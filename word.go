package textrank

// lexicalCoOccurenceFactor is used when linking text nodes containing words,
// used for keyword ranking. It should be in the range [2, 10].
const lexicalCoOccurenceFactor = 2

// RankWords ranks the words in the given text based on the TextRank algorithm
// and returned a list of the ranked words in descending order or score.
func RankWords(text string, iterations int) []string {
	graph := &textgraph{}

	// Setup graph.
	for _, word := range cleanText(text) {
		if _, ok := stopwords[word]; ok {
			continue
		}
		graph.addNode(word, nodeInitialScore)
	}
	linkWords(graph, tokenizeSentences(text))

	// Score sentence nodes.
	for _, node := range *graph {
		node.Score = scoreNode(node, iterations)
	}
	return graph.normalize()
}

// linkWords links word nodes within a graph.
func linkWords(tg *textgraph, sentences []string) *textgraph {
	for _, sentence := range sentences {
		words := cleanText(sentence)

		// Connect nodes based on co-occurrence.
		for i := 0; i <= len(words)-lexicalCoOccurenceFactor; i++ {
			cooccurring := words[i : i+lexicalCoOccurenceFactor]
			for _, wordA := range cooccurring {
				nodeA := tg.getNode(wordA)
				if nodeA == nil {
					continue
				}

				for _, wordB := range cooccurring {
					nodeB := tg.getNode(wordB)
					if nodeB == nil {
						continue
					}

					// Prevent nodes being reflexive.
					if wordA == wordB {
						continue
					}
					nodeA.Edges = append(nodeA.Edges, nodeB)
				}
			}
		}
	}
	return tg
}

// cleanText tokenizes the text into words and cleans it by removing stopwords.
func cleanText(text string) []string {
	tokens := []string{}
	for _, word := range tokenizeWords(text) {
		if _, ok := stopwords[word]; ok {
			continue
		}
		tokens = append(tokens, word)
	}
	return tokens
}
