package graph_test
/*
import (
	"fmt"

	"github.com/soniakeys/graph"
)

func ExamplePrim_Scan() {
	g := graph.WeightedAdjacencyList{
		0: {{1, 3}, {2, 5}},
		1: {{0, 3}, {2, 4}},
		2: {{0, 5}, {1, 4}},
		3: {{4, 2}},
		4: {{3, 2}},
	}
	uw := g.Unweighted()
	ud, _, _ := uw.Undirected()
	fmt.Println("Undirected:", ud)
	rep, nNodes := uw.ConnectedComponents()
	fmt.Println("Connected component representatives:", rep)
	p := graph.NewPrim(g)
	nSpanned := p.Span(rep[0])
	fmt.Println("Spanned:", nSpanned == nNodes[0])
	fmt.Println("Node  From  Weight")
	for n, pe := range p.Result.Paths {
//		if pe.Len > 0 {
			fmt.Printf("%4d %5d %7.1f %#v\n",
				n, pe.From.From, pe.From.ArcWeight, pe)
//		}
	}
	for n, nbs := range g {
		fmt.Println(n, nbs)
	}
	// Output:
	// Undirected: true
	// Connected component representatives: [0 3]
	// :P
}
*/
