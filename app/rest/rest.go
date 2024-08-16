package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"

	cmn "github.com/denniswon/tcex/app/common"
	cfg "github.com/denniswon/tcex/app/config"
	d "github.com/denniswon/tcex/app/data"
	"github.com/denniswon/tcex/app/db"
	ps "github.com/denniswon/tcex/app/pubsub"
	"github.com/denniswon/tcex/app/rest/graph/generated"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/denniswon/tcex/app/rest/graph"
)

// RunHTTPServer - Holds definition for all REST API(s) to be exposed
func RunHTTPServer(_db *gorm.DB, _status *d.StatusHolder, _redisClient *redis.Client) {

	respondWithJSON := func(data []byte, c *gin.Context) {
		if data != nil {
			c.Data(http.StatusOK, "application/json", data)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "JSON encoding failed",
		})
	}

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

	grp := router.Group("/v1")

	{

		// For checking the service's syncing status
		grp.GET("/synced", func(c *gin.Context) {

			currentOrderNumber := _status.GetLatestOrderNumber()
			orderCountInDB := _status.OrderCountInDB()
			remaining := (currentOrderNumber + 1) - orderCountInDB
			elapsed := _status.ElapsedTime()

			status := fmt.Sprintf("%.2f %%", (float64(orderCountInDB)/float64(currentOrderNumber+1))*100)
			eta := "0s"
			if remaining > 0 {
				eta = (time.Duration((elapsed.Seconds()/float64(_status.Done()))*float64(remaining)) * time.Second).String()
			}

			c.JSON(http.StatusOK, gin.H{
				"synced":    status,
				"processed": _status.Done(),
				"elapsed":   elapsed.String(),
				"eta":       eta,
				"status":	_status.State,
			})

		})

		// Query order data using order hash/ number/ order number range ( 10 at max )
		grp.GET("/order", func(c *gin.Context) {

			number := c.Query("number")

			// Order number based single order retrieval request handler
			if number != "" {

				_num, err := cmn.ParseNumber(number)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"msg": "Bad order number",
					})
					return
				}

				if order := db.GetOrderByNumber(_db, _num); order != nil {
					respondWithJSON(order.ToJSON(), c)
					return
				}

				c.JSON(http.StatusNotFound, gin.H{
					"msg": "Not found",
				})
				return
			}

			// Order number range based query
			// At max 10 orders at a time to be returned
			fromOrder := c.Query("fromOrder")
			toOrder := c.Query("toOrder")

			if fromOrder != "" && toOrder != "" {

				_from, _to, err := cmn.RangeChecker(fromOrder, toOrder, cfg.GetOrderNumberRange())
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"msg": "Bad order number range",
					})
					return
				}

				if orders := db.GetOrdersByNumberRange(_db, _from, _to); orders != nil {
					respondWithJSON(orders.ToJSON(), c)
					return
				}

				c.JSON(http.StatusNotFound, gin.H{
					"msg": "Not found",
				})
				return
			}

			// Query orders by timestamp range, at max 60 seconds of timestamp
			// can be mentioned, otherwise request to be rejected
			fromTime := c.Query("fromTime")
			toTime := c.Query("toTime")

			if fromTime != "" && toTime != "" {

				_from, _to, err := cmn.RangeChecker(fromTime, toTime, cfg.GetTimeRange())
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"msg": "Bad order time range",
					})
					return
				}

				if orders := db.GetOrdersByTimeRange(_db, _from, _to); orders != nil {
					respondWithJSON(orders.ToJSON(), c)
					return
				}

				c.JSON(http.StatusNotFound, gin.H{
					"msg": "Not found",
				})
				return
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Bad query param(s)",
			})

		})
	}

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
			Topics:     make(map[string]map[string]*ps.SubscriptionRequest),
			Consumers:  make(map[string]ps.Consumer),
			Client:     _redisClient,
			Connection: conn,
			DB:         _db,
			ConnLock:   &connLock,
			TopicLock:  &topicLock,
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

			// Validating incoming request on websocket subscription channel
			if !req.Validate(&pubsubManager) {
				// -- Critical section of code begins
				//
				// Attempting to write to shared network connection
				connLock.Lock()

				if err := conn.WriteJSON(&ps.SubscriptionResponse{Code: 0, Message: "Bad Payload"}); err != nil {
					log.Printf("[!] Failed to write message : %s\n", err.Error())
				}

				connLock.Unlock()
				// -- ends here
				break
			}

			// Attempting to subscribe to/ unsubscribe from this topic
			switch req.Type {
			case "subscribe":
				pubsubManager.Subscribe(&req)
			case "unsubscribe":
				pubsubManager.Unsubscribe(&req)
			}

		}

	})

	router.POST("/v1/graphql",
		// Attempting to pass router context, which so that some job can
		// be done if needed to (i.e. logging, stats, etc.) before delivering requested piece of data to client
		func(c *gin.Context) {
			ctx := context.WithValue(c.Request.Context(), "RouterContextInGraphQL", c)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		},

		func(c *gin.Context) {

			gql := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
				Resolvers: &graph.Resolver{},
			}))

			if gql == nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "Failed to handle graphQL query",
				})
				return
			}

			gql.ServeHTTP(c.Writer, c.Request)

		})

	router.GET("/v1/graphql-playground", func(c *gin.Context) {

		gpg := playground.Handler("tcex", "/v1/graphql")

		if gpg == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Failed to create graphQL playground",
			})
			return
		}

		gpg.ServeHTTP(c.Writer, c.Request)

	})

	router.Run(fmt.Sprintf(":%s", cfg.Get("PORT")))
}
