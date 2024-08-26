package data

import (
	"encoding/json"
	"fmt"
	"log"
)

// Kline - OHLCV data for the orders for a specific granularity
type Kline struct {
	Timestamp						int64 	`json:"timestamp"` 		// bucket start time in unix timestamp
	Low 								float64 `json:"low"`  			  // lowest price during the bucket interval
	High 								float64 `json:"high"`  				// highest price during the bucket interval
	Open								float64 `json:"open"`  				// opening price (first trade) in the bucket interval
	Close 							float64 `json:"close"`  			// closing price (last trade) in the bucket interval
	Volume 							int64 	`json:"volume"`  			// volume of trading activity during the bucket interval
	Granularity         uint16  `json:"granularity"`	// granularity field is in "seconds"
}

// MarshalBinary - Implementing binary marshalling function, to be invoked
// by redis before publishing data on channel
func (k *Kline) MarshalBinary() ([]byte, error) {
	return json.Marshal(k)
}

// MarshalJSON - Custom JSON encoder
func (k *Kline) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"timestamp":%d,"low":%f,"high":%f,"open":%f,"close":%f,"volume":%d,"granularity":%d}`,
		k.Timestamp,
		k.Low,
		k.High,
		k.Open,
		k.Close,
		k.Volume,
		k.Granularity,
	)), nil
}

// ToJSON - Encodes into JSON, to be supplied when queried for order kline data
func (k *Kline) ToJSON() []byte {
	data, err := json.Marshal(k)
	if err != nil {
		log.Printf("[!] Failed to encode order kline data to JSON : %s\n", err.Error())
		return nil
	}

	return data
}
