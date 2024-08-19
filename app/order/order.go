package order

import (
	q "github.com/denniswon/tcex/app/queue"
	d "github.com/denniswon/tcex/app/redis"
)

// ProcessOrders - Processes orders batch
func ProcessOrders(order *d.Order, redis *d.RedisInfo, queue *q.OrderReplayQueue) bool {

	// -- 3 step pub/sub attempt
	//
	// Attempting to publish whole order data to redis pubsub channel

	// 1. Asking queue whether we need to publish order or not
	if !queue.CanPublish(order.Number) {
		return false
	}

	// 2. Attempting to publish order on Pub/Sub topic
	if !PublishOrder(order, redis) {
		return false
	}

	// 3. Marking this order as published
	if !queue.Published(order.Number) {
		return false
	}

	return true
}
