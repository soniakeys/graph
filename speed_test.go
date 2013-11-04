// Copyright 2013 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package ed

import (
	"math"
	"math/rand"
	"sort"
	"testing"
)

type xy struct {
	x, y float64
	n    *Node
}
type xyList []xy

func (l xyList) Len() int           { return len(l) }
func (l xyList) Less(i, j int) bool { return l[i].n.Name < l[j].n.Name }
func (l xyList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

// generate a random graph
func r(nNodes, nEdges int) (all xyList, start, end *Node) {
	s := rand.New(rand.NewSource(59))
	// generate unique node names
	nameMap := map[string]bool{}
	name := make([]byte, int(math.Log(float64(nNodes))/3)+1)
	for len(nameMap) < nNodes {
		for i := range name {
			name[i] = 'A' + byte(s.Intn(26))
		}
		nameMap[string(name)] = true
	}
	// sort for repeatability
	nodes := make(xyList, nNodes)
	i := 0
	for n := range nameMap {
		nodes[i].n = &Node{Name: n}
		i++
	}
	sort.Sort(nodes)
	// now assign random coordinates.
	for i := range nodes {
		nodes[i].x = s.Float64()
		nodes[i].y = s.Float64()
	}
	// generate edges.
	for i := 0; i < nEdges; {
		n1 := &nodes[s.Intn(nNodes)]
		n2 := &nodes[s.Intn(nNodes)]
		dist := math.Hypot(n2.x-n1.x, n2.y-n1.y)
		if dist > s.Float64()*math.Sqrt2 {
			continue
		}
		n1.n.Neighbors = append(n1.n.Neighbors, Neighbor{dist, n2.n})
		switch i {
		case 0:
			start = n1.n
		case 1:
			end = n2.n
		}
		i++
	}
	return nodes, start, end
}

func Benchmark100(b *testing.B) {
	// 100 nodes
	all, start, end := r(100, 200)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ShortestPath(start, end)
		for _, n := range all {
			n.n.Reset()
		}
	}
}

func Benchmark1e3(b *testing.B) {
	// 1000 nodes
	all, start, end := r(1000, 3000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ShortestPath(start, end)
		for _, n := range all {
			n.n.Reset()
		}
	}
}

func Benchmark1e4(b *testing.B) {
	// 10k nodes
	all, start, end := r(1e4, 5e4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ShortestPath(start, end)
		for _, n := range all {
			n.n.Reset()
		}
	}
}

func Benchmark1e5(b *testing.B) {
	// 100k nodes
	all, start, end := r(1e5, 1e6)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ShortestPath(start, end)
		for _, n := range all {
			n.n.Reset()
		}
	}
}
