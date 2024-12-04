package domain

type Rate struct {
	Ask       float32 `json:"ask"`
	Bid       float32 `json:"bid"`
	Timestamp int64   `json:"timestamp"`
}
