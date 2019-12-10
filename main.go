package main

import (
	"fmt"
	"log"
	"math/rand"
)

//---------------------------------------------Helpers-------------------------//

// dataLen(data) counts the amount of non-nil elements in a node's data array
func dataLen(data [2]*int) int {
	count := 0
	for i := 0; i < 2; i++ {
		if data[i] != nil {
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
	for i := 0; i < 4; i++ {
		if children[i] != nil {
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
}

// Tree.Insert(data) inserts the given data into the tree.
func (t *Tree) Insert(data int) {
	if t.Contains(data) == false {
		n := t.FindLeaf(data)
		if n != nil {
			n.insert(data)
			t.refreshRoot()
		} else {
			t.root = &node{}
			t.root.insert(data)
		}
	}
}

func (t *Tree) InsertMany(data []int) {
	for i := 0; i < len(data); i++ {
		t.Insert(data[i])
	}
}

// Tree.FindNode(data) initiatiates the recursive node.findNode(data) method.
func (t *Tree) FindLeaf(data int) *node {
	if t.root != nil {
		return t.root.findLeaf(data)
	} else {
		return nil
	}
}

// Tree.Contains(data) initiates the recursive node.contains(data) method.
func (t *Tree) Contains(data int) bool {
	if t.root != nil {
		return t.root.contains(data)
	}
	return false
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
			// if data is smaller than stored data
			// move stored data to the right and insert
			n.data[0] = &data
			n.data[1] = &stored
		} else {
			// otherwise, insert data to the right
			n.data[1] = &data
		}
		break

	case 2:
		// if node is already a three node
		left := *n.data[0]
		right := *n.data[1]
		if left > data {
			// data becomes left and left gets pushed up
			n.data[0] = &data
			n.toParent(left)
		} else if right < data {
			// data becomes right and right gets pushed up
			n.data[1] = &data
			n.toParent(right)
		} else {
			// data gets pushed up
			n.toParent(data)
		}
		break
	}
}

// node.findNode(data) recursive method to find the the leaf that data should belong to.
func (n *node) findLeaf(data int) *node {
	hasChildren := n.children[0] != nil

	if data < *n.data[0] {
		if hasChildren {
			return n.children[0].findLeaf(data)
		}
	}

	// switch to handle two and three nodes slighlty differently
	switch dataLen(n.data) {
	case 1:
		if hasChildren {
			return n.children[1].findLeaf(data)
		}
		break

	case 2:
		if data < *n.data[1] {
			if hasChildren {
				return n.children[1].findLeaf(data)
			}
		} else {
			if hasChildren {
				return n.children[2].findLeaf(data)
			}
		}
		break

	}
	// node has been found
	return n
}

// node.contains(data) recursive method returns bool if data exists.
func (n *node) contains(data int) bool {
	hasChildren := n.children[0] != nil

	if data == *n.data[0] {
		return true
	} else if hasChildren && data < *n.data[0] {
		return n.children[0].contains(data)
	}

	// switch to handle two and three nodes slighlty differently
	switch dataLen(n.data) {
	case 1:
		if hasChildren {
			return n.children[1].contains(data)
		}
		break

	case 2:
		if data == *n.data[1] {
			return true
		} else if hasChildren {
			if data < *n.data[1] {
				return n.children[1].contains(data)
			} else {
				return n.children[2].contains(data)
			}
		}
		break
	}
	return false
}

// node.toParent(data) pushes data to the node's parent
func (n *node) toParent(data int) {
	if n.parent == nil {
		n.parent = &node{children: [4]*node{n}}
	}
	n.split()
	n.parent.insert(data)
}

func (n *node) kidnap(c *node) {
	for i := 0; i < 4; i++ {
		if n.children[i] == nil {
			n.children[i] = c
			c.parent = n
			return
		}
	}
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

func (n *node) childIndex(c *node) int {
	for i := 0; i < 3; i++ {
		if n.children[i] == c {
			return i
		}
	}
	log.Fatal("Child not present in parent")
	return -1
}

func (n *node) removeChild(c *node) int {
	index := n.childIndex(c)
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
		left.kidnap(n.children[0])
		left.kidnap(n.children[1])
		right.kidnap(n.children[2])
		right.kidnap(n.children[3])
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

	for i := 0; i < 123420; i++ {
		t.Insert(i)
	}

	for i := 0; i < 123420; i++ {
		if t.Contains(i) == false {
			fmt.Printf("Missing number in tree: %d\n", i)
		}
	}

	fmt.Println("the root after a bunch of inserts:")
	fmt.Println(t.root.toString(true))

	t = &Tree{}

	rands := randomNumbers(123420, -1)

	t.InsertMany(rands)

	for i := 0; i < len(rands); i++ {
		if t.Contains(rands[i]) == false {
			fmt.Printf("Missing number in tree: %d\n", rands[i])
		}
	}

	fmt.Println("the tree after a bunch of random inserts:")
	fmt.Println(t.root.toString(true))
	// t.PrintTree()
}
