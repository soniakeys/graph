package ed

import ()

type BreadthFirst struct {
	Graph  AdjacencyList
	Start  int
	Paths  []PathEnd
	MaxLen int
}

type PathEnd struct {
	From int
	Len  int
}

func NewBreadthFirst(g AdjacencyList) *BreadthFirst {
	return &BreadthFirst{
		Graph: g,
		Paths: make([]PathEnd, len(g)),
	}
}

func (b *BreadthFirst) Path(start, end int) []int {
	b.search(start, end)
	return b.PathTo(end)
}

func (b *BreadthFirst) AllPaths(start int) int {
	return b.search(start, -1)
}

func (b *BreadthFirst) search(start, end int) int {
	for n := range b.Paths {
		b.Paths[n] = PathEnd{From: -1, Len: 0}
	}
	b.Start = start
	b.MaxLen = 1
	b.Paths[start].Len = b.MaxLen
	if start == end {
		return -1
	}
	nReached := 1 // accumulated for a return value
	// the frontier consists of nodes all at the same level
	frontier := []int{start}
	for {
		b.MaxLen++
		var next []int
		for _, n := range frontier {
			for _, nb := range b.Graph[n] {
				if b.Paths[nb].Len == 0 {
					b.Paths[nb] = PathEnd{From: n, Len: b.MaxLen}
					if nb == end {
						return -1
					}
					next = append(next, nb)
					nReached++
				}
			}
		}
		if len(next) == 0 {
			break
		}
		frontier = next
	}
	b.MaxLen--
	return nReached
}

func (b *BreadthFirst) PathTo(end int) []int {
	n := b.Paths[end].Len
	if n == 0 {
		return nil
	}
	p := make([]int, n)
	for {
		n--
		p[n] = end
		if n == 0 {
			return p
		}
		end = b.Paths[end].From
	}
}

/*
type BreadthFirst2 struct {
	Out, In [][]int
	M int
	Level []int
	FromTree []FromNode
}

func NewBF(out, in [][]int, m int) *BreadthFirst {



type func Visitor(int) bool

//func BreadthFirst(out, in [][]int, n, m int) (D []int) {
func (bf *BreadthFirst) Search(start int, v Visitor) int {
	// source defined by the problem to be vertex 1
	source := 1
	d0 := make([]int, len(out)) // 0 element unused.  return value is d0[1:]
	for i := range d0 {
		d0[i] = -1
	}
	lNum := 0 // level number
	d0[source] = lNum

	frontier := []int{source} // verices all at the same level
	mf := len(out[source])    // number of arcs leading out from frontier
	ctb := m / 10             // threshold to change from top-down to bottom-up
	k14 := 14 * m / n         // 14 * mean degree
	cbt := n / k14            // threshold to change from bottom-up to top-down

	for {
		lNum++
		frontier, mf = topDown(lNum, out, frontier, d0)
		if len(frontier) == 0 {
			break
		}
		if mf > ctb {
			// switch to bottom up!
		} else {
			// stick with top down
			continue
		}
		// convert
		nf := 0 // number of vertices on the frontier
		fBits := &big.Int{}
		for _, v := range frontier {
			fBits.SetBit(fBits, v, 1)
			nf++
		}
	bottomUpLoop:
		lNum++
		fBits, nf = bottomUp(lNum, in, fBits, d0)
		if fBits.BitLen() == 0 {
			break
		}
		if nf < cbt {
			// switch back to top down!
		} else {
			// stick with bottom up
			goto bottomUpLoop
		}
		// convert
		mf = 0
		frontier = frontier[:0]
		for v := 1; v <= n; v++ {
			if fBits.Bit(v) == 1 {
				frontier = append(frontier, v)
				mf += len(out[v])
			}
		}
	}
	return d0[1:] // drop unused 0 element
}

func bottomUp(lNum int, in [][]int, frontier *big.Int, d0 []int) (next *big.Int, nNext int) {
	next = &big.Int{}
	for v := 1; v < len(in); v++ {
		if d0[v] < 0 {
			for _, nb := range in[v] {
				if frontier.Bit(nb) == 1 {
					d0[v] = lNum
					next.SetBit(next, v, 1)
					nNext++
					break
				}
			}
		}
	}
	return
}*/
