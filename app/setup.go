package app

import (
	"context"
	"log"

	c "github.com/denniswon/tcex/app/client"
	cfg "github.com/denniswon/tcex/app/config"
	q "github.com/denniswon/tcex/app/queue"
	"github.com/go-redis/redis/v8"
)

// Setting ground up i.e. acquiring resources required & determining with
// some basic checks whether we can proceed to next step or not
func bootstrap(configFile string) (*c.OrderClient, *redis.Client, *q.OrderReplayQueue) {

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

	// order processor queue
	_queue := q.New()

	// orders client for fetching orders from the input file
	_orderClient := c.NewOrderClient()


	return _orderClient, _redisClient, _queue
}
