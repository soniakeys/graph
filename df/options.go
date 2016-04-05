// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package df

import (
	"math/rand"

	"github.com/soniakeys/graph"
)

type config struct {
	iterateFrom   func(n graph.NI)
	bits          *graph.Bits
	nodeVisitor   graph.NodeVisitor
	okNodeVisitor graph.OkNodeVisitor
	visitOk       *bool
	rand          *rand.Rand
}

// Bits specifies a graph.Bits value to record visited nodes.
//
// The search algorithm controls the search using a graph.Bits.  If this
// function is used, argument b will be used as the controlling value.
// Bits are set as a side effect of the algorithm.  Bits should generally
// be zero initially.  Non-zero bits will limit the search.
func Bits(b *graph.Bits) func(*config) {
	return func(c *config) { c.bits = b }
}

// NodeVisitor specifies a visitor function to call at each node.
//
// See also OkNodeVisitor.
func NodeVisitor(v graph.NodeVisitor) func(*config) {
	return func(c *config) {
		c.nodeVisitor = v
	}
}

// OkNodeVisitor specifies a visitor function to perform some test at each node
// and return a boolean result.
//
// As long as v return a result of true, the search progresses to traverse all
// nodes, and ok is ultimately set to true.
//
// If v returns false, the search terminates immediately and ok is set to false.
//
// The ok pointer can be nil if the bool result is not needed.
//
// See also NodeVisitor.
func OkNodeVisitor(v graph.OkNodeVisitor, ok *bool) func(*config) {
	return func(c *config) {
		c.okNodeVisitor = v
		c.visitOk = ok
	}
}

// Rand specifies to traverse edges from each visited node in random order.
func Rand(r *rand.Rand) func(*config) {
	return func(c *config) { c.rand = r }
}
