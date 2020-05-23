package graphite

const (
	formatParam = "format"
	fromParam   = "from"
	untilParam  = "until"
	queryParam  = "query"
)

type dto struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Leaf int    `json:"leaf"`
}
