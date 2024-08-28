package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// KlineConsumer - To be subscribed to `kline` topic using this consumer handle
// and client connected using websocket needs to be delivered this piece of data
type KlineConsumer struct {
	Client     *redis.Client
	Request    *SubscriptionRequest
	Connection *websocket.Conn
	PubSub     *redis.PubSub
	ConnLock   *sync.Mutex
	TopicLock  *sync.RWMutex
}

// Subscribe - Subscribe to `kline` channel
func (k *KlineConsumer) Subscribe() {
	k.PubSub = k.Client.Subscribe(context.Background(), k.Request.ID)
}

// Listen - Listener function, which keeps looping in infinite loop
// and reads data from subcribed channel, which also gets delivered to client application
func (k *KlineConsumer) Listen() {

	for {

		msg, err := k.PubSub.ReceiveTimeout(context.Background(), time.Second)
		if err != nil {
			continue
		}

		switch m := msg.(type) {

		case *redis.Subscription:

			// Pubsub broker informed we've been unsubscribed from this topic
			if m.Kind == "unsubscribe" {
				return
			}

			k.SendData(&SubscriptionResponse{
				Code:    1,
				ID:      m.Channel,
				Message: fmt.Sprintf("Subscribed to `%s`", m.Channel),
			})

		case *redis.Message:
			k.Send(m.Payload)

		}

	}
}

// Send - Tries to deliver subscribed kline data to client application
// connected over websocket
func (k *KlineConsumer) Send(msg string) {

	if strings.Contains(msg, "request_id") {
		k.SendEOF(msg)
		return

	}

	var kline struct {
		Timestamp						uint64 	`json:"timestamp"` 		// bucket start time in unix timestamp
		Low 								float32 `json:"low"`  			  // lowest price during the bucket interval
		High 								float32 `json:"high"`  				// highest price during the bucket interval
		Open								float32 `json:"open"`  				// opening price (first trade) in the bucket interval
		Close 							float32 `json:"close"`  			// closing price (last trade) in the bucket interval
		Volume 							int64 	`json:"volume"`  			// volume of trading activity during the bucket interval
		Turnover						float64 `json:"turnover"`			// total usd volume of trading activity during the bucket interval
		Granularity         uint16  `json:"granularity"`	// granularity field is in "seconds"
	}

	_msg := []byte(msg)

	err := json.Unmarshal(_msg, &kline)
	if err != nil {

		log.Printf("[!] Failed to decode published kline data to JSON : %s\n", err.Error())

		return
	}

	k.SendData(&kline)
}

// SendEOF - Tries to deliver eof data to client application
// connected over websocket
func (k *KlineConsumer) SendEOF(msg string) {

	var eof struct{
		RequestID	string `json:"request_id"`
	}

	_msg := []byte(msg)

	err := json.Unmarshal(_msg, &eof)

	if err != nil {

		log.Printf("[!] Failed to decode published eof data to JSON : %s\n", err.Error())

		return
	}

	k.SendData(&eof)
	log.Printf("Published EOF for request %s\n", eof.RequestID)
}

// SendData - Sending message to client application, connected over websocket
//
// If failed, we're going to remove subscription & close websocket
// connection ( connection might be already closed though )
func (k *KlineConsumer) SendData(data interface{}) bool {

	// -- Critical section of code begins
	//
	// Attempting to write to a network resource,
	// shared among multiple go routines
	k.ConnLock.Lock()
	defer k.ConnLock.Unlock()

	if err := k.Connection.WriteJSON(data); err != nil {
		log.Printf("[!] Failed to deliver kline data for request %s : %s\n", k.Request.ID, err.Error())
		return false
	}

	return true
}

// Unsubscribe - Unsubscribe from kline data publishing event this client has subscribed to
func (k *KlineConsumer) Unsubscribe() {

	if k.PubSub == nil {
		log.Printf("[!] Bad attempt to unsubscribe from `kline` topic\n")
		return
	}

	if err := k.PubSub.Unsubscribe(context.Background(), k.Request.ID); err != nil {
		log.Printf("[!] Failed to unsubscribe from topic %s : %s\n", k.Request.ID, err.Error())
		return
	}

	resp := &SubscriptionResponse{
		Code:    1,
		Message: fmt.Sprintf("Unsubscribed from `%s`", k.Request.ID),
	}

	// -- Critical section of code begins
	//
	// Attempting to write to a network resource,
	// shared among multiple go routines
	k.ConnLock.Lock()
	defer k.ConnLock.Unlock()

	if err := k.Connection.WriteJSON(resp); err != nil {

		log.Printf("[!] Failed to deliver unsubscription confirmation for request %s : %s\n", k.Request.ID, err.Error())
		return

	}
}
