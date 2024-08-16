package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"

	cmn "github.com/denniswon/tcex/app/common"
	cfg "github.com/denniswon/tcex/app/config"
	_db "github.com/denniswon/tcex/app/db"
	"github.com/denniswon/tcex/app/rest/graph/generated"
	"github.com/denniswon/tcex/app/rest/graph/model"
	"github.com/ethereum/go-ethereum/common"
)

func (r *queryResolver) OrderByHash(ctx context.Context, hash string) (*model.Order, error) {
	if !(strings.HasPrefix(hash, "0x") && len(hash) == 66) {
		return nil, errors.New("Bad Order Hash")
	}

	return getGraphQLCompatibleOrder(ctx, _db.GetOrderByHash(db, common.HexToHash(hash)))
}

func (r *queryResolver) OrderByNumber(ctx context.Context, number string) (*model.Order, error) {
	_number, err := strconv.ParseUint(number, 10, 64)
	if err != nil {
		return nil, errors.New("Bad Order Number")
	}

	return getGraphQLCompatibleOrder(ctx, _db.GetOrderByNumber(db, _number))
}

func (r *queryResolver) OrdersByNumberRange(ctx context.Context, from string, to string) ([]*model.Order, error) {
	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetOrderNumberRange())
	if err != nil {
		return nil, errors.New("Bad Order Number Range")
	}

	return getGraphQLCompatibleOrders(ctx, _db.GetOrdersByNumberRange(db, _from, _to))
}

func (r *queryResolver) OrdersByTimeRange(ctx context.Context, from string, to string) ([]*model.Order, error) {
	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetTimeRange())
	if err != nil {
		return nil, errors.New("Bad Order Timestamp Range")
	}

	return getGraphQLCompatibleOrders(ctx, _db.GetOrdersByTimeRange(db, _from, _to))
}

func (r *queryResolver) Transaction(ctx context.Context, hash string) (*model.Transaction, error) {
	if !(strings.HasPrefix(hash, "0x") && len(hash) == 66) {
		return nil, errors.New("Bad Transaction Hash")
	}

	return getGraphQLCompatibleTransaction(ctx, _db.GetTransactionByHash(db, common.HexToHash(hash)), true)
}

func (r *queryResolver) TransactionCountByOrderHash(ctx context.Context, hash string) (int, error) {
	if !(strings.HasPrefix(hash, "0x") && len(hash) == 66) {
		return 0, errors.New("Bad Order Hash")
	}

	count := int(_db.GetTransactionCountByOrderHash(db, common.HexToHash(hash)))

	// Attempting to calculate byte form of number
	// so that we can keep track of how much data was transferred
	// to client
	_count := make([]byte, 4)
	binary.LittleEndian.PutUint32(_count, uint32(count))

	return count, nil
}

func (r *queryResolver) TransactionsByOrderHash(ctx context.Context, hash string) ([]*model.Transaction, error) {
	if !(strings.HasPrefix(hash, "0x") && len(hash) == 66) {
		return nil, errors.New("Bad Order Hash")
	}

	return getGraphQLCompatibleTransactions(ctx, _db.GetTransactionsByOrderHash(db, common.HexToHash(hash)))
}

func (r *queryResolver) TransactionCountByOrderNumber(ctx context.Context, number string) (int, error) {
	_number, err := strconv.ParseUint(number, 10, 64)
	if err != nil {
		return 0, errors.New("Bad Order Number")
	}

	count := int(_db.GetTransactionCountByOrderNumber(db, _number))

	// Attempting to calculate byte form of number
	// so that we can keep track of how much data was transferred
	// to client
	_count := make([]byte, 4)
	binary.LittleEndian.PutUint32(_count, uint32(count))

	return count, nil
}

func (r *queryResolver) TransactionsByOrderNumber(ctx context.Context, number string) ([]*model.Transaction, error) {
	_number, err := strconv.ParseUint(number, 10, 64)
	if err != nil {
		return nil, errors.New("Bad Order Number")
	}

	return getGraphQLCompatibleTransactions(ctx, _db.GetTransactionsByOrderNumber(db, _number))
}

