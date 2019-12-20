package graph

import (
	"container/heap"
	"strconv"
	"strings"
	"sync"
)

// VertexID is an identity for each vertex.
type VertexID int

// Graph implements an adjacency list.
type Graph struct {
	vertices map[VertexID]*vertex
}

// Vertex of graph. It is safe to copy Vertex.
type Vertex struct {
	vertex *vertex
	// Value stores user defined values.
	Value *interface{}
}

type vertex struct {
	id        VertexID
	outgoing  []edge
	container Vertex
}

type edge struct {
	to     *vertex
	weight int
}

// distanceHeap implements min-heap interface for algorithm operations.
type distanceHeap []edge

func (id VertexID) String() string {
	return strconv.Itoa(int(id))
}

// NewVertex returns a new vertex which is ready to use.
func (g *Graph) NewVertex() Vertex {
	newVertex := &vertex{id: vertexIDGenrator()}
	newVertex.container = Vertex{
		vertex: newVertex,
		Value:  new(interface{}),
	}
	g.vertices[newVertex.id] = newVertex
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
	newEdge := edge{to: to.vertex, weight: weight}
	from.vertex.outgoing = append(from.vertex.outgoing, newEdge)
}

// ShortestPaths returns shortest paths from source to each other vertices.
// The value of map[Vertex]int is negative if the Vertex is unreachable
// from the source.
func (g *Graph) ShortestPaths(source Vertex) map[Vertex]int {
	var dists = make(map[Vertex]int)
	for _, v := range g.vertices {
		weight := -1
		if v == source.vertex {
			weight = 0
		}
		dists[v.container] = weight
	}

	distHeap := distanceHeap(make([]edge, 0, len(dists)))
	for v, w := range dists {
		distHeap = append(distHeap, edge{
			to:     v.vertex,
			weight: w,
		})
	}
	heap.Init(&distHeap)

	// Dijkstra's algorithm
	for i := 0; i < len(dists); i++ {
		closestEdge := heap.Pop(&distHeap).(edge)
		if closestEdge.weight < 0 {
			break
		}
		for _, e := range closestEdge.to.outgoing {
			weight := closestEdge.weight + e.weight
			if dists[e.to.container] < 0 || weight < dists[e.to.container] {
				dists[e.to.container] = weight
				distHeap.update(e.to, weight)
			}
		}
	}
	return dists
}

func (g *Graph) String() string {
	var result strings.Builder
	for _, v := range g.vertices {
		result.WriteString(v.String() + "\n")
	}
	return result.String()
}

// ID returns VertexID(int).
func (v Vertex) ID() VertexID {
	return v.vertex.id
}

func (v *vertex) removeEdge(dest VertexID) (removed bool) {
	for i := 0; i < len(v.outgoing); i++ {
		if v.outgoing[i].to.id == dest {
			v.outgoing[i] = v.outgoing[len(v.outgoing)-1]
			v.outgoing = v.outgoing[:len(v.outgoing)-1]
			removed = true
		}
	}
	return
}

func (v *vertex) String() string {
	format := func(id VertexID) string {
		return "[" + id.String() + "]"
	}
	var result strings.Builder
	result.WriteString(format(v.id))
	if len(v.outgoing) <= 0 {
		return result.String()
	}

	result.WriteString(" -> ")
	var IDs = make([]string, 0, len(v.outgoing))
	for _, e := range v.outgoing {
		IDs = append(IDs, format(e.to.id))
	}
	result.WriteString(strings.Join(IDs, ", "))
	return result.String()
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

func (d distanceHeap) update(v *vertex, weight int) {
	for i, e := range d {
		if v == e.to {
			d[i].weight = weight
			heap.Fix(&d, i)
			return
		}
	}
}

var vertexIDGenrator = func() func() VertexID {
	var vertexIDLast VertexID
	var vertexIDLock sync.Mutex
	return func() VertexID {
		vertexIDLock.Lock()
		defer vertexIDLock.Unlock()
		vertexIDLast++
		return vertexIDLast
	}
}()

// New returns initialized Graph.
func New() *Graph {
	return &Graph{vertices: make(map[VertexID]*vertex)}
}
