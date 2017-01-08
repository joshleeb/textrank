# TextRank

The TextRank algorithm implemented in Go. This particular implementation is suited for sentence extraction and is based off the paper by [Rada Mihalcea and Paul Tarau (2004)](https://web.eecs.umich.edu/~mihalcea/papers/mihalcea.emnlp04.pdf).

## Example

```go
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

	// Iterating 30 times was chosen based on the convergence curves in Figure 1
	// of "TextRank: Bringing Order into Texts" by Rada Mihalcea and Paul Tarau,
	// 2004 - https://web.eecs.umich.edu/~mihalcea/papers/mihalcea.emnlp04.pdf
	words := textrank.RankWords(text, 30)[:5]
	fmt.Println(words)

	// Iterating 5 times was chosen based on the convergence curves in Figure 1
	// of "TextRank: Bringing Order into Texts" by Rada Mihalcea and Paul Tarau,
	// 2004 - https://web.eecs.umich.edu/~mihalcea/papers/mihalcea.emnlp04.pdf
	sentences := textrank.RankSentences(text, 5)
	for _, sentence := range sentences[:5] {
		fmt.Println("\n" + sentence)
	}
}
```
