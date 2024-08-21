package pubsub

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// SubscriptionManager - Higher level abstraction to be used
// by websocket connection acceptor, for subscribing to topics
type SubscriptionManager struct {
	Topics     	map[string]*SubscriptionRequest
	Consumers  	map[string]Consumer
	Redis   	 	*redis.Client
	Connection 	*websocket.Conn
	ConnLock   	*sync.Mutex
	TopicLock  	*sync.RWMutex
}

// Subscribe - Websocket connection manager can reliably call
// this function when ever it receives one valid subscription request
// with out worrying about how will it be handled
func (s *SubscriptionManager) Subscribe(req *SubscriptionRequest) {

	s.TopicLock.Lock()
	defer s.TopicLock.Unlock()

	_, ok := s.Topics[req.ID]
	if !ok {

		s.Topics[req.ID] = req
		s.Consumers[req.ID] = NewOrderConsumer(s.Redis, req, s.Connection, s.ConnLock, s.TopicLock)

	}

	s.Consumers[req.ID].SendData(
		&SubscriptionResponse{
			Code:    	1,
			Message: 	fmt.Sprintf("Subscription request for replay : `%s` (`x%d`)", req.Filename, req.ReplayRate),
			ID:    		req.ID,
		})
}

// Unsubscribe - Websocket connection manager can reliably call
// this to unsubscribe from topic for this client
func (s *SubscriptionManager) Unsubscribe(req *SubscriptionRequest) {

	s.TopicLock.Lock()
	defer s.TopicLock.Unlock()

	_, ok := s.Topics[req.ID]
	if !ok {
		return
	}

	delete(s.Topics, req.ID)

	s.Consumers[req.ID].SendData(
		&SubscriptionResponse{
			Code:    1,
			Message: fmt.Sprintf("Unsubscribed from `%s`", req.ID),
		})

	s.Consumers[req.ID].Unsubscribe()
	delete(s.Consumers, req.ID)
}
