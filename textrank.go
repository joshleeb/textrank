package textrank

import (
	"math"
	"sort"
	"strings"

	"github.com/neurosnap/sentences"
	"github.com/neurosnap/sentences/english"
)

const dampeningFactor = 0.85

// Rank ranks the sentences in the given text based on the TextRank algorithm
// and returned a list of the ranked sentences in descending order or score.
func Rank(text string, iterations int) []string {
	sentences := tokenize(text)
	graph := newGraph(sentences)
	ranked := []string{}

	// Iterating 5 times was chosen based on the convergence curves in Figure 1
	// of "TextRank: Bringing Order into Texts" by Rada Mihalcea and Paul Tarau,
	// 2004 - https://web.eecs.umich.edu/~mihalcea/papers/mihalcea.emnlp04.pdf
	for _, node := range *graph {
		node.Score = scoreNode(node, iterations)
	}

	sort.Sort(sort.Reverse(graph))
	for _, node := range *graph {
		ranked = append(ranked, node.Data)
	}
	return ranked
}

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

// tokenize tokenises the text into sentences.
func tokenize(text string) []string {
	tokenizer, _ := english.NewSentenceTokenizer(nil)

	var sentences []string
	for _, token := range tokenizer.Tokenize(text) {
		sentence := strings.TrimSpace(token.Text)
		if sentence != "" {
			sentences = append(sentences, sentence)
		}
	}
	return sentences
}

// similarity calculates the similarity between two sentences, normalizing for
// sentence length.
func similarity(a, b string) float64 {
	tokenizer := english.NewWordTokenizer(sentences.NewPunctStrings())

	tokensA := tokenizer.Tokenize(a, false)
	tokensB := tokenizer.Tokenize(b, false)

	if len(tokensA) == 0 || len(tokensB) == 0 {
		return 0
	}

	similarWords := make(map[string]bool)
	for _, tokenA := range tokensA {
		wordA := strings.TrimSuffix(strings.ToLower(tokenA.Tok), ",")

		for _, tokenB := range tokensB {
			wordB := strings.TrimSuffix(strings.ToLower(tokenB.Tok), ",")

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
