package pubsub

import (
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// Consumer - Order, transaction & event consumers need to implement these methods
type Consumer interface {
	Subscribe()
	Listen()
	Send(msg string)
	SendData(data interface{}) bool
	Unsubscribe()
}

// NewOrderConsumer - Creating one new order data consumer, which will subscribe to order
// topic & listen for data being published on this channel, which will eventually be
// delivered to client application over websocket connection
func NewOrderConsumer(client *redis.Client, request *SubscriptionRequest, conn *websocket.Conn, connLock *sync.Mutex, topicLock *sync.RWMutex) *OrderConsumer {
	consumer := OrderConsumer{
		Client:     client,
		Request:   request,
		Connection: conn,
		ConnLock:   connLock,
		TopicLock:  topicLock,
	}

	consumer.Subscribe()
	go consumer.Listen()

	return &consumer
}
