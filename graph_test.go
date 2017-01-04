package textrank

import "testing"

func TestNewGraph(t *testing.T) {
	graph := *newGraph([]string{"A", "B", "A", "C", "B"})

	if len(graph) != 3 {
		t.Errorf("graph length = %d, expected 3", len(graph))
	}
	if graph[0].Text != "A" {
		t.Errorf("graph[0] = %s, expected A", graph[0].Text)
	}
	if graph[1].Text != "B" {
		t.Errorf("graph[0] = %s, expected B", graph[1].Text)
	}
	if graph[2].Text != "C" {
		t.Errorf("graph[0] = %s, expected C", graph[2].Text)
	}
}

func TestAddNode(t *testing.T) {
	graph := textgraph{}
	graph.addNode("some-text", 1)

	if len(graph) != 1 {
		t.Errorf("graph length = %d, expected 1", len(graph))
	}
}

func TestGetExistingNode(t *testing.T) {
	graph := textgraph{}
	graph.addNode("A", 1)
	graph.addNode("B", 2)

	node := graph.getNode("B")
	if node.Text != "B" {
		t.Errorf("node text = %s, expected B", node.Text)
	}
}

func TestGetNonExistingNode(t *testing.T) {
	graph := textgraph{}

	node := graph.getNode("A")
	if node != nil {
		t.Errorf("node = %#v, expected nil", node)
	}
}

func TestNormalizeGraph(t *testing.T) {
	graph := textgraph{}
	graph.addNode("A", 1)
	graph.addNode("B", 2)

	norm := graph.normalize()
	if len(norm) != 2 {
		t.Errorf("normalized length = %d, expected 2", len(norm))
	}
	if !eqStringSlices(norm, []string{"B", "A"}) {
		t.Errorf("normalized = %#v, expected [B, A]", norm)
	}
}
