package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// OrderConsumer - To be subscribed to `order` topic using this consumer handle
// and client connected using websocket needs to be delivered this piece of data
type OrderConsumer struct {
	Client     *redis.Client
	Request    *SubscriptionRequest
	Connection *websocket.Conn
	PubSub     *redis.PubSub
	ConnLock   *sync.Mutex
	TopicLock  *sync.RWMutex
}

// Subscribe - Subscribe to `order` channel
func (b *OrderConsumer) Subscribe() {
	b.PubSub = b.Client.Subscribe(context.Background(), b.Request.ID)
}

// Listen - Listener function, which keeps looping in infinite loop
// and reads data from subcribed channel, which also gets delivered to client application
func (b *OrderConsumer) Listen() {

	for {

		msg, err := b.PubSub.ReceiveTimeout(context.Background(), time.Second)
		if err != nil {
			continue
		}

		switch m := msg.(type) {

		case *redis.Subscription:

			// Pubsub broker informed we've been unsubscribed from this topic
			if m.Kind == "unsubscribe" {
				return
			}

			b.SendData(&SubscriptionResponse{
				Code:    1,
				ID:      m.Channel,
				Message: fmt.Sprintf("Subscribed to `%s`", m.Channel),
			})

		case *redis.Message:
			b.Send(m.Payload)

		}

	}
}

// Send - Tries to deliver subscribed order data to client application
// connected over websocket
func (b *OrderConsumer) Send(msg string) {

	var order struct {
		Price               string `json:"price"`
		Quantity            uint64 `json:"quantity"`
		Aggressor           string `json:"aggressor"`
		Timestamp           int64  `json:"timestamp"`
	}

	_msg := []byte(msg)

	err := json.Unmarshal(_msg, &order)
	if err != nil {
		log.Printf("[!] Failed to decode published order data to JSON : %s\n", err.Error())
		return
	}

	b.SendData(&order)
}

// SendData - Sending message to client application, connected over websocket
//
// If failed, we're going to remove subscription & close websocket
// connection ( connection might be already closed though )
func (b *OrderConsumer) SendData(data interface{}) bool {

	// -- Critical section of code begins
	//
	// Attempting to write to a network resource,
	// shared among multiple go routines
	b.ConnLock.Lock()
	defer b.ConnLock.Unlock()

	if err := b.Connection.WriteJSON(data); err != nil {
		log.Printf("[!] Failed to deliver order data for request %s : %s\n", b.Request.ID, err.Error())
		return false
	}

	return true
}

// Unsubscribe - Unsubscribe from order data publishing event this client has subscribed to
func (b *OrderConsumer) Unsubscribe() {

	if b.PubSub == nil {
		log.Printf("[!] Bad attempt to unsubscribe from `order` topic\n")
		return
	}

	if err := b.PubSub.Unsubscribe(context.Background(), b.Request.ID); err != nil {
		log.Printf("[!] Failed to unsubscribe from topic %s : %s\n", b.Request.ID, err.Error())
		return
	}

	resp := &SubscriptionResponse{
		Code:    1,
		Message: fmt.Sprintf("Unsubscribed from `%s`", b.Request.ID),
	}

	// -- Critical section of code begins
	//
	// Attempting to write to a network resource,
	// shared among multiple go routines
	b.ConnLock.Lock()
	defer b.ConnLock.Unlock()

	if err := b.Connection.WriteJSON(resp); err != nil {

		log.Printf("[!] Failed to deliver unsubscription confirmation for request %s : %s\n", b.Request.ID, err.Error())
		return

	}
}
