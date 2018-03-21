// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package alt

import (
	"math/rand"

	"github.com/soniakeys/bits"
	"github.com/soniakeys/graph"
)

type config struct {
	start          graph.NI
	arcVisitor     func(n graph.NI, x int)
	iterateFrom    func(n graph.NI)
	levelVisitor   func(l int, n []graph.NI)
	nodeVisitor    func(n graph.NI)
	okArcVisitor   func(n graph.NI, x int) bool
	okNodeVisitor  func(n graph.NI) bool
	okLevelVisitor func(l int, n []graph.NI) bool
	rand           *rand.Rand
	visBits        *bits.Bits
	pathBits       *bits.Bits
	fromList       *graph.FromList
}

// A TraverseOption specifies an option for a breadth first or depth first
// traversal.
//
// Values of this type are returned by various TraverseOption constructor
// functions.  These constructors take optional values to be used
// in a traversal and wrap them in the TraverseOption type.  This type
// is actually a function.  The BreadthFirst and DepthFirst traversal
// methods call these functions in order, to initialize state that controls
// the traversal.
type TraverseOption func(*config)

// ArcVisitor specifies a visitor function to call at each arc.
//
// See also OkArcVisitor.
func ArcVisitor(v func(n graph.NI, x int)) TraverseOption {
	return func(c *config) {
		c.arcVisitor = v
	}
}

// From specifies a graph.FromList to populate.
func From(f *graph.FromList) TraverseOption {
	return func(c *config) {
		c.fromList = f
	}
}

// LevelVisitor specifies a visitor function to call at each level or depth.
//
// The level visitor function is called before any node or arc visitor
// functions.
//
// See also OkLevelVisitor.
func LevelVisitor(v func(level int, nodes []graph.NI)) TraverseOption {
	return func(c *config) {
		c.levelVisitor = v
	}
}

// NodeVisitor specifies a visitor function to call at each node.
//
// The node visitor function is called before any arc visitor functions.
//
// See also OkNodeVisitor.
func NodeVisitor(v func(graph.NI)) TraverseOption {
	return func(c *config) {
		c.nodeVisitor = v
	}
}

// OkArcVisitor specifies a visitor function to perform some test at each arc
// and return a boolean result.
//
// As long as v return a result of true, the traverse progresses to traverse all
// arcs.
//
// If v returns false, the traverse terminates immediately.
//
// See also ArcVisitor.
func OkArcVisitor(v func(n graph.NI, x int) bool) TraverseOption {
	return func(c *config) {
		c.okArcVisitor = v
	}
}

// OKLevelVisitor specifies a visitor function to call at each level or depth,
// returning a boolean result
//
// As long as v returns a result of true, the traverse progresses to traverse
// all nodes.  If v returns false, the traverse terminates immediately.
//
// The level visitor function is called before any node or arc visitor
// functions.
//
// See also LevelVisitor.
func OkLevelVisitor(v func(level int, nodes []graph.NI) bool) TraverseOption {
	return func(c *config) {
		c.okLevelVisitor = v
	}
}

// OkNodeVisitor specifies a visitor function to perform some test at each node
// and return a boolean result.
//
// As long as v returns a result of true, the traverse progresses to traverse
// all nodes.  If v returns false, the traverse terminates immediately.
//
// The node visitor function is called before any arc visitor functions.
//
// See also NodeVisitor.
func OkNodeVisitor(v func(graph.NI) bool) TraverseOption {
	return func(c *config) {
		c.okNodeVisitor = v
	}
}

// PathBits specifies a bits.Bits value for nodes of the path to the
// currently visited node.
//
// A use for PathBits is identifying back arcs in a traverse.
//
// Unlike Visited, PathBits are zeroed at the start of a traverse.
func PathBits(b *bits.Bits) TraverseOption {
	return func(c *config) { c.pathBits = b }
}

// Rand specifies to traverse edges from each visited node in random order.
func Rand(r *rand.Rand) TraverseOption {
	return func(c *config) { c.rand = r }
}

// Visited specifies a bits.Bits value to record visited nodes.
//
// For each node visited, the corresponding bit is set to 1.  Other bits
// are not modified.
//
// The traverse algorithm controls the traverse using a bits.Bits.  If this
// function is used, argument b will be used as the controlling value.
//
// Bits are not zeroed at the start of a traverse, so the initial Bits value
// passed in should generally be zero.  Non-zero bits will limit the traverse.
func Visited(b *bits.Bits) TraverseOption {
	return func(c *config) { c.visBits = b }
}

