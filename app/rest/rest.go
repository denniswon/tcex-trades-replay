package rest

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gin-contrib/cors"

	c "github.com/denniswon/tcex/app/client"
	cfg "github.com/denniswon/tcex/app/config"
	o "github.com/denniswon/tcex/app/order"
	ps "github.com/denniswon/tcex/app/pubsub"
	q "github.com/denniswon/tcex/app/queue"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// RunHTTPServer - Holds definition for all REST API(s) to be exposed
func RunHTTPServer(_orderClient *c.OrderClient, _queue *q.OrderReplayQueue, _redisClient *redis.Client) {

	// Checking if webserver in production mode or not
	checkIfInProduction := func() bool {
		return strings.ToLower(cfg.Get("Production")) == "yes"
	}

	// Running in production/ debug mode depending upon
	// config specified in .env file
	if checkIfInProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	// enabled cors
	router.Use(cors.Default())

	router.GET("/v1/ws", func(c *gin.Context) {

		// Setting read & write buffer size
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {

			log.Printf("[!] Failed to upgrade to websocket : %s\n", err.Error())
			return

		}

		// Registering websocket connection closing, to be executed when leaving
		// this function order
		defer conn.Close()

		// To be used for concurrent safe access of
		// underlying network socket
		connLock := sync.Mutex{}
		// To be used for concurrent safe access of subscribed
		// topic's associative array
		topicLock := sync.RWMutex{}

		// Log it when closing connection
		defer func() {

			log.Printf("[] Closing websocket connection\n",)

		}()



		// All topic subscription/ unsubscription requests
		// to handled by this higher layer abstraction
		pubsubManager := ps.SubscriptionManager{
			Topics:     	make(map[string]*ps.SubscriptionRequest),
			Consumers:  	make(map[string]ps.Consumer),
			Redis:     		_redisClient,
			Connection: 	conn,
			ConnLock:   	&connLock,
			TopicLock:  	&topicLock,
		}

		// Unsubscribe from all pubsub topics ( 3 at max ) when returning from
		// this execution scope
		defer func() {

			topicLock.Lock()
			defer topicLock.Unlock()

			for _, v := range pubsubManager.Consumers {
				v.Unsubscribe()
			}

		}()

		// Client communication handling logic
		for {

			var req ps.SubscriptionRequest

			if err := conn.ReadJSON(&req); err != nil {

				log.Printf("[!] Failed to read message : %s\n", err.Error())
				break

			}

			// Attempting to subscribe to/ unsubscribe from this topic
			switch req.Type {
			case "subscribe":
				go o.SubscribeToNewOrders(_orderClient, _queue)
				pubsubManager.Subscribe(&req)
			case "unsubscribe":
				pubsubManager.Unsubscribe(&req)
			}

		}

	})

	router.Run(fmt.Sprintf(":%s", cfg.Get("PORT")))
}
