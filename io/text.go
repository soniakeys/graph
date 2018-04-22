// Copyright 2018 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Package graph/io provides graph readers and writers.
package io

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/soniakeys/graph"
)

// ArcDir specifies whether to consider all arcs, or for the case of
// undirected graphs, to consider only a single arc for each edge.
type ArcDir int

const (
	All   ArcDir = iota // all directed arcs
	Upper               // for undirected, only arcs where to >= from
	Lower               // for undirected, only arcs where to <= from
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

	// FrDelim is a delimiter following from-node, for Sparse and Edge formats.
	// Write methods write ": " if FrDelim is blank.  For read behavior
	// see description at Type Text.
	FrDelim string

	// ToDelim separates nodes of to-list, for Sparse format.
	// Write methods write " " if ToDelim is blank.  For read behavior see
	// description at Type Text.
	ToDelim string

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
//
// In many cases it will be simpler to just write a struct literal intiializing
// fields as needed.
func NewText() *Text {
	return &Text{Comment: "//"}
}

// ReadAdjacencyList reads text data and returns an  AdjacencyList.
//
// Fields of the receiver Text define how the text data is interpreted.
// See documentation of the Text struct.
//
// ReadAdjacencyList reads to EOF.
//
// On successful read, a valid AdjacencyList is returned with error = nil.
// In addition, with Text.MapNames true, the method returns a list
// of node names indexed by NI and the reverse mapping of NI by name.
func (t Text) ReadAdjacencyList(r io.Reader) (
	graph.AdjacencyList, []string, map[string]graph.NI, error) {
	switch t.Format {
	case Sparse:
		return t.readALSparse(r)
	case Dense:
		return t.readALDense(r)
	case Arcs:
		return t.readALArcs(r)
	}
	return nil, nil, nil, fmt.Errorf("format %d invalid", t.Format)
}

// call before any read or write of NIs (with MapNames false.)
// Note top level API methods generally take non-pointer receivers and so
// operations like this do not modify the caller's struct, they just supply
// a default over the zero value of the struct field.  Note also it's called
// inside of sep() so methods that call sep don't need to call it.
func (t *Text) fixBase() {
	if t.Base == 0 {
		t.Base = 10
	}
}

// package var for most common case of base 10 delimiter.
var bx10 = regexp.MustCompile("[^0-9-+]+")

// return regular expression delimiting numbers of given base.
// This supports the behavior documented "When MapNames is false..."
// for type Text above.
func (t *Text) sep() *regexp.Regexp {
	t.fixBase()
	if t.Base == 10 {
		return bx10
	}
	expr := ""
	if t.Base <= 10 {
		expr = fmt.Sprintf("[^0-%d-+]+", t.Base-1)
	} else {
		expr = fmt.Sprintf("[^0-9a-%c-+]+", 'a'+t.Base-11)
	}
	return regexp.MustCompile(expr)
}

func (t Text) readALDense(r io.Reader) (
	g graph.AdjacencyList, name []string, ni map[string]graph.NI, err error) {
	if t.MapNames {
		return nil, nil, nil,
			fmt.Errorf("name translation not valid for dense format")
	}
	sep := t.sep()
	for b := bufio.NewReader(r); ; {
		f, err := t.readSplitInts(b, sep)
		if err != nil {
			if err != io.EOF {
				return nil, nil, nil, err
			}
			return g, nil, nil, nil
		}
		to, err := parseNIs(f, t.Base)
		if err != nil {
			return nil, nil, nil, err
		}
		g = append(g, to)
	}
}

// read and split a single line of integers.  sep must be compiled based on
// the integer base.
// if the last line is missing the newline, it returns the data and err == nil.
// a subsequent call will return io.EOF.
func (t *Text) readSplitInts(r *bufio.Reader, sep *regexp.Regexp) (
	f []string, err error) {
	s, err := r.ReadString('\n') // but allow last line without \n
	if err != nil {
		if err != io.EOF || s == "" {
			return
		}
	}
	if s == "\n" { // fast path for blank line
		return
	}
	if t.Comment > "" {
		if i := strings.Index(s, t.Comment); i >= 0 {
			s = s[:i]
		}
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return // almost as fast
	}
	// s is non-empty at this point
	f = sep.Split(s, -1)
	// some minimal possibilities:
	// single sep: ":" -> ["", ""], the empty strings before and after sep
	// single non-sep: "0" -> ["0"], a non-empty string
	if f[0] == "" {
		f = f[1:] // toss "" before any leading separator
	}
	last := len(f) - 1
	if f[last] == "" {
		f = f[:last] // toss "" after any trailing separator, perhaps a ","
	}
	return
}

// parse a slice of strings expected to contain valid NIs.
func parseNIs(f []string, base int) (n []graph.NI, err error) {
	if len(f) > 0 {
		n = make([]graph.NI, len(f))
		for x, s := range f {
			i, err := strconv.ParseInt(s, base, graph.NIBits)
			if err != nil {
				return nil, err
			}
			n[x] = graph.NI(i)
		}
	}
	return
}

func (t Text) readALSparse(r io.Reader) (
	g graph.AdjacencyList, name []string, ni map[string]graph.NI, err error) {
	if t.MapNames {
		return t.readALSparseNames(r)
	}
	sep := t.sep()
	b := bufio.NewReader(r)
	for {
		f, err := t.readSplitInts(b, sep)
		if err != nil {
			if err != io.EOF {
				return nil, nil, nil, err
			}
			return g, nil, nil, nil
		}
		if len(f) == 0 {
			continue
		}
		fr, err := strconv.ParseInt(f[0], t.Base, graph.NIBits)
		if err != nil {
			return nil, nil, nil, err
		}
		for int(fr) >= len(g) {
			g = append(g, nil)
		}
		if len(f) == 1 {
			continue // from-NI with no to-list is allowed.
		}
		to, err := parseNIs(f[1:], t.Base)
		if err != nil {
			return nil, nil, nil, err
		}
		if g[fr] == nil {
			g[fr] = to
		} else {
			g[fr] = append(g[fr], to...)
		}
	}
}

func (t Text) readALSparseNames(r io.Reader) (
	g graph.AdjacencyList, name []string, ni map[string]graph.NI, err error) {
	ni = map[string]graph.NI{}
	getNI := func(s string) graph.NI {
		n, ok := ni[s]
		if !ok {
			n = graph.NI(len(g))
			g = append(g, nil)
			name = append(name, s)
			ni[s] = n
		}
		return n
	}
	split := t.splitALSparseNames()
	s := ""
	b := bufio.NewReader(r)
	for line := 1; ; line++ {
		s, err = t.readStripComment(b)
		if err != nil {
			if err != io.EOF {
				return nil, nil, nil, err
			}
			err = nil
			return
		}
		fs, ts := split(s)
		if fs == "" {
			continue
		}
		fr := getNI(fs)
		if len(ts) == 0 {
			continue
		}
		if g[fr] == nil {
			to := make([]graph.NI, len(ts))
			for i, s := range ts {
				to[i] = getNI(s)
			}
			g[fr] = to
		} else {
			for _, s := range ts {
				g[fr] = append(g[fr], getNI(s))
			}
		}
	}
}

func (t *Text) splitALSparseNames() func(string) (string, []string) {
	// simplest case, no delimiters, fields can be split all at once
	if t.FrDelim > "" || t.ToDelim > "" {
		return func(s string) (string, []string) {
			f := strings.Fields(s)
			if len(f) == 0 {
				return "", nil
			}
			return f[0], f[1:]
		}
	}
	// otherwise fr must be split first, then to.
	frIndex := func(s string) int {
		return strings.IndexFunc(s, unicode.IsSpace)
	}
	fdLen := 0
	if t.FrDelim > "" {
		fd := strings.TrimSpace(t.FrDelim)
		if fd == "" {
			fd = t.FrDelim
		}
		fdLen = len(fd)
		frIndex = func(s string) int {
			return strings.Index(s, fd)
		}
	}
	toSplit := strings.Fields
	if t.ToDelim > "" {
		td := strings.TrimSpace(t.ToDelim)
		if td == "" {
			td = t.ToDelim
		}
		toSplit = func(s string) []string {
			to := strings.Split(s, td)
			for i, ti := range to {
				to[i] = strings.TrimSpace(ti)
			}
			return to
		}
	}
	return func(s string) (string, []string) {
		s = strings.TrimLeftFunc(s, unicode.IsSpace)
		x := frIndex(s)
		if x < 0 {
			return "", nil
		}
		return s[:x], toSplit(s[x+fdLen:])
	}
}

func (t *Text) readStripComment(r *bufio.Reader) (s string, err error) {
	s, err = r.ReadString('\n') // but allow last line without \n
	if err != nil {
		if err != io.EOF || s == "" {
			return
		}
	}
	if s == "\n" { // fast path for blank line
		return "", nil
	}
	if t.Comment > "" {
		if i := strings.Index(s, t.Comment); i >= 0 {
			s = s[:i]
		}
	}
	return
}

func (t *Text) readSplitFields(r *bufio.Reader) (f []string, err error) {
	s, err := t.readStripComment(r)
	if err == nil && s > "" {
		f = strings.Fields(s)
	}
	return
}

func (t Text) readALArcs(r io.Reader) (
	g graph.AdjacencyList, name []string, ni map[string]graph.NI, err error) {
	if t.MapNames {
		return t.readALArcNames(r)
	}
	sep := t.sep()
	var max graph.NI
	e := map[int][]graph.NI{} // full graph with to-lists as multisets.
	for b := bufio.NewReader(r); ; {
		f, err := t.readSplitInts(b, sep)
		if err != nil {
			if err != io.EOF {
				return nil, nil, nil, err
			}
			// normal return
			g = make(graph.AdjacencyList, max+1)
			for fr := range g {
				g[fr] = e[fr]
			}
			return g, nil, nil, nil
		}
		if len(f) == 0 {
			continue
		}
		if len(f) != 2 {
			return nil, nil, nil, fmt.Errorf("Arc must have two nodes")
		}
		a, err := parseNIs(f, t.Base)
		if err != nil {
			return nil, nil, nil, err
		}
		fr, to := a[0], a[1]
		if fr < 0 {
			return nil, nil, nil, fmt.Errorf("invalid from: %d", fr)
		}
		if fr > max {
			max = fr
		}
		if to > max {
			max = to
		}
		e[int(fr)] = append(e[int(fr)], graph.NI(to))
	}
}

// ReadArcNames reads a graph from the simple text format written by
// WriteArcNames.
//
// The format is similar to that of ReadArcNIs but reads an input graph of
// named nodes rather than NIs.  Unique node names are accumulated
// in a list parallel to the returned graph, naming each node.  A map, with
// the reverse mapping from name to NI is also accumulated.
//
// An empty string is not allowed as a node name. A final newline is not
// required.
//
// If argument comment is non-empty, it specifies a string marking end-of-line
// comments.
func (t Text) readALArcNames(r io.Reader) (
	g graph.AdjacencyList, name []string, ni map[string]graph.NI, err error) {
	ni = map[string]graph.NI{}
	getNI := func(s string) graph.NI {
		n, ok := ni[s]
		if !ok {
			n = graph.NI(len(g))
			g = append(g, nil)
			name = append(name, s)
			ni[s] = n
		}
		return n
	}
	delIndex := func(s string) int { return strings.Index(s, t.FrDelim) }
	if t.FrDelim == "" {
		delIndex = func(s string) int {
			return strings.IndexFunc(s, unicode.IsSpace)
		}
	}
	b := bufio.NewReader(r)
	s := ""
	for line := 1; ; line++ {
		s, err = b.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, nil, nil, fmt.Errorf("line %d: %v", line, err)
			}
			if s == "" {
				err = nil
				return
			}
			// allow last line without \n
		}
		if t.Comment > "" {
			if i := strings.Index(s, t.Comment); i >= 0 {
				s = s[:i]
			}
		}
		s = strings.TrimSpace(s)
		if s == "" {
			continue // allow and ignore blank line
		}
		i := delIndex(s)
		if i < 0 {
			return nil, nil, nil,
				fmt.Errorf("line %d: t.FrDelimiter required", line)
		}
		fn := strings.TrimSpace(s[:i])
		if fn == "" {
			return nil, nil, nil,
				fmt.Errorf("line %d: blank name not allowed", line)
		}
		tn := strings.TrimSpace(s[i+len(t.FrDelim):])
		if tn == "" {
			return nil, nil, nil,
				fmt.Errorf("line %d: blank name not allowed", line)
		}
		fr := getNI(fn)
		g[fr] = append(g[fr], getNI(tn))
	}
}

