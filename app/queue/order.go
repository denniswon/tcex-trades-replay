package queue

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	d "github.com/denniswon/tcex/app/data"
	ps "github.com/denniswon/tcex/app/pubsub"
	"github.com/go-redis/redis/v8"
)

type Order struct {
	RequestId 		string
	OrderNumber 	uint64
	ExecuteTime 	int64
}

func (o *Order) String() string {
	return fmt.Sprintf(`{"request_id":%s,"order_number":%d,"execute_time":%d}`,
		o.RequestId,
		o.OrderNumber,
		o.ExecuteTime,
	)
}

func (o *Order) ID() string {
	return fmt.Sprintf("%s:%d", o.RequestId, o.OrderNumber)
}


type Error struct {
	RequestId 	string
	Err   	 		error
}

type FileRef struct {
	File *os.File
	RC   uint64
}

// Client defines typed wrappers for the Ethereum RPC API.
type OrderQueue struct {
	stopped     			bool
	requests 					map[string]*ps.SubscriptionRequest
	files 						map[string]*FileRef
	requestChannel 		chan string
	orderChannel 			chan Order
	stopChannel 			chan string
	errorChannel		 	chan Error
	redis    					*redis.Client
	mutex 						*sync.RWMutex
}

// NewClient creates a client that uses the given RPC client.
func NewOrderQueue(_redis *redis.Client) *OrderQueue {
	client := &OrderQueue{
		stopped:					false,
		stopChannel: 			make(chan string),
		errorChannel: 		make(chan Error),
		orderChannel: 		make(chan Order),
		requestChannel: 	make(chan string),
		files: 						make(map[string]*FileRef),
		requests: 				make(map[string]*ps.SubscriptionRequest),
		redis:    				_redis,
		mutex:      			&sync.RWMutex{},
	}
	return client
}

func (q *OrderQueue) Put(request *ps.SubscriptionRequest) bool {

	if !request.Validate() {
		q.Error(request.ID, fmt.Errorf("invalid request"))
		return false
	}

	q.requests[request.ID] = request
	q.requestChannel <- request.ID

	return true
}

func (q *OrderQueue) Remove(requestId string) {

	q.mutex.RLock()
	defer q.mutex.RUnlock()

 	request := q.requests[requestId]
	if request != nil {
		delete(q.requests, requestId)
	}

	if q.files[requestId] == nil {
		return
	}

	if q.files[requestId].RC == 1 {
		q.files[requestId].File.Close()
		delete(q.files, requestId)
	} else {
		q.files[requestId].RC--
	}

}

func (q *OrderQueue) Start(ctx context.Context) {

	log.Println("Worker Started")

	for {
		select {

		case <-ctx.Done():
			q.Close()
			return

		case requestId := <-q.requestChannel:
			if err := q.HandleRequest(requestId); err != nil {
				q.Error(requestId, err)
			}

		case <-q.stopChannel:

			log.Println("Order Client Stopped.")

			q.stopChannel <- "Stop"
			return

		}
	}
}

func (q *OrderQueue) HandleRequest(requestId string) error {
	request, ok := q.requests[requestId]

	if !ok {
		return fmt.Errorf("missing request for request id: %s", requestId)
	}

	log.Printf("Reading input file for request id: %s\n", request.String())

	if q.files[request.Filename] == nil {
		file, err := os.Open(request.Filename)
		if err != nil {
			log.Fatalf("Error opening file: %s\n", err.Error())
			return err
		}
		q.files[request.Filename] = &FileRef{
			File: file,
			RC:   1,
		}
	} else {
		q.files[request.Filename].RC++
	}

	return q.Run(request)
}

func (q *OrderQueue) Run(request *ps.SubscriptionRequest) error {
	if q.files[request.Filename] == nil {
		return fmt.Errorf("missing file: %s", request.Filename)
	}

	fref := q.files[request.Filename]
	scanner := bufio.NewScanner(fref.File)

	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		return advance, token, err
	}
	scanner.Split(scanLines)

	var orderNumber uint64 = 0
	var currTime int64 = time.Now().Unix()
	var indexTime int64 = 0
	keys := []string{}
	values := [][]byte{}
	orders := []Order{}

	for scanner.Scan() {

		if q.IsStopped() {
			return nil
		}

		var order d.Order
		err := json.Unmarshal([]byte(scanner.Text()), &order)
		if err != nil {
			log.Fatalf("Failed to decode order data to JSON : %s\n", err.Error())
			return err
		}

		keys = append(keys, fmt.Sprintf("%s:%d", request.ID, orderNumber))
		values = append(values, order.ToJSON())

		if indexTime == 0 {
			indexTime = order.Timestamp
		}

		orders = append(orders, Order {
			RequestId: request.ID,
			OrderNumber: orderNumber,
			ExecuteTime: currTime + (order.Timestamp - indexTime),
		})

		orderNumber++

	}

	_, err := q.redis.MSet(context.Background(), keys, values, 0).Result()
	if err != nil {
		log.Fatalf("Failed to cache order for request %s order number %d: %s\n",
			request.ID, orderNumber, err.Error(),
		)
		return err
	}

	for _, order := range orders {
		q.orderChannel <- order
	}

	return nil
}

func (q *OrderQueue) Stop() {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if q.stopped {
		return
	}

	q.stopped = true

	q.stopChannel <- "Stop"
	<-q.stopChannel
}

func (q *OrderQueue) Close() {
	if !q.IsStopped() {
		q.Stop()
	}

	q.mutex.RLock()
	defer q.mutex.RUnlock()

	for k := range q.requests {
		delete(q.requests, k)
	}

	for k := range q.files {
		_ = q.files[k].File.Close()
		delete(q.files, k)
	}

	close(q.stopChannel)
	close(q.errorChannel)
	close(q.orderChannel)
	close(q.requestChannel)
}

func (q *OrderQueue) Error(requestId string, err error) {
	request, ok := q.requests[requestId]
	if ok {
		delete(q.requests, requestId)
		delete(q.files, request.Filename)
	}

	q.errorChannel <- Error {
		RequestId: requestId,
		Err: err,
	}
}

func (q *OrderQueue) Order() chan Order {
	return q.orderChannel
}

func (q *OrderQueue) Err() chan Error {
	return q.errorChannel
}

func (q *OrderQueue) IsStopped() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return q.stopped
}
