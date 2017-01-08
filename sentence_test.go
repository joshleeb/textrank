package textrank

import (
	"fmt"
	"math"
	"testing"
)

func TestLinkSentences(t *testing.T) {
	cases := map[string]struct {
		Sentences [2]string
		Linked    bool
	}{
		"linked":         {[2]string{"a b c d e f", "b c d e f g"}, true},
		"not linked":     {[2]string{"a b c d e f", "g, h, i, j, k, l"}, false},
		"small sentence": {[2]string{"a b c", "a b c"}, false},
	}

	for k, tc := range cases {
		graph := textgraph{}

		// Setup graph with nodes.
		for _, sentence := range tc.Sentences {
			graph.addNode(sentence, 0)
		}

		linkSentences(&graph)
		if tc.Linked {
			if len(graph[0].Edges) != 1 {
				t.Errorf("%s: graph[0] edges = %d, expected %d",
					k, len(graph[0].Edges), 1)
			}
			if len(graph[1].Edges) != 1 {
				t.Errorf("%s: graph[1] edges = %d, expected %d",
					k, len(graph[0].Edges), 1)
			}
		} else {
			if len(graph[0].Edges) != 0 {
				for _, edge := range graph[0].Edges {
					fmt.Println(edge.Text)
				}
				t.Errorf("%s: graph[0] edges = %d, expected %d",
					k, len(graph[0].Edges), 0)
			}
			if len(graph[1].Edges) != 0 {
				t.Errorf("%s: graph[1] edges = %d, expected %d",
					k, len(graph[0].Edges), 0)
			}
		}
	}
}

func TestSentenceSimilarity(t *testing.T) {
	allowedDelta := 0.005

	cases := map[string]struct {
		Sentences          [2]string
		ExpectedSimilarity float64
	}{
		"both empty sentences": {[2]string{"", ""}, 0},
		"case insensitive":     {[2]string{"A b C d E f", "a b c d e f"}, 1.116},
		"empty sentence":       {[2]string{"", "a b c"}, 0},
		"full similarity":      {[2]string{"a b c d e f", "a b c d e f"}, 1.116},
		"no similarity":        {[2]string{"a b c d e f", "g h i j k l"}, 0},
		"small sentence":       {[2]string{"a b c", "a b c"}, 0},
		"some similarity":      {[2]string{"a b c d e f", "a b d f x z"}, 0.558},
	}

	for k, tc := range cases {
		similarity := sentenceSimilarity(tc.Sentences[0], tc.Sentences[1])
		if math.Abs(tc.ExpectedSimilarity-similarity) > allowedDelta {
			t.Errorf("%s: similarity = %f expected = %f",
				k, similarity, tc.ExpectedSimilarity)
		}
	}
}