// ReadLabeledAdjacencyList reads text data and returns a LabeledAdjacencyList.
//
// Fields of the receiver Text define how the text data is interpreted.
// See documentation of the Text struct.
//
// ReadLabeledAdjacencyList reads to EOF.
//
// On successful read, a valid AdjacencyList is returned with error = nil.
// In addition, with Text.MapNames true, the method returns a list
// of node names indexed by NI and the reverse mapping of NI by name.
func (t Text) ReadLabeledAdjacencyList(r io.Reader) (g graph.LabeledAdjacencyList, err error) {
	sep := t.sep()
	for b := bufio.NewReader(r); ; {
		to, err := t.readHalf(b, sep)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return g, nil
		}
		g = append(g, to)
	}
}

// read and parse a single line of a to-list.  if the last line is missing
// the newline, it returns the data and err == nil.  a subsequent call will
// return io.EOF.
func (t *Text) readHalf(r *bufio.Reader, sep *regexp.Regexp) (to []graph.Half, err error) {
	s, err := r.ReadString('\n') // but allow last line without \n
	if err != nil {
		if err != io.EOF || s == "" {
			return
		}
	}
	if s == "\n" { // fast path for blank line
		return
	}
	// s is non-empty at this point
	f := sep.Split(s, -1)
	// some minimal possibilities:
	// single sep: "\n" -> ["", ""], the empty strings before and after sep
	// single non-sep: "0" -> ["0"], a non-empty string
	if f[0] == "" {
		f = f[1:] // toss "" before leading any separator
	}
	last := len(f) - 1
	if f[last] == "" {
		f = f[:last] // toss "" after trailing separator, usually the \n
	}
	if len(f)%2 != 0 {
		return nil, fmt.Errorf("odd data")
	}
	to = make([]graph.Half, len(f)/2)
	y := 0
	for x := range to {
		ni, err := strconv.ParseInt(f[y], t.Base, graph.NIBits)
		if err != nil {
			return nil, err
		}
		y++
		li, err := strconv.ParseInt(f[y], t.Base, graph.NIBits)
		if err != nil {
			return nil, err
		}
		y++
		to[x] = graph.Half{graph.NI(ni), graph.LI(li)}
	}
	return to, nil
}

