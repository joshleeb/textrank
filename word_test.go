package textrank

import (
	"strings"
	"testing"
)

func TestLinkWords(t *testing.T) {
	graph := textgraph{}
	sentence := "1 2 3 4 5"
	links := map[string][]string{
		"1": []string{"2"},
		"2": []string{"1", "3"},
		"3": []string{"2", "4"},
		"4": []string{"3", "5"},
		"5": []string{"4"},
	}

	// Setup graph with nodes.
	for _, word := range strings.Split(sentence, " ") {
		graph.addNode(word, 0)
	}

	linkWords(&graph, []string{sentence})
	for _, node := range graph {
		nodeEdges := []string{}
		for _, edge := range node.Edges {
			nodeEdges = append(nodeEdges, edge.Text)
		}

		for word, wordLinks := range links {
			if node.Text == word && !eqStringSlices(nodeEdges, wordLinks) {
				t.Errorf("%s: edges = %v, expected %v",
					word, nodeEdges, wordLinks)
			}
		}
	}
}

func TestCleanText(t *testing.T) {
	cases := map[string]struct {
		Text         string
		ExpectedText []string
	}{
		"stopwords":       {"is the a", []string{}},
		"stopwords cases": {"iS ThE A", []string{}},
		"mixed":           {"the fox jumped", []string{"fox", "jumped"}},
	}

	for k, tc := range cases {
		cleaned := cleanText(tc.Text)
		if !eqStringSlices(cleaned, tc.ExpectedText) {
			t.Errorf("%s: cleaned = %v, expected %v",
				k, cleaned, tc.ExpectedText)
		}
	}
}
