// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

// Package df provides a paramertized depth-first search.  A single variadic
// function, Search, takes options in the form of configuration functions.
package df

import (
	"errors"

	"github.com/soniakeys/graph"
)

// Search performs a depth-first search or traversal of graph g starting at
// node start.
//
// Options controlling the search are specified with configuration functions
// defined in this package.
//
// A non-nil error indicates some problem initializing the search, such as
// an invalid graph type or options.
func Search(g interface{}, start graph.NI, options ...func(*config)) (err error) {
	cf := &config{}
	for _, o := range options {
		o(cf)
	}
	if cf.bits == nil {
		cf.bits = &graph.Bits{}
	}
	var f func(start graph.NI) bool
	switch t := g.(type) {
	case graph.AdjacencyList:
		f, err = cf.adjSearchFunc(t)
	case graph.LabeledAdjacencyList:
		f, err = cf.labSearchFunc(t)
	default:
		return errors.New("invalid graph type")
	}
	if err == nil {
		f(start)
	}
	return
}

type dff struct {
	visited func(graph.NI) bool
	recurse func(graph.NI) bool
}

func (f *dff) search(n graph.NI) bool {
	return f.visited(n) || f.recurse(n)
}

func (cf *config) adjSearchFunc(g graph.AdjacencyList) (func(graph.NI) bool, error) {
	f := &dff{}
	search := f.search
	f.visited = cf.visitedFunc()
	f.recurse = cf.adjRecurseFunc(g, search)
	return search, nil
}

func (cf *config) visitedFunc() func(graph.NI) bool {
	b := cf.bits
	return func(n graph.NI) (t bool) {
		if b.Bit(n) != 0 {
			return true
		}
		b.SetBit(n, 1)
		return false
	}
}

/*
	if v := cf.okNodeVisitor; v != nil {
		return func(n graph.NI) bool {
			return t(n) || !v(n)
		}
	}
}
*/

func (cf *config) adjRecurseFunc(g graph.AdjacencyList, search func(graph.NI) bool) func(graph.NI) bool {
	if r := cf.rand; r != nil {
		if v := cf.okNodeVisitor; v != nil {
			if cf.visitOk != nil {
				*cf.visitOk = true
			}
			return func(n graph.NI) bool {
				if !v(n) {
					if cf.visitOk != nil {
						*cf.visitOk = false
					}
					return false
				}
				to := g[n]
				for _, i := range r.Perm(len(to)) {
					if !search(to[i]) {
						return false
					}
				}
				return true
			}
		} // else just rand, no visitor
		return func(n graph.NI) bool {
			to := g[n]
			for _, i := range r.Perm(len(to)) {
				if !search(to[i]) {
					return false
				}
			}
			return true

		}
	}
	// else no rand
	if v := cf.okNodeVisitor; v != nil {
		if cf.visitOk != nil {
			*cf.visitOk = true
		}
		return func(n graph.NI) bool {
			if !v(n) {
				if cf.visitOk != nil {
					*cf.visitOk = false
				}
				return false
			}
			for _, to := range g[n] {
				if !search(to) {
					return false
				}
			}
			return true
		}
	}
	// else no rand, no visitor
	return func(n graph.NI) bool {
		for _, to := range g[n] {
			if !search(to) {
				return false
			}
		}
		return true
	}
}

func (cf *config) labSearchFunc(g graph.LabeledAdjacencyList) (func(start graph.NI) bool, error) {
	return nil, nil
}
