// Copyright 2014 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package graph

import (
	"math/big"

)

// cg_dir_al.go is code generated from cg_dir.go by directive in graph.go.
// Editing cg_dir.go is okay.
// DO NOT EDIT cg_dir_al.go.

// Balanced returns true if for every node in g, in-degree equals out-degree.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g DirectedLabeled) Balanced() bool {
    for n, in := range g.InDegree() {
        if in != len(g.LabeledAdjacencyList[n]) {
            return false
        }
    }
    return true
}

// Cyclic, for directed graphs, determines if g contains cycles.
//
// Cyclic returns true if g contains at least one cycle.
// Cyclic returns false if g is acyclic.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g DirectedLabeled) Cyclic() (cyclic bool, fr NI, to Half) {
    fr, to.To = -1, -1
    var temp, perm big.Int
    var df func(NI)
    df = func(n NI) {
        switch {
        case temp.Bit(int(n)) == 1:
            cyclic = true
            return
        case perm.Bit(int(n)) == 1:
            return
        }
        temp.SetBit(&temp, int(n), 1)
        for _, nb := range g.LabeledAdjacencyList[n] {
            df(nb.To)
            if cyclic {
                if fr < 0 {
                    fr, to = n, nb
                }
                return
            }
        }
        temp.SetBit(&temp, int(n), 0)
        perm.SetBit(&perm, int(n), 1)
    }
    for n := range g.LabeledAdjacencyList {
        if perm.Bit(n) == 1 {
            continue
        }
        if df(NI(n)); cyclic { // short circuit as soon as a cycle is found
            break
        }
    }
    return
}

// FromList transposes a graph into a FromList, encoding a tree or forest.
//
// Returns FromList, -1 as long as g represents a tree.
//
// If any node n in g is the to-node of more than one node, then g does not
// represent a tree and the method returns (nil, n).
//
// There are equivalent labeled and unlabeled versions of this method.
func (g DirectedLabeled) FromList() (*FromList, NI) {
    // init paths
    paths := make([]PathEnd, len(g.LabeledAdjacencyList))
    for i := range paths {
        paths[i].From = -1
    }
    // init leaves
    var leaves big.Int
    OneBits(&leaves, len(g.LabeledAdjacencyList))
    // iterate over arcs, setting from pointers and and marking non-leaves.
    for fr, to := range g.LabeledAdjacencyList {
        for _, to := range to {
            if paths[to.To].From >= 0 {
                return nil, to.To
            }
            paths[to.To].From = NI(fr)
            leaves.SetBit(&leaves, fr, 0)
        }
    }
    // f to set path lengths
    var leng func(NI) int
    leng = func(n NI) int {
        if l := paths[n].Len; l > 0 {
            return l
        }
        fr := paths[n].From
        if fr < 0 {
            paths[n].Len = 1
            return 1
        }
        l := 1 + leng(fr)
        paths[n].Len = l
        return l
    }
    // for each leaf, trace path to set path length, accumulate max
    maxLen := 0
    for i := range paths {
        if leaves.Bit(i) == 1 {
            if l := leng(NI(i)); l > maxLen {
                maxLen = l
            }
        }
    }
    return &FromList{paths, leaves, maxLen}, -1
}

// InDegree computes the in-degree of each node in g
//
// There are equivalent labeled and unlabeled versions of this method.
func (g DirectedLabeled) InDegree() []int {
    ind := make([]int, len(g.LabeledAdjacencyList))
    for _, nbs := range g.LabeledAdjacencyList {
        for _, nb := range nbs {
            ind[nb.To]++
        }
    }
    return ind
}

// IsTree identifies trees in directed graphs.
//
// IsTree returns true if the subgraph reachable from
// root is a tree.  It does not validate that the entire graph is a tree.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g DirectedLabeled) IsTree(root NI) bool {
    var v big.Int
    var df func(NI) bool
    df = func(n NI) bool {
        if v.Bit(int(n)) == 1 {
            return false
        }
        v.SetBit(&v, int(n), 1)
        for _, to := range g.LabeledAdjacencyList[n] {
            if !df(to.To) {
                return false
            }
        }
        return true
    }
    return df(root)
}

