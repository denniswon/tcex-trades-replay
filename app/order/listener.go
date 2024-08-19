package order

import (
	"log"
	"runtime"

	c "github.com/denniswon/tcex/app/client"
	cfg "github.com/denniswon/tcex/app/config"
	q "github.com/denniswon/tcex/app/queue"
	d "github.com/denniswon/tcex/app/redis"
	"github.com/gammazero/workerpool"
)

// SubscribeToNewOrders - Listen for new order header available, then fetch order content
// including all transactions in different worker
func SubscribeToNewOrders(orderClient *c.OrderClient, redis *d.RedisInfo, queue *q.OrderReplayQueue) {
	ordersChan := make(chan d.Orders)

	orderClient.Start(ordersChan)
	defer orderClient.Stop()

	// Creating a job queue of size `#-of CPUs present in machine` * concurrency factor where order fetching requests are submitted to
	// There is no upper limit on the number of tasks queued, other than the limits of system resources
	// If the number of inbound tasks is too many to even queue for pending processing, then we should distribute workload over multiple systems,
	// and/or storing input for pending processing in intermediate storage such as a distributed message queue, etc.
	wp := workerpool.New(runtime.NumCPU() * int(cfg.GetConcurrencyFactor()))
	defer wp.Stop()

	for {
		select {
		case err := <-orderClient.Err():

			log.Fatalf("Order client stopped : %s\n", err.Error())

		case orders := <-ordersChan:

			// Receive orders as a batch and that job submitted in job queue
			//
			// Putting it in a different function scope so that job submitter gets its own copy of orders,
			// otherwise it might get wrong info, if new orders batch is received very soon & this job is not yet submitted
			func(orders d.Orders, _queue *q.OrderReplayQueue) {

				// Next order which can be attempted to be checked
				// while finally considering it confirmed & put into DB
				if nxt, ok := _queue.PublishedNext(); ok {

					log.Printf("Processing order %d\n", nxt)

					// Note, we are taking `next` variable's copy in local scope of closure, so that during
					// iteration over queue elements, none of them get missed, becuase in a concurrent system,
					// previous `next` can be overwritten by new `next` & we can end up missing a order
					func(_oldestOrder uint64, _queue *q.OrderReplayQueue) {

						wp.Submit(func() {

							ProcessOrders(order, redis, queue)

						})

					}(nxt, _queue)

				}

				wp.Submit(func() {

					_queue.Put(&orders)

				})

			}(orders, queue)

		}
	}
}
