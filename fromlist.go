// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// fromlist.go

package graph

import "math/big"

// FromList represents a rooted tree (or forest) where each node is associated
// with a half arc identifying an arc "from" another node.
//
// Other terms for this data structure include "parent list",
// "predecessor list", "in-tree", "inverse arborescence", and
// "spaghetti stack."
//
// The Paths member represents the tree structure.  Leaves and MaxLen are
// not always needed.  Where Leaves is used it serves as a bitmap where
// Leave.Bit() == 1 for each leaf of the tree.  Where MaxLen is used it is
// provided primarily as a convenience for functions that might want to
// anticipate the maximum path length that would be encountered traversing
// the tree.
//
// Various graph search functions use a FromList to returns search results.
// For a start node of a search, From will be -1 and Len will be 1. For other
// nodes reached by the search, From represents a half arc in a path back to
// start and Len represents the number of nodes in the path.  For nodes not
// reached by the search, From will be -1 and Len will be 0.
//
// A single FromList can also represent a forest.  In this case paths from
// all leaves do not return to a single root node, but multiple root nodes.
type FromList struct {
	Paths  []PathEnd // tree representation
	Leaves big.Int   // leaves of tree
	MaxLen int       // length of longest path, max of all PathEnd.Len values
}

// PathEnd associates a half arc and a path length.
//
// A PathEnd list is an element type of FromList.
type PathEnd struct {
	From NI  // a "from" half arc, the node the arc comes from
	Len  int // number of nodes in path from start
}

// NewFromList creates a FromList object of given order.
func NewFromList(n int) FromList {
	return FromList{Paths: make([]PathEnd, n)}
}

// reset initializes a FromList in preparation for a search.  Search methods
// will call this function and you don't typically call it from application
// code.
func (f *FromList) reset() {
	for n := range f.Paths {
		f.Paths[n] = PathEnd{From: -1, Len: 0}
	}
	f.Leaves = big.Int{}
	f.MaxLen = 0
}

// BoundsOk validates the "from" values in the list.
//
// Negative values are allowed as they indicate root nodes.
//
// BoundsOk returns true when all from values are less than len(t).
// Otherwise it returns false and a node with a from value >= len(t).
func (f *FromList) BoundsOk() (ok bool, n NI) {
	for n, e := range f.Paths {
		if int(e.From) >= len(f.Paths) {
			return false, NI(n)
		}
	}
	return true, -1
}

// CommonAncestor returns the common ancestor of a and b.
//
// It returns -1 if a or b are invalid node numbers.
func (f *FromList) CommonAncestor(a, b NI) NI {
	p := f.Paths
	if a < 0 || b < 0 || a >= NI(len(p)) || b >= NI(len(p)) {
		return -1
	}
	if p[a].Len < p[b].Len {
		a, b = b, a
	}
	for bl := p[b].Len; p[a].Len > bl; {
		a = p[a].From
	}
	for a != b {
		a = p[a].From
		b = p[b].From
	}
	return a
}

// PathTo decodes a FromList, recovering a single path.
//
// The path is returned as a list of nodes where the first element will be
// a root node and the last element will be the specified end node.
//
// Only the Paths member of the receiver is used.  Other members of the
// FromList do not need to be valid, however the MaxLen member can be useful
// for allocating argument p.
//
// Argument p can provide the result slice.  If p has capacity for the result
// it will be used, otherwise a new slice is created for the result.
//
// See also function PathTo.
func (f *FromList) PathTo(end NI, p []NI) []NI {
	return PathTo(f.Paths, end, p)
}

// PathTo decodes a single path from a PathEnd list.
//
// The path is returned as a list of nodes where the first element will be
// a root node and the last element will be the specified end node.
//
// Argument p can provide the result slice.  If p has capacity for the result
// it will be used, otherwise a new slice is created for the result.
//
// See also method FromList.PathTo.
func PathTo(paths []PathEnd, end NI, p []NI) []NI {
	n := paths[end].Len
	if n == 0 {
		return nil
	}
	if cap(p) >= n {
		p = p[:n]
	} else {
		p = make([]NI, n)
	}
	for {
		n--
		p[n] = end
		if n == 0 {
			return p
		}
		end = paths[end].From
	}
}

// RecalcLeaves recomputes the Leaves member of f.
func (f *FromList) RecalcLeaves() {
	p := f.Paths
	lv := &f.Leaves
	OneBits(lv, len(p))
	for n := range f.Paths {
		if fr := p[n].From; fr >= 0 {
			lv.SetBit(lv, int(fr), 0)
		}
	}
}

// RecalcLen recomputes Len for each path end, and recomputes MaxLen.
//
// RecalcLen relies on the Leaves member being valid.  If it is not known
// to be valid, call RecalcLeaves before calling RecalcLen.
func (f *FromList) RecalcLen() {
	p := f.Paths
	var setLen func(NI) int
	setLen = func(n NI) int {
		switch {
		case p[n].Len > 0:
			return p[n].Len
		case p[n].From < 0:
			p[n].Len = 1
			return 1
		}
		l := 1 + setLen(p[n].From)
		p[n].Len = l
		return l
	}
	for n := range f.Paths {
		p[n].Len = 0
	}
	f.MaxLen = 0
	lv := &f.Leaves
	for n := range f.Paths {
		if lv.Bit(n) == 1 {
			if l := setLen(NI(n)); l > f.MaxLen {
				f.MaxLen = l
			}
		}
	}
}

