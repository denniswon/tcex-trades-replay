package pubsub

import (
	"fmt"

	"github.com/google/uuid"
)

// SubscriptionRequest
type SubscriptionRequest struct {
	ID       		string `json:"id"`
	Filename   	string `json:"filename"`
	ReplayRate 	float32 `json:"replay_rate"`
	Type    		string `json:"type"`
}

func (req SubscriptionRequest) Generate() SubscriptionRequest {

	if req.ID == "" {
		req.ID = uuid.New().String()
	}

	if req.Filename == "" {
		req.Filename = "trades.txt"
	}

	if req.ReplayRate == 0.0 {
		req.ReplayRate = 60.0
	}

	return req
}


func (req *SubscriptionRequest) Validate() bool {
	return req.Filename != "" && req.ReplayRate > 0.0 && req.ID != ""
}

func (req *SubscriptionRequest) String() string {
	return fmt.Sprintf(`{"request_id":%s,"filename":%s,"replay_rate":%f}`,
		req.ID,
		req.Filename,
		req.ReplayRate,
	)
}


// SubscriptionResponse
type SubscriptionResponse struct {
	Code    uint   `json:"code"`
	ID    	string `json:"id"`
	Message string `json:"msg"`
}
