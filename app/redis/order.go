package data

import (
	"encoding/json"
	"fmt"
	"log"
)

// Order - Order related info to be delivered to client in this format
type Order struct {
	RequestId           string  `json:"request_id"`
	Price               float64 `json:"price"`
	Quantity            uint64  `json:"quantity"`
	Aggressor           string  `json:"aggressor"`
	Timestamp           uint64  `json:"timestamp"`
}

// MarshalBinary - Implementing binary marshalling function, to be invoked
// by redis before publishing data on channel
func (b *Order) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

// MarshalJSON - Custom JSON encoder
func (b *Order) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"request_id":%s,"price":%.2f,"quantity":%d,"aggressor":%s,"timestamp":%d}`,
		b.RequestId,
		b.Price,
		b.Quantity,
		b.Aggressor,
		b.Timestamp,
	)), nil
}

// ToJSON - Encodes into JSON, to be supplied when queried for order data
func (b *Order) ToJSON() []byte {
	data, err := json.Marshal(b)
	if err != nil {
		log.Printf("[!] Failed to encode order data to JSON : %s\n", err.Error())
		return nil
	}

	return data
}

// Orders - A set of orders to be held, extracted from DB query result
// also to be supplied to client in JSON encoded form
type Orders struct {
	Orders []*Order `json:"orders"`
}

// ToJSON - Encoding into JSON, to be invoked when delivering query result to client
func (b *Orders) ToJSON() []byte {
	data, err := json.Marshal(b)
	if err != nil {
		log.Printf("[!] Failed to encode order data to JSON : %s\n", err.Error())
		return nil
	}

	return data
}
