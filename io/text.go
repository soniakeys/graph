// Copyright 2018 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Package graph/io provides graph readers and writers.
package io

import (
	"github.com/soniakeys/graph"
)

// ArcDir specifies whether to consider all arcs, or only arcs that would
// be in the upper or lower triangle of an adjacency matrix representation.
//
// For the case of undirected graphs, the effect of Upper or Lower is to
// specify that the text representation does not contain reciprocal arcs
// but contains only a single arc for each undirected edge.
type ArcDir int

const (
	All   ArcDir = iota // all directed arcs
	Upper               // only arcs where to >= from
	Lower               // only arcs where to <= from
)

// Format defines the fundamental format of text data.
type Format int

const (
	// Sparse is the default.  Each line has a from-node followed by a list
	// of to-nodes. A node with no to-nodes does not need to be listed unless
	// nodes are being stored as NIs and it is the maximum value NI.  In this
	// case it goes on a line by itself to preserve graph order (number of
	// nodes.)  If multiple lines have the same from-node, the to-nodes are
	// concatenated.
	Sparse Format = iota

	// Dense is only meaningful when nodes are being stored as NIs.  The node
	// NI is implied by line number.  There is a line for each node, to-nodes
	// or not.
	Dense

	// Arc format is not actually adjacency list, but arc list or edge list.
	// There are exactly two nodes per line, a from-node and a to-node.
	Arcs
)

// Type Text defines options for reading and writing simple text formats.
//
// The zero value is valid and usable by all methods.  Writing a graph with
// a zero value Text writes a default format that is readable with a zero
// value Text.
//
// Read methods delimit text based on a combination of Text fields:
//
// When MapNames is false, text data is parsed as numeric NIs and LIs in
// the numeric base of field Base, although using 10 if Base is 0.  IDs are
// delimited in this case by any non-empty string of characters that cannot be
// in an integer of the specified base.  For example in the case of base 10
// IDs, any string of characters not in "+-0123456789" will delimit IDs.
//
// When MapNames is true, delimiting depends on FrDelim and ToDelim.  If
// the whitespace-trimmed value of FrDelim or ToDelim is non-empty, text data
// is split on the trimmed value.  If the delimiter is non-empty but
// all whitespace, text data is split on the untrimmed value.  In either
// of these cases, individual fields are then trimmed of leading and trailing
// whitespace.  If FrDelim or ToDelim are empty strings, input text data is
// delimited by non-empty strings of whitespace.
type Text struct {
	Format  Format // Fundamental format of text representation
	Comment string // End of line comment delimiter

	// FrDelim is a delimiter following from-node, for Sparse and Arc formats.
	// Write methods write ": " if FrDelim is blank.  For read behavior
	// see description at Type Text.
	FrDelim string

	// ToDelim separates nodes of to-list, for Sparse and Dense formats.
	// Write methods write " " if ToDelim is blank.  For read behavior see
	// description at Type Text.
	ToDelim string

	// HalfDelim separates the NI and LI of a half arc in a labeled adjacency
	// list.  WriteLabeledAdjacencyList writes " " if HalfDelim is blank.
	HalfDelim string

	// Open and Close surround the NI LI pair of a half arc in a labeled
	// adjacency list.  WriteLabeledAdjacenyList writes "(" and ")" if both
	// are blank.
	Open, Close string

	// Base is the numeric base for NIs and LIs.  Methods pass this to strconv
	// functions and so values should be in the range 2-36, although they will
	// pass 10 to strconv functions when Base is 0.
	Base int

	// MapNames true means to consider node text to be symbolic rather
	// than numeric NIs.  Read methods assign numeric graph NIs as data is
	// read, and return the mapping between node names and NIs.
	MapNames bool

	// A non-nil NodeName is used by write methods to translate NIs to
	// strings and write the strings as symbolic node names rather than numeric
	// NIs.
	NodeName func(graph.NI) string

	// WriteArcs can specify to write only a single arc of an undirected
	// graph.  See definition of ArcDir.
	WriteArcs ArcDir
}

// NewText is a small convenience constructor.
//
// It simply returns &Text{Comment: "//"}.
func NewText() *Text {
	return &Text{Comment: "//"}
}
