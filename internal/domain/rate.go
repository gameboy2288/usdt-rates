package domain

type Ask struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   Type   `json:"type"`
}

type Rates struct {
	Timestamp int64 `json:"timestamp"`
	Asks      []Ask `json:"asks"`
	Bids      []Ask `json:"bids"`
}

type Type string

const (
	Limit Type = "limit"
)
