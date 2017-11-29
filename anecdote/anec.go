// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

// Anecdotal timings
package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/soniakeys/graph"
	"github.com/soniakeys/graph/alt"
)

func main() {
	fmt.Println("Anecdotal timings")
	fmt.Println(runtime.GOOS, runtime.GOARCH)
	random()
	prop()
	trav()
	allpairs()
	sssp()
	shortestone()
	eulerian()
	parallel()
	cycles()
}

func h(n int) string {
	switch {
	case n < 1000:
		return fmt.Sprint(n)
	case n < 1e4:
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	case n < 1e6:
		return fmt.Sprint(n/1000, "K")
	case n < 1e7:
		return fmt.Sprintf("%.1fM", float64(n)/1e6)
	default:
		return fmt.Sprint(n/1e6, "M")
	}
}

func random() {
	fmt.Println("\nRandom graph generation")
	fmt.Println("(arcs/edges lists arcs for directed graphs, edges for undirected)")
	fmt.Println("Graph type                Nodes  Arcs/edges           Time")
	for _, tc := range []func() (string, int, int){
		ChungLuSmall, ChungLuLarge,
		EucSmall, EucLarge,
		GeoSmall, GeoLarge,
		GnpDSmall, GnpDLarge, GnpUSmall, GnpULarge,
		GnmDSmall, GnmDLarge, GnmUSmall, GnmULarge,
		Gnm3USmall, Gnm3ULarge,
		KronDSmall, KronDLarge, KronUSmall, KronULarge,
	} {
		t := time.Now()
		g, n, a := tc()
		d := time.Now().Sub(t)
		fmt.Printf("%-25s %5s %5s %20s\n", g, h(n), h(a), d)
	}
}

func prop() {
	fmt.Println("\nProperties")
	fmt.Println("Method                 Graph                          Time")
	for _, tc := range []func() (string, string){
		CCSmall, CCLarge,
		SCCKosarajuSmall, SCCKosarajuLarge,
		SCCTarjanSmall, SCCTarjanLarge,
		SCCPathSmall, SCCPathLarge,
		SCCPearceSmall, SCCPearceLarge,
		SCCEucSmall, SCCEucLarge,
	} {
		t := time.Now()
		m, g := tc()
		d := time.Now().Sub(t)
		fmt.Printf("%-22s %-22s %12s\n", m, g, d)
	}
}

func trav() {
	fmt.Println("\nTraversal")
	fmt.Println("Method                 Graph                                          Time")
	for _, tc := range []func() (string, string){
		DFSmall, DFLarge, BFSmall, BFLarge,
		AltBFSmall, AltBFLarge,
	} {
		t := time.Now()
		m, g := tc()
		d := time.Now().Sub(t)
		fmt.Printf("%-22s %-38s %12s\n", m, g, d)
	}

}

var geoSmallEnd graph.NI
var geoLargeEnd graph.NI
var geoSmallHeuristic func(graph.NI) float64
var geoLargeHeuristic func(graph.NI) float64

func allpairs() {
	fmt.Println("\nShortest path all pairs")
	fmt.Println("Method                 Graph                                          Time")
	for _, tc := range []func() (string, string){
		FloydEuc, FloydGeo,
	} {
		t := time.Now()
		m, g := tc()
		d := time.Now().Sub(t)
		fmt.Printf("%-22s %-38s %12s\n", m, g, d)
	}
}

func sssp() {
	fmt.Println("\nSingle source shortest path")
	fmt.Println("Method                 Graph                                          Time")
	for _, tc := range []func() (string, string){
		BellmanSmall,
		DijkstraAllSmall, DijkstraAllLarge,
	} {
		t := time.Now()
		m, g := tc()
		d := time.Now().Sub(t)
		fmt.Printf("%-22s %-38s %12s\n", m, g, d)
	}

}

