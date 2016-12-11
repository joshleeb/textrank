package main

// graph is a graph of nodes.
type graph []*node

// node is a node of the graph.
type node struct {
	Data  string
	Links []*link
}

type link struct {
	To     *node
	Weight float64
}

// newGraph creates a new graph from sentences.
// NOTE: we are assuming that there are no duplicate sentences.
func newGraph(sentences []string) *graph {
	g := &graph{}
	seen := make(map[int]map[int]bool) // to prevent computing similarity twice

	// Add nodes.
	for i, sentence := range sentences {
		g.AddNode(sentence)
		seen[i] = make(map[int]bool)
	}

	// Connect nodes.
	for a, nodeA := range *g {
		for b, nodeB := range *g {
			if _, ok := seen[a][b]; ok {
				continue
			}
			seen[a][b] = true
			seen[b][a] = true

			similar := similarity(nodeA.Data, nodeB.Data)
			if similar >= 1 {
				nodeA.Links = append(nodeA.Links, &link{
					To: nodeB, Weight: similar,
				})
				nodeB.Links = append(nodeB.Links, &link{
					To: nodeA, Weight: similar,
				})
			}
		}
	}
	return g
}

func (g *graph) AddNode(data string) {
	*g = append(*g, &node{
		Data: data, Links: []*link{},
	})
}
