package textrank

import "sort"

// nodeInitialScore is the initial score of the node.
const nodeInitialScore = 1

// dampeningFactor is used when scoring text nodes, and is in the range (0, 1].
const dampeningFactor = 0.85

// textgraph is a textrank graph.
type textgraph []*node

// node is a node in the textrank graph.
type node struct {
	Text  string
	Edges []*node
	Score float64
}

// newTextGraph creates a new graph (undirected and unweighted) from provided
// text.
func newGraph(tokens []string) *textgraph {
	newGraph := &textgraph{}
	seenNodes := make(map[string]bool) // prevent duplication
	for _, token := range tokens {
		if _, ok := seenNodes[token]; ok {
			continue
		}
		newGraph.addNode(token, nodeInitialScore)
		seenNodes[token] = true
	}
	return newGraph
}

// addNode adds a node to the graph, giving it a random score in the range
// [0.0, 1.0).
func (tg *textgraph) addNode(text string, initialScore float64) *node {
	newNode := &node{Text: text, Edges: []*node{}, Score: initialScore}
	*tg = append(*tg, newNode)
	return newNode
}

// getNode gets a node from the graph and `nil` if it doesn't exist.
func (tg *textgraph) getNode(text string) *node {
	for _, node := range *tg {
		if node.Text == text {
			return node
		}
	}
	return nil
}

// scoreNode calculates the voting score for a given node.
func scoreNode(n *node, iterations int) float64 {
	if iterations == 0 {
		return 0
	}

	var successiveScore float64
	for _, edges := range n.Edges {
		successiveScore += scoreNode(edges, iterations-1)
	}
	return (1 - dampeningFactor) + dampeningFactor*successiveScore
}

// normalize normalizes the graph into a list containing the string of each
// node, ordered in descending order by the score of the node.
func (tg *textgraph) normalize() []string {
	ranked := []string{}
	sort.Sort(sort.Reverse(tg))
	for _, node := range *tg {
		ranked = append(ranked, node.Text)
	}
	return ranked
}

// Len returns the number of nodes in the graph.
func (tg textgraph) Len() int {
	return len(tg)
}

// Less returns true if the `Score` in the supplied slice at index `i` is less
// than the `Score` at index `j`.
func (tg textgraph) Less(i, j int) bool {
	return tg[i].Score < tg[j].Score
}

// Swap swaps elements at the indexes of `i` and `j` in the provided graph.
func (tg textgraph) Swap(i, j int) {
	tg[i], tg[j] = tg[j], tg[i]
}
