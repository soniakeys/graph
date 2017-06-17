// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package traverse

import (
	"errors"

	"github.com/soniakeys/graph"
)

// BreadthFirst traverses a directed or undirected graph in breadth first order.
//
// Argument start is the start node for the traversal.  If r is nil, nodes are
// visited in deterministic order.  If a random number generator is supplied,
// nodes at each level are visited in random order.
//
// Argument f can be nil if you have no interest in the FromList path result.
// If FromList f is non-nil, the method populates f.Paths and sets f.MaxLen.
// It does not set f.Leaves.  For convenience argument f can be a zero value
// FromList.  If f.Paths is nil, the FromList is initialized first.  If f.Paths
// is non-nil however, the FromList is  used as is.  The method uses a value of
// PathEnd.Len == 0 to indentify unvisited nodes.  Existing non-zero values
// will limit the traversal.
//
// Traversal calls the visitor function v for each node starting with node
// start.  If v returns true, traversal continues.  If v returns false, the
// traversal terminates immediately.  PathEnd Len and From values are updated
// before calling the visitor function.
//
// On return f.Paths and f.MaxLen are set but not f.Leaves.
//
// Returned is the number of nodes visited and ok = true if the traversal
// ran to completion or ok = false if it was terminated by the visitor
// function returning false.
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also alt.BreadthFirst, a direction optimizing algorithm.
func BreadthFirst(g interface{}, start graph.NI, options ...Option) error {
	cf := &config{start: start}
	for _, o := range options {
		o(cf)
	}
	switch t := g.(type) {
	case graph.AdjacencyList:
		cf.bfAdj(t)
	case graph.Directed:
		cf.bfAdj(t.AdjacencyList)
	case graph.Undirected:
		cf.bfAdj(t.AdjacencyList)
	case graph.LabeledAdjacencyList:
		cf.bfLab(t)
	case graph.LabeledDirected:
		cf.bfLab(t.LabeledAdjacencyList)
	case graph.LabeledUndirected:
		cf.bfLab(t.LabeledAdjacencyList)
	default:
		return errors.New("invalid graph type for BreadthFirst")
	}
	return nil
}

func (cf *config) bfAdj(g graph.AdjacencyList) {
	f := cf.fromList
	switch {
	case f == nil:
		e := graph.NewFromList(len(g))
		f = &e
	case f.Paths == nil:
		*f = graph.NewFromList(len(g))
	}
	rp := f.Paths
	// the frontier consists of nodes all at the same level
	frontier := []graph.NI{cf.start}
	level := 1
	// assign path when node is put on frontier,
	rp[cf.start] = graph.PathEnd{Len: level, From: -1}
	for {
		f.MaxLen = level
		level++
		var next []graph.NI
		if cf.rand == nil {
			for _, n := range frontier {
				// visit nodes as they come off frontier
				if cf.nodeVisitor != nil {
					cf.nodeVisitor(n)
				}
				if cf.okNodeVisitor != nil {
					if !cf.okNodeVisitor(n) {
						return
					}
				}
				for _, nb := range g[n] {
					if rp[nb].Len == 0 {
						next = append(next, nb)
						rp[nb] = graph.PathEnd{From: n, Len: level}
					}
				}
			}
		} else { // take nodes off frontier at random
			for _, i := range cf.rand.Perm(len(frontier)) {
				n := frontier[i]
				// remainder of block same as above
				if cf.nodeVisitor != nil {
					cf.nodeVisitor(n)
				}
				if cf.okNodeVisitor != nil {
					if !cf.okNodeVisitor(n) {
						return
					}
				}
				for _, nb := range g[n] {
					if rp[nb].Len == 0 {
						next = append(next, nb)
						rp[nb] = graph.PathEnd{From: n, Len: level}
					}
				}
			}
		}
		if len(next) == 0 {
			break
		}
		frontier = next
	}
}

func (cf *config) bfLab(g graph.LabeledAdjacencyList) {
	f := cf.fromList
	switch {
	case f == nil:
		e := graph.NewFromList(len(g))
		f = &e
	case f.Paths == nil:
		*f = graph.NewFromList(len(g))
	}
	rp := f.Paths
	// the frontier consists of nodes all at the same level
	frontier := []graph.NI{cf.start}
	level := 1
	// assign path when node is put on frontier,
	rp[cf.start] = graph.PathEnd{Len: level, From: -1}
	for {
		f.MaxLen = level
		level++
		var next []graph.NI
		if cf.rand == nil {
			for _, n := range frontier {
				// visit nodes as they come off frontier
				if cf.nodeVisitor != nil {
					cf.nodeVisitor(n)
				}
				if cf.okNodeVisitor != nil {
					if !cf.okNodeVisitor(n) {
						return
					}
				}
				for _, nb := range g[n] {
					if rp[nb.To].Len == 0 {
						next = append(next, nb.To)
						rp[nb.To] = graph.PathEnd{From: n, Len: level}
					}
				}
			}
		} else { // take nodes off frontier at random
			for _, i := range cf.rand.Perm(len(frontier)) {
				n := frontier[i]
				// remainder of block same as above
				if cf.nodeVisitor != nil {
					cf.nodeVisitor(n)
				}
				if cf.okNodeVisitor != nil {
					if !cf.okNodeVisitor(n) {
						return
					}
				}
				for _, nb := range g[n] {
					if rp[nb.To].Len == 0 {
						next = append(next, nb.To)
						rp[nb.To] = graph.PathEnd{From: n, Len: level}
					}
				}
			}
		}
		if len(next) == 0 {
			break
		}
		frontier = next
	}
}
