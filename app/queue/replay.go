package queue

import (
	"log"
	"sync"
	"time"
)

// You don’t have unlimited resource on your machine, the minimal size of a goroutine object is 2 KB,
// when you spawn too many goroutine, your machine will quickly run out of memory and the CPU will keep processing
// the task until it reach the limit. By using limited pool of workers and keep the task on the queue,
// we can reduce the burst of CPU and memory since the task will wait on the queue until the the worker pull the task.
type Status struct {
	Order             		Order
	Inserted            	bool // 1. Order data inserted to queue
	Published           	bool // 2. Pub/Sub publishing
}

type PutRequest struct {
	Order  				Order
	ResponseChan 	chan bool
}

// Request - Any request to be placed into queue's channels in this form
// client can also receive response/ confirmation over channel that they specify
type Request struct {
	Order  				string
	ResponseChan 	chan bool
}


// Next - Order to be processed next, asked by sending this request
type Next struct {
	ResponseChan chan struct {
		Status bool
		Order  string
		Time   int64
		EOF    bool
	}
}

// ReplayQueue - concurrent safe queue to be interacted with before attempting to replay any order
type ReplayQueue struct {
	Orders                map[string]*Status
	PutChan               chan PutRequest
	CanPublishChan        chan Request
	PublishedChan         chan Request
	PublishNextChan    		chan Next
	CleanUpChan           chan bool
	stopChannel 					chan string
	mutex 								*sync.RWMutex
	stopped               bool
}

// New - Getting new instance of queue, to be invoked during setting up application
func NewReplayQueue() *ReplayQueue {

	return &ReplayQueue{
		Orders:                make(map[string]*Status),
		PutChan:               make(chan PutRequest, 128),
		CanPublishChan:        make(chan Request, 128),
		PublishedChan:         make(chan Request, 128),
		PublishNextChan:     	 make(chan Next, 1),
		stopChannel:           make(chan string, 1),
		mutex:                 &sync.RWMutex{},
	}

}

// Put - Client is supposed to be invoking this method
// when it's interested in putting new order to processing queue
func (q *ReplayQueue) Put(order Order) bool {

	resp := make(chan bool)
	req := PutRequest{
		Order:  			order,
		ResponseChan: resp,
	}

	q.PutChan <- req

	return <-resp

}

func (q *ReplayQueue) CanPublish(order string) bool {

	resp := make(chan bool)
	req := Request{
		Order:  			order,
		ResponseChan: resp,
	}

	q.CanPublishChan <- req
	return <-resp

}

// Published - Asks queue manager to mark that this order has been
// successfully published on Pub/Sub topic
//
// Future order processing attempts (if any), are supposed to be
// avoiding doing this, if already done successfully
func (q *ReplayQueue) Published(order string) bool {

	resp := make(chan bool)
	req := Request{
		Order:  			order,
		ResponseChan: resp,
	}

	q.PublishedChan <- req
	return <-resp

}

// PublishNext - Next order that can be published
func (q *ReplayQueue) PublishNext() (string, int64, bool, bool) {

	resp := make(chan struct {
		Status bool
		Order  string
		Time   int64
		EOF    bool
	})
	req := Next{ResponseChan: resp}

	q.PublishNextChan <- req

	v := <-resp
	return v.Order, v.Time, v.EOF, v.Status

}

func (q *ReplayQueue) CleanUp() {
	q.CleanUpChan <- true
}

// Stop - You're supposed to be stopping this method as an
// independent go routine
func (q *ReplayQueue) Stop() {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if q.stopped {
		return
	}

	q.stopped = true

	q.stopChannel <- "Stop"
	<-q.stopChannel

}

// Start - You're supposed to be starting this method as an
// independent go routine, with will listen on multiple channels
// & respond back over provided channel ( by client )
func (q *ReplayQueue) Start() {

	for {
		select {

		case <-q.stopChannel:

			log.Println("Stopping replay queue")

			q.stopChannel <- "Stop"
			return

		case req := <-q.PutChan:

			// Once a order is inserted into processing queue, don't overwrite its history with some new request
			if _, ok := q.Orders[req.Order.ID()]; ok {

				req.ResponseChan <- false
				break

			}

			q.Orders[req.Order.ID()] = &Status { Order: req.Order, Inserted: true }
			req.ResponseChan <- true

		case req := <-q.CanPublishChan:

			order, ok := q.Orders[req.Order]
			if !ok {
				req.ResponseChan <- false
				break
			}

			req.ResponseChan <- !order.Published

		case req := <-q.PublishedChan:
			// Worker go rountine marks this order has been published
			//
			// If not, it'll be marked so & no future attempt
			// should try to publish it again over Pub/Sub

			order, ok := q.Orders[req.Order]
			if !ok {
				req.ResponseChan <- false
				break
			}

			order.Published = true
			req.ResponseChan <- true

		case nxt := <-q.PublishNextChan:

			var selected string
			var found bool

			var min int64
			// TODO: use min heap for this
			for k := range q.Orders {

				if q.Orders[k].Published {
					delete(q.Orders, k)
					continue
				}

				if min == 0 {
					min = q.Orders[k].Order.ExecuteTime
				}

				if q.Orders[k].Order.ExecuteTime <= time.Now().UnixMicro() {
					if min >= q.Orders[k].Order.ExecuteTime || (min == q.Orders[k].Order.ExecuteTime && k < selected) {
						selected = k
						min = q.Orders[k].Order.ExecuteTime
						found = true
					}
				}

			}

			if !found {

				nxt.ResponseChan <- struct {
					Status 	bool
					Order 	string
					Time  	int64
					EOF   	bool
				}{
					Status: false,
				}
				break

			}

			nxt.ResponseChan <- struct {
				Status 	bool
				Order 	string
				Time   	int64
				EOF   	bool
			}{
				Status: true,
				Order: selected,
				Time: min,
				EOF: q.Orders[selected].Order.EOF,
			}

		}
	}

}
