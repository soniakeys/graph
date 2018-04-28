// Copyright 2018 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package io_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/io"
)

type allErr struct{}

func (allErr) Read([]byte) (int, error) {
	return 0, errors.New("always error")
}

func TestReadALArcNames(t *testing.T) {
	// test a read error
	tx := io.Text{Format: io.Arcs, MapNames: true}
	if _, _, _, err := tx.ReadAdjacencyList(allErr{}); err == nil {
		t.Fatal("readALArcNames allowed read error")
	}

	// normal operation
	r := bytes.NewBufferString(`
a c 
b
a d
a d
c d`)
	g, names, m, err := tx.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal(err)
	}
	want := graph.AdjacencyList{
		0: {1, 3, 3},
		1: {3},
		3: {},
	}
	if !g.Equal(want) {
		for fr, to := range g {
			t.Logf("%d %d", fr, to)
		}
		t.Fatal("g")
	}
	if !reflect.DeepEqual(names, []string{"a", "c", "b", "d"}) {
		t.Fatal("names")
	}
	if graph.OrderMap(m) != "map[a:0 b:2 c:1 d:3]" {
		t.Fatal("map")
	}

	// test blank from
	tx.FrDelim = ":"
	r = bytes.NewBufferString(`:z`)
	if _, _, _, err = tx.ReadAdjacencyList(r); err == nil {
		t.Fatal("readALArcNames allowed blank from")
	}
}

func TestArcNameSplitter(t *testing.T) {
	// with FrDelim
	r := bytes.NewBufferString(`
a->c
b
a->d
a->d
c->d`)
	tx := io.Text{Format: io.Arcs, MapNames: true, FrDelim: "->"}
	g, names, m, err := tx.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal(err)
	}
	want := graph.AdjacencyList{
		0: {1, 3, 3},
		1: {3},
		3: {},
	}
	if !g.Equal(want) {
		for fr, to := range g {
			t.Logf("%d %d", fr, to)
		}
		t.Fatal("g")
	}
	if !reflect.DeepEqual(names, []string{"a", "c", "b", "d"}) {
		t.Fatal("names")
	}
	if graph.OrderMap(m) != "map[a:0 b:2 c:1 d:3]" {
		t.Fatal("map")
	}

	// with whitespace FrDelim
	r = bytes.NewBufferString(`
a	c
b
a	d
a	d
c	d`)
	tx.FrDelim = "\t"
	g, names, m, err = tx.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal(err)
	}
	if !g.Equal(want) {
		for fr, to := range g {
			t.Logf("%d %d", fr, to)
		}
		t.Fatal("g")
	}
	if !reflect.DeepEqual(names, []string{"a", "c", "b", "d"}) {
		t.Fatal("names")
	}
	if graph.OrderMap(m) != "map[a:0 b:2 c:1 d:3]" {
		t.Fatal("map")
	}
}

