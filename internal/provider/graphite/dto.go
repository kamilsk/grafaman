package graphite

const (
	formatParam = "format"
	fromParam   = "from"
	untilParam  = "until"
	queryParam  = "query"
)

type dto struct {
	ID   string `json:"id,omitempty"`
	Text string `json:"text,omitempty"`
	Leaf int    `json:"leaf,omitempty"`
}
