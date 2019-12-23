package graph

import (
	"sort"
	"testing"
)

func TestPrintGraph(t *testing.T) {
	graph := New(Directional)
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

func TestRemoveEdges(t *testing.T) {
	graph := New(Directional)
	vertices := []Vertex{
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
	}
	graph.AddEdge(vertices[0], vertices[2], 1)
	graph.AddEdge(vertices[0], vertices[2], 2)
	graph.AddEdge(vertices[2], vertices[4], 2)
	graph.AddEdge(vertices[3], vertices[5], 4)
	graph.AddEdge(vertices[4], vertices[5], 5)
	graph.AddEdge(vertices[5], vertices[1], 6)
	t.Log(graph.String())

	graph.RemoveEdges(vertices[0], vertices[2])
	t.Log(graph.String())
}

func TestShortestPaths(t *testing.T) {
	printResult := func(result map[Vertex]Path) {
		var keys []Vertex
		for v := range result {
			keys = append(keys, v)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].ID() < keys[j].ID()
		})
		for _, v := range keys {
			t.Logf("[%s]: weight(%d): %v",
				v.ID(), result[v].Distance(), result[v])
		}
	}

	graph := New(Directional)
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
	printResult(shortestPaths)

	t.Log("")
	graph = New(Directional)
	vertices = []Vertex{
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
	}
	graph.AddEdge(vertices[0], vertices[1], 4)
	graph.AddEdge(vertices[0], vertices[2], 5)
	graph.AddEdge(vertices[0], vertices[4], 3)
	graph.AddEdge(vertices[1], vertices[3], 1)
	graph.AddEdge(vertices[2], vertices[4], 4)
	graph.AddEdge(vertices[2], vertices[5], 5)
	graph.AddEdge(vertices[2], vertices[7], 3)
	graph.AddEdge(vertices[3], vertices[0], 2)
	graph.AddEdge(vertices[3], vertices[5], 5)
	graph.AddEdge(vertices[3], vertices[6], 2)
	graph.AddEdge(vertices[4], vertices[6], 4)
	graph.AddEdge(vertices[5], vertices[6], 1)
	graph.AddEdge(vertices[5], vertices[7], 2)
	graph.AddEdge(vertices[6], vertices[7], 2)
	graph.AddEdge(vertices[7], vertices[0], 5)

	shortestPaths = graph.ShortestPaths(vertices[0])
	printResult(shortestPaths)
}

func TestBidirectionalGraph(t *testing.T) {
	printResult := func(result map[Vertex]Path) {
		var keys []Vertex
		for v := range result {
			keys = append(keys, v)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].ID() < keys[j].ID()
		})
		for _, v := range keys {
			t.Logf("[%s]: weight(%d): %v",
				v.ID(), result[v].Distance(), result[v])
		}
	}

	t.Log("Directional Graph")
	graph := New(Directional)
	vertices := []Vertex{
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
	}
	graph.AddEdge(vertices[0], vertices[1], 1)
	graph.AddEdge(vertices[1], vertices[2], 2)

	shortestPaths := graph.ShortestPaths(vertices[1])
	printResult(shortestPaths)

	t.Log("Bidirectional Graph")
	graph = New(Bidirectional)
	vertices = []Vertex{
		graph.NewVertex(),
		graph.NewVertex(),
		graph.NewVertex(),
	}
	graph.AddEdge(vertices[0], vertices[1], 1)
	graph.AddEdge(vertices[1], vertices[2], 2)

	shortestPaths = graph.ShortestPaths(vertices[1])
	printResult(shortestPaths)
}

func BenchmarkNewVertex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		graph := New(Directional)
		for j := 0; j < 1000; j++ {
			graph.NewVertex()
		}
	}
}
