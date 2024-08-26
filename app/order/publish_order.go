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

// PublishKline - Attempts to publish kline data to Redis pubsub channel
func PublishKline(orderId string, kline *d.Kline, redis *redis.Client) bool {

	if kline == nil {
		return false
	}

	tokens := strings.Split(orderId, ":")
	if len(tokens) != 2 {
		log.Printf("Unexpected order id %s\n", orderId)
		return false
	}

	requestId := tokens[0]
	if err := redis.Publish(context.Background(), requestId, kline).Err(); err != nil {

		log.Printf("Failed to publish kline data for order %s : %s\n", orderId, err.Error())
		return false

	}

	log.Printf("📎 Published kline data for order %s\n", orderId)
	return true

}


// PublishEOF - Attempts to publish replay eof to Redis pubsub channel
func PublishEOF(orderId string, redis *redis.Client) bool {

	tokens := strings.Split(orderId, ":")
	if len(tokens) != 2 {
		log.Printf("Unexpected order id %s\n", orderId)
		return false
	}

	requestId := tokens[0]
	eof := d.EOF{
		RequestID: requestId,
	}
	if err := redis.Publish(context.Background(), requestId, &eof).Err(); err != nil {

		log.Printf("Failed to publish eof %s : %s\n", requestId, err.Error())
		return false

	}

	log.Printf("📎 Published eof for request %s\n", requestId)
	return true

}
