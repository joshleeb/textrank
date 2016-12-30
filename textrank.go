package textrank

import (
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/neurosnap/sentences/english"
)

// dampeningFactor is used when scoring nodes.
// It should be in the range (0.0, 1.0].
const dampeningFactor = 0.85

// lexicalCoOccurenceFactor is used when linking nodes containing words.
// It should be in the range [2, 10].
const lexicalCoOccurenceFactor = 2

// minWordSentence is the minimum number of words a sentence can have to become
// a node in the graph.
const minWordSentence = 5

// tokenizeWordsReplacePunctRe is a RegExp that replaces punctuation with spaces
// when tokenizing text into words.
var tokenizeWordsReplacePunctRe = regexp.MustCompile(`[.,\/!&;:=\-_]`)

// tokenizeWordsRemovePunctRe is a RegExp that removes punctuation when
// tokenizing text into words.
var tokenizeWordsRemovePunctRe = regexp.MustCompile(`[#$%\^\*{}~()\?\'\"]`)

// RankSentences ranks the sentences in the given text based on the TextRank
// algorithm and returned a list of the ranked sentences in descending order or
// score.
func RankSentences(text string, iterations int) []string {
	sentences := TokenizeSentences(text)
	graph := newSentenceGraph(sentences)
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

// RankWords ranks the rods in the given text based on the TextRank algorithm
// and returned a list of the ranked words in descending order or score.
func RankWords(text string, iterations int) []string {
	words := TokenizeWords(text)

	// Remove stopwords.
	cleanWords := []string{}
	for _, word := range words {
		if _, ok := stopwords[word]; ok {
			continue
		}
		cleanWords = append(cleanWords, word)
	}

	graph := newWordGraph(cleanWords)
	ranked := []string{}

	// Iterating 30 times was chosen based on the convergence curves in Figure 1
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

	rawSentences := []string{}
	for _, token := range tokenizer.Tokenize(text) {
		token := strings.TrimSpace(token.Text)
		if token != "" {
			rawSentences = append(rawSentences, strings.TrimSuffix(token, "."))
		}
	}

	// Often text will contain a sentence that finished with a period followed
	// by a sentence startin with a capital letter, and zero or more spaces in
	// between.  Since the order of these sentences doesn't matter for TextRank
	// we can go through the sentences again and split any that match this
	// format.
	re := regexp.MustCompile("\\.\\s*[A-Z]")
	sentences := []string{}
	for _, token := range rawSentences {
		sentences = append(sentences, splitByRegexp(token, re)...)
	}
	return sentences
}

// TokenizeWords tokenizes the text into words.
func TokenizeWords(text string) []string {
	text = strings.ToLower(
		tokenizeWordsReplacePunctRe.ReplaceAllString(text, " "))
	text = strings.ToLower(
		tokenizeWordsRemovePunctRe.ReplaceAllString(text, ""))

	words := []string{}
	for _, word := range strings.Split(text, " ") {
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
		successiveScore += scoreNode(linked, iterations-1)
	}
	return (1 - dampeningFactor) + dampeningFactor*successiveScore
}

// similarity calculates the similarity between two sentences, normalizing for
// sentence length.
func similarity(a, b string) float64 {
	tokensA := TokenizeWords(a)
	tokensB := TokenizeWords(b)

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

// splitByRegexp splits a specified string by the regex provided.
func splitByRegexp(text string, re *regexp.Regexp) []string {
	indexes := re.FindAllStringIndex(text, -1)
	prevStart := 0
	results := make([]string, len(indexes)+1)
	for i, element := range indexes {
		results[i] = text[prevStart:element[0]]
		prevStart = element[1] - 1
	}
	results[len(indexes)] = text[prevStart:len(text)]
	return results
}
