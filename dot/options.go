package dot

import (
	"strconv"

	"github.com/soniakeys/graph"
)

// AttrVal represents the dot format concept of an attribute-value pair.
type AttrVal struct {
	Attr string
	Val  string
}

// Config holds options that control the dot output.
//
// See Overview/Scheme for an overview of how this works.  Generally you will
// not set members of a Config struct directly.  There is an option function
// for each member.  To set a member, pass the option function as an optional
// argument to a Write or String function.
type Config struct {
	Directed     bool
	EdgeLabel    func(graph.LI) string
	GraphAttr    []AttrVal
	Indent       string
	Isolated     bool
	NodeLabel    func(graph.NI) string
	UndirectArcs bool
}

// Defaults holds a package default Config struct.
//
// Defaults is copied as the first configuration step.  See Overview/Scheme.
var Defaults = Config{
	Directed:  true,
	EdgeLabel: func(l graph.LI) string { return strconv.Itoa(int(l)) },
	Indent:    "  ",
	NodeLabel: func(n graph.NI) string { return strconv.Itoa(int(n)) },
}

// Directed specifies whether to write a dot format directected or undirected
// graph.
//
// The default, Directed(true), specifies a dot format directected graph,
// or "digraph."  In this case the Write or String function outputs each arc
// of the graph as a dot format directed edge.
//
// Directed(false) specifies a dot format undirected graph or simply a "graph."
// In this case the Write or String function requires that all arcs between
// distinct nodes occur in reciprocal pairs.  For each pair the function
// outputs a single edge in dot format.
func Directed(d bool) func(*Config) {
	return func(c *Config) { c.Directed = d }
}

// EdgeLabel specifies a function to generate edge label strings for the
// dot format given the arc label integers of graph package.
//
// The default function is simply strconv.Itoa of the graph package arc label.
func EdgeLabel(f func(graph.LI) string) func(*Config) {
	return func(c *Config) { c.EdgeLabel = f }
}

// GraphAttr adds a dot format graph attribute.
//
// Graph attributes are held in a slice, and so are ordered.  This function
// updates the value of the last matching attribute if it exists, or adds a
// new attribute to the end of the list.
func GraphAttr(attr, val string) func(*Config) {
	return func(c *Config) {
		for i := len(c.GraphAttr) - 1; i >= 0; i-- {
			if c.GraphAttr[i].Attr == attr {
				c.GraphAttr[i].Val = val
				return
			}
		}
		c.GraphAttr = append(c.GraphAttr, AttrVal{attr, val})
	}
}

// Indent specifies an indent string for the body of the dot format.
//
// The default is two spaces.
func Indent(i string) func(*Config) {
	return func(c *Config) { c.Indent = i }
}

// Isolated specifies whether to include isolated nodes.
//
// An isolated node has no arcs in or out.  By default, isolated = false,
// isolated nodes are not included in the dot output.
//
// Isolated(true) will include isolated nodes.
func Isolated(i bool) func(*Config) {
	return func(c *Config) { c.Isolated = i }
}

// NodeLabel specifies a function to generate node label strings for the
// dot format given the node integers of graph package.
//
// The default function is simply strconv.Itoa of the graph package node
// integer.
func NodeLabel(f func(graph.NI) string) func(*Config) {
	return func(c *Config) { c.NodeLabel = f }
}

// UndirectArcs, for the WeightedEdgeList graph type, specifies to write
// each element of the edge list as a dot file undirected edge.
//
// Note that while Directed(false) requires arcs to occur in reciprocal pairs,
// UndirectArcs(true) does not, and does not collapse reciprocal arc pairs to
// single dot format edges.
//
// See WriteWeightedEdgeList for more detail.
func UndirectArcs(u bool) func(*Config) {
	return func(c *Config) { c.UndirectArcs = u }
}
