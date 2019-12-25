package graph

import (
	"fmt"
	"strconv"
	"strings"
)

// VertexID is an identity for each vertex.
type VertexID int

type vertexible interface {
	fmt.Stringer

	id() VertexID
	accessor() Vertex
	addEdge(vertexible, int)
	removeEdge(VertexID)
	edges() []edge
}

type vertex struct {
	VertexID
	outgoing  []edge
	container Vertex
}

// Vertex of graph. It is safe to copy Vertex.
type Vertex struct {
	vertex vertexible
	// Value stores user defined values.
	Value *interface{}
}

func (id VertexID) String() string {
	return strconv.Itoa(int(id))
}

func (v *vertex) id() VertexID { return v.VertexID }

func (v *vertex) accessor() Vertex { return v.container }

func (v *vertex) addEdge(new vertexible, weight int) {
	v.outgoing = append(v.outgoing, edge{
		to:     new,
		weight: weight,
	})
}

func (v *vertex) removeEdge(dest VertexID) {
	for i := 0; i < len(v.outgoing); i++ {
		if v.outgoing[i].to.id() == dest {
			v.outgoing[i] = v.outgoing[len(v.outgoing)-1]
			v.outgoing = v.outgoing[:len(v.outgoing)-1]
			i--
		}
	}
}

func (v *vertex) edges() []edge { return v.outgoing }

func (v *vertex) String() string {
	format := func(id VertexID) string {
		return "[" + id.String() + "]"
	}
	var result strings.Builder
	result.WriteString(format(v.VertexID))
	if len(v.outgoing) <= 0 {
		return result.String()
	}

	result.WriteString(" -> ")
	var IDs = make([]string, 0, len(v.outgoing))
	for _, e := range v.outgoing {
		IDs = append(IDs, format(e.to.id()))
	}
	result.WriteString(strings.Join(IDs, ", "))
	return result.String()
}

// ID returns VertexID(int).
func (v Vertex) ID() VertexID {
	return v.vertex.id()
}
