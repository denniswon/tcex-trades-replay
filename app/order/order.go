package order

import (
	d "github.com/denniswon/tcex/app/data"
	q "github.com/denniswon/tcex/app/queue"
	"github.com/go-redis/redis/v8"
)

// PublishReplayOrder - Attempts to process order data from Redis pubsub channel
func PublishReplayOrder(orderId string, order *d.Order, queue *q.ReplayQueue, redis *redis.Client) bool {

	// -- 3 step pub/sub attempt

	// 1. Asking queue whether we need to publish order or not
	if !queue.CanPublish(orderId) {
		return false
	}

	// 2. Attempting to publish order on Pub/Sub topic
	if !PublishOrder(orderId, order, redis) {
		return false
	}

	// 3. Marking this order as published
	if !queue.Published(orderId) {
		return false
	}

	return true
}

// PublishReplayKline - Attempts to process kline data from Redis pubsub channel
func PublishReplayKline(orderId string, kline *d.Kline, queue *q.ReplayQueue, redis *redis.Client) bool {

	// -- 3 step pub/sub attempt

	// 1. Asking queue whether we need to publish order or not
	if !queue.CanPublish(orderId) {
		return false
	}

	// 2. Attempting to publish order on Pub/Sub topic
	if !PublishKline(orderId, kline, redis) {
		return false
	}

	// 3. Marking this order as published
	if !queue.Published(orderId) {
		return false
	}

	return true
}

// PublishReplayEOF - Attempts to notify the client of the EOF for replay
func PublishReplayEOF(orderId string, queue *q.ReplayQueue, redis *redis.Client) bool {

	// -- 3 step pub/sub attempt

	// 1. Asking queue whether we need to publish order or not
	if !queue.CanPublish(orderId) {
		return false
	}

	// 2. Attempting to publish replay EOF on Pub/Sub topic
	if !PublishEOF(orderId, redis) {
		return false
	}

	// 3. Marking this order as published
	if !queue.Published(orderId) {
		return false
	}

	return true
}