func shortestone() {
	fmt.Println("\nSingle shortest path")
	fmt.Println("Method                 Graph                                          Time")
	// pick end nodes about .7 distant from node 0
	p1 := geoSmallPos[0]
	nearestSmall, nearestLarge := 2., 2.
	var geoSmallEndPos struct{ X, Y float64 }
	var geoLargeEndPos struct{ X, Y float64 }
	for l := 1; l < len(geoSmallPos); l++ {
		p2 := geoSmallPos[1]
		d := math.Abs(.7 - math.Hypot(p2.X-p1.X, p2.Y-p1.Y))
		if d < nearestSmall {
			nearestSmall = d
			geoSmallEnd = graph.NI(l)
			geoSmallEndPos = p2
		}
		p2 = geoLargePos[1]
		d = math.Abs(.7 - math.Hypot(p2.X-p1.X, p2.Y-p1.Y))
		if d < nearestLarge {
			nearestLarge = d
			geoLargeEnd = graph.NI(l)
			geoLargeEndPos = p2
		}
	}
	// and define heuristics for AStar
	geoSmallHeuristic = func(from graph.NI) float64 {
		p := &geoSmallPos[from]
		return math.Hypot(geoSmallEndPos.X-p.X, geoSmallEndPos.Y-p.Y)
	}
	geoLargeHeuristic = func(from graph.NI) float64 {
		p := &geoLargePos[from]
		return math.Hypot(geoLargeEndPos.X-p.X, geoLargeEndPos.Y-p.Y)
	}
	for _, tc := range []func() (string, string){
		Dijkstra1Small, Dijkstra1Large,
		AStarASmall, AStarALarge, AStarMSmall, AStarMLarge,
	} {
		t := time.Now()
		m, g := tc()
		d := time.Now().Sub(t)
		fmt.Printf("%-22s %-38s %12s\n", m, g, d)
	}
}

func eulerian() {
	fmt.Println("\nEulerian cycles")
	dirEuSmall := dirEu(1e3)
	dirEuLarge := dirEu(1e5)
	uEuSmall := uEu(1e3)
	uEuLarge := uEu(1e5)
	for _, r := range []euResult{
		dirEuTest(dirEuSmall),  // non-destructive first
		dirEuDTest(dirEuSmall), // this consumes dirEuSmall
		dirEuTest(dirEuLarge),
		dirEuDTest(dirEuLarge),
		uEuTest(uEuSmall),
		uEuDTest(uEuSmall),
		uEuTest(uEuLarge),
		uEuDTest(uEuLarge),
	} {
		fmt.Printf("%-22s %-38s %12s\n", r.method, r.tag, r.d)
	}
}

func parallel() {
	fmt.Println("\nParallel arc tests")
	fmt.Println("Graph       Sort           Map  Best")
	for i, g := range []interface{}{
		chungLuSmall,
		chungLuLarge,
		eucLarge,
		gnpDSmall,
		gnpDLarge,
		gnpUSmall,
		gnpULarge,
		gnmDSmall,
		gnmDLarge,
		gnmUSmall,
		gnmULarge,
		kronDSmall,
		kronDLarge,
		kronUSmall,
		kronULarge,
	} {
		var tm, ts time.Duration
		switch h := g.(type) {
		case graph.Directed:
			t := time.Now()
			alt.AnyParallelMap(h.AdjacencyList)
			tm = time.Now().Sub(t)
			t = time.Now()
			h.AnyParallel()
			ts = time.Now().Sub(t)
		case graph.Undirected:
			t := time.Now()
			alt.AnyParallelMap(h.AdjacencyList)
			tm = time.Now().Sub(t)
			t = time.Now()
			h.AnyParallel()
			ts = time.Now().Sub(t)
		}
		best := "Map"
		if ts < tm {
			best = "Sort"
		}
		fmt.Printf("%2d  %12s  %12s  %s\n", i, ts, tm, best)
	}
}

func cycles() {
	fmt.Println("\nDirectedCycles")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	g := graph.GnmDirected(18, 90, r)
	g.ShuffleArcLists(r)
	fmt.Println("Method                 Graph                                          Time")
	for _, tc := range []struct {
		f func(graph.Directed, func([]graph.NI) bool)
		s string
	}{
		{graph.Directed.Cycles, "Johnson"},
		{alt.TarjanCycles, "Tarjan"},
	} {
		t := time.Now()
		tc.f(g, func(emit []graph.NI) bool { return true })
		d := time.Now().Sub(t)
		fmt.Printf("%-22s %-38s %12s\n", tc.s, "Gnm 18 90", d)
	}
}
