package app

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"

	c "github.com/denniswon/tcex/app/client"
	cfg "github.com/denniswon/tcex/app/config"
)

// Connect to order client
func getOrderClient(orderNumber uint64) *c.OrderClient {
	client := c.NewOrderClient()
	if client != nil {
		log.Fatalln("[!] Failed to initialize order client")
	}
	return client
}

// Creates connection to Redis server & returns that handle to be used for further communication
func getRedisClient() *redis.Client {

	var options *redis.Options

	// If password is given in config file
	if cfg.Get("RedisPassword") != "" {

		options = &redis.Options{
			Network:  cfg.Get("RedisConnection"),
			Addr:     cfg.Get("RedisAddress"),
			Password: cfg.Get("RedisPassword"),
			DB:       0,
		}

	} else {
		// If password is not given, attempting to connect with out it
		//
		// Though this is not recommended
		options = &redis.Options{
			Network: cfg.Get("RedisConnection"),
			Addr:    cfg.Get("RedisAddress"),
			DB:      0,
		}

	}

	_redis := redis.NewClient(options)
	// Checking whether connection was successful or not
	if err := _redis.Ping(context.Background()).Err(); err != nil {
		return nil
	}

	return _redis
}
