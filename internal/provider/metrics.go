package provider

type Node struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Leaf int    `json:"leaf"`
}

type Nodes []Node

func (nodes Nodes) Len() int {
	return len(nodes)
}

func (nodes Nodes) Less(i, j int) bool {
	return nodes[i].ID < nodes[j].ID
}

func (nodes Nodes) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}

func (nodes Nodes) OnlyLeafs() Nodes {
	filtered := make(Nodes, 0, 8)
	for _, node := range nodes {
		if node.Leaf == 1 {
			filtered = append(filtered, node)
		}
	}
	return filtered
}
