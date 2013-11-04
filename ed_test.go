// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed_test

import (
	"fmt"

	"github.com/soniakeys/ed"
)

// edgeData struct for simple specification of example data
type edgeData struct {
	n1, n2 string
	l      float64
}

// example data
var (
	exampleEdges = []edgeData{
		{"a", "b", 7},
		{"a", "c", 9},
		{"a", "f", 14},
		{"b", "c", 10},
		{"b", "d", 15},
		{"c", "d", 11},
		{"c", "f", 2},
		{"d", "e", 6},
		{"e", "f", 9},
	}
	exampleStart = "a"
	exampleEnd   = "e"
)

// linkGraph constructs a linked representation of example data.
func linkGraph(g []edgeData, start, end string) (allNodes map[string]*ed.Node, startNode, endNode *ed.Node) {
	all := map[string]*ed.Node{}
	// one pass over data to collect nodes
	for _, e := range g {
		if all[e.n1] == nil {
			all[e.n1] = &ed.Node{Name: e.n1}
		}
		if all[e.n2] == nil {
			all[e.n2] = &ed.Node{Name: e.n2}
		}
	}
	// second pass to link neighbors
	for _, ge := range g {
		n1 := all[ge.n1]
		n1.Neighbors = append(n1.Neighbors, ed.Neighbor{ge.l, all[ge.n2]})
	}
	return all, all[start], all[end]
}

func ExampleShortestPath() {
	// construct linked representation of example data
	allNodes, startNode, endNode :=
		linkGraph(exampleEdges, exampleStart, exampleEnd)
	// echo initial conditions
	fmt.Printf("Directed graph with %d nodes, %d edges\n",
		len(allNodes), len(exampleEdges))
	// run Dijkstra's shortest path algorithm
	p, l := ed.ShortestPath(startNode, endNode)
	if p == nil {
		fmt.Println("No path from start node to end node")
		return
	}
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// repeat search
	for _, n := range allNodes {
		n.Reset()
	}
	p, l = ed.ShortestPath(startNode, endNode)
	if p == nil {
		fmt.Println("No path from start node to end node")
		return
	}
	fmt.Println("Shortest path:", p)
	fmt.Println("Path length:", l)
	// Output:
	// Directed graph with 6 nodes, 9 edges
	// Shortest path: [{0 a} {9 c} {11 d} {6 e}]
	// Path length: 26
	// Shortest path: [{0 a} {9 c} {11 d} {6 e}]
	// Path length: 26
}
