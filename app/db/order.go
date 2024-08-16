package db

import (
	"errors"
	"log"

	d "github.com/denniswon/tcex/app/data"
	q "github.com/denniswon/tcex/app/queue"
	"gorm.io/gorm"
)

// StoreOrder - Persisting order data in database,
// if data already not stored
//
// Also checks equality with existing data, if mismatch found,
// updated with latest data
func StoreOrder(dbWOTx *gorm.DB, order *Order, status *d.StatusHolder, queue *q.OrderReplayQueue) error {

	if order == nil {
		return errors.New("empty order received while attempting to persist")
	}

	// -- Starting DB transaction
	return dbWOTx.Transaction(func(dbWTx *gorm.DB) error {

		orderInserted := false

		persistedOrder := GetOrder(dbWTx, order.Number)
		if persistedOrder == nil {

			if err := PutOrder(dbWTx, order); err != nil {
				return err
			}

			orderInserted = true

		} else if !persistedOrder.SimilarTo(order) {

			log.Printf("[!] Order %d already present in DB, similar ❌\n", order.Number)

			// cascaded deletion !
			if err := DeleteOrder(dbWTx, order.Number); err != nil {
				return err
			}

			if err := PutOrder(dbWTx, order); err != nil {
				return err
			}

			orderInserted = true

		} else {

			log.Printf("[+] Order %d already present in DB, similar \n", order.Number)
			return nil

		}

		// If we've really inserted a new order into database,
		// count will get updated
		if orderInserted && status != nil {
			status.IncrementOrdersInserted()
		}

		return nil

	})
	// -- Ending DB transaction
}

// GetOrder - Fetch order by number, from database
func GetOrder(_db *gorm.DB, number uint64) *Order {
	var order Order

	if err := _db.Where("number = ?", number).First(&order).Error; err != nil {
		return nil
	}

	return &order
}

// PutOrder - Persisting fetched order
func PutOrder(dbWTx *gorm.DB, order *Order) error {

	return dbWTx.Create(order).Error

}

// DeleteOrder - Delete order entry, identified by order number, while
// cascading all dependent entries ( i.e. in transactions/ events table )
func DeleteOrder(dbWTx *gorm.DB, number uint64) error {

	return dbWTx.Where("number = ?", number).Delete(&Order{}).Error

}

// UpdateOrder - Updating already existing order
func UpdateOrder(dbWTx *gorm.DB, order *Order) error {

	return dbWTx.Model(&Order{}).Where("number = ?", order.Number).Updates(map[string]interface{}{
		"timestamp":   	order.Timestamp,
		"aggressor":		order.Aggressor,
		"price":      	order.Price,
		"quantity":   	order.Quantity,
	}).Error

}
