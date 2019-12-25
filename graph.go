// Package graph implements several graph algorithms like Dijkstra, A*, etc.
package graph

import (
	"container/heap"
	"fmt"
	"strings"
	"sync"
)

// Type enumerates type of Graph.
type Type int8

const (
	// Directional graph. A -> B != B -> A.
	Directional Type = iota
	// Bidirectional graph. A <-> B.
	Bidirectional
)

// Graph implements an adjacency list. You should create a Graph by calling
// New(Type) function. As Graph doest have any location or coordinate, Graph
// cannot use A* algorithm. If you want A* algorithm, use CGraph instead.
type Graph struct {
	Type
	vertices   map[VertexID]vertexible
	generateID func() VertexID
}

type edge struct {
	to     vertexible
	weight int
}

// Path implements specific path to a vertex.
type Path struct {
	edges []edge
}

// distanceHeap implements min-heap interface for algorithm operations.
type distanceHeap []edge

// NewVertex returns a new vertex which is ready to use.
func (g *Graph) NewVertex() Vertex {
	newVertex := &vertex{VertexID: g.generateID()}
	newVertex.container = Vertex{
		vertex: newVertex,
		Value:  new(interface{}),
	}
	g.vertices[newVertex.VertexID] = newVertex
	return newVertex.container
}

// RemoveVertex removes the vertex with 'id'. Then edges that point to
// the removed vertex are also removed. Returns true if the vertex is removed.
func (g *Graph) RemoveVertex(id VertexID) bool {
	if _, ok := g.vertices[id]; !ok {
		return false
	}
	delete(g.vertices, id)
	for _, v := range g.vertices {
		v.removeEdge(id)
	}
	return true
}

// AddEdge adds a new edge with weight from Vertex to Vertex.
func (g *Graph) AddEdge(from, to Vertex, weight int) {
	from.vertex.addEdge(to.vertex, weight)
	if g.Type == Bidirectional {
		to.vertex.addEdge(from.vertex, weight)
	}
}

// RemoveEdges removes all edges from 'from' to 'to'. If Type of g is
// Directional, the other edges (from 'to' to 'from') are not removed.
func (g *Graph) RemoveEdges(from, to Vertex) {
	from.vertex.removeEdge(to.ID())
	if g.Type == Bidirectional {
		to.vertex.removeEdge(from.ID())
	}
}

func (g *Graph) dijkstra(src vertexible, handler func(vertexible, Path) bool) {
	var shortestPaths = make(map[vertexible]Path)
	for _, v := range g.vertices {
		if v != src {
			shortestPaths[v] = Path{}
		}
	}

	var distHeap = distanceHeap(make([]edge, 0, len(g.vertices)))
	for _, v := range g.vertices {
		weight := -1
		if v == src {
			weight = 0
		}
		distHeap = append(distHeap, edge{
			to:     v,
			weight: weight,
		})
	}
	heap.Init(&distHeap)

	entireSize := len(distHeap)
	for i := 0; i < entireSize; i++ {
		closestEdge := heap.Pop(&distHeap).(edge)
		if closestEdge.weight < 0 {
			break
		}
		// stop finding paths if handler returns false.
		if closestEdge.to != src &&
			!handler(closestEdge.to, shortestPaths[closestEdge.to]) {
			break
		}
		for _, e := range closestEdge.to.edges() {
			if e.to == src {
				continue
			}
			newW := closestEdge.weight + e.weight
			oldW := shortestPaths[e.to].Distance()
			if oldW < 0 || newW < oldW {
				fixedPath := shortestPaths[closestEdge.to]
				fixedPath.addEdge(e)
				shortestPaths[e.to] = fixedPath
				distHeap.update(e.to, newW)
			}
		}
	}
}

// ShortestPath returns shortest path p from src to dest. You can check whether
// the path exists by checking p.Destination() or p.Distance(). As g
// cannot be applied A* algorithm, ShortestPath uses Dijkstra's one instead.
func (g *Graph) ShortestPath(src, dest Vertex) (p Path) {
	g.dijkstra(src.vertex, func(v vertexible, shortest Path) bool {
		if v == dest.vertex {
			p = shortest
			return false
		}
		return true
	})
	return
}

// ShortestPaths returns shortest paths from source to every vertices which
// are reachable from source.
func (g *Graph) ShortestPaths(source Vertex) map[Vertex]Path {
	var dists = make(map[Vertex]Path)
	g.dijkstra(source.vertex, func(v vertexible, p Path) bool {
		dists[v.accessor()] = p
		return true
	})
	return dists
}

func (g *Graph) String() string {
	var result strings.Builder
	for _, v := range g.vertices {
		result.WriteString(v.String() + "\n")
	}
	return result.String()
}

func (e edge) String() string {
	return fmt.Sprintf("->%d [%v]", e.weight, e.to.id())
}

// Destination returns the destination of current path. ok is false if
// the Path does not include any Vertex yet.
func (p Path) Destination() (dest Vertex, ok bool) {
	if len(p.edges) == 0 {
		return dest, false
	}
	return p.edges[len(p.edges)-1].to.accessor(), true
}

// Distance returns a total distance to the destination. Returns negative
// number if the Path does not include any Vertex.
func (p Path) Distance() int {
	if len(p.edges) == 0 {
		return -1
	}
	var distance int
	for _, e := range p.edges {
		distance += e.weight
	}
	return distance
}

// IterateEdge iterates edges of p sequentially. If handler returns false,
// iteration will stop.
func (p Path) IterateEdge(handler func(to Vertex, weight int) bool) {
	for _, e := range p.edges {
		if !handler(e.to.accessor(), e.weight) {
			break
		}
	}
}

func (p Path) String() string {
	var builder strings.Builder
	for i, e := range p.edges {
		builder.WriteString(e.String())
		if i != len(p.edges)-1 {
			builder.WriteByte(' ')
		}
	}
	return builder.String()
}

func (p *Path) addEdge(target edge) error {
	if target.weight < 0 {
		return fmt.Errorf("cannot add edge with a negative weight")
	}
	p.edges = append(p.edges, target)
	return nil
}

func (d distanceHeap) Len() int { return len(d) }

func (d distanceHeap) Less(i, j int) bool {
	wi, wj := d[i].weight, d[j].weight
	// treat negative numbers as if it is greater than any positive numbers.
	if wi < 0 {
		return false
	} else if wj < 0 {
		return true
	}
	return wi < wj
}

func (d distanceHeap) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

func (d *distanceHeap) Push(x interface{}) {
	*d = append(*d, x.(edge))
}

func (d *distanceHeap) Pop() interface{} {
	old := *d
	size := len(old)
	popped := old[size-1]
	*d = old[:size-1]
	return popped
}

func (d distanceHeap) update(v vertexible, weight int) {
	for i, e := range d {
		if v == e.to {
			d[i].weight = weight
			heap.Fix(&d, i)
			return
		}
	}
}

// New returns initialized Graph.
func New(t Type) *Graph {
	return &Graph{
		Type:     t,
		vertices: make(map[VertexID]vertexible),
		generateID: func() func() VertexID {
			var vertexIDLast VertexID
			var vertexIDLock sync.Mutex
			return func() VertexID {
				vertexIDLock.Lock()
				defer vertexIDLock.Unlock()
				vertexIDLast++
				return vertexIDLast
			}
		}(),
	}
}
