// Logical satisfiability
//
// Currently just a 2-sat function.
package sat

import "github.com/soniakeys/graph"

type VarSense int

const (
	Bare    VarSense = 0
	Negated VarSense = 1
)

type CNF2 struct {
	X, Y   string
	XS, YS VarSense
}

func TwoSat(exp []CNF2) map[string]bool {
	// scan clauses for unique variable names
	varN := map[string]graph.NI{}
	for _, c := range exp {
		varN[c.X] = 0
		varN[c.Y] = 0
	}
	// map numbers
	vars := make([]string, len(varN))
	var n graph.NI
	for v := range varN {
		vars[n] = v
		varN[v] = n
		n++
	}
	// build graph
	g := make(graph.AdjacencyList, 2*len(vars))
	for _, c := range exp {
		// var z bare is node 2z, negated is 2z+1
		nx := varN[c.X]<<1 | graph.NI(c.XS^1)
		ny := varN[c.Y]<<1 | graph.NI(c.YS)
		g[nx] = append(g[nx], ny) // ^X -> Y
		ny ^= 1
		g[ny] = append(g[ny], nx^1) // ^Y -> X
	}
	// test strongly connected components
	p := make([]int, len(vars))
	m := map[string]bool{}
	// (Tarjan returns components in reverse topological order, which is
	// the needed processing order here.)
	g.Tarjan(func(c []graph.NI) bool {
		for i := range p {
			p[i] = 0
		}
		for _, nn := range c {
			vn := nn >> 1          // recover variable number from node number
			vs := VarSense(nn & 1) // recover variable sense from node number
			p[vn] |= 1 << uint(vs)
			if p[vn] == 3 { // bare and negated in same scc
				m = nil // nil return means no solution
				return false
			}
			str := vars[vn]
			if _, ok := m[str]; !ok {
				// var not set yet,
				// set value bare variable to true, negated to false
				m[str] = vs == Bare
			}
		}
		return true
	})
	return m
}