// WriteAdjacencyList writes an adjacency list as text.
//
// Fields of the receiver Text define how the text data is formatted.
// See documentation of the Text struct.
//
// Returned is number of bytes written and error.
func (t Text) WriteAdjacencyList(g graph.AdjacencyList, w io.Writer) (
	int, error) {
	switch t.Format {
	case Dense:
		return t.writeAdjacencyListDense(g, w)
	case Sparse:
		return t.writeAdjacencyListSparse(g, w)
	case Arcs:
		return t.writeArcs(g, w)
	}
	return 0, fmt.Errorf("format %d unimplemente", t.Format)
}

// WriteAdjacencyList writes an adjacency list in a simple text format.
//
// A line is written for each node, consisting of the to-list, the NIs
// formatted as space separated base 10 numbers.  The NI of the from-node
// is not written but is implied by the line number.
//
// Returned is number of bytes written and error.
func (t Text) writeAdjacencyListDense(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	if t.WriteArcs == Upper {
		return t.writeUpperDense(g, w)
	}
	t.fixBase()
	b := bufio.NewWriter(w)
	var c int
	for _, to := range g {
		if len(to) > 0 {
			c, err = b.WriteString(strconv.FormatInt(int64(to[0]), t.Base))
			n += c
			if err != nil {
				return
			}
			for _, to := range to[1:] {
				if err = b.WriteByte(' '); err != nil {
					return
				}
				c, err = b.WriteString(strconv.FormatInt(int64(to), t.Base))
				n += c + 1
				if err != nil {
					return
				}
			}
		}
		if err = b.WriteByte('\n'); err != nil {
			return
		}
		n++
	}
	b.Flush()
	return
}