// Tarjan identifies strongly connected components in a directed graph using
// Tarjan's algorithm.
//
// The method calls the emit argument for each component identified.  Each
// component is a list of nodes.  A property of the algorithm is that
// components are emitted in reverse topological order of the condensation.
// (See https://en.wikipedia.org/wiki/Strongly_connected_component#Definitions
// for description of condensation.)
//
// There are equivalent labeled and unlabeled versions of this method.
//
// See also TarjanForward and TarjanCondensation.
func (g DirectedLabeled) Tarjan(emit func([]NI) bool) {
    // See "Depth-first search and linear graph algorithms", Robert Tarjan,
    // SIAM J. Comput. Vol. 1, No. 2, June 1972.
    //
    // Implementation here from Wikipedia pseudocode,
    // http://en.wikipedia.org/w/index.php?title=Tarjan%27s_strongly_connected_components_algorithm&direction=prev&oldid=647184742
    var indexed, stacked big.Int
    index := make([]int, len(g.LabeledAdjacencyList))
    lowlink := make([]int, len(g.LabeledAdjacencyList))
    x := 0
    var S []NI
    var sc func(NI) bool
    sc = func(n NI) bool {
        index[n] = x
        indexed.SetBit(&indexed, int(n), 1)
        lowlink[n] = x
        x++
        S = append(S, n)
        stacked.SetBit(&stacked, int(n), 1)
        for _, nb := range g.LabeledAdjacencyList[n] {
            if indexed.Bit(int(nb.To)) == 0 {
                if !sc(nb.To) {
                    return false
                }
                if lowlink[nb.To] < lowlink[n] {
                    lowlink[n] = lowlink[nb.To]
                }
            } else if stacked.Bit(int(nb.To)) == 1 {
                if index[nb.To] < lowlink[n] {
                    lowlink[n] = index[nb.To]
                }
            }
        }
        if lowlink[n] == index[n] {
            var c []NI
            for {
                last := len(S) - 1
                w := S[last]
                S = S[:last]
                stacked.SetBit(&stacked, int(w), 0)
                c = append(c, w)
                if w == n {
                    if !emit(c) {
                        return false
                    }
                    break
                }
            }
        }
        return true
    }
    for n := range g.LabeledAdjacencyList {
        if indexed.Bit(n) == 0 && !sc(NI(n)) {
            return
        }
    }
}

// TarjanForward returns strongly connected components. 
// 
// It returns components in the reverse order of Tarjan, for situations
// where a forward topological ordering is easier.
func (g DirectedLabeled) TarjanForward() [][]NI {
    var r [][]NI
    g.Tarjan(func(c []NI) bool {
        r = append(r, c)
        return true
    })
    scc := make([][]NI, len(r))
    last := len(r) - 1
    for i, ci := range r {
        scc[last-i] = ci
    }
    return scc
}   
    
// TarjanCondensation returns strongly connected components and their
// condensation graph.
//  
// Components are ordered in a forward topological ordering.
func (g DirectedLabeled) TarjanCondensation() (scc [][]NI, cd AdjacencyList) { 
    scc = g.TarjanForward()
    cd = make(AdjacencyList, len(scc)) // return value
    cond := make([]NI, len(g.LabeledAdjacencyList))         // mapping from g node to cd node
    for cn := len(scc) - 1; cn >= 0; cn-- {
        c := scc[cn] 
        for _, n := range c {
            cond[n] = NI(cn) // map g node to cd node
        }
        var tos []NI  // list of 'to' nodes
        var m big.Int // tos map
        m.SetBit(&m, cn, 1)
        for _, n := range c {
            for _, to := range g.LabeledAdjacencyList[n] {
                if ct := cond[to.To]; m.Bit(int(ct)) == 0 {
                    m.SetBit(&m, int(ct), 1)
                    tos = append(tos, ct)
                }
            }
        }
        cd[cn] = tos
    }
    return
}

