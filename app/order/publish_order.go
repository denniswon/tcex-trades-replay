package order

import (
	"context"
	"log"

	d "github.com/denniswon/tcex/app/data"
	"github.com/go-redis/redis/v8"
)

// PublishOrder - Attempts to publish order data to Redis pubsub channel
func PublishOrder(orderId string, order *d.Order, redis *redis.Client) bool {

	if order == nil {
		return false
	}

	if err := redis.Publish(context.Background(), orderId, order).Err(); err != nil {

		log.Printf("Failed to publish order %d : %s\n", orderId, err.Error())
		return false

	}

	log.Printf("📎 Published order %d\n", orderId)
	return true

}
