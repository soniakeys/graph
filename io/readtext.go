// Copyright 2018 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Package graph/io provides graph readers and writers.
package io

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/soniakeys/graph"
)

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

func (t Text) readALSparse(r io.Reader) (
	g graph.AdjacencyList, name []string, ni map[string]graph.NI, err error) {
	if t.MapNames {
		return t.readALSparseNames(r)
	}
	sep, err := t.sep()
	if err != nil {
		return nil, nil, nil, err
	}
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
			panic(fmt.Sprintf("in readALSparse: %v", err))
		}
		for int(fr) >= len(g) {
			g = append(g, nil)
		}
		if len(f) == 1 {
			continue // from-NI with no to-list is allowed.
		}
		to := parseNIs(f[1:], t.Base)
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
	split := t.sparseNameSplitter()
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
			if len(ts) > 0 {
				return nil, nil, nil, errors.New("blank node name")
			}
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

func (t *Text) sparseNameSplitter() func(string) (string, []string) {
	// simplest case, no delimiters, fields can be split all at once
	if t.FrDelim == "" && t.ToDelim == "" {
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
			nb := 0
			for _, ti := range to {
				if ti = strings.TrimSpace(ti); ti > "" {
					to[nb] = ti
					nb++
				}
			}
			return to[:nb]
		}
	}
	return func(s string) (string, []string) {
		s = strings.TrimLeftFunc(s, unicode.IsSpace)
		x := frIndex(s)
		if x < 0 {
			return "", toSplit(s)
		}
		return strings.TrimRightFunc(s[:x], unicode.IsSpace),
			toSplit(s[x+fdLen:])
	}
}

func (t Text) readALDense(r io.Reader) (
	g graph.AdjacencyList, name []string, ni map[string]graph.NI, err error) {
	if t.MapNames {
		return nil, nil, nil,
			fmt.Errorf("name translation not valid for dense format")
	}
	sep, err := t.sep()
	if err != nil {
		return nil, nil, nil, err
	}
	for b := bufio.NewReader(r); ; {
		f, err := t.readSplitInts(b, sep)
		if err != nil {
			if err != io.EOF {
				return nil, nil, nil, err
			}
			return g, nil, nil, nil
		}
		g = append(g, parseNIs(f, t.Base))
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
	return f, nil
}

// parse a slice of strings expected to contain valid NIs.  The slice can be
// empty, but any strings present must parse by strconv in the specified base.
// This is a "MustParse" function.
// an error from ParseInt cannot be a user error, it must be bug, so panic
func parseNIs(f []string, base int) (n []graph.NI) {
	if len(f) > 0 {
		n = make([]graph.NI, len(f))
		for x, s := range f {
			i, err := strconv.ParseInt(s, base, graph.NIBits)
			if err != nil {
				panic(err)
			}
			n[x] = graph.NI(i)
		}
	}
	return
}

// normal return: two non-empty strings.
// second string empty is okay, means just define node name.
// empty first and non-empty second should be considered an error.
func (t *Text) arcNameSplitter() func(string) (string, string) {
	index := func(s string) int {
		return strings.IndexFunc(s, unicode.IsSpace)
	}
	fdLen := 0
	if t.FrDelim > "" {
		fd := strings.TrimSpace(t.FrDelim)
		if fd == "" {
			fd = t.FrDelim
		}
		fdLen = len(fd)
		index = func(s string) int {
			return strings.Index(s, fd)
		}
	}
	return func(s string) (string, string) {
		s = strings.TrimLeftFunc(s, unicode.IsSpace)
		x := index(s)
		if x < 0 {
			return strings.TrimRightFunc(s, unicode.IsSpace), ""
		}
		return strings.TrimRightFunc(s[:x], unicode.IsSpace),
			strings.TrimSpace(s[x+fdLen:])
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
	return s, nil
}

func (t Text) readALArcs(r io.Reader) (
	g graph.AdjacencyList, name []string, ni map[string]graph.NI, err error) {
	if t.MapNames {
		return t.readALArcNames(r)
	}
	sep, err := t.sep()
	if err != nil {
		return nil, nil, nil, err
	}
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
		if len(f) > 2 {
			return nil, nil, nil, fmt.Errorf("Arc can only have two nodes")
		}
		a := parseNIs(f, t.Base)
		fr := a[0]
		if fr < 0 {
			return nil, nil, nil, fmt.Errorf("invalid from: %d", fr)
		}
		if fr > max {
			max = fr
		}
		if len(f) == 2 {
			to := a[1]
			if to > max {
				max = to
			}
			e[int(fr)] = append(e[int(fr)], graph.NI(to))
		}
	}
}

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
	split := t.arcNameSplitter()
	b := bufio.NewReader(r)
	s := ""
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
			if len(ts) > 0 {
				return nil, nil, nil, errors.New("blank from-node")
			}
			continue
		}
		fr := getNI(fs)
		if len(ts) > 0 {
			g[fr] = append(g[fr], getNI(ts))
		}
	}
}

// package var for most common case of base 10 delimiter.
var bx10 = regexp.MustCompile("[^0-9-+]+")

// return regular expression delimiting numbers of given base.
// This supports the behavior documented "When MapNames is false..."
// for type Text above.
func (t *Text) sep() (*regexp.Regexp, error) {
	if err := t.fixBase(); err != nil {
		return nil, err
	}
	if t.Base == 10 {
		return bx10, nil
	}
	expr := ""
	if t.Base <= 10 {
		expr = fmt.Sprintf("[^0-%d-+]+", t.Base-1)
	} else {
		expr = fmt.Sprintf("[^0-9a-%c-+]+", 'a'+t.Base-11)
	}
	return regexp.MustCompile(expr), nil
}

// call before any read or write of NIs (with MapNames false.)
// Note top level API methods generally take non-pointer receivers and so
// operations like this do not modify the caller's struct, they just supply
// a default over the zero value of the struct field.  Note also it's called
// inside of sep() so methods that call sep don't need to call it.
func (t *Text) fixBase() error {
	if t.Base == 0 {
		t.Base = 10
		return nil
	}
	if t.Base < 2 || t.Base > 36 {
		return errors.New("invalid Text.Base")
	}
	return nil
}
