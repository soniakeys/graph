package graph

import (
	"testing"
)

func Test_zL(t *testing.T) {
	//     -1      2
	//  0----->1------>2--\
	//  ^     ^^\     ^|   \2
	// 3|  -4/ ||1   / |-1  v
	//  |   /  ||   /  |    3
	//  |  /   ||  /   |   /
	//  | /  -2|| /-2  |  /-1
	//  |/     \v/     v /
	//  4<------5<-----6<
	//      -1      1
	g := LabeledDirected{LabeledAdjacencyList{
		0: {{To: 1, Label: -1}},
		1: {{To: 2, Label: 2}, {To: 5, Label: 1}},
		2: {{To: 3, Label: 2}, {To: 6, Label: -1}},
		3: {{To: 6, Label: -1}},
		4: {{To: 0, Label: 3}, {To: 1, Label: -4}},
		5: {{To: 1, Label: -2}, {To: 2, Label: -2}, {To: 4, Label: -1}},
		6: {{To: 5, Label: 1}},
	}}
	a := g.LabeledAdjacencyList
	tr, _ := g.UnlabeledTranspose()
	F := path{4, []Half{{To: 0, Label: 3}}}
	R := map[arc]bool{{5, Half{To: 2, Label: -2}}: true}
	GFR := gfr(a, tr.AdjacencyList, F, R)
	t.Log("GFR:")
	for fr, to := range GFR {
		t.Log(fr, " ", to)
	}
	wf := func(l LI) float64 { return float64(l) }
	t.Log("wPath ", wPath(F, wf))
	zL(GFR, F, wf, wPath(F, wf))
}
