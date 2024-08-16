package order

import (
	"context"
	"log"

	d "github.com/denniswon/tcex/app/data"
	"github.com/denniswon/tcex/app/db"
)

// PublishOrder - Attempts to publish order data to Redis pubsub channel
func PublishOrder(order *db.PackedOrder, redis *d.RedisInfo) bool {

	if order == nil {
		return false
	}

	_order := &d.Order{
		Hash:                order.Order.Hash,
		Number:              order.Order.Number,
		Time:                order.Order.Time,
		ParentHash:          order.Order.ParentHash,
		Difficulty:          order.Order.Difficulty,
		GasUsed:             order.Order.GasUsed,
		GasLimit:            order.Order.GasLimit,
		Nonce:               order.Order.Nonce,
		Miner:               order.Order.Miner,
		Size:                order.Order.Size,
		StateRootHash:       order.Order.StateRootHash,
		UncleHash:           order.Order.UncleHash,
		TransactionRootHash: order.Order.TransactionRootHash,
		ReceiptRootHash:     order.Order.ReceiptRootHash,
		ExtraData:           order.Order.ExtraData,
	}

	if err := redis.Client.Publish(context.Background(), redis.OrderPublishTopic, _order).Err(); err != nil {

		log.Printf("Failed to publish order %d : %s\n", order.Order.Number, err.Error())
		return false

	}

	log.Printf("📎 Published order %d\n", order.Order.Number)

	return PublishTxs(order.Order.Number, order.Transactions, redis)
}
