package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/neurosnap/sentences"
	"github.com/neurosnap/sentences/english"
)

func main() {
	bytes, _ := ioutil.ReadAll(os.Stdin)
	text := string(bytes)

	sentences := tokenize(text)
	graph := newGraph(sentences)

	// Iterating 5 times was chosen based on the convergence curves in Figure 1
	// of "TextRank: Bringing Order into Texts" by Rada Mihalcea and Paul Tarau,
	// 2004 - https://web.eecs.umich.edu/~mihalcea/papers/mihalcea.emnlp04.pdf
	for _, node := range *graph {
		node.Score = scoreNode(node, 5)
	}

	sort.Sort(sort.Reverse(graph))
	for _, node := range (*graph)[:4] {
		fmt.Println(node.Data + "\n")
	}
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