// BreadthFirst traverses a directed or undirected graph in breadth first order.
//
// Argument start is the start node for the traversal.  Argument opt can be
// any number of values returned by a supported TraverseOption function.
//
// Supported:
//
//   ArcVisitor
//   From
//   LevelVisitor
//   NodeVisitor
//   OkArcVisitor
//   OkLevelVisitor
//   OkNodeVisitor
//   Rand
//   Visited
//
// Unsupported:
//
//   PathBits
//
// See also alt.BreadthFirst2, a direction optimizing breadth first algorithm.
func BreadthFirst(g graph.AdjacencyList, start graph.NI, opt ...TraverseOption) {
	cf := &config{start: start}
	for _, o := range opt {
		o(cf)
	}
	// either visBits or fromList are suitable for recording visited nodes.
	// if neither is specified as an option, allocate bits.
	b := cf.visBits
	var rp []graph.PathEnd
	if cf.fromList == nil {
		if b == nil {
			n := bits.New(len(g))
			b = &n
		}
	} else {
		if cf.fromList.Paths == nil {
			*cf.fromList = graph.NewFromList(len(g))
		}
		rp = cf.fromList.Paths
	}
	fillBits := b != nil
	fillPath := rp != nil
	// not-visited test
	nvis := func(n graph.NI) bool { return rp[n].Len == 0 }
	// the frontier consists of nodes all at the same level
	frontier := []graph.NI{cf.start}
	level := 1
	var next []graph.NI
	// fill bits/path when node is put on frontier
	if fillBits {
		if b.Bit(int(cf.start)) == 1 {
			return
		}
		b.SetBit(int(cf.start), 1)
		nvis = func(n graph.NI) bool { return b.Bit(int(n)) == 0 }
	}
	if fillPath {
		rp[cf.start] = graph.PathEnd{Len: level, From: -1}
	}
	visitNode := func(n graph.NI) bool {
		// visit nodes as they come off frontier
		if cf.nodeVisitor != nil {
			cf.nodeVisitor(n)
		}
		if cf.okNodeVisitor != nil {
			if !cf.okNodeVisitor(n) {
				return false
			}
		}
		for x, nb := range g[n] {
			if cf.arcVisitor != nil {
				cf.arcVisitor(n, x)
			}
			if cf.okArcVisitor != nil {
				if !cf.okArcVisitor(n, x) {
					return false
				}
			}
			if nvis(nb) {
				next = append(next, nb)
				if fillBits {
					b.SetBit(int(nb), 1)
				}
				if fillPath {
					rp[nb] = graph.PathEnd{From: n, Len: level}
				}
			}
		}
		return true
	}
	for {
		if cf.rand != nil {
			cf.rand.Shuffle(len(frontier), func(i, j int) {
				frontier[i], frontier[j] = frontier[j], frontier[i]
			})
		}
		if cf.levelVisitor != nil {
			cf.levelVisitor(level, frontier)
		}
		if cf.okLevelVisitor != nil && !cf.okLevelVisitor(level, frontier) {
			return
		}
		if fillPath {
			cf.fromList.MaxLen = level
		}
		level++
		for _, n := range frontier {
			if !visitNode(n) {
				return
			}
		}
		if len(next) == 0 {
			break
		}
		frontier, next = next, nil
	}
}

// DepthFirst traverses a directed or undirected graph in depth first order.
//
// Argument start is the start node for the traversal.  Argument opt can be
// any number of values returned by a supported TraverseOption function.
//
// Supported:
//
//   NodeVisitor
//   OkNodeVisitor
//   ArcVisitor
//   OkArcVisitor
//   Visited
//   PathBits
//   Rand
//
// Unsupported:
//
//   From
//   LevelVisitor
//   OkLevelVisitor
func DepthFirst(g graph.AdjacencyList, start graph.NI, options ...TraverseOption) {
	cf := &config{start: start}
	for _, o := range options {
		o(cf)
	}
	b := cf.visBits
	if b == nil {
		n := bits.New(len(g))
		b = &n
	} else if b.Bit(int(cf.start)) != 0 {
		return
	}
	if cf.pathBits != nil {
		cf.pathBits.ClearAll()
	}
	var df func(graph.NI) bool
	df = func(n graph.NI) bool {
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
				if b.Bit(int(to)) != 0 {
					continue
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
				if b.Bit(int(to[x])) != 0 {
					continue
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
