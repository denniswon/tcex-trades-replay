package queue

import (
	"context"
	"time"

	d "github.com/denniswon/tcex/app/data"
)

// You don’t have unlimited resource on your machine, the minimal size of a goroutine object is 2 KB,
// when you spawn too many goroutine, your machine will quickly run out of memory and the CPU will keep processing
// the task until it reach the limit. By using limited pool of workers and keep the task on the queue,
// we can reduce the burst of CPU and memory since the task will wait on the queue until the the worker pull the task.

// Order - Keeps track of single order i.e. how many
// times attempted till date, last attempted to process
// whether order data has been published on pubsub topic or not,
// is order processing currently
type Order struct {
	Inserted            	bool // 1. Order data inserted to queue
	Published           	bool // 2. Pub/Sub publishing
	ExecuteTime       		time.Time
}

// Request - Any request to be placed into queue's channels in this form
// client can also receive response/ confirmation over channel that they specify
type Request struct {
	Order  				uint64
	Timestamp 		time.Time
	ResponseChan 	chan bool
}

// Next - Order to be processed next, asked by sending this request
// when receptor detects so, will attempt to find out what should be next processed and
// send that order number is response over channel specified by client
type Next struct {
	ResponseChan chan struct {
		Status bool
		Order  uint64
	}
}

// Stat - Clients can query how many orders present in queue currently
type Stat struct {
	ResponseChan chan StatResponse
}

// StatResponse - Statistics of queue to be responded back to client in this form
type StatResponse struct {
	Published 					uint64
	Inserted   					uint64
	Total               uint64
}

// OrderReplayQueue - concurrent safe queue to be interacted with before attempting to process any order
type OrderReplayQueue struct {
	Orders                map[uint64]*Order
	TotalInserted         uint64
	Total                 uint64
	PutChan               chan Request
	CanPublishChan        chan Request
	PublishedChan         chan Request
	InsertedChan          chan Request
	StatChan              chan Stat
	PublishedNextChan     chan Next
	Delay                 time.Duration
}

// New - Getting new instance of queue, to be invoked during setting up application
func New(startingWith uint64) *OrderReplayQueue {

	return &OrderReplayQueue{
		Orders:                make(map[uint64]*Order),
		TotalInserted:         0,
		Total:                 0,
		PutChan:               make(chan Request, 128),
		InsertedChan:          make(chan Request, 128),
		CanPublishChan:        make(chan Request, 128),
		PublishedChan:         make(chan Request, 128),
		StatChan:              make(chan Stat, 1),
		PublishedNextChan:     make(chan Next, 1),
		Delay:                 0 * time.Millisecond,
	}

}

// Put - Client is supposed to be invoking this method
// when it's interested in putting new order to processing queue
//
// If responded with `true`, they're good to go with execution of
// processing of this order
//
// If this order is already put into queue, it'll ask client
// to not proceed with this number
func (b *OrderReplayQueue) Put(orders *d.Orders) {

	first := true
	for _, order := range orders.Orders {
		if first {

			first = false
		}
		if b.Orders[order.Number] != nil {
			continue
		}

		resp := make(chan bool)
		req := Request{
			Order:  			order.Number,
			ResponseChan: resp,
		}

		b.PutChan <- req
	}

}

// CanPublish - Before any client attempts to publish any order
// on Pub/Sub topic, they're supposed to be invoking this method
// to check whether they're eligible of publishing or not
//
// NOTE: if any other client has already published it, we'll avoid redoing it
func (b *OrderReplayQueue) CanPublish(orderNumber uint64) bool {

	resp := make(chan bool)
	req := Request{
		Order:  			orderNumber,
		ResponseChan: resp,
	}

	b.CanPublishChan <- req
	return <-resp

}

// Published - Asks queue manager to mark that this order has been
// successfully published on Pub/Sub topic
//
// Future order processing attempts (if any), are supposed to be
// avoiding doing this, if already done successfully
func (b *OrderReplayQueue) Published(orderNumber uint64) bool {

	resp := make(chan bool)
	req := Request{
		Order:  			orderNumber,
		ResponseChan: resp,
	}

	b.PublishedChan <- req
	return <-resp

}

// Inserted - Marking this order has been inserted into DB (not updation, it's insertion)
func (b *OrderReplayQueue) Inserted(orderNumber uint64) bool {

	resp := make(chan bool)
	req := Request{
		Order:  			orderNumber,
		ResponseChan: resp,
	}

	b.InsertedChan <- req
	return <-resp

}

// Stat - Client's are supposed to be invoking this abstracted method
// for checking queue status
func (b *OrderReplayQueue) Stat() StatResponse {

	resp := make(chan StatResponse)
	req := Stat{ResponseChan: resp}

	b.StatChan <- req
	return <-resp

}

// PublishedNext - Next order that can be published
func (b *OrderReplayQueue) PublishedNext() (uint64, bool) {

	resp := make(chan struct {
		Status bool
		Order  uint64
	})
	req := Next{ResponseChan: resp}

	b.PublishedNextChan <- req

	v := <-resp
	return v.Order, v.Status

}

// CanBeConfirmed - Checking whether given order number has reached
// finality as per given user set preference, then it can be attempted
// to be checked again & finally entered into storage
func (b *OrderReplayQueue) CanBeConfirmed(orderNumber uint64) bool {

	order, ok := b.Orders[orderNumber]
	if !ok {
		return false
	}

	return time.Now().UTC().After(order.TimeInserted.Add(b.Delay))

}

