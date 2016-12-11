package main

import (
	"math/rand"
	"time"
)

// graph is a graph of nodes.
type graph []*node

// node is a node of the graph.
type node struct {
	Data  string
	Links []*node
	Score float64
}

// Len returns the number of nodes in the graph.
func (g graph) Len() int {
	return len(g)
}

// Less returns true if the `Score` in the supplied slice at index `i` is less
// than the `Score` at index `j`.
func (g graph) Less(i, j int) bool {
	return g[i].Score < g[j].Score
}

// Swap swaps elements at the indexes of `i` and `j` in the provided graph.
func (g graph) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
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
			if similar > 0 {
				nodeA.Links = append(nodeA.Links, nodeB)
				nodeB.Links = append(nodeB.Links, nodeA)
			}
		}
	}
	return g
}

// AddNode adds a node to the graph, giving it a random score in the range
// [0.0, 1.0).
func (g *graph) AddNode(data string) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	*g = append(*g, &node{
		Data:  data,
		Links: []*node{},
		Score: random.Float64(),
	})
}
