package data

import (
	"encoding/json"
	"fmt"
	"log"
)

// EOF - Replay EOF info to be delivered to client in this format
type EOF struct {
	RequestID string `json:"request_id"`
}

// MarshalBinary - Implementing binary marshalling function, to be invoked
// by redis before publishing data on channel
func (e *EOF) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

// MarshalJSON - Custom JSON encoder
func (e *EOF) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"request_id":%q}`,
		e.RequestID,
	)), nil
}

// ToJSON - Encodes into JSON, to be supplied when queried for eof data
func (e *EOF) ToJSON() []byte {
	data, err := json.Marshal(e)
	if err != nil {
		log.Printf("[!] Failed to encode EOF data to JSON : %s\n", err.Error())
		return nil
	}

	return data
}
