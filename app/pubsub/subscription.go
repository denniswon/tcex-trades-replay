package pubsub

// SubscriptionRequest
type SubscriptionRequest struct {
	ID       		string `json:"id"`
	Filename   	string `json:"filename"`
	ReplayRate 	uint64 `json:"replay_rate"`
	Type    		string `json:"type"`
}

// SubscriptionResponse
type SubscriptionResponse struct {
	Code    uint   `json:"code"`
	ID    	string `json:"id"`
	Message string `json:"msg"`
}
