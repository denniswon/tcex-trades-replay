package config

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
)

// Read - Reading .env file content, during application start up
func Read(file string) error {
	viper.SetConfigFile(file)

	return viper.ReadInConfig()
}

// Get - Get config value by key
func Get(key string) string {
	return viper.GetString(key)
}

func GetBatchSize() uint64 {

	batchSize := Get("BatchSize")
	if batchSize == "" {
		return 125
	}

	parsedBatchSize, err := strconv.ParseUint(batchSize, 10, 64)
	if err != nil {
		log.Printf("[!] Failed to parse batch size : %s\n", err.Error())
		return 1
	}

	return parsedBatchSize
}

// GetConcurrencyFactor - Reads concurrency factor specified in `.env` file, during deployment
// and returns that number as unsigned integer
func GetConcurrencyFactor() uint64 {

	factor := Get("ConcurrencyFactor")
	if factor == "" {
		return 1
	}

	parsedFactor, err := strconv.ParseUint(factor, 10, 64)
	if err != nil {
		log.Printf("[!] Failed to parse concurrency factor : %s\n", err.Error())
		return 1
	}

	return parsedFactor
}

// GetOrderConfirmations - Number of order confirmations required
// before considering that order to be finalized, and can be persisted
// in a permanent data store
func GetOrderConfirmations() uint64 {

	confirmationCount := Get("OrderConfirmations")
	if confirmationCount == "" {
		return 0
	}

	parsedConfirmationCount, err := strconv.ParseUint(confirmationCount, 10, 64)
	if err != nil {
		log.Printf("[!] Failed to parse order confirmations : %s\n", err.Error())
		return 0
	}

	return parsedConfirmationCount
}

// GetOrderNumberRange - Returns how many orders can be queried at a time
// when performing range based queries from client side
func GetOrderNumberRange() uint64 {

	orderRange := Get("OrderRange")
	if orderRange == "" {
		return 100
	}

	parsedOrderRange, err := strconv.ParseUint(orderRange, 10, 64)
	if err != nil {
		log.Printf("[!] Failed to parse order range : %s\n", err.Error())
		return 100
	}

	return parsedOrderRange
}

// GetTimeRange - Returns what's the max time span that can be used while performing query
// from client side, in terms of second
func GetTimeRange() uint64 {

	timeRange := Get("TimeRange")
	if timeRange == "" {
		return 3600
	}

	parsedTimeRange, err := strconv.ParseUint(timeRange, 10, 64)
	if err != nil {
		log.Printf("[!] Failed to parse time range : %s\n", err.Error())
		return 3600
	}

	return parsedTimeRange
}
