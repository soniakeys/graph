// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

// A single variadic function, DepthFirst, takes options in the form of configuration functions.
package search

import (
	"errors"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

// DepthFirst performs a depth-first search or traversal of graph g starting at
// node start.
//
// Options controlling the search are specified with configuration functions
// defined in this package.
//
// A non-nil error indicates some problem initializing the search, such as
// an invalid graph type or options.
func DepthFirst(g interface{}, start graph.NI, options ...func(*config)) error {
	cf := &config{start: start}
	for _, o := range options {
		o(cf)
	}
	switch t := g.(type) {
	case graph.AdjacencyList:
		cf.adjFunc(t)
	case graph.Directed:
		cf.adjFunc(t.AdjacencyList)
	case graph.Undirected:
		cf.adjFunc(t.AdjacencyList)
	case graph.LabeledAdjacencyList:
		cf.labFunc(t)
	case graph.LabeledDirected:
		cf.labFunc(t.LabeledAdjacencyList)
	case graph.LabeledUndirected:
		cf.labFunc(t.LabeledAdjacencyList)
	default:
		return errors.New("invalid graph type for DepthFirst")
	}
	return nil
}

func (cf *config) adjFunc(g graph.AdjacencyList) {
	b := cf.visited
	if b == nil {
		n := bits.New(len(g))
		b = &n
	}
	var df func(graph.NI) bool
	df = func(n graph.NI) bool {
		if b.Bit(int(n)) != 0 {
			return true
		}
		b.SetBit(int(n), 1)
		if cf.pathBits != nil {
			cf.pathBits.SetBit(int(n), 1)
		}

		if cf.nodeVisitor != nil {
			cf.nodeVisitor(n)
		}
		if cf.okNodeVisitor != nil {
			if !cf.okNodeVisitor(n) {
				return false
			}
		}

		if cf.rand == nil {
			for x, to := range g[n] {
				if cf.arcVisitor != nil {
					cf.arcVisitor(n, x)
				}
				if cf.okArcVisitor != nil {
					if !cf.okArcVisitor(n, x) {
						return false
					}
				}
				if !df(to) {
					return false
				}
			}
		} else {
			to := g[n]
			for _, x := range cf.rand.Perm(len(to)) {
				if cf.arcVisitor != nil {
					cf.arcVisitor(n, x)
				}
				if cf.okArcVisitor != nil {
					if !cf.okArcVisitor(n, x) {
						return false
					}
				}
				if !df(to[x]) {
					return false
				}
			}
		}
		if cf.pathBits != nil {
			cf.pathBits.SetBit(int(n), 0)
		}
		return true
	}
	df(cf.start)
}

func (cf *config) labFunc(g graph.LabeledAdjacencyList) {
}
