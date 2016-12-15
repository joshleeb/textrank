package textrank

import (
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/neurosnap/sentences"
	"github.com/neurosnap/sentences/english"
)

const dampeningFactor = 0.85

// minWordSentence is the minimum number of words a sentence can have to become
// a node in the graph.
const minWordSentence = 5

// Rank ranks the sentences in the given text based on the TextRank algorithm
// and returned a list of the ranked sentences in descending order or score.
func Rank(text string, iterations int) []string {
	sentences := TokenizeSentences(text)
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

// TokenizeSentences tokenises the text into sentences.
func TokenizeSentences(text string) []string {
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

// TokenizeWords tokenizes the text into words.
func TokenizeWords(text string) []string {
	tokenizer := english.NewWordTokenizer(sentences.NewPunctStrings())

	var words []string
	for _, token := range tokenizer.Tokenize(text, false) {
		word := strings.TrimSpace(token.Tok)
		if word != "" {
			words = append(words, word)
		}
	}
	return words
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

// similarity calculates the similarity between two sentences, normalizing for
// sentence length.
func similarity(a, b string) float64 {
	punctRe := regexp.MustCompile(`[.,\/#!$%\^&\*;:{}=\-_~()]`)
	stopwords := getStopwords()

	tokensA := TokenizeWords(a)
	tokensB := TokenizeWords(b)

	if len(tokensA) < minWordSentence || len(tokensB) < minWordSentence {
		return 0
	}

	similarWords := make(map[string]bool)
	for _, tokenA := range tokensA {
		wordA := strings.ToLower(punctRe.ReplaceAllString(tokenA, ""))

		// Ignore stopwords. Only need to check wordA because if wordA is not a
		// stopword and wordB is a stopword, then they are not going to match.
		if _, ok := stopwords[wordA]; ok {
			continue
		}

		for _, tokenB := range tokensB {
			wordB := strings.ToLower(punctRe.ReplaceAllString(tokenB, ""))

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
