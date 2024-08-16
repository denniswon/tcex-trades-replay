package order

import (
	"log"
	"time"

	d "github.com/denniswon/tcex/app/data"
	"github.com/denniswon/tcex/app/db"
	q "github.com/denniswon/tcex/app/queue"
	"gorm.io/gorm"
)

// ProcessOrders - Processes orders batch
func ProcessOrders(order *d.Order, _db *gorm.DB, redis *d.RedisInfo, queue *q.OrderReplayQueue, status *d.StatusHolder, startingAt time.Time) bool {

	// Closure managing publishing whole order data i.e. order header, txn(s), event logs on redis pubsub channel
	pubsubWorker := func() (*db.Order, bool) {

	// Constructing order data to published & persisted
	packedOrder := BuildPackedOrder(order)

	// -- 3 step pub/sub attempt
	//
	// Attempting to publish whole order data to redis pubsub channel

	// 1. Asking queue whether we need to publish order or not
	if !queue.CanPublish(order.NumberU64()) {
		return packedOrder, true
	}

	// 2. Attempting to publish order on Pub/Sub topic
	if !PublishOrder(packedOrder, redis) {
		return nil, false
	}

	// 3. Marking this order as published
	if !queue.Published(order.NumberU64()) {
		return nil, false
	}

		// -- done, with publishing on Pub/Sub topic

		return packedOrder, true

	}

	packedOrder, ok := pubsubWorker()
	if !ok {
		return false
	}

	// If order doesn't contain any tx, we'll attempt to persist only order
	if err := db.StoreOrder(_db, packedOrder, status, queue); err != nil {

		log.Printf("Failed to process order %d : %s\n", order.NumberU64(), err.Error())
		return false

	}

	return true
}
