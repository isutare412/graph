package graph

import (
	"strconv"
	"strings"
	"sync"
)

// VertexID is an identity for each vertex
type VertexID int

// Graph implements an adjacency list.
type Graph struct {
	vertices map[VertexID]*vertex
}

// Vertex of graph.
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
func (g *Graph) AddEdge(from *Vertex, to *Vertex, weight int) {
	newEdge := edge{to: to.vertex, weight: weight}
	from.vertex.outgoing = append(from.vertex.outgoing, newEdge)
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
		return "[" + v.id.String() + "]"
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