// Start - You're supposed to be starting this method as an
// independent go routine, with will listen on multiple channels
// & respond back over provided channel ( by client )
func (b *OrderReplayQueue) Start(ctx context.Context) {

	for {
		select {

		case <-ctx.Done():
			return

		case req := <-b.PutChan:

			// Once a order is inserted into processing queue, don't overwrite its history with some new request
			if _, ok := b.Orders[req.Order]; ok {

				req.ResponseChan <- false
				break

			}

			b.Orders[req.Order] = &Order{
				Inserted: 					true,
				Published: 					false,
				TimeInserted:       time.Now().UTC(),
				Delay:              00,
			}
			req.ResponseChan <- true

		case req := <-b.CanPublishChan:

			order, ok := b.Orders[req.Order]
			if !ok {
				req.ResponseChan <- false
				break
			}

			req.ResponseChan <- !order.Published

		case req := <-b.PublishedChan:
			// Worker go rountine marks this order has been published
			//
			// If not, it'll be marked so & no future attempt
			// should try to publish it again over Pub/Sub

			order, ok := b.Orders[req.Order]
			if !ok {
				req.ResponseChan <- false
				break
			}

			order.Published = true
			req.ResponseChan <- true

		case req := <-b.InsertedChan:
			// Increments how many orders were inserted into DB

			order, ok := b.Orders[req.Order]
			if !ok {
				req.ResponseChan <- false
				break
			}

			b.TotalInserted++

			order.ConfirmedDone = b.CanBeConfirmed(req.Order)

			order.ResetDelay()
			order.SetTimeInserted()

			req.ResponseChan <- true

			req.ResponseChan <- true

		case req := <-b.ConfirmedFailedChan:

			order, ok := b.Orders[req.Order]
			if !ok {
				req.ResponseChan <- false
				break
			}

			order.ConfirmedProgress = false
			order.SetDelay()

			req.ResponseChan <- true

		case req := <-b.ConfirmedDoneChan:

			order, ok := b.Orders[req.Order]
			if !ok {
				req.ResponseChan <- false
				break
			}

			order.ConfirmedProgress = false
			order.ConfirmedDone = true

			req.ResponseChan <- true

		case nxt := <-b.UnconfirmedNextChan:

			// This is the order number which should be processed by requester client
			var selected uint64
			var found bool

			for k := range b.Orders {

				if b.Orders[k].ConfirmedDone || b.Orders[k].ConfirmedProgress {
					continue
				}

				if b.Orders[k].UnconfirmedDone || b.Orders[k].UnconfirmedProgress {
					continue
				}

				if b.Orders[k].CanAttempt() {
					selected = k
					found = true

					break
				}

			}

			if !found {

				// As we've failed to find any order which can be processed now
				// ask client to try again later. When to come back is upto client
				nxt.ResponseChan <- struct {
					Status bool
					Number uint64
				}{
					Status: false,
				}
				break

			}

			// Updated when last this order was attempted to be processed
			b.Orders[selected].SetTimeInserted()
			b.Orders[selected].UnconfirmedProgress = true

			// Asking client to proceed with processing of this order
			nxt.ResponseChan <- struct {
				Status bool
				Number uint64
			}{
				Status: true,
				Number: selected,
			}

		case nxt := <-b.ConfirmedNextChan:

			var selected uint64
			var found bool

			for k := range b.Orders {

				if b.Orders[k].ConfirmedDone || b.Orders[k].ConfirmedProgress {
					continue
				}

				if !b.Orders[k].UnconfirmedDone {
					continue
				}

				if b.Orders[k].CanAttempt() && b.CanBeConfirmed(k) {
					selected = k
					found = true

					break
				}

			}

			if !found {

				nxt.ResponseChan <- struct {
					Status bool
					Number uint64
				}{
					Status: false,
				}
				break

			}

			b.Orders[selected].SetTimeInserted()
			b.Orders[selected].ConfirmedProgress = true

			nxt.ResponseChan <- struct {
				Status bool
				Number uint64
			}{
				Status: true,
				Number: selected,
			}

		case req := <-b.StatChan:

			// Returning back how many orders currently living
			// in order processor queue & in what state
			var stat StatResponse

			for k := range b.Orders {

				if b.Orders[k].UnconfirmedProgress {
					stat.UnconfirmedProgress++
					continue
				}

				if b.Orders[k].UnconfirmedProgress == b.Orders[k].UnconfirmedDone {
					stat.UnconfirmedWaiting++
					continue
				}

				if b.Orders[k].ConfirmedProgress {
					stat.ConfirmedProgress++
					continue
				}

				if b.Orders[k].ConfirmedProgress == b.Orders[k].ConfirmedDone {
					stat.ConfirmedWaiting++
					continue
				}

			}

			stat.Total = b.Total
			req.ResponseChan <- stat

		case <-time.After(time.Duration(100) * time.Millisecond):

			// Finding out which orders are confirmed & we're good to clean those up
			for k := range b.Orders {

				if b.Orders[k].ConfirmedDone {
					delete(b.Orders, k)
					b.Total++ // Successfully processed #-of orders
				}

			}

		}
	}

}
