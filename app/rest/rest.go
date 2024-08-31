package rest

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/google/uuid"

	cfg "github.com/denniswon/tcex/app/config"
	ps "github.com/denniswon/tcex/app/pubsub"
	q "github.com/denniswon/tcex/app/queue"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// RunHTTPServer - Holds definition for all REST API(s) to be exposed
func RunHTTPServer(_queue *q.RequestQueue, _redis *redis.Client, tempDir string) {

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20

	// enabled cors
	router.Use(cors.Default())

	grp := router.Group("/v1")

	{

		// For checking the service's syncing status
		grp.POST("/upload", func(c *gin.Context) {

			// single file
			file, _ := c.FormFile("file")

			log.Printf("Uploading File: %s (size: %d)\n", file.Filename, file.Size)

			_filepath := filepath.Join(tempDir, file.Filename)

			header := ps.UploadHeader{
				ID:       uuid.New().String(),
				Filepath: _filepath,
				Size:     file.Size,
			}

			// Check if file already exists
			if _, err := os.Stat(_filepath); err == nil {

				log.Printf("Upload file already exists : %s %d\n", header.Filepath, header.Size)

				c.JSON(http.StatusOK, header)

				return

			}

			c.SaveUploadedFile(file, header.Filepath)

			c.JSON(http.StatusOK, header)

		})

	}

	router.GET("/v1/ws", func(c *gin.Context) {

		// Setting read & write buffer size
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// TODO: You should never blindly trust any Origin by return true.
			// Have the function range over a list of accepted origins.
			CheckOrigin: func(r *http.Request) bool { return true },
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

			log.Printf("[] Closing websocket connection\n")

		}()

		// All topic subscription/ unsubscription requests
		// to handled by this higher layer abstraction
		pubsubManager := ps.SubscriptionManager{
			Topics:     make(map[string]*ps.SubscriptionRequest),
			Consumers:  make(map[string]ps.Consumer),
			Redis:      _redis,
			Connection: conn,
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
				continue

			}

			// Attempting to subscribe to/ unsubscribe from this topic
			switch req.Type {

			case "subscribe":

				// Validating incoming request on websocket subscription channel
				if !req.Validate() {
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

				_queue.Put(&req)
				pubsubManager.Subscribe(&req)

			case "unsubscribe":
				_queue.Remove(req.ID)
				pubsubManager.Unsubscribe(&req)

			}

		}

		for {

			err := <-_queue.Err()
			log.Printf("[!] Failed to process order %s : %s\n", err.RequestId, err.Err.Error())
			pubsubManager.Unsubscribe(pubsubManager.Topics[err.RequestId])

		}

	})

	router.Run(fmt.Sprintf(":%s", cfg.GetPort()))
}
