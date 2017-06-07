// Copyright 2016 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package df

import (
	"math/rand"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

type config struct {
	arcVisitor    func(n graph.NI, x int)
	iterateFrom   func(n graph.NI)
	nodeVisitor   graph.NodeVisitor
	okArcVisitor  func(n graph.NI, x int) bool
	okNodeVisitor graph.OkNodeVisitor
	pathBits      *bits.Bits
	rand          *rand.Rand
	visited       *bits.Bits
}

// ArcVisitor specifies a visitor function to call at each arc.
//
// See also OkArcVisitor.
func ArcVisitor(v func(n graph.NI, x int)) func(*config) {
	return func(c *config) {
		c.arcVisitor = v
	}
}

// NodeVisitor specifies a visitor function to call at each node.
//
// See also OkNodeVisitor.
func NodeVisitor(v graph.NodeVisitor) func(*config) {
	return func(c *config) {
		c.nodeVisitor = v
	}
}

// OkArcVisitor specifies a visitor function to perform some test at each arc
// and return a boolean result.
//
// As long as v return a result of true, the search progresses to traverse all
// arcs.
//
// If v returns false, the search terminates immediately.
//
// See also ArcVisitor.
func OkArcVisitor(v func(n graph.NI, x int) bool) func(*config) {
	return func(c *config) {
		c.okArcVisitor = v
	}
}

// OkNodeVisitor specifies a visitor function to perform some test at each node
// and return a boolean result.
//
// As long as v return a result of true, the search progresses to traverse all
// nodes.
//
// If v returns false, the search terminates immediately.
//
// See also NodeVisitor.
func OkNodeVisitor(v graph.OkNodeVisitor) func(*config) {
	return func(c *config) {
		c.okNodeVisitor = v
	}
}

// PathBits specifies a bits.Bits value for nodes of the path to the
// currently visited node.
//
// A use for PathBits is identifying back arcs in a search.
//
// Unlike Visited, PathBits are zeroed at the start of a search.
func PathBits(b *bits.Bits) func(*config) {
	return func(c *config) { c.pathBits = b }
}

// Rand specifies to traverse edges from each visited node in random order.
func Rand(r *rand.Rand) func(*config) {
	return func(c *config) { c.rand = r }
}

// Visited specifies a bits.Bits value to record visited nodes.
//
// For each node visited, the corresponding bit is set to 1.  Other bits
// are not modified.
//
// The search algorithm controls the search using a bits.Bits.  If this
// function is used, argument b will be used as the controlling value.
//
// Bits are not zeroed at the start of a search, so the initial Bits value
// passed in should generally be zero.  Non-zero bits will limit the search.
func Visited(b *bits.Bits) func(*config) {
	return func(c *config) { c.visited = b }
}
