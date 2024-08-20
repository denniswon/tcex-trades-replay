package order

import (
	"context"
	"encoding/json"
	"log"
	"runtime"

	cfg "github.com/denniswon/tcex/app/config"
	d "github.com/denniswon/tcex/app/data"
	q "github.com/denniswon/tcex/app/queue"
	"github.com/gammazero/workerpool"
	"github.com/go-redis/redis/v8"
)

// ProcessOrderReplays
func ProcessOrderReplays(queue *q.PublishQueue, redis *redis.Client) {

	// Creating a job queue of size `#-of CPUs present in machine` * concurrency factor where order fetching requests are submitted to
	// There is no upper limit on the number of tasks queued, other than the limits of system resources
	// If the number of inbound tasks is too many to even queue for pending processing, then we should distribute workload over multiple systems,
	// and/or storing input for pending processing in intermediate storage such as a distributed message queue, etc.
	wp := workerpool.New(runtime.NumCPU() * int(cfg.GetConcurrencyFactor()))
	defer wp.Stop()

	for {
		wp.Submit(func() {
			for {
				if order, ok := queue.PublishNext(); ok {
					log.Printf("Publishing order %d\n", order)

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

					PublishReplayOrder(order, &_order, queue, redis)
				}
			}
		})
	}
}
