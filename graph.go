package graph

import (
	"container/heap"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// VertexID is an identity for each vertex.
type VertexID int

// Graph implements an adjacency list.
type Graph struct {
	vertices   map[VertexID]*vertex
	generateID func() VertexID
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

// Path implements specific path to a vertex.
type Path struct {
	edges []edge
}

// distanceHeap implements min-heap interface for algorithm operations.
type distanceHeap []edge

func (id VertexID) String() string {
	return strconv.Itoa(int(id))
}

// NewVertex returns a new vertex which is ready to use.
func (g *Graph) NewVertex() Vertex {
	newVertex := &vertex{id: g.generateID()}
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
func (g *Graph) ShortestPaths(source Vertex) map[Vertex]Path {
	var dists = make(map[Vertex]Path)
	for _, v := range g.vertices {
		if v != source.vertex {
			dists[v.container] = Path{}
		}
	}

	distHeap := distanceHeap(make([]edge, 0, len(dists)+1))
	distHeap = append(distHeap, edge{
		to:     source.vertex,
		weight: 0,
	})
	for v := range dists {
		distHeap = append(distHeap, edge{
			to:     v.vertex,
			weight: -1,
		})
	}
	heap.Init(&distHeap)

	// Dijkstra's algorithm
	entireSize := len(distHeap)
	for i := 0; i < entireSize; i++ {
		closestEdge := heap.Pop(&distHeap).(edge)
		if closestEdge.weight < 0 {
			break
		}
		for _, e := range closestEdge.to.outgoing {
			if e.to == source.vertex {
				continue
			}
			newW := closestEdge.weight + e.weight
			oldW := dists[e.to.container].Distance()
			if oldW < 0 || newW < oldW {
				fixedPath := dists[closestEdge.to.container]
				fixedPath.addEdge(e)
				dists[e.to.container] = fixedPath
				distHeap.update(e.to, newW)
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

func (e edge) String() string {
	return fmt.Sprintf("->%d [%v]", e.weight, e.to.id)
}

// Destination returns the destination of current path. ok is false if
// the Path does not include any Vertex.
func (p Path) Destination() (dest Vertex, ok bool) {
	if len(p.edges) == 0 {
		return dest, false
	}
	return p.edges[len(p.edges)-1].to.container, true
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

func (d distanceHeap) update(v *vertex, weight int) {
	for i, e := range d {
		if v == e.to {
			d[i].weight = weight
			heap.Fix(&d, i)
			return
		}
	}
}

// New returns initialized Graph.
func New() *Graph {
	return &Graph{
		vertices: make(map[VertexID]*vertex),
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