// WriteAdjacencyListNIs writes an adjacency list in a simple text format.
//
// The format is similar to Go keyed composite literals.  Each line has a
// from-NI as a "key" followed by ": " and a list of to-NIs.
// Nodes with no to-nodes generate no output, except a maximum NI with no
// to-nodes will be written with an empty to-list.
//
// Returned is number of bytes written and error.
func (t Text) writeAdjacencyListSparse(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	if t.FrDelim == "" {
		t.FrDelim = ": "
	}
	if t.ToDelim == "" {
		t.ToDelim = " "
	}
	if t.WriteArcs == Upper {
		return t.writeUpperSparse(g, w)
	}
	writeLast := t.NodeName == nil
	if writeLast {
		t.fixBase()
		t.NodeName = func(n graph.NI) string {
			return strconv.FormatInt(int64(n), t.Base)
		}
	}
	b := bufio.NewWriter(w)
	var c int
	last := len(g) - 1
	for i, to := range g {
		switch {
		case len(to) > 0:
			fr := graph.NI(i)
			c, err = b.WriteString(t.NodeName(fr))
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(t.FrDelim)
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(t.NodeName(to[0]))
			n += c
			if err != nil {
				return
			}
			for _, to := range to[1:] {
				c, err = b.WriteString(t.ToDelim)
				n += c
				if err != nil {
					return
				}
				c, err = b.WriteString(t.NodeName(to))
				n += c
				if err != nil {
					return
				}
			}
		case i == last && writeLast:
			fr := graph.NI(i)
			c, err = b.WriteString(t.NodeName(fr))
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(
				strings.TrimRightFunc(t.FrDelim, unicode.IsSpace))
			n += c
			if err != nil {
				return
			}
		default:
			continue // without writing \n
		}
		if err = b.WriteByte('\n'); err != nil {
			return
		}
		n++
	}
	b.Flush()
	return
}