func (r *queryResolver) TransactionCountFromAccountByNumberRange(ctx context.Context, account string, from string, to string) (int, error) {
	if !(strings.HasPrefix(account, "0x") && len(account) == 42) {
		return 0, errors.New("Bad Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetOrderNumberRange())
	if err != nil {
		return 0, errors.New("Bad Order Number Range")
	}

	count := int(_db.GetTransactionCountFromAccountByOrderNumberRange(db, common.HexToAddress(account), _from, _to))

	// Attempting to calculate byte form of number
	// so that we can keep track of how much data was transferred
	// to client
	_count := make([]byte, 4)
	binary.LittleEndian.PutUint32(_count, uint32(count))

	return count, nil
}

func (r *queryResolver) TransactionsFromAccountByNumberRange(ctx context.Context, account string, from string, to string) ([]*model.Transaction, error) {
	if !(strings.HasPrefix(account, "0x") && len(account) == 42) {
		return nil, errors.New("Bad Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetOrderNumberRange())
	if err != nil {
		return nil, errors.New("Bad Order Number Range")
	}

	return getGraphQLCompatibleTransactions(ctx, _db.GetTransactionsFromAccountByOrderNumberRange(db, common.HexToAddress(account), _from, _to))
}

func (r *queryResolver) TransactionCountFromAccountByTimeRange(ctx context.Context, account string, from string, to string) (int, error) {
	if !(strings.HasPrefix(account, "0x") && len(account) == 42) {
		return 0, errors.New("Bad Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetTimeRange())
	if err != nil {
		return 0, errors.New("Bad Order Timestamp Range")
	}

	count := int(_db.GetTransactionCountFromAccountByOrderTimeRange(db, common.HexToAddress(account), _from, _to))

	// Attempting to calculate byte form of number
	// so that we can keep track of how much data was transferred
	// to client
	_count := make([]byte, 4)
	binary.LittleEndian.PutUint32(_count, uint32(count))

	return count, nil
}

func (r *queryResolver) TransactionsFromAccountByTimeRange(ctx context.Context, account string, from string, to string) ([]*model.Transaction, error) {
	if !(strings.HasPrefix(account, "0x") && len(account) == 42) {
		return nil, errors.New("Bad Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetTimeRange())
	if err != nil {
		return nil, errors.New("Bad Order Timestamp Range")
	}

	return getGraphQLCompatibleTransactions(ctx, _db.GetTransactionsFromAccountByOrderTimeRange(db, common.HexToAddress(account), _from, _to))
}

func (r *queryResolver) TransactionCountToAccountByNumberRange(ctx context.Context, account string, from string, to string) (int, error) {
	if !(strings.HasPrefix(account, "0x") && len(account) == 42) {
		return 0, errors.New("Bad Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetOrderNumberRange())
	if err != nil {
		return 0, errors.New("Bad Order Number Range")
	}

	count := int(_db.GetTransactionCountToAccountByOrderNumberRange(db, common.HexToAddress(account), _from, _to))

	// Attempting to calculate byte form of number
	// so that we can keep track of how much data was transferred
	// to client
	_count := make([]byte, 4)
	binary.LittleEndian.PutUint32(_count, uint32(count))

	return count, nil
}

func (r *queryResolver) TransactionsToAccountByNumberRange(ctx context.Context, account string, from string, to string) ([]*model.Transaction, error) {
	if !(strings.HasPrefix(account, "0x") && len(account) == 42) {
		return nil, errors.New("Bad Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetOrderNumberRange())
	if err != nil {
		return nil, errors.New("Bad Order Number Range")
	}

	return getGraphQLCompatibleTransactions(ctx, _db.GetTransactionsToAccountByOrderNumberRange(db, common.HexToAddress(account), _from, _to))
}

func (r *queryResolver) TransactionCountToAccountByTimeRange(ctx context.Context, account string, from string, to string) (int, error) {
	if !(strings.HasPrefix(account, "0x") && len(account) == 42) {
		return 0, errors.New("Bad Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetTimeRange())
	if err != nil {
		return 0, errors.New("Bad Order Timestamp Range")
	}

	count := int(_db.GetTransactionCountToAccountByOrderTimeRange(db, common.HexToAddress(account), _from, _to))

	// Attempting to calculate byte form of number
	// so that we can keep track of how much data was transferred
	// to client
	_count := make([]byte, 4)
	binary.LittleEndian.PutUint32(_count, uint32(count))

	return count, nil
}

func (r *queryResolver) TransactionsToAccountByTimeRange(ctx context.Context, account string, from string, to string) ([]*model.Transaction, error) {
	if !(strings.HasPrefix(account, "0x") && len(account) == 42) {
		return nil, errors.New("Bad Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetTimeRange())
	if err != nil {
		return nil, errors.New("Bad Order Timestamp Range")
	}

	return getGraphQLCompatibleTransactions(ctx, _db.GetTransactionsToAccountByOrderTimeRange(db, common.HexToAddress(account), _from, _to))
}

func (r *queryResolver) TransactionCountBetweenAccountsByNumberRange(ctx context.Context, fromAccount string, toAccount string, from string, to string) (int, error) {
	if !(strings.HasPrefix(fromAccount, "0x") && len(fromAccount) == 42) {
		return 0, errors.New("Bad From Account Address")
	}

	if !(strings.HasPrefix(toAccount, "0x") && len(toAccount) == 42) {
		return 0, errors.New("Bad To Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetOrderNumberRange())
	if err != nil {
		return 0, errors.New("Bad Order Number Range")
	}

	count := int(_db.GetTransactionCountBetweenAccountsByOrderNumberRange(db, common.HexToAddress(fromAccount), common.HexToAddress(toAccount), _from, _to))

	// Attempting to calculate byte form of number
	// so that we can keep track of how much data was transferred
	// to client
	_count := make([]byte, 4)
	binary.LittleEndian.PutUint32(_count, uint32(count))

	return count, nil
}

func (r *queryResolver) TransactionsBetweenAccountsByNumberRange(ctx context.Context, fromAccount string, toAccount string, from string, to string) ([]*model.Transaction, error) {
	if !(strings.HasPrefix(fromAccount, "0x") && len(fromAccount) == 42) {
		return nil, errors.New("Bad From Account Address")
	}

	if !(strings.HasPrefix(toAccount, "0x") && len(toAccount) == 42) {
		return nil, errors.New("Bad To Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetOrderNumberRange())
	if err != nil {
		return nil, errors.New("Bad Order Number Range")
	}

	return getGraphQLCompatibleTransactions(ctx, _db.GetTransactionsBetweenAccountsByOrderNumberRange(db, common.HexToAddress(fromAccount), common.HexToAddress(toAccount), _from, _to))
}

func (r *queryResolver) TransactionCountBetweenAccountsByTimeRange(ctx context.Context, fromAccount string, toAccount string, from string, to string) (int, error) {
	if !(strings.HasPrefix(fromAccount, "0x") && len(fromAccount) == 42) {
		return 0, errors.New("Bad From Account Address")
	}

	if !(strings.HasPrefix(toAccount, "0x") && len(toAccount) == 42) {
		return 0, errors.New("Bad To Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetTimeRange())
	if err != nil {
		return 0, errors.New("Bad Order Timestamp Range")
	}

	count := int(_db.GetTransactionCountBetweenAccountsByOrderTimeRange(db, common.HexToAddress(fromAccount), common.HexToAddress(toAccount), _from, _to))

	// Attempting to calculate byte form of number
	// so that we can keep track of how much data was transferred
	// to client
	_count := make([]byte, 4)
	binary.LittleEndian.PutUint32(_count, uint32(count))

	return count, nil
}

func (r *queryResolver) TransactionsBetweenAccountsByTimeRange(ctx context.Context, fromAccount string, toAccount string, from string, to string) ([]*model.Transaction, error) {
	if !(strings.HasPrefix(fromAccount, "0x") && len(fromAccount) == 42) {
		return nil, errors.New("Bad From Account Address")
	}

	if !(strings.HasPrefix(toAccount, "0x") && len(toAccount) == 42) {
		return nil, errors.New("Bad To Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetTimeRange())
	if err != nil {
		return nil, errors.New("Bad Order Timestamp Range")
	}

	return getGraphQLCompatibleTransactions(ctx, _db.GetTransactionsBetweenAccountsByOrderTimeRange(db, common.HexToAddress(fromAccount), common.HexToAddress(toAccount), _from, _to))
}

func (r *queryResolver) ContractsCreatedFromAccountByNumberRange(ctx context.Context, account string, from string, to string) ([]*model.Transaction, error) {
	if !(strings.HasPrefix(account, "0x") && len(account) == 42) {
		return nil, errors.New("Bad Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetOrderNumberRange())
	if err != nil {
		return nil, errors.New("Bad Order Number Range")
	}

	return getGraphQLCompatibleTransactions(ctx, _db.GetContractCreationTransactionsFromAccountByOrderNumberRange(db, common.HexToAddress(account), _from, _to))
}

func (r *queryResolver) ContractsCreatedFromAccountByTimeRange(ctx context.Context, account string, from string, to string) ([]*model.Transaction, error) {
	if !(strings.HasPrefix(account, "0x") && len(account) == 42) {
		return nil, errors.New("Bad Account Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetTimeRange())
	if err != nil {
		return nil, errors.New("Bad Order Timestamp Range")
	}

	return getGraphQLCompatibleTransactions(ctx, _db.GetContractCreationTransactionsFromAccountByOrderTimeRange(db, common.HexToAddress(account), _from, _to))
}

func (r *queryResolver) TransactionFromAccountWithNonce(ctx context.Context, account string, nonce string) (*model.Transaction, error) {
	if !(strings.HasPrefix(account, "0x") && len(account) == 42) {
		return nil, errors.New("Bad Account Address")
	}

	_nonce, err := strconv.ParseUint(nonce, 10, 64)
	if err != nil {
		return nil, errors.New("Bad Account Nonce")
	}

	return getGraphQLCompatibleTransaction(ctx, _db.GetTransactionFromAccountWithNonce(db, common.HexToAddress(account), _nonce), true)
}

func (r *queryResolver) EventsFromContractByNumberRange(ctx context.Context, contract string, from string, to string) ([]*model.Event, error) {
	if !(strings.HasPrefix(contract, "0x") && len(contract) == 42) {
		return nil, errors.New("Bad Contract Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetOrderNumberRange())
	if err != nil {
		return nil, errors.New("Bad Order Number Range")
	}

	return getGraphQLCompatibleEvents(ctx, _db.GetEventsFromContractByOrderNumberRange(db, common.HexToAddress(contract), _from, _to))
}

func (r *queryResolver) EventsFromContractByTimeRange(ctx context.Context, contract string, from string, to string) ([]*model.Event, error) {
	if !(strings.HasPrefix(contract, "0x") && len(contract) == 42) {
		return nil, errors.New("Bad Contract Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetTimeRange())
	if err != nil {
		return nil, errors.New("Bad Order Timestamp Range")
	}

	return getGraphQLCompatibleEvents(ctx, _db.GetEventsFromContractByOrderTimeRange(db, common.HexToAddress(contract), _from, _to))
}

func (r *queryResolver) EventsByOrderHash(ctx context.Context, hash string) ([]*model.Event, error) {
	if !(strings.HasPrefix(hash, "0x") && len(hash) == 66) {
		return nil, errors.New("Bad Order Hash")
	}

	return getGraphQLCompatibleEvents(ctx, _db.GetEventsByOrderHash(db, common.HexToHash(hash)))
}

func (r *queryResolver) EventsByTxHash(ctx context.Context, hash string) ([]*model.Event, error) {
	if !(strings.HasPrefix(hash, "0x") && len(hash) == 66) {
		return nil, errors.New("Bad Transaction Hash")
	}

	return getGraphQLCompatibleEvents(ctx, _db.GetEventsByTransactionHash(db, common.HexToHash(hash)))
}

func (r *queryResolver) EventsFromContractWithTopicsByNumberRange(ctx context.Context, contract string, from string, to string, topics []string) ([]*model.Event, error) {
	if !(strings.HasPrefix(contract, "0x") && len(contract) == 42) {
		return nil, errors.New("Bad Contract Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetOrderNumberRange())
	if err != nil {
		return nil, errors.New("Bad Order Number Range")
	}

	return getGraphQLCompatibleEvents(ctx, _db.GetEventsFromContractWithTopicsByOrderNumberRange(db, common.HexToAddress(contract), _from, _to, cmn.CreateEventTopicMap(FillUpTopicArray(topics))))
}

func (r *queryResolver) EventsFromContractWithTopicsByTimeRange(ctx context.Context, contract string, from string, to string, topics []string) ([]*model.Event, error) {
	if !(strings.HasPrefix(contract, "0x") && len(contract) == 42) {
		return nil, errors.New("Bad Contract Address")
	}

	_from, _to, err := cmn.RangeChecker(from, to, cfg.GetTimeRange())
	if err != nil {
		return nil, errors.New("Bad Order Timestamp Range")
	}

	return getGraphQLCompatibleEvents(ctx, _db.GetEventsFromContractWithTopicsByOrderTimeRange(db, common.HexToAddress(contract), _from, _to, cmn.CreateEventTopicMap(FillUpTopicArray(topics))))
}

func (r *queryResolver) LastXEventsFromContract(ctx context.Context, contract string, x int) ([]*model.Event, error) {
	if !(strings.HasPrefix(contract, "0x") && len(contract) == 42) {
		return nil, errors.New("Bad Contract Address")
	}

	if !(x <= 50) {
		return nil, errors.New("Too Many Events Requested")
	}

	return getGraphQLCompatibleEvents(ctx, _db.GetLastXEventsFromContract(db, common.HexToAddress(contract), x))
}

func (r *queryResolver) EventByOrderHashAndLogIndex(ctx context.Context, hash string, index string) (*model.Event, error) {
	if !(strings.HasPrefix(hash, "0x") && len(hash) == 66) {
		return nil, errors.New("Bad Order Hash")
	}

	_index, err := strconv.ParseUint(index, 10, 64)
	if err != nil {
		return nil, errors.New("Bad Log Index")
	}

	return getGraphQLCompatibleEvent(ctx, _db.GetEventByOrderHashAndLogIndex(db, common.HexToHash(hash), uint(_index)), true)
}

func (r *queryResolver) EventByOrderNumberAndLogIndex(ctx context.Context, number string, index string) (*model.Event, error) {
	_number, err := strconv.ParseUint(number, 10, 64)
	if err != nil {
		return nil, errors.New("Bad Order Number")
	}

	_index, err := strconv.ParseUint(index, 10, 64)
	if err != nil {
		return nil, errors.New("Bad Log Index")
	}

	return getGraphQLCompatibleEvent(ctx, _db.GetEventByOrderNumberAndLogIndex(db, _number, uint(_index)), true)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) TransactionByHash(ctx context.Context, hash string) (*model.Transaction, error) {
	if !(strings.HasPrefix(hash, "0x") && len(hash) == 66) {
		return nil, errors.New("Bad Transaction Hash")
	}

	return getGraphQLCompatibleTransaction(ctx, _db.GetTransactionByHash(db, common.HexToHash(hash)), true)
}
