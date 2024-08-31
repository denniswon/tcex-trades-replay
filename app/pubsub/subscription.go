package pubsub

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

// UploadHeader
type UploadHeader struct {
	ID       string `json:"id"`
	Filepath string `json:"filepath"`
	Size     int64  `json:"size"` // size of the file if uploaded
}

func (header *UploadHeader) Generate() *UploadHeader {

	if header.ID == "" {
		header.ID = uuid.New().String()
	}

	return header
}

func (header *UploadHeader) String() string {

	return fmt.Sprintf(`{"id":%s,"filepath":%s,"size":%d}`,
		header.ID,
		header.Filepath,
		header.Size,
	)

}

// SubscriptionRequest
type SubscriptionRequest struct {
	ID          string  `json:"id"`
	Filename    string  `json:"filename"`
	ReplayRate  float32 `json:"replay_rate"`
	Type        string  `json:"type"`
	Name        string  `json:"name"` // "order" or "kline"
	Granularity uint16  `json:"granularity"`
}

func (req *SubscriptionRequest) Generate() *SubscriptionRequest {

	if req.ID == "" {
		req.ID = uuid.New().String()
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

	// Check if file exists
	if _, err := os.Stat(req.Filename); err != nil {

		log.Printf("Request input file does not exist : %s\n", req.Filename)

		ret = false

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
