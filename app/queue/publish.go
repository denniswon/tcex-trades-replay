package queue

import (
	"context"
	"time"
)

// You don’t have unlimited resource on your machine, the minimal size of a goroutine object is 2 KB,
// when you spawn too many goroutine, your machine will quickly run out of memory and the CPU will keep processing
// the task until it reach the limit. By using limited pool of workers and keep the task on the queue,
// we can reduce the burst of CPU and memory since the task will wait on the queue until the the worker pull the task.
type Status struct {
	Order             		Order
	Inserted            	bool // 1. Order data inserted to queue
	CanPublish           	bool // 2. Replay the order
	Published           	bool // 3. Pub/Sub publishing
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
	}
}

// PublishQueue - concurrent safe queue to be interacted with before attempting to replay any order
type PublishQueue struct {
	Orders                map[string]*Status
	PutChan               chan PutRequest
	CanPublishChan        chan Request
	PublishedChan         chan Request
	PublishNextChan    		chan Next
}

// New - Getting new instance of queue, to be invoked during setting up application
func NewPublishQueue() *PublishQueue {

	return &PublishQueue{
		Orders:                make(map[string]*Status),
		PutChan:               make(chan PutRequest, 128),
		CanPublishChan:        make(chan Request, 128),
		PublishedChan:         make(chan Request, 128),
		PublishNextChan:     	 make(chan Next, 1),
	}

}

// Put - Client is supposed to be invoking this method
// when it's interested in putting new order to processing queue
func (q *PublishQueue) Put(order Order) bool {

	resp := make(chan bool)
	req := PutRequest{
		Order:  			order,
		ResponseChan: resp,
	}

	q.PutChan <- req

	return <-resp

}

// CanPublish - Before any client attempts to publish any order
// on Pub/Sub topic, they're supposed to be invoking this method
// to check whether they're eligible of publishing or not
//
// NOTE: if any other client has already published it, we'll avoid redoing it
func (q *PublishQueue) CanPublish(order string) bool {

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
func (q *PublishQueue) Published(order string) bool {

	resp := make(chan bool)
	req := Request{
		Order:  			order,
		ResponseChan: resp,
	}

	q.PublishedChan <- req
	return <-resp

}

// PublishNext - Next order that can be published
func (q *PublishQueue) PublishNext() (string, bool) {

	resp := make(chan struct {
		Status bool
		Order  string
	})
	req := Next{ResponseChan: resp}

	q.PublishNextChan <- req

	v := <-resp
	return v.Order, v.Status

}

// Start - You're supposed to be starting this method as an
// independent go routine, with will listen on multiple channels
// & respond back over provided channel ( by client )
func (q *PublishQueue) Start(ctx context.Context) {

	for {
		select {

		case <-ctx.Done():
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

		case req := <-q.CanPublishChan:

			order, ok := q.Orders[req.Order]
			if !ok {
				req.ResponseChan <- false
				break
			}

			order.CanPublish = true
			req.ResponseChan <- true

		case nxt := <-q.PublishNextChan:

			var selected string
			var found bool

			for k := range q.Orders {

				if q.Orders[k].Published || !q.Orders[k].CanPublish {
					continue
				}

				if q.Orders[k].Order.ExecuteTime > time.Now().Unix() {
					selected = k
					found = true

					break
				}

			}

			if !found {

				nxt.ResponseChan <- struct {
					Status bool
					Order string
				}{
					Status: false,
				}
				break

			}

			nxt.ResponseChan <- struct {
				Status bool
				Order string
				}{
				Status: true,
				Order: selected,
			}

		case <-time.After(time.Duration(100) * time.Millisecond):

			// Finding out which orders are confirmed & we're good to clean those up
			for k := range q.Orders {

				if q.Orders[k].Published {
					delete(q.Orders, k)
				}

			}

		}
	}

}
