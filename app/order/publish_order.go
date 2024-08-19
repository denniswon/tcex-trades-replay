package order

import (
	"context"
	"log"

	d "github.com/denniswon/tcex/app/redis"
)

// PublishOrder - Attempts to publish order data to Redis pubsub channel
func PublishOrder(order *d.Order, redis *d.RedisInfo) bool {

	if order == nil {
		return false
	}

	if err := redis.Client.Publish(context.Background(), redis.OrderPublishTopic, order).Err(); err != nil {

		log.Printf("Failed to publish order %d : %s\n", order.Number, err.Error())
		return false

	}

	log.Printf("📎 Published order %d\n", order.Number)
	return true

}
