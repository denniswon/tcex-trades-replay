package app

import (
	"context"
	"log"
	"sync"

	c "github.com/denniswon/tcex/app/client"
	cfg "github.com/denniswon/tcex/app/config"
	d "github.com/denniswon/tcex/app/data"
	"github.com/denniswon/tcex/app/db"
	q "github.com/denniswon/tcex/app/queue"
	"github.com/denniswon/tcex/app/rest/graph"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Setting ground up i.e. acquiring resources required & determining with
// some basic checks whether we can proceed to next step or not
func bootstrap(configFile string) (*c.OrderClient, *redis.Client, *d.RedisInfo, *gorm.DB, *d.StatusHolder, *q.OrderReplayQueue) {

	err := cfg.Read(configFile)
	if err != nil {
		log.Fatalf("[!] Failed to read `.env` : %s\n", err.Error())
	}

	_redisClient := getRedisClient()
	if _redisClient == nil {
		log.Fatalf("[!] Failed to connect to Redis Server\n")
	}
	if err := _redisClient.FlushAll(context.Background()).Err(); err != nil {
		log.Printf("[!] Failed to flush all keys from redis : %s\n", err.Error())
	}

	_db := db.Connect()

	// Passing db handle to graph for resolving graphQL queries
	graph.GetDatabaseConnection(_db)

	_status := &d.StatusHolder{
		State: &d.SyncState{
			OrderCountAtStartUp:     db.GetOrderCount(_db),
			MaxOrderNumberAtStartUp: db.GetCurrentOrderNumber(_db),
		},
		Mutex: &sync.RWMutex{},
	}

	_redisInfo := &d.RedisInfo{
		Client:            _redisClient,
		OrderPublishTopic: "order",
	}

	// order processor queue
	_queue := q.New(db.GetCurrentOrderNumber(_db))

	_orderClient := getOrderClient(db.GetCurrentOrderNumber(_db))

	return _orderClient, _redisClient, _redisInfo, _db, _status, _queue
}
