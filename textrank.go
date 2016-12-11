package main

const dampeningFactor = 0.85

// scoreNode calculates the voting score for a given node.
func scoreNode(n *node, iterations int) float64 {
	if iterations == 0 {
		return 0
	}

	var successiveScore float64
	for _, linked := range n.Links {
		if len(n.Links) == 0 {
			continue
		}
		successiveScore += scoreNode(linked, iterations-1)
	}
	return (1 - dampeningFactor) + dampeningFactor*successiveScore
}
