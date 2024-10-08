package app

import (
	"context"
	"log"
	"os"

	cfg "github.com/denniswon/tcex/app/config"
	q "github.com/denniswon/tcex/app/queue"
	"github.com/go-redis/redis/v8"
)

// Setting ground up i.e. acquiring resources required & determining with
// some basic checks whether we can proceed to next step or not
func bootstrap(configFile string) (*q.RequestQueue, *q.ReplayQueue, *redis.Client) {

	err := cfg.Read(configFile)
	if err != nil {
		log.Fatalf("[!] Failed to read `.env` : %s\n", err.Error())
	}

	_redis := getRedisClient()
	if _redis == nil {
		log.Fatalf("[!] Failed to connect to Redis Server\n")
	}
	if err := _redis.FlushAll(context.Background()).Err(); err != nil {
		log.Printf("[!] Failed to flush all keys from redis : %s\n", err.Error())
	}

	// orders queue for fetching orders from the input file
	requestQueue := q.NewRequestQueue(_redis)
	// order replay publishing queue
	replayQueue := q.NewReplayQueue()

	// Create a temporary directory for file uploads
	tempDir, err := os.MkdirTemp("", "uploads-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	return requestQueue, replayQueue, _redis
}
