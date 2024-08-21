package order

import (
	"context"
	"encoding/json"
	"log"

	d "github.com/denniswon/tcex/app/data"
	q "github.com/denniswon/tcex/app/queue"
	"github.com/gammazero/workerpool"
	"github.com/go-redis/redis/v8"
)

// ProcessOrderReplays
func ProcessOrderReplays(ctx context.Context, requestQueue *q.RequestQueue, replayQueue *q.ReplayQueue, redis *redis.Client) {

	orderChan := make(chan q.Order)

	// first start the request queue as a separate go routine
	go requestQueue.Start(orderChan)

	// second start the replay queue as a separate go routine
	go replayQueue.Start()

	// TODO create a job queue of size `#-of CPUs present in machine` * concurrency factor for worker pool

	// There is no upper limit on the number of tasks queued, other than the limits of system resources
	// If the number of inbound tasks is too many to even queue for pending processing, then we should distribute workload over multiple systems,
	// and/or storing input for pending processing in intermediate storage such as a distributed message queue, etc.
	wp := workerpool.New(/* runtime.NumCPU() * int(cfg.GetConcurrencyFactor()) */ 1)
	defer wp.Stop()

	for {

		wp.Submit(func() {

			for {

				select {

				case <-ctx.Done():

					log.Println("Exiting order replay publisher")

					requestQueue.Close()
					replayQueue.Stop()
					if orderChan != nil {
						close(orderChan)
						orderChan = nil
					}

					return

				case order := <-orderChan:

					log.Println("Submitting order for replay", order)

					replayQueue.Put(order)

				default:
					if order, extime, ok := replayQueue.PublishNext(); ok {
						// retrieve the cached order data
						encoded, err := redis.Get(context.Background(), order).Result()
						if err != nil {
							log.Fatalf("Failed to retrieve cached order %s : %s\n", order, err.Error())
							continue
						}

						_order := d.Order{}
						err = json.Unmarshal([]byte(encoded), &_order)
						if err != nil {
							log.Fatalf("Failed to unmarshal cached order %s : %s\n", order, err.Error())
							continue
						}

						log.Printf(
							"Publishing order %s at time %d (order timestamp : %d)\n",
							order, extime, _order.Timestamp,
						)

						if ok := PublishReplayOrder(order, &_order, replayQueue, redis); !ok {
							log.Fatalf("Failed to publish replay order %s : %s\n", order, err.Error())
							continue
						}

						res, err := redis.Del(context.Background(), order).Result()
						if err != nil {
							log.Printf("[!%d] Failed to delete cached order %s from redis : %s\n", res, order, err.Error())
						}
					}

				}

			}

		})

	}
}