// ReRoot reorients the tree containing n to make n the root node.
//
// It keeps the tree connected by "reversing" the path from n to the old root.
//
// After ReRoot, the Leaves and Len members are invalid.
// Call RecalcLeaves or RecalcLen as needed.
func (f *FromList) ReRoot(n NI) {
	p := f.Paths
	fr := p[n].From
	if fr < 0 {
		return
	}
	p[n].From = -1
	for {
		ff := p[fr].From
		p[fr].From = n
		if ff < 0 {
			return
		}
		n = fr
		fr = ff
	}
}

// Root finds the root of a node in a FromList.
func (f *FromList) Root(n NI) NI {
	for p := f.Paths; ; {
		fr := p[n].From
		if fr < 0 {
			return n
		}
		n = fr
	}
}

// Transpose constructs the directed graph corresponding to FromList f
// but with arcs in the opposite direction.  That is, from roots toward leaves.
//
// See FromList.TransposeRoots for a version that also accumulates and returns
// information about the roots.
func (f *FromList) Transpose() Directed {
	g := make(AdjacencyList, len(f.Paths))
	for n, p := range f.Paths {
		if p.From == -1 {
			continue
		}
		g[p.From] = append(g[p.From], NI(n))
	}
	return Directed{g}
}

// TransposeLabeled constructs the directed labeled graph corresponding
// to FromList f but with arcs in the opposite direction.  That is, from
// roots toward leaves.
//
// The argument labels can be nil.  In this case labels are generated matching
// the path indexes.  This corresponds to the "to", or child node.
//
// If labels is non-nil, it must be the same length as f.Paths and is used
// to look up label numbers by the path index.
//
// See FromList.TransposeLabeledRoots for a version that also accumulates
// and returns information about the roots.
func (f *FromList) TransposeLabeled(labels []LI) DirectedLabeled {
	g := make(LabeledAdjacencyList, len(f.Paths))
	for n, p := range f.Paths {
		if p.From == -1 {
			continue
		}
		l := LI(n)
		if labels != nil {
			l = labels[n]
		}
		g[p.From] = append(g[p.From], Half{NI(n), l})
	}
	return DirectedLabeled{g}
}

// TransposeLabeledRoots constructs the labeled directed graph corresponding
// to FromList f but with arcs in the opposite direction.  That is, from
// roots toward leaves.
//
// TransposeLabeledRoots also returns a count of roots of the resulting forest
// and a bitmap of the roots.
//
// The argument labels can be nil.  In this case labels are generated matching
// the path indexes.  This corresponds to the "to", or child node.
//
// If labels is non-nil, it must be the same length as t.Paths and is used
// to look up label numbers by the path index.
//
// See FromList.TransposeLabeled for a simpler verstion that returns the
// forest only.
func (f *FromList) TransposeLabeledRoots(labels []LI) (forest DirectedLabeled, nRoots int, roots big.Int) {
	p := f.Paths
	nRoots = len(p)
	OneBits(&roots, len(p))
	g := make(LabeledAdjacencyList, len(p))
	for n, p := range f.Paths {
		if p.From == -1 {
			continue
		}
		l := LI(n)
		if labels != nil {
			l = labels[n]
		}
		g[p.From] = append(g[p.From], Half{NI(n), l})
		if roots.Bit(n) == 1 {
			roots.SetBit(&roots, n, 0)
			nRoots--
		}
	}
	return DirectedLabeled{g}, nRoots, roots
}

// TransposeRoots constructs the directed graph corresponding to FromList f
// but with arcs in the opposite direction.  That is, from roots toward leaves.
//
// TransposeRoots also returns a count of roots of the resulting forest and
// a bitmap of the roots.
//
// See FromList.Transpose for a simpler verstion that returns the forest only.
func (f *FromList) TransposeRoots() (forest Directed, nRoots int, roots big.Int) {
	p := f.Paths
	nRoots = len(p)
	OneBits(&roots, len(p))
	g := make(AdjacencyList, len(p))
	for n, e := range p {
		if e.From == -1 {
			continue
		}
		g[e.From] = append(g[e.From], NI(n))
		if roots.Bit(n) == 1 {
			roots.SetBit(&roots, n, 0)
			nRoots--
		}
	}
	return Directed{g}, nRoots, roots
}

// Undirected contructs the undirected graph corresponding to the FromList.
func (f *FromList) Undirected() Undirected {
	g := make(AdjacencyList, len(f.Paths))
	for n, p := range f.Paths {
		if p.From == -1 {
			continue
		}
		g[n] = append(g[n], p.From)
		if p.From != NI(n) {
			g[p.From] = append(g[p.From], NI(n))
		}
	}
	return Undirected{g}
}

// UndirectedLabeled contructs the corresponding undirected graph with
// edge labels.
//
// The argument labels can be nil.  In this case labels are generated matching
// the path indexes.  This corresponds to the to, or child node.
//
// If labels is non-nil, it must be the same length as t.Paths and is used
// to look up label numbers by the path index.
func (f *FromList) UndirectedLabeled(labels []LI) UndirectedLabeled {
	g := make(LabeledAdjacencyList, len(f.Paths))
	for n, p := range f.Paths {
		if p.From == -1 {
			continue
		}
		l := LI(n)
		if labels != nil {
			l = labels[n]
		}
		g[n] = append(g[n], Half{To: p.From, Label: l})
		if p.From != NI(n) {
			g[p.From] = append(g[p.From], Half{NI(n), l})
		}
	}
	return UndirectedLabeled{g}
}