// WriteArcNIs writes arcs of an adjacency list in a simple text format.
//
// Each line has an arc, a from-NI and to-NI, separated with a space.
//
// Returned is number of bytes written and error.
func (t Text) writeArcs(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	t.fixBase()
	if t.FrDelim == "" {
		t.FrDelim = " "
	}
	if t.NodeName == nil {
		t.NodeName = func(n graph.NI) string {
			return strconv.FormatInt(int64(n), t.Base)
		}
	}
	b := bufio.NewWriter(w)
	var c int
	for fr, to := range g {
		for _, to := range to {
			c, err = b.WriteString(t.NodeName(graph.NI(fr)))
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(t.FrDelim)
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(t.NodeName(to))
			n += c
			if err != nil {
				return
			}
			if err = b.WriteByte('\n'); err != nil {
				return
			}
			n++
		}
	}
	b.Flush()
	return
}

// WriteUpper writes the "upper triangle" of an adjacency list.
//
// For an adjacency list representing an undirected graph, this writes
// an arc for each undirected edge but omits reciprocal arcs.
//
// A line is written for each node, consisting of to-NIs greater than or
// equal to the from-NI.  The from-NI is not written but is implied by the
// line number.
//
// Returned is number of bytes written and error.
func (t Text) writeUpperDense(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	t.fixBase()
	b := bufio.NewWriter(w)
	var c int
	for i, to := range g {
		fr := graph.NI(i)
		one := false
		for _, to := range to {
			if to >= fr {
				if one {
					if err = b.WriteByte(' '); err != nil {
						return
					}
					n++
				} else {
					one = true
				}
				c, err = b.WriteString(strconv.FormatInt(int64(to), t.Base))
				n += c
				if err != nil {
					return
				}
			}
		}
		if err = b.WriteByte('\n'); err != nil {
			return
		}
		n++
	}
	b.Flush()
	return
}

