package graph

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/denniswon/tcex/app/data"
	"github.com/denniswon/tcex/app/rest/graph/model"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

var db *gorm.DB

// GetDatabaseConnection - Passing already connected database handle to this package,
// so that it can be used for handling database queries for resolving graphQL queries
func GetDatabaseConnection(conn *gorm.DB) {
	db = conn
}

func routerContextFromGraphQLContext(ctx context.Context) (*gin.Context, error) {

	ginContext := ctx.Value("RouterContextInGraphQL")
	if ginContext == nil {
		return nil, errors.New("Failed to retrieve router context")
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		return nil, errors.New("Type assert of router context failed")
	}

	return gc, nil
}

// Converting order data to graphQL compatible data structure
func getGraphQLCompatibleOrder(ctx context.Context, order *data.Order) (*model.Order, error) {

	if order == nil {
		return nil, errors.New("Found nothing")
	}

	extraData := ""
	if _h := hex.EncodeToString(order.ExtraData); _h != "" {
		extraData = fmt.Sprintf("0x%s", _h)
	}

	return &model.Order{
		Hash:            order.Hash,
		Number:          fmt.Sprintf("%d", order.Number),
		Time:            fmt.Sprintf("%d", order.Time),
		ParentHash:      order.ParentHash,
		Difficulty:      order.Difficulty,
		GasUsed:         fmt.Sprintf("%d", order.GasUsed),
		GasLimit:        fmt.Sprintf("%d", order.GasLimit),
		Nonce:           order.Nonce,
		Miner:           order.Miner,
		Size:            order.Size,
		StateRootHash:   order.StateRootHash,
		UncleHash:       order.UncleHash,
		TxRootHash:      order.TransactionRootHash,
		ReceiptRootHash: order.ReceiptRootHash,
		ExtraData:       extraData,
	}, nil
}

// Converting order array to graphQL compatible data structure
func getGraphQLCompatibleOrders(ctx context.Context, orders *data.Orders) ([]*model.Order, error) {
	if orders == nil {
		return nil, errors.New("Found nothing")
	}

	if !(len(orders.Orders) > 0) {
		return nil, errors.New("Found nothing")
	}

	_orders := make([]*model.Order, len(orders.Orders))

	for k, v := range orders.Orders {
		_v, _ := getGraphQLCompatibleOrder(ctx, v)
		_orders[k] = _v
	}

	return _orders, nil
}

// Converting transaction data to graphQL compatible data structure
func getGraphQLCompatibleTransaction(ctx context.Context, tx *data.Transaction, bookKeeping bool) (*model.Transaction, error) {
	if tx == nil {
		return nil, errors.New("Found nothing")
	}

	data := ""
	if _h := hex.EncodeToString(tx.Data); _h != "" {
		data = fmt.Sprintf("0x%s", _h)
	}

	if !strings.HasPrefix(tx.Contract, "0x") {
		return &model.Transaction{
			Hash:      tx.Hash,
			From:      tx.From,
			To:        tx.To,
			Contract:  "",
			Value:     tx.Value,
			Data:      data,
			Gas:       fmt.Sprintf("%d", tx.Gas),
			GasPrice:  tx.GasPrice,
			Cost:      tx.Cost,
			Nonce:     fmt.Sprintf("%d", tx.Nonce),
			State:     fmt.Sprintf("%d", tx.State),
			OrderHash: tx.OrderHash,
		}, nil
	}

	return &model.Transaction{
		Hash:      tx.Hash,
		From:      tx.From,
		To:        "",
		Contract:  tx.Contract,
		Value:     tx.Value,
		Data:      data,
		Gas:       fmt.Sprintf("%d", tx.Gas),
		GasPrice:  tx.GasPrice,
		Cost:      tx.Cost,
		Nonce:     fmt.Sprintf("%d", tx.Nonce),
		State:     fmt.Sprintf("%d", tx.State),
		OrderHash: tx.OrderHash,
	}, nil
}

// Converting transaction array to graphQL compatible data structure
func getGraphQLCompatibleTransactions(ctx context.Context, tx *data.Transactions) ([]*model.Transaction, error) {
	if tx == nil {
		return nil, errors.New("Found nothing")
	}

	if !(len(tx.Transactions) > 0) {
		return nil, errors.New("Found nothing")
	}

	_tx := make([]*model.Transaction, len(tx.Transactions))

	for k, v := range tx.Transactions {
		_v, _ := getGraphQLCompatibleTransaction(ctx, v, false)
		_tx[k] = _v
	}

	return _tx, nil
}

// Converting event data to graphQL compatible data structure
func getGraphQLCompatibleEvent(ctx context.Context, event *data.Event, bookKeeping bool) (*model.Event, error) {
	if event == nil {
		return nil, errors.New("Found nothing")
	}

	data := ""
	if _h := hex.EncodeToString(event.Data); _h != "" && _h != strings.Repeat("0", 64) {
		data = fmt.Sprintf("0x%s", _h)
	}

	return &model.Event{
		Origin:    event.Origin,
		Index:     fmt.Sprintf("%d", event.Index),
		Topics:    getTopicSignaturesAsStringSlice(event.Topics),
		Data:      data,
		TxHash:    event.TransactionHash,
		OrderHash: event.OrderHash,
	}, nil
}

// Converting event array to graphQL compatible data structure
func getGraphQLCompatibleEvents(ctx context.Context, events *data.Events) ([]*model.Event, error) {
	if events == nil {
		return nil, errors.New("Found nothing")
	}

	if !(len(events.Events) > 0) {
		return nil, errors.New("Found nothing")
	}

	_events := make([]*model.Event, len(events.Events))

	for k, v := range events.Events {
		_v, _ := getGraphQLCompatibleEvent(ctx, v, false)
		_events[k] = _v
	}

	return _events, nil
}

func getTopicSignaturesAsStringSlice(topics pq.StringArray) []string {
	_tmp := make([]string, len(topics))

	for k, v := range topics {
		_tmp[k] = v
	}

	return _tmp
}

// FillUpTopicArray - Creates a topic signature array of length
// 4, while putting all elements passed from graphQL query & appending
// empty strings, in remaining places
func FillUpTopicArray(topics []string) []string {

	if len(topics) == 4 {
		return topics
	}

	result := make([]string, 0, 4)
	result = append(result, topics...)

	i := 0
	target := 4 - len(topics)

	for i < target {

		result = append(result, "")
		i++

	}

	return result
}