func TestReadALSparseNames(t *testing.T) {
	// test a read error
	tx := io.Text{MapNames: true}
	if _, _, _, err := tx.ReadAdjacencyList(allErr{}); err == nil {
		t.Fatal("readALSparse allowed read error")
	}

	// test empty to-list, additional arcs from a node
	r := bytes.NewBufferString(`
a b c
a c
b
d e`)
	g, names, m, err := tx.ReadAdjacencyList(r)
	// result should be same as example without the b line,
	// tediously checked here
	if err != nil {
		t.Fatalf("no error expected.  got %v", err)
	}
	if !reflect.DeepEqual(names, []string{"a", "b", "c", "d", "e"}) {
		t.Fatal("names: ", names)
	}
	if graph.OrderMap(m) != "map[a:0 b:1 c:2 d:3 e:4]" {
		t.Fatal("map: ", m)
	}
	want := graph.AdjacencyList{
		0: {1, 2, 2},
		3: {4},
		4: {},
	}
	if !g.Equal(want) {
		t.Fatal("g wrong")
	}

	// test delimiters.  setting both hits most code in splitALSparseNames
	tx.FrDelim = "->"
	tx.ToDelim = ":"
	r = bytes.NewBufferString(`
a->b:c
d->e`)
	g, names, m, err = tx.ReadAdjacencyList(r)
	if err != nil {
		t.Fatalf("delims. no error expected.  got %v", err)
	}
	if !reflect.DeepEqual(names, []string{"a", "b", "c", "d", "e"}) {
		t.Fatal("delims. names: ", names)
	}
	if graph.OrderMap(m) != "map[a:0 b:1 c:2 d:3 e:4]" {
		t.Fatal("delims. map: ", m)
	}
	want[0] = []graph.NI{1, 2}
	if !g.Equal(want) {
		t.Fatal("delims. g wrong")
	}

	// setting just ToDelim exercises default for FrDelim
	tx.FrDelim = ""
	r = bytes.NewBufferString(`
a  b : c
d	e`)
	g, names, m, err = tx.ReadAdjacencyList(r)
	if err != nil {
		t.Fatalf("ToDelim. no error expected.  got %v", err)
	}
	if !reflect.DeepEqual(names, []string{"a", "b", "c", "d", "e"}) {
		t.Fatal("ToDelim names: ", names)
	}
	if graph.OrderMap(m) != "map[a:0 b:1 c:2 d:3 e:4]" {
		t.Fatal("ToDelim map: ", m)
	}
	if !g.Equal(want) {
		t.Fatal("ToDelimn g wrong")
	}

	// set both to whitespace
	tx.FrDelim = "	"  // tab
	tx.ToDelim = "  " // two spaces
	r = bytes.NewBufferString(`
a	b   c
d	    e f`)
	g, names, m, err = tx.ReadAdjacencyList(r)
	if err != nil {
		t.Fatalf("ToDelim. no error expected.  got %v", err)
	}
	if !reflect.DeepEqual(names, []string{"a", "b", "c", "d", "e f"}) {
		t.Fatal("ToDelim names: ", names)
	}
	if graph.OrderMap(m) != "map[a:0 b:1 c:2 d:3 e f:4]" {
		t.Fatal("ToDelim map: ", m)
	}
	if !g.Equal(want) {
		t.Fatal("ToDelimn g wrong")
	}

	// only whitespace in front of ToDelim is error
	r = bytes.NewBufferString(`  	b`)
	_, _, _, err = tx.ReadAdjacencyList(r)
	if err == nil {
		t.Fatalf("just ws error expected")
	}
}

func TestReadALSparse(t *testing.T) {
	// a few more tests for coverage:
	// test invalid base
	_, _, _, err := io.Text{Base: 1}.ReadAdjacencyList(nil)
	if err == nil {
		t.Fatal("readALSparse allowed invalid base")
	}
	// test a read error
	_, _, _, err = io.Text{}.ReadAdjacencyList(allErr{})
	if err == nil {
		t.Fatal("readALSparse allowed read error")
	}
	// additional arcs on a node
	want := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	r := bytes.NewBufferString(`
0: 2 1
0: 1
2: 1`)
	got, _, _, err := io.Text{FrDelim: ":"}.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal("readSplitInts test: ", err)
	}
	if !got.Equal(want) {
		for fr, to := range got {
			t.Log(fr, " ", to)
		}
		t.Fail()
	}
}

func TestReadALDense(t *testing.T) {
	// test normal operation
	want := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	r := bytes.NewBufferString(`2 1 1

1`)
	got, _, _, err := io.Text{Format: io.Dense}.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal("readALDense: ", err)
	}
	if !got.Equal(want) {
		for fr, to := range got {
			t.Log(fr, " ", to)
		}
		t.Fail()
	}
	// test invalid base
	tx := io.Text{Format: io.Dense, Base: 1}
	if _, _, _, err = tx.ReadAdjacencyList(nil); err == nil {
		t.Fatal("readALDense allowed invalid base")
	}
	// test MapNames invalid
	tx = io.Text{Format: io.Dense, MapNames: true}
	if _, _, _, err := tx.ReadAdjacencyList(nil); err == nil {
		t.Fatal("readALDense allowed MapNames")
	}
	// test a read error
	_, _, _, err = io.Text{Format: io.Dense}.ReadAdjacencyList(allErr{})
	if err == nil {
		t.Fatal("readALDense allowed read error")
	}
}

