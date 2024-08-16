package db

import (
	"log"
	"math"

	"github.com/denniswon/tcex/app/data"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// GetAllOrderNumbersInRange - Returns all order numbers in given range, both inclusive
func GetAllOrderNumbersInRange(db *gorm.DB, from uint64, to uint64) []uint64 {

	var orders []uint64
	rangeFrom := math.Min(float64(from), float64(to))
	rangeTo := math.Max(float64(from), float64(to))
	if err := db.Model(&Order{}).Where("number >= ? and number <= ?", rangeFrom, rangeTo).Order("number asc").Select("number").Find(&orders).Error; err != nil {

		log.Printf("[!] Failed to fetch order numbers by range : %s\n", err.Error())
		return nil

	}

	return orders
}

// GetCurrentOldestOrderNumber - Fetches what's lowest order number present in database,
// which denotes if it's not 0, from here we can start syncing again, until we reach 0
func GetCurrentOldestOrderNumber(db *gorm.DB) uint64 {
	var number uint64

	if err := db.Raw("select min(number) from orders").Scan(&number).Error; err != nil {
		return 0
	}

	return number
}

// GetCurrentOrderNumber - Returns highest order number, which got processed
// by the service
func GetCurrentOrderNumber(db *gorm.DB) uint64 {
	var number uint64

	if err := db.Raw("select max(number) from orders").Scan(&number).Error; err != nil {
		return 0
	}

	return number
}

// GetOrderCount - Returns how many orders currently present in database
//
// Caution : As we're dealing with very large tables
// ( with row count  ~ 10M & increasing 1 row every 2 seconds )
// this function needs to be least frequently, otherwise due to full table
// scan it'll cost us a lot
//
// Currently only using during application start up
//
// All other order count calculation requirements can be fulfilled by
// using in-memory program state holder
func GetOrderCount(db *gorm.DB) uint64 {
	var number int64

	if err := db.Model(&Order{}).Count(&number).Error; err != nil {
		return 0
	}

	return uint64(number)
}

// GetOrderByHash - Given orderhash finds out order related information
//
// If not found, returns nil
func GetOrderByHash(db *gorm.DB, hash common.Hash) *data.Order {
	var order data.Order

	if res := db.Model(&Order{}).Where("hash = ?", hash.Hex()).First(&order); res.Error != nil {
		return nil
	}

	return &order
}

// GetOrderByNumber - Fetch order using order number
//
// If not found, returns nil
func GetOrderByNumber(db *gorm.DB, number uint64) *data.Order {
	var order data.Order

	if res := db.Model(&Order{}).Where("number = ?", number).First(&order); res.Error != nil {
		return nil
	}

	return &order
}

// GetOrdersByNumberRange - Given order numbers as range, it'll extract out those orders
// by number, while returning them in ascendically sorted form in terms of order numbers
//
// Note : Can return at max 10 orders in a single query
//
// If more orders are requested, simply to be rejected
// In that case, consider splitting them such that they satisfy criteria
func GetOrdersByNumberRange(db *gorm.DB, from uint64, to uint64) *data.Orders {
	var orders []*data.Order

	if res := db.Model(&Order{}).Where("number >= ? and number <= ?", from, to).Order("number asc").Find(&orders); res.Error != nil {
		return nil
	}

	return &data.Orders{
		Orders: orders,
	}
}

// GetOrdersByTimeRange - Given time range ( of 60 sec span at max ), returns orders
// mined in that time span
//
// If asked to find out orders in time span larger than 60 sec, simply drops query request
func GetOrdersByTimeRange(db *gorm.DB, from uint64, to uint64) *data.Orders {
	var orders []*data.Order

	if res := db.Model(&Order{}).Where("timestamp >= ? and timestamp <= ?", from, to).Order("number asc").Find(&orders); res.Error != nil {
		return nil
	}

	return &data.Orders{
		Orders: orders,
	}
}
