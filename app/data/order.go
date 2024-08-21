package data

import (
	"encoding/json"
	"fmt"
	"log"
)

// Order - Order related info to be delivered to client in this format
type Order struct {
	Price               string `json:"price"`
	Quantity            uint64 `json:"quantity"`
	Aggressor           string `json:"aggressor"`
	Timestamp           int64  `json:"timestamp"`
}

// MarshalBinary - Implementing binary marshalling function, to be invoked
// by redis before publishing data on channel
func (b *Order) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

// MarshalJSON - Custom JSON encoder
func (b *Order) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"price":%q,"quantity":%d,"aggressor":%q,"timestamp":%d}`,
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
