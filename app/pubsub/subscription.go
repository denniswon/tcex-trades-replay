package pubsub

import (
	"fmt"

	"github.com/google/uuid"
)

// SubscriptionRequest
type SubscriptionRequest struct {
	ID          string  `json:"id"`
	Filename    string  `json:"filename"`
	Size        int     `json:"size"` // size of the file if uploaded
	ReplayRate  float32 `json:"replay_rate"`
	Type        string  `json:"type"`
	Name        string  `json:"name"` // "order" or "kline"
	Granularity uint16  `json:"granularity"`
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

	if req.Name == "kline" && req.Granularity == 0 {
		req.Granularity = 60
	}

	return req
}

func (req *SubscriptionRequest) Validate() bool {
	ret := req.Filename != "" && req.ReplayRate > 0.0 && req.ID != "" && req.Name != ""
	if req.Name == "kline" {
		ret = ret && req.Granularity > 0
	}
	return ret
}

func (req *SubscriptionRequest) String() string {
	if req.Name == "kline" {
		return fmt.Sprintf(`{"request_id":%s,"filename":%s,"replay_rate":%f,"name":%s,"granularity":%d}`,
			req.ID,
			req.Filename,
			req.ReplayRate,
			req.Name,
			req.Granularity,
		)
	}

	return fmt.Sprintf(`{"request_id":%s,"filename":%s,"replay_rate":%f,"name":%s}`,
		req.ID,
		req.Filename,
		req.ReplayRate,
		req.Name,
	)
}

// SubscriptionResponse
type SubscriptionResponse struct {
	Code    uint   `json:"code"`
	ID      string `json:"id"`
	Message string `json:"msg"`
}
