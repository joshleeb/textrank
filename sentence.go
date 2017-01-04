package textrank

import (
	"math"
	"strings"
)

// minWordSentence is the minimum number of words a sentence can have to become
// a node in the graph.
const minWordSentence = 5

// RankSentences ranks the sentences in the given text based on the TextRank
// algorithm and returned a list of the ranked sentences in descending order or
// score.
func RankSentences(text string, iterations int) []string {
	graph := &textgraph{}

	// Setup graph.
	seenNodes := make(map[string]bool) // prevent duplication
	for _, token := range tokenizeSentences(text) {
		if _, ok := seenNodes[token]; ok {
			continue
		}
		graph.addNode(token, nodeInitialScore)
		seenNodes[token] = true
	}
	linkSentences(graph)

	// Score sentence nodes.
	for _, node := range *graph {
		node.Score = scoreNode(node, iterations)
	}
	return graph.normalize()
}

// linkSentences links sentence nodes within a graph.
func linkSentences(tg *textgraph) *textgraph {
	seenEdges := make(map[string]map[string]bool) // prevent duplication
	for _, nodeA := range *tg {
		seenEdges[nodeA.Text] = make(map[string]bool)

		for _, nodeB := range *tg {
			_, okA := seenEdges[nodeA.Text][nodeB.Text]
			_, okB := seenEdges[nodeB.Text][nodeA.Text]
			if okA || okB {
				continue
			}
			seenEdges[nodeA.Text][nodeB.Text] = true

			// Connect nodes based on similarity.
			if sentenceSimilarity(nodeA.Text, nodeB.Text) > 1 {
				nodeA.Edges = append(nodeA.Edges, nodeB)
				nodeB.Edges = append(nodeB.Edges, nodeA)
			}
		}
	}
	return tg
}

// sentenceSimilarity calculates the similarity between two sentences,
// normalizing for sentence length.
func sentenceSimilarity(a, b string) float64 {
	tokensA := tokenizeWords(a)
	tokensB := tokenizeWords(b)

	if len(tokensA) < minWordSentence || len(tokensB) < minWordSentence {
		return 0
	}

	similarWords := make(map[string]bool)
	for _, tokenA := range tokensA {
		wordA := strings.ToLower(tokenA)

		// Ignore stopwords. Only need to check wordA because if wordA is not a
		// stopword and wordB is a stopword, then they are not going to match.
		if _, ok := stopwords[wordA]; ok {
			continue
		}

		for _, tokenB := range tokensB {
			wordB := strings.ToLower(tokenB)
			if strings.Compare(wordA, wordB) == 0 {
				similarWords[wordA] = true
			}
		}
	}

	numSimilarWords := float64(len(similarWords))
	numWordsMult := float64(len(tokensA) * len(tokensB))

	if numWordsMult == 1 {
		return 0
	}

	return numSimilarWords / math.Log(numWordsMult)
}