// Topological, for directed acyclic graphs, computes a topological sort of g.
//
// For an acyclic graph, return value ordering is a permutation of node numbers
// in topologically sorted order and cycle will be nil.  If the graph is found
// to be cyclic, ordering will be nil and cycle will be the path of a found
// cycle.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g DirectedLabeled) Topological() (ordering, cycle []NI) {
    ordering = make([]NI, len(g.LabeledAdjacencyList))
    i := len(ordering)
    var temp, perm big.Int
    var cycleFound bool
    var cycleStart NI
    var df func(NI)
    df = func(n NI) {
        switch {
        case temp.Bit(int(n)) == 1:
            cycleFound = true
            cycleStart = n
            return
        case perm.Bit(int(n)) == 1:
            return
        }
        temp.SetBit(&temp, int(n), 1)
        for _, nb := range g.LabeledAdjacencyList[n] {
            df(nb.To)
            if cycleFound {
                if cycleStart >= 0 {
                    // a little hack: orderng won't be needed so repurpose the
                    // slice as cycle.  this is read out in reverse order
                    // as the recursion unwinds.
                    x := len(ordering) - 1 - len(cycle)
                    ordering[x] = n
                    cycle = ordering[x:]
                    if n == cycleStart {
                        cycleStart = -1
                    }
                }
                return
            }
        }
        temp.SetBit(&temp, int(n), 0)
        perm.SetBit(&perm, int(n), 1)
        i--
        ordering[i] = n
    }
    for n := range g.LabeledAdjacencyList {
        if perm.Bit(n) == 1 {
            continue
        }
        df(NI(n))
        if cycleFound {
            return nil, cycle
        }
    }
    return ordering, nil
}
// TopologicalKahn, for directed acyclic graphs, computes a topological sort of g.
//
// For an acyclic graph, return value ordering is a permutation of node numbers
// in topologically sorted order and cycle will be nil.  If the graph is found
// to be cyclic, ordering will be nil and cycle will be the path of a found
// cycle.
//
// This function is based on the algorithm by Arthur Kahn and requires the
// transpose of g be passed as the argument.
//
// There are equivalent labeled and unlabeled versions of this method.
func (g DirectedLabeled) TopologicalKahn(tr Directed) (ordering, cycle []NI) {
    // code follows Wikipedia pseudocode.
    var L, S []NI
    // rem for "remaining edges," this function makes a local copy of the
    // in-degrees and consumes that instead of consuming an input.
    rem := make([]int, len(g.LabeledAdjacencyList))
    for n, fr := range tr.AdjacencyList {
        if len(fr) == 0 {
            // accumulate "set of all nodes with no incoming edges"
            S = append(S, NI(n))
        } else {
            // initialize rem from in-degree
            rem[n] = len(fr)
        }
    }
    for len(S) > 0 {
        last := len(S) - 1 // "remove a node n from S"
        n := S[last]
        S = S[:last]
        L = append(L, n) // "add n to tail of L"
        for _, m := range g.LabeledAdjacencyList[n] {
            // WP pseudo code reads "for each node m..." but it means for each
            // node m *remaining in the graph.*  We consume rem rather than
            // the graph, so "remaining in the graph" for us means rem[m] > 0.
            if rem[m.To] > 0 {
                rem[m.To]--         // "remove edge from the graph"
                if rem[m.To] == 0 { // if "m has no other incoming edges"
                    S = append(S, m.To) // "insert m into S"
                }
            }
        }
    }
    // "If graph has edges," for us means a value in rem is > 0.
    for c, in := range rem {
        if in > 0 {
            // recover cyclic nodes
            for _, nb := range g.LabeledAdjacencyList[c] {
                if rem[nb.To] > 0 {
                    cycle = append(cycle, NI(c))
                    break
                }
            }
        }
    }
    if len(cycle) > 0 {
        return nil, cycle
    }
    return L, nil
}
