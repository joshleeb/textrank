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
	text, _ := ioutil.ReadAll(os.Stdin)
	sentences := textrank.Rank(string(text), 5)
	for _, sentence := range sentences {
		fmt.Println(sentence)
	}
}
```
