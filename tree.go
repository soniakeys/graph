package graph

import "math/rand"

// RandomUTree, random unrooted binary tree, returned as a parent list.
//
// An unrooted binary tree, as used in phylogeny, can be represented as
// an undirected graph where leaves are nodes with exactly one neighbor
// and internal nodes have exactly 3.  (It's binary because after following
// one edge into an internal node there are two other ways out.)
//
// RandomUTree generates a random unrooted binary tree but returns the
// result as a parent list, a more compact representation than the full
// undirected graph adjacency list.
//
// The unrooted binary tree for n leaves has n-2 internal nodes. The
// root node is not stored so len(parentList) == nLeaves+nLeaves-3.
// The leaves correspond to nodes 0:nLeaves and the unrepresented root
// corresponds to len(parentList).
//
// The unrepresented root is only a root for the purpose of the parent list
// representation.  It has no significance in the corresponding unrooted
// binary tree.
//
// Example code:
//   n := 5
//   pl := graph.RandomUTree(n)
//   fmt.Println(n+n-2, "nodes,")
//   fmt.Println(len(pl), "in parent list:")
//   for i, p := range pl {
//       fmt.Printf("%d: parent %d\n", i, p)
//   }
//
// Output:
//
//   8 nodes,
//   7 in parent list:
//   0: parent 7
//   1: parent 5
//   2: parent 6
//   3: parent 5
//   4: parent 6
//   5: parent 7
//   6: parent 7
//
// Diagram:
//
//          7
//        / | \
//       /  /  6
//      /  /   |\
//     /  5    | \
//    /  / \   |  \
//   0  1   3  2   4
//
// Equivalently:
//
//       1
//        \
//         5 - 3
//        /
//   0 - 7
//        \
//         6 - 2
//        /
//       4
//
func RandomUTree(nLeaves int) (parentList []int) {
	// allocate space for whole tree except root
	parentList = make([]int, nLeaves+nLeaves-3)
	// initial tree has three leaves and the internal root
	root := len(parentList)
	parentList[0] = root
	parentList[1] = root
	parentList[2] = root
	// now add remaining nodes.  on each iteration add one each internal
	// node and parent node.  pick random edge from parent list,
	// embed a new internal node in this edge, also connect a new leaf.
	// new edges are from new leaf to new internal node and from
	// new internal node to parent.
	for newLeaf := 3; newLeaf < nLeaves; newLeaf++ {
		i := nLeaves + newLeaf - 3     // new internal node
		l1 := rand.Intn(newLeaf*2 - 3) // (range is number of existing edges)
		if l1 >= newLeaf {
			l1 += nLeaves - newLeaf // skip to range of internal nodes
		}
		l2 := parentList[l1]
		parentList[l1] = i
		parentList[i] = l2      // new edge
		parentList[newLeaf] = i // new edge
	}
	return
}