// WriteUpperNames writes the "upper triangle" of an adjacency list.
//
// Returned is number of bytes written and error.
/*
func (t Text) WriteUpperNames(g graph.AdjacencyList, w io.Writer,
	format func(graph.NI) string) (n int, err error) {
	return t.writeUpper(g, w, format, false)
}

// WriteUpperNIs writes the "upper triangle" of an adjacency list.
//
// For an adjacency list representing an undirected graph, this writes
// an arc for each undirected edge but omits reciprocal arcs.
//
// A line is written for each from-node that has to-nodes >= the from-node.
// The format is similar to Go keyed composite literals.  Each line has the
// from-NI as a "key" followed by ": " and a list of to-NIs.  Nodes with no
// to-nodes >= the from-node generate no output, except a maximum NI with no
// to-nodes will be written with an empty to-list.
//
// Returned is number of bytes written and error.
func (t Text) WriteUpperNIs(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	t.fixBase()
	return t.writeUpper(g, w, func(n graph.NI) string {
		return strconv.FormatInt(int64(n), t.Base)
	}, true)
}
*/

func (t Text) writeUpperSparse(g graph.AdjacencyList, w io.Writer) (
	n int, err error) {
	b := bufio.NewWriter(w)
	writeLast := t.NodeName == nil
	if writeLast {
		t.fixBase()
		t.NodeName = func(n graph.NI) string {
			return strconv.FormatInt(int64(n), t.Base)
		}
	}
	var c int
	last := len(g) - 1
	for i, to := range g {
		fr := graph.NI(i)
		one := false
		for _, to := range to {
			if to >= fr {
				if !one {
					one = true
					c, err = b.WriteString(t.NodeName(fr))
					n += c
					if err != nil {
						return
					}
					c, err = b.WriteString(t.FrDelim)
					n += c
					if err != nil {
						return
					}
				} else {
					c, err = b.WriteString(t.ToDelim)
					n += c
					if err != nil {
						return
					}
				}
				c, err = b.WriteString(t.NodeName(to))
				n += c
				if err != nil {
					return
				}
			}
		}
		if writeLast && i == last && !one {
			one = true
			c, err = b.WriteString(t.NodeName(fr))
			n += c
			if err != nil {
				return
			}
			c, err = b.WriteString(
				strings.TrimRightFunc(t.FrDelim, unicode.IsSpace))
			n += c
			if err != nil {
				return
			}
		}
		if one {
			if err = b.WriteByte('\n'); err != nil {
				return
			}
			n++
		}
	}
	b.Flush()
	return
}

// WriteLabeledAdjacencyList writes a lableed adjacency list as text.
//
// Fields of the receiver Text define how the text data is formatted.
// See documentation of the Text struct.
//
// Returned is number of bytes written and error.
func (t Text) WriteLabeledAdjacencyList(g graph.LabeledAdjacencyList,
	w io.Writer) (n int, err error) {
	t.fixBase()
	b := bufio.NewWriter(w)
	var c int
	for _, to := range g {
		if len(to) > 0 {
			c, err = b.WriteString(strconv.FormatInt(int64(to[0].To), t.Base))
			n += c
			if err != nil {
				return
			}
			if err = b.WriteByte(' '); err != nil {
				return
			}
			c, err = b.WriteString(strconv.FormatInt(int64(to[0].Label), t.Base))
			n += c + 1
			if err != nil {
				return
			}
			for _, to := range to[1:] {
				if err = b.WriteByte(' '); err != nil {
					return
				}
				c, err = b.WriteString(strconv.FormatInt(int64(to.To), t.Base))
				n += c + 1
				if err != nil {
					return
				}
				if err = b.WriteByte(' '); err != nil {
					return
				}
				c, err = b.WriteString(strconv.FormatInt(int64(to.Label), t.Base))
				n += c + 1
				if err != nil {
					return
				}
			}
		}
		if err = b.WriteByte('\n'); err != nil {
			return
		}
		n++
	}
	b.Flush()
	return
}
