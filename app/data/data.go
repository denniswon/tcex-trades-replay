package data

import (
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// SyncState - Whether the service is synced with orderchain or not
type SyncState struct {
	Done                    uint64
	StartedAt               time.Time
	OrderCountAtStartUp     uint64
	MaxOrderNumberAtStartUp uint64
	NewOrdersInserted       uint64
	LatestOrderNumber       uint64
}

// OrderCountInDB - Orders currently present in database
func (s *SyncState) OrderCountInDB() uint64 {
	return s.OrderCountAtStartUp + s.NewOrdersInserted
}

// StatusHolder - Keeps track of progress. To be delivered when `/v1/synced` is queried
type StatusHolder struct {
	State *SyncState
	Mutex *sync.RWMutex
}

// MaxOrderNumberAtStartUp - thread safe read latest order number at the time of service start
// To determine whether a missing order related notification needs to be sent on a pubsub channel or not
func (s *StatusHolder) MaxOrderNumberAtStartUp() uint64 {

	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return s.State.MaxOrderNumberAtStartUp
}

// SetStartedAt - Sets started at time
func (s *StatusHolder) SetStartedAt() {

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.State.StartedAt = time.Now().UTC()
}

// IncrementOrdersInserted - thread safe increments number of orders inserted into DB since start
func (s *StatusHolder) IncrementOrdersInserted() {

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.State.NewOrdersInserted++
}

// IncrementOrdersProcessed - thread safe increments number of orders processed by after it started
func (s *StatusHolder) IncrementOrdersProcessed() {

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.State.Done++
}

// OrderCountInDB - thread safe reads currently present orders in db
func (s *StatusHolder) OrderCountInDB() uint64 {

	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return s.State.OrderCountInDB()
}

// ElapsedTime - thread safe uptime of the service
func (s *StatusHolder) ElapsedTime() time.Duration {

	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return time.Now().UTC().Sub(s.State.StartedAt)
}

// Done - thread safe  #-of Orders processed during uptime i.e. after last time it started
func (s *StatusHolder) Done() uint64 {

	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return s.State.Done
}

// GetLatestOrderNumber - thread safe read latest order number
func (s *StatusHolder) GetLatestOrderNumber() uint64 {

	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return s.State.LatestOrderNumber
}

// SetLatestOrderNumber - thread safe write latest order number
func (s *StatusHolder) SetLatestOrderNumber(num uint64) {

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.State.LatestOrderNumber = num
}

// RedisInfo
type RedisInfo struct {
	Client *redis.Client
	OrderPublishTopic string
}

// ResultStatus
type ResultStatus struct {
	Success uint64
	Failure uint64
}

// Total - Returns total count of operations which were supposed to be performed
//
// To check whether all go routines have sent their status i.e. completed their tasks or not
func (r ResultStatus) Total() uint64 {
	return r.Success + r.Failure
}

// Job - For running a order fetching job
type Job struct {
	DB     *gorm.DB
	Redis  *RedisInfo
	Order  uint64
	Status *StatusHolder
}