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


type RequestError struct {
	RequestId 	string
	Err   	 		error
}

func (m *RequestError) Error() string {
	return m.Err.Error()
}

type FileRef struct {
	File *os.File
	RC   uint64
}

// Client defines typed wrappers for the Ethereum RPC API.
type RequestQueue struct {
	stopped     			bool
	requests 					map[string]*ps.SubscriptionRequest
	files 						map[string]*FileRef
	requestChannel 		chan string
	stopChannel 			chan string
	orderChannel 			chan Order
	errorChannel		 	chan RequestError
	redis    					*redis.Client
	mutex 						*sync.RWMutex
}

// NewClient creates a client that uses the given RPC client.
func NewRequestQueue(_redis *redis.Client) *RequestQueue {
	client := &RequestQueue{
		stopped:					false,
		stopChannel:      make(chan string, 1),
		errorChannel: 		make(chan RequestError),
		requestChannel: 	make(chan string),
		files: 						make(map[string]*FileRef),
		requests: 				make(map[string]*ps.SubscriptionRequest),
		redis:    				_redis,
		mutex:      			&sync.RWMutex{},
	}
	return client
}

func (q *RequestQueue) Put(request *ps.SubscriptionRequest) bool {

	if !request.Validate() {
		q.Error(request.ID, fmt.Errorf("invalid request"))
		return false
	}

	q.requests[request.ID] = request
	q.requestChannel <- request.ID

	return true
}

func (q *RequestQueue) Remove(requestId string) {

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

func (q *RequestQueue) Start(orderChannel chan Order) {

	log.Println("Request queue started")
	q.orderChannel = orderChannel

	for {
		select {

		case requestId := <-q.requestChannel:
			if err := q.HandleRequest(requestId); err != nil {
				q.Error(requestId, err)
			}

		case <-q.stopChannel:

			log.Println("Stopping request queue")

			q.stopChannel <- "Stop"
			return

		}
	}
}

func (q *RequestQueue) HandleRequest(requestId string) error {
	request, ok := q.requests[requestId]

	if !ok {
		return fmt.Errorf("missing request for request id : %s", requestId)
	}

	log.Printf("Reading input file for request id : %s\n", request.String())

	if q.files[request.Filename] == nil {
		file, err := os.Open(request.Filename)
		if err != nil {
			log.Fatalf("Error opening file : %s\n", err.Error())
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

func (q *RequestQueue) Run(request *ps.SubscriptionRequest) error {
	if q.files[request.Filename] == nil {
		return fmt.Errorf("missing file : %s", request.Filename)
	}

	fref := q.files[request.Filename]
	fref.File.Seek(0, 0)
	scanner := bufio.NewScanner(fref.File)

	var orderNumber uint64 = 0
	var currTime int64 = time.Now().UnixMilli()
	var indexTime int64 = 0
	var pairs []interface{}
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

		pairs = append(pairs, fmt.Sprintf("%s:%d", request.ID, orderNumber))
		pairs = append(pairs, order.ToJSON())

		if indexTime == 0 {
			indexTime = order.Timestamp
		}

		log.Printf("!!!!!!! %f %d\n", float32(order.Timestamp - indexTime) / request.ReplayRate, int64(float32(order.Timestamp - indexTime) / request.ReplayRate))
		orders = append(orders, Order {
			RequestId: request.ID,
			OrderNumber: orderNumber,
			ExecuteTime: currTime + int64(float32(order.Timestamp - indexTime) / request.ReplayRate),
		})

		orderNumber++

	}

	if len(pairs) > 0 {
		_, err := q.redis.MSet(context.Background(), pairs...).Result()
		if err != nil {
			log.Fatalf("Failed to cache order for request %s order number %d : %s\n",
				request.ID, orderNumber, err.Error(),
			)
			return err
		}
	}

	for _, order := range orders {
		q.orderChannel <- order
	}

	return nil
}

func (q *RequestQueue) Stop() {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if q.stopped {
		return
	}

	q.stopped = true
	q.orderChannel = nil

	q.stopChannel <- "Stop"
	<-q.stopChannel
}

func (q *RequestQueue) Close() {
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

	if q.stopChannel != nil {
		close(q.stopChannel)
		q.stopChannel = nil
	}
	if q.errorChannel != nil {
		close(q.errorChannel)
		q.errorChannel = nil
	}
	if q.requestChannel != nil {
		close(q.requestChannel)
		q.requestChannel = nil
	}
}

func (q *RequestQueue) Error(requestId string, err error) {
	request, ok := q.requests[requestId]
	if ok {
		delete(q.requests, requestId)
		delete(q.files, request.Filename)
	}

	q.errorChannel <- RequestError {
		RequestId: requestId,
		Err: err,
	}
}

func (q *RequestQueue) Err() chan RequestError {
	return q.errorChannel
}

func (q *RequestQueue) IsStopped() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return q.stopped
}
