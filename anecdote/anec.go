// Copyright 2017 Sonia Keys
// License MIT: http://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println("Anecdotal timings")
	fmt.Println(runtime.GOOS, runtime.GOARCH)
	fmt.Println("\nRandom graph generation")
	fmt.Println("(arcs/edges lists arcs for directed graphs, edges for undirected)")
	fmt.Println("Graph type           nodes  arcs/edges           time")
	for _, tc := range []func() (string, int, int){
		ChungLuSmall, ChungLuLarge,
		/*
		EucSmall, EucLarge,
			GeoSmall, GeoLarge,
			GnpDSmall, GnpDLarge, GnpUSmall, GnpULarge,
			GnmDSmall, GnmDLarge, GnmUSmall, GnmULarge, // Gnm3...
			KronDSmall, KronDLarge, KronUSmall, KronULarge,
		*/
	} {
		t := time.Now()
		g, n, a := tc()
		d := time.Now().Sub(t)
		fmt.Printf("%-20s %5s %5s %20s\n", g, h(n), h(a), d)
	}

	fmt.Println("\nProperties")
	fmt.Println("Method                Graph                      time")
	for _, tc := range []func() (string, string){
		CCSmall, CCLarge,
	} {
		t := time.Now()
		m, g := tc()
		d := time.Now().Sub(t)
		fmt.Printf("%-21s %-16s %14s\n", m, g, d)
	}

	fmt.Println("\nTraversal")
	fmt.Println("Method                Graph                                 time")
	for _, tc := range []func() (string, string){
		DFSmall, DFLarge, BFSmall, BFLarge,
	} {
		t := time.Now()
		m, g := tc()
		d := time.Now().Sub(t)
		fmt.Printf("%-21s %-33s %14s\n", m, g, d)
	}

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
