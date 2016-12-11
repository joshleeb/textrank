package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/joshleeb/textrank"
)

func main() {
	bytes, _ := ioutil.ReadAll(os.Stdin)
	text := string(bytes)

	// Iterating 5 times was chosen based on the convergence curves in Figure 1
	// of "TextRank: Bringing Order into Texts" by Rada Mihalcea and Paul Tarau,
	// 2004 - https://web.eecs.umich.edu/~mihalcea/papers/mihalcea.emnlp04.pdf
	sentences := textrank.Rank(text, 5)
	for _, sentence := range sentences {
		fmt.Println("\n" + sentence)
	}
}
