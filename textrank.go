package textrank

// dampeningFactor is used when scoring text nodes, and is in the range (0, 1].
const dampeningFactor = 0.85

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
