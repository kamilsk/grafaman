package graphite

const (
	fromKey  = "from"
	queryKey = "query"
)

type dto struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Leaf int    `json:"leaf"`
}
