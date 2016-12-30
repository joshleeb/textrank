package textrank

const nodeInitialScore = 1

// node is a node of the graph.
type node struct {
	Data  string
	Links []*node
	Score float64
}

// graph is a graph of nodes.
type graph []*node

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

// newSentenceGraph creates a new graph from sentences. For a graph of
// sentences, the order of the `sentences` array is not important.
func newSentenceGraph(sentences []string) *graph {
	g := &graph{}
	seen := make(map[int]map[int]bool) // to prevent computing similarity twice

	// Add nodes.
	for i, sentence := range sentences {
		g.addNode(sentence, nodeInitialScore)
		seen[i] = make(map[int]bool)
	}

	// Connect nodes based on similarity.
	for a, nodeA := range *g {
		for b, nodeB := range *g {
			if _, ok := seen[a][b]; ok {
				continue
			}
			seen[a][b] = true
			seen[b][a] = true

			similar := similarity(nodeA.Data, nodeB.Data)
			if similar > 1 {
				nodeA.Links = append(nodeA.Links, nodeB)
				nodeB.Links = append(nodeB.Links, nodeA)
			}
		}
	}
	return g
}

// newWordGraph creates a new graph from words. For a graph of words, the order
// of the `words` array are important.
func newWordGraph(words []string) *graph {
	g := &graph{}

	// Connect nodes based on co-occurrence.
	for i := 0; i <= len(words)-lexicalCoOccurenceFactor; i++ {
		cooccurring := words[i : i+lexicalCoOccurenceFactor]
		for _, wordA := range cooccurring {
			nodeA := g.getNode(wordA)
			if nodeA == nil {
				nodeA = g.addNode(wordA, nodeInitialScore)
			}

			for _, wordB := range cooccurring {
				nodeB := g.getNode(wordB)
				if nodeB == nil {
					nodeB = g.addNode(wordB, nodeInitialScore)
				}

				// We don't want nodes to be reflexive.
				if wordA == wordB {
					continue
				}
				nodeA.Links = append(nodeA.Links, nodeB)
			}
		}
	}
	return g
}

// addNode adds a node to the graph, giving it a random score in the range
// [0.0, 1.0).
func (g *graph) addNode(data string, score float64) *node {
	newNode := &node{
		Data:  data,
		Links: []*node{},
		Score: score,
	}
	*g = append(*g, newNode)
	return newNode
}

// getNode gets a node from a graph that has the specified data.
func (g *graph) getNode(data string) *node {
	for _, n := range *g {
		if n.Data == data {
			return n
		}
	}
	return nil
}
