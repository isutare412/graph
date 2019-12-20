package graph

import (
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
	graph.AddEdge(vertices[0], vertices[2], 1)
	graph.AddEdge(vertices[2], vertices[4], 2)
	graph.AddEdge(vertices[3], vertices[5], 4)
	graph.AddEdge(vertices[4], vertices[5], 5)
	graph.AddEdge(vertices[5], vertices[1], 6)
	t.Log(graph.String())
}

func TestShortestPaths(t *testing.T) {
	graph := New()
	vertices := []Vertex{
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
	}
	graph.AddEdge(vertices[0], vertices[2], 1)
	graph.AddEdge(vertices[2], vertices[4], 2)
	graph.AddEdge(vertices[2], vertices[5], 3)
	graph.AddEdge(vertices[3], vertices[5], 4)
	graph.AddEdge(vertices[4], vertices[5], 5)
	graph.AddEdge(vertices[5], vertices[1], 6)

	shortestPaths := graph.ShortestPaths(vertices[0])
	for v, weight := range shortestPaths {
		t.Logf("[%s]: weight %d", v.ID(), weight)
	}
}

func BenchmarkNewVertex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		graph := New()
		for j := 0; j < 1000; j++ {
			graph.NewVertex()
		}
	}
}
