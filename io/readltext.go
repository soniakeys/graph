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

	"github.com/soniakeys/graph"
)

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
	sep, err := t.sep()
	if err != nil {
		return nil, err
	}
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
