package main

//go:generate go run DijkstraAllPaths.go

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/soniakeys/graph"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ns := []int{16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384,
		32768, 65536, 131072, 262144, 524288, 1048576, 2097152, 4194304}
	rep := 5
	pts := make([]plotter.XYer, len(ns))
	for i, n := range ns {
		g, _, wt, err := graph.LabeledEuclidean(n, n*10, 4, 100, r)
		if err != nil {
			log.Fatal(err)
		}
		d := graph.NewDijkstra(g.LabeledAdjacencyList,
			func(l graph.LI) float64 { return wt[l] })
		xys := make(plotter.XYs, rep)
		pts[i] = xys
		for j := 0; j < rep; j++ {
			var start graph.NI
			for {
				start = graph.NI(r.Intn(n))
				if len(g.LabeledAdjacencyList[start]) > 0 {
					break
				}
			}
			f := func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					d.AllPaths(start)
					d.Reset()
				}
			}
			b := testing.Benchmark(f)
			fmt.Printf("n=%4d, run %d: %v\n", n, j, b)
			xys[j].X = float64(n)
			xys[j].Y = float64(b.NsPerOp()) * .001
		}
	}
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.Title.Text = "Dijkstra.AllPaths, Directed Graph\nArc Size MA = 10N"
	p.X.Label.Text = "Graph Order, N"
	p.Y.Label.Text = "µsec"
	nln := plotter.NewFunction(func(n float64) float64 {
		ln := math.Log(n)
		return (200*n + 40*n*ln + 6*n*ln*ln) * .001
	})
	nln.Color = color.RGBA{B: 127, A: 255}
	p.Add(nln)
	mmm, err := plotutil.NewErrorPoints(meanMinMax, pts...)
	if err != nil {
		log.Fatal(err)
	}
	if err = plotutil.AddYErrorBars(p, mmm); err != nil {
		log.Fatal(err)
	}
	p.X.Scale = plot.LogScale{}
	p.Y.Scale = plot.LogScale{}
	p.X.Tick.Marker = plot.LogTicks{}
	p.Y.Tick.Marker = plot.LogTicks{}
	err = p.Save(4*vg.Inch, 4*vg.Inch, "DijkstraAllPaths.svg")
	if err != nil {
		log.Fatal(err)
	}
}

func meanMinMax(vs []float64) (mean, lowerr, higherr float64) {
	low := math.Inf(1)
	high := math.Inf(-1)
	for _, v := range vs {
		mean += v
		low = math.Min(v, low)
		high = math.Max(v, high)
	}
	mean /= float64(len(vs))
	return mean, mean - low, high - mean
}
