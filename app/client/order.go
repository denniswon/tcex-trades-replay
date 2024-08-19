package client

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	d "github.com/denniswon/tcex/app/redis"
)

type Request struct {
	RequestId 	string
	Filename 		string
	ReplayRate 	uint64
}

func (b *Request) Validate() bool {
	return b.Filename != "" && b.ReplayRate > 0 && b.RequestId != ""
}

func (b *Request) String() string {
	return fmt.Sprintf(`{"request_id":%s,"filename":%s,"replay_rate":%d}`,
		b.RequestId,
		b.Filename,
		b.ReplayRate,
	)
}

type RequestError struct {
	RequestId string
	Error   	  error
}

// Client defines typed wrappers for the Ethereum RPC API.
type OrderClient struct {
	stopped     			bool
	requests 					map[string]Request
	files 						map[string]*os.File
	requestChannel 		chan string
	orderChannel 			chan d.Order
	stopChannel 			chan string
	errorChannel		 	chan RequestError
	mutex 						*sync.RWMutex
}

// NewClient creates a client that uses the given RPC client.
func NewOrderClient() *OrderClient {
	client := &OrderClient{
		stopped:					false,
		stopChannel: 			make(chan string),
		errorChannel: 		make(chan RequestError),
		orderChannel: 		make(chan d.Order, 128),
		requests: 				make(map[string]Request),
		requestChannel: 	make(chan string, 128),
		mutex:      			&sync.RWMutex{},
	}
	return client
}

func (c *OrderClient) Put(request Request) bool {
	if c.IsStopped()|| !request.Validate() {
		return false
	}
	c.requests[request.RequestId] = request
	c.requestChannel <- request.RequestId
	return true
}

func (c *OrderClient) Run() {

	log.Println("Worker Started")

	for {
		select {
		case requestId := <-c.requestChannel:
			if err := c.HandleRequest(requestId); err != nil {
				c.Err(requestId, err)
			}
			// This breaks out of the select, not the for loop.
			break
		case <-c.stopChannel:
			c.stopChannel <- "Stop"
			return
		}
	}
}

func (c *OrderClient) HandleRequest(requestId string) error {
	request, ok := c.requests[requestId]

	if !ok {
		return errors.New(fmt.Sprintf("Missing request for request id: %s", requestId))
	}

	log.Printf("Reading input file for request id: %s\n", request.String())

	if c.files[request.Filename] == nil {
		file, err := os.Open(request.Filename)
		if err != nil {
			log.Fatalf("Error opening file: %s\n", err.Error())
			return err
		}
		c.files[request.Filename] = file
	}

	return c.ProcessRequest(request)
}

func (c *OrderClient) ProcessRequest(request Request) error {
	if c.files[request.Filename] == nil {
		return errors.New(fmt.Sprintf("Missing file: %s", request.Filename))
	}

	file := c.files[request.Filename]
	scanner := bufio.NewScanner(file)

	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		return advance, token, err
	}
	scanner.Split(scanLines)

	for scanner.Scan() {

		if c.IsStopped() {
			return nil
		}

		var order d.Order
		err := json.Unmarshal([]byte(scanner.Text()), &order)
		if err != nil {
			log.Printf("Failed to decode order data to JSON : %s\n", err.Error())
			return err
		}

		order.RequestId = request.RequestId
		c.orderChannel <- order
	}

	return nil
}

func (c *OrderClient) Stop() {
	if c.IsStopped() {
		return
	}

	c.stopped = true

	c.stopChannel <- "Stop"
	<-c.stopChannel
}

func (c *OrderClient) Close() {
	if !c.IsStopped() {
		return
	}


	for k := range c.requests {
		delete(c.requests, k)
	}

	for k := range c.files {
		_ = c.files[k].Close()
		delete(c.files, k)
	}

	close(c.stopChannel)
	close(c.errorChannel)
	close(c.orderChannel)
	close(c.requestChannel)
}

func (c *OrderClient) Err(requestId string, err error) {
	request, ok := c.requests[requestId]
	if !ok {
		return
	}

	delete(c.requests, requestId)
	delete(c.files, request.Filename)

	c.errorChannel <- RequestError {
		RequestId: requestId,
		Error: err,
	}
}

func (c *OrderClient) OrderChannel() chan d.Order {
	return c.orderChannel
}

func (c *OrderClient) ErrorChannel() chan RequestError {
	return c.errorChannel
}

func (c *OrderClient) IsStopped() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.stopped
}