func TestReadALArcs(t *testing.T) {
	// test a read error
	tx := io.Text{Format: io.Arcs}
	if _, _, _, err := tx.ReadAdjacencyList(allErr{}); err == nil {
		t.Fatal("readALDense allowed read error")
	}

	// test normal operation
	want := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	r := bytes.NewBufferString(`
0 2
0 1
0 1
2 1`)
	got, _, _, err := tx.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal("readALArcs: ", err)
	}
	if !got.Equal(want) {
		for fr, to := range got {
			t.Log(fr, " ", to)
		}
		t.Fail()
	}

	// 3 IDs?
	r = bytes.NewBufferString(`1 2 3`)
	if _, _, _, err = tx.ReadAdjacencyList(r); err == nil {
		t.Fatal("expected error for 3 IDs")
	}

	// bad from
	r = bytes.NewBufferString(`
1 0
-1`)
	if _, _, _, err = tx.ReadAdjacencyList(r); err == nil {
		t.Fatal("readALArcs allowed invalid from")
	}

	// test invalid base
	tx.Base = 1
	if _, _, _, err = tx.ReadAdjacencyList(nil); err == nil {
		t.Fatal("readALArcs allowed invalid base")
	}
}

func TestReadAdjacencyList(t *testing.T) {
	// test bad format
	_, _, _, err := io.Text{Format: -1}.ReadAdjacencyList(nil)
	if err == nil {
		t.Fatal("ReadAdacencyList no err from bad Format")
	}
}

func TestWriteAdjacencyList(t *testing.T) {
	// test bad format
	_, err := io.Text{Format: -1}.WriteAdjacencyList(nil, nil)
	if err == nil {
		t.Fatal("ReadAdacencyList no err from bad Format")
	}
}

func TestSep(t *testing.T) {
	// test a base < 10
	want := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	r := bytes.NewBufferString(`
0:829181
2:91`)
	got, _, _, err := io.Text{Base: 8}.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal("sep: ", err)
	}
	if !got.Equal(want) {
		for fr, to := range got {
			t.Log(fr, " ", to)
		}
		t.Fail()
	}
	// test a base > 10
	const zz = 36*36 - 1
	want = graph.AdjacencyList{
		zz: {0},
	}
	r = bytes.NewBufferString(`zz: 0`)
	got, _, _, err = io.Text{Base: 36}.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal("sep: ", err)
	}
	if !got.Equal(want) {
		for fr, to := range got {
			t.Log(fr, " ", to)
		}
		t.Fail()
	}
}

func TestReadSplitInts(t *testing.T) {
	// a couple of cases for test coverage
	want := graph.AdjacencyList{
		0: {2, 1, 1},
		2: {1},
	}
	r := bytes.NewBufferString(` 
:
0: 2 1 1
2: 1`)
	got, _, _, err := io.Text{FrDelim: ":"}.ReadAdjacencyList(r)
	if err != nil {
		t.Fatal("readSplitInts test: ", err)
	}
	if !got.Equal(want) {
		for fr, to := range got {
			t.Log(fr, " ", to)
		}
		t.Fail()
	}
}

func TestWriteALDenseTriangle(t *testing.T) {
        //   0
        //  / \\
        // 1   2--\
        //      \-/
        var g graph.Undirected
        g.AddEdge(0, 1)
        g.AddEdge(0, 2)
        g.AddEdge(0, 2)
        g.AddEdge(2, 2)

	// test upper
	var b bytes.Buffer
        tx := io.Text{Format: io.Dense, WriteArcs: io.Upper}
        n, err := tx.WriteAdjacencyList(g.AdjacencyList, &b)
	if err != nil {
		t.Fatal(err)
	}
	if n != 9 {
		t.Fatal("got", n, "bytes")
	}
	got := b.String()
	if got != "1 2 2\n\n2\n" {
		t.Fatalf("got %q", got)
	}
	// test lower
	b.Reset()
	tx.WriteArcs = io.Lower
        n, err = tx.WriteAdjacencyList(g.AdjacencyList, &b)
	if err != nil {
		t.Fatal(err)
	}
	if n != 9 {
		t.Fatal("got", n, "bytes")
	}
	got = b.String()
	if got != "\n0\n0 0 2\n" {
		t.Fatalf("got %q", got)
	}
}

