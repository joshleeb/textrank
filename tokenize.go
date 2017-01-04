package textrank

import (
	"regexp"
	"strings"

	"github.com/neurosnap/sentences/english"
)

// tokenizeWordsReplacePunctRe is a RegExp that replaces punctuation with spaces
// when tokenizing text into words.
var tokenizeWordsReplacePunctRe = regexp.MustCompile(`[.,\/!&;:=\-_]`)

// tokenizeWordsRemovePunctRe is a RegExp that removes punctuation when
// tokenizing text into words.
var tokenizeWordsRemovePunctRe = regexp.MustCompile(`[#$%\^\*{}~()\?\'\"]`)

// tokenizeSentences tokenises the text into sentences.
func tokenizeSentences(text string) []string {
	tokenizer, _ := english.NewSentenceTokenizer(nil)

	sentences := []string{}
	for _, token := range tokenizer.Tokenize(text) {
		token := strings.TrimSpace(token.Text)
		if token != "" {
			sentences = append(sentences, strings.TrimSuffix(token, "."))
		}
	}
	return sentences
}

// tokenizeWords tokenizes the text into words.
func tokenizeWords(text string) []string {
	text = strings.ToLower(tokenizeWordsReplacePunctRe.ReplaceAllString(text, " "))
	text = strings.ToLower(tokenizeWordsRemovePunctRe.ReplaceAllString(text, ""))

	words := []string{}
	for _, word := range strings.Split(text, " ") {
		if word != "" {
			words = append(words, word)
		}
	}
	return words
}
