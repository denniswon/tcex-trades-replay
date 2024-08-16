package order

import (
	d "github.com/denniswon/tcex/app/data"
	"github.com/denniswon/tcex/app/db"
)

// BuildPackedOrder - Builds struct holding whole order data i.e.
// order header, order body i.e. tx(s), event log(s)
func BuildPackedOrder(order *d.Order) *db.Order {

	packedOrder := &db.Order{
		Number:              order.Number,
		Price:               order.Price,
		Timestamp:           order.Timestamp,
		Aggressor:           order.Aggressor,
		Quantity:            order.Quantity,
	}

	return packedOrder
}
