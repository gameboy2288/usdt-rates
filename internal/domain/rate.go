package domain

type Rate struct {
	Ask       float64 `json:"ask"`
	Bid       float64 `json:"bid"`
	Timestamp int64   `json:"timestamp"`
}
