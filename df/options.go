// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package df

import (
	"math/rand"

	"github.com/soniakeys/graph"
)

type config struct {
	iterateFrom func(n graph.NI)
	bits        *graph.Bits
	visitor     graph.Visitor
	visitOk     *bool
	rand        *rand.Rand
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

// Rand specifies to traverse edges from each visited node in random order.
func Rand(r *rand.Rand) func(*config) {
	return func(c *config) { c.rand = r }
}

// Visitor specifies a visitor function to call at each node.
//
// As long as v return true, the search progresses to traverse all nodes
// reachable from start, and ok is ultimately set to true.
//
// If the visitor function returns false, the search terminates immediately
// and ok is set to false.
//
// The ok pointer can be nil if the bool result is not needed.
func Visitor(v graph.Visitor, ok *bool) func(*config) {
	return func(c *config) {
		c.visitor = v
		c.visitOk = ok
	}
}
