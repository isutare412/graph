package graph

import (
	"fmt"
	"testing"
)

func TestPrintGraph(t *testing.T) {
	graph := New()
	vertices := []Vertex{
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
	}
	graph.AddEdge(&vertices[0], &vertices[2], 1)
	graph.AddEdge(&vertices[2], &vertices[4], 2)
	graph.AddEdge(&vertices[3], &vertices[5], 3)
	graph.AddEdge(&vertices[4], &vertices[5], 4)
	graph.AddEdge(&vertices[5], &vertices[1], 5)
	fmt.Println(graph.String())
}

func BenchmarkNewVertex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		graph := New()
		for j := 0; j < 1000; j++ {
			graph.NewVertex()
		}
	}
}
