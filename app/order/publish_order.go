package order

import (
	"context"
	"log"
	"strings"

	d "github.com/denniswon/tcex/app/data"
	"github.com/go-redis/redis/v8"
)

// PublishOrder - Attempts to publish order data to Redis pubsub channel
func PublishOrder(orderId string, order *d.Order, redis *redis.Client) bool {

	if order == nil {
		return false
	}

	tokens := strings.Split(orderId, ":")
	if len(tokens) != 2 {
		log.Printf("Unexpected order id %s\n", orderId)
		return false
	}

	requestId := tokens[0]
	if err := redis.Publish(context.Background(), requestId, order).Err(); err != nil {

		log.Printf("Failed to publish order %s : %s\n", orderId, err.Error())
		return false

	}

	log.Printf("📎 Published order %s\n", orderId)
	return true

}
