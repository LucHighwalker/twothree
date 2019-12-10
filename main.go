package main

import (
	"fmt"
	"math/rand"
)

//---------------------------------------------Helpers-------------------------//

// dataLen(data) counts the amount of non-nil elements in a node's data array
func dataLen(data [2]*int) int {
	count := 0
	for _, d := range data {
		if d != nil {
			count++
		} else {
			break
		}
	}
	return count
}

// childLen(children) counts the amount of non-nil elements in a node's children array
func childLen(children [4]*node) int {
	count := 0
	for _, c := range children {
		if c != nil {
			count++
		} else {
			break
		}
	}
	return count
}

func randomNumbers(count int, max int) []int {
	seed := true
	if max < 1 {
		seed = false
	}

	rands := []int{}
	for i := 0; i < count; i++ {
		if seed {
			rands = append(rands, rand.Intn(max))
		} else {
			rands = append(rands, rand.Int())
		}
	}
	return rands
}

//---------------------------------------------Tree----------------------------//

// Tree struct for a two-three tree
type Tree struct {
	root *node
	Size int
}

// Tree.Insert(data) inserts the given data into the tree.
func (t *Tree) Insert(data int) {
	if n, exists := t.FindNode(data); exists == false {
		if n != nil {
			n.insert(data)
			t.refreshRoot()
		} else {
			t.root = &node{}
			t.root.insert(data)
		}
		t.Size++
	}
}

func (t *Tree) InsertMany(data []int) {
	for _, d := range data {
		t.Insert(d)
	}
}

// Tree.Contains(data) initiates the recursive node.contains(data) method.
func (t *Tree) FindNode(data int) (*node, bool) {
	if t.root != nil {
		return t.root.findNode(data)
	}
	return nil, false
}

// Tree.refreshRoot() makes sure that the root is the top level node.
// If it isn't the top level node, the root is assigned to the top level node
func (t *Tree) refreshRoot() {
	if t.root.parent == nil {
		return
	} else {
		t.root = t.root.parent
	}
}

func (t *Tree) PrintTree() {
	if t.root == nil {
		fmt.Println("tree is empty")
		return
	} else {
		t.root.visit(stringifyNode)
	}
}

//---------------------------------------------Node----------------------------//

// node struct for a two-three tree
type node struct {
	data     [2]*int
	children [4]*node
	parent   *node
}

// node.insert(data) inserts the data into the node.
func (n *node) insert(data int) {
	switch dataLen(n.data) {
	case 0:
		// if node has no data
		n.data[0] = &data
		break

	case 1:
		stored := *n.data[0]
		if data < stored {
			// move existing data to the right
			n.data[0] = &data
			n.data[1] = &stored
		} else {
			n.data[1] = &data
		}
		break

	case 2:
		// if node is already a three node
		left := *n.data[0]
		right := *n.data[1]
		if left > data {
			n.data[0] = &data
			n.toParent(left)
		} else if right < data {
			n.data[1] = &data
			n.toParent(right)
		} else {
			n.toParent(data)
		}
		break
	}
}

// node.contains(data) recursive method returns bool if data exists.
func (n *node) findNode(data int) (*node, bool) {
	hasChildren := n.children[0] != nil

	if data == *n.data[0] {
		return n, true
	} else if hasChildren && data < *n.data[0] {
		return n.children[0].findNode(data)
	}

	// switch to handle two and three nodes slighlty differently
	switch dataLen(n.data) {
	case 1:
		if hasChildren {
			return n.children[1].findNode(data)
		}
		break

	case 2:
		if data == *n.data[1] {
			return n, true
		} else if hasChildren {
			if data < *n.data[1] {
				return n.children[1].findNode(data)
			} else {
				return n.children[2].findNode(data)
			}
		}
		break
	}
	return n, false
}

// node.toParent(data) pushes data to the node's parent
func (n *node) toParent(data int) {
	if n.parent == nil {
		n.parent = &node{children: [4]*node{n}}
	}
	n.split()
	n.parent.insert(data)
}

func (n *node) adopt(left, right *node) {
	n.children[0] = left
	n.children[1] = right
	left.parent = n
	right.parent = n
}

func (n *node) pushChildren(left, right *node, location int) {
	switch location {
	case 0:
		n.children[0] = left
		if n.children[2] != nil {
			n.children[3] = n.children[2]
		}
		if n.children[1] != nil {
			n.children[2] = n.children[1]
		}
		n.children[1] = right
		break

	case 1:
		n.children[1] = left
		if n.children[2] != nil {
			n.children[3] = n.children[2]
		}
		n.children[2] = right
		break

	case 2:
		n.children[2] = left
		n.children[3] = right
		break
	}
}

func (n *node) pushChildrenLeft(left *node, right *node) {
	n.children[0] = left
	if n.children[2] != nil {
		n.children[3] = n.children[2]
	}
	if n.children[1] != nil {
		n.children[2] = n.children[1]
	}
	n.children[1] = right
}

func (n *node) pushChildrenMid(left *node, right *node) {
	n.children[1] = left
	if n.children[2] != nil {
		n.children[3] = n.children[2]
	}
	n.children[2] = right
}

func (n *node) pushChildrenRight(left *node, right *node) {
	n.children[2] = left
	n.children[3] = right
}

func (n *node) removeChild(c *node) int {
	index := -1
	for i, child := range n.children {
		if child == c {
			index = i
			break
		}
	}
	n.children[index] = nil
	return index
}

// node.split() splits the node into 2 nodes
func (n *node) split() {
	left := &node{parent: n.parent}
	right := &node{parent: n.parent}

	left.insert(*n.data[0])
	right.insert(*n.data[1])

	if n.children[0] != nil {
		left.adopt(n.children[0], n.children[1])
		right.adopt(n.children[2], n.children[3])
	}

	index := n.parent.removeChild(n)
	n.parent.pushChildren(left, right, index)
}

func stringifyNode(n *node) {
	fmt.Println("----node----")
	fmt.Println(n.toString(true))
}

func (n *node) visit(f func(n *node)) {
	f(n)
	childCount := childLen(n.children)
	for i := 0; i < childCount; i++ {
		n.children[i].visit(f)
	}
}

// node.toString() returns a string representation of the node
func (n *node) toString(children bool) string {
	var s string

	switch dataLen(n.data) {
	case 2:
		s = fmt.Sprintf("[ %d | %d ]", *n.data[0], *n.data[1])
		break

	case 1:
		s = fmt.Sprintf("[ %d | <nil> ]", *n.data[0])
		break

	case 0:
		s = "[ <nil> | <nil> ]"

	}

	if children {
		switch childLen(n.children) {
		case 3:
			s += "\n/   |   \\"
			s += "\n" + n.children[0].toString(false)
			s += "\n" + n.children[1].toString(false)
			s += "\n" + n.children[2].toString(false)
			break
		case 2:
			s += "\n/     \\"
			s += "\n" + n.children[0].toString(false)
			s += "\n" + n.children[1].toString(false)
			break
		case 1:
			s += "\n   |   "
			s += "\n" + n.children[0].toString(false)
		}
	}

	return s
}

func main() {
	t := &Tree{}

	rands := randomNumbers(9000, 123420)

	t.InsertMany(rands)

	for _, r := range rands {
		if _, e := t.FindNode(r); e == false {
			fmt.Printf("Missing number in tree: %d\n", r)
		}
	}

	fmt.Println("The root after a bunch of random inserts:")
	fmt.Println(t.root.toString(true))
}
