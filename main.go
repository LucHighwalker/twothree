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
func childLen(children [3]*node) int {
	count := 0
	for i := 0; i < 3; i++ {
		if children[i] != nil {
			count++
		} else {
			break
		}
	}
	return count
}

func randomNumbers(count int, max int) []int {
	seed := false
	if max < 1 {
		seed = true
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
	if t.root.contains(data) == false {
		n := t.FindNode(data)
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
func (t *Tree) FindNode(data int) *node {
	if t.root != nil {
		return t.root.findNode(data)
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

//---------------------------------------------Node----------------------------//

// node struct for a two-three tree
type node struct {
	data     [2]*int
	children [3]*node
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

// node.findNode(data) recursive method to find the node that data should belong to.
func (n *node) findNode(data int) *node {
	hasChildren := childLen(n.children) > 0

	if data < *n.data[0] {
		if hasChildren {
			return n.children[0].findNode(data)
		}
	}

	// switch to handle two and three nodes slighlty differently
	switch dataLen(n.data) {
	case 1:
		if hasChildren {
			return n.children[1].findNode(data)
		}
		break

	case 2:
		if data < *n.data[1] {
			if hasChildren {
				return n.children[1].findNode(data)
			}
		} else {
			if hasChildren {
				return n.children[2].findNode(data)
			}
		}
		break

	}
	// node has been found
	return n
}

// node.contains(data) recursive method to find specific location of data.
func (n *node) contains(data int) bool {
	hasChildren := childLen(n.children) > 0

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
		n.parent = &node{children: [3]*node{n}}
	}
	n.parent.insert(data)
	n.parent.validate()
}

// node.validate() ensures that the node is a valid two-three node.
func (n *node) validate() {
	childLength := childLen(n.children)

	// switch to handle two and three nodes slighlty differently
	switch dataLen(n.data) {
	case 1:
		if childLength == 2 {
			// a 2 node with 2 children is valid
			return
		} else {
			if dataLen(n.children[0].data) == 2 {
				// split the child if it is a 3 node
				left, right := n.children[0].split()
				n.children = [3]*node{left, right}
				return
			}
		}
		break

	case 2:
		if childLength == 3 {
			// a 3 node with 3 children is valid
			return
		} else if childLength == 2 {
			if dataLen(n.children[0].data) == 2 {
				// if the left node is a 3 node, split it
				left, middle := n.children[0].split()
				n.children = [3]*node{left, middle, n.children[1]}
				return
			} else if dataLen(n.children[1].data) == 2 {
				// if the right node is a 3 node, split it
				middle, right := n.children[1].split()
				n.children = [3]*node{n.children[0], middle, right}
				return
			}
		}
		break
	}

	log.Fatalf("\nCould not validate node:\n%s\nSomething seems to be wrong....\n", n.toString(true))
}

// node.split() splits the node into 2 nodes and returns them
func (n *node) split() (*node, *node) {
	if n.parent == nil {
		log.Fatal("Cannot split a node without a parent")
	}
	if dataLen(n.data) < 2 {
		log.Fatal("Cannot split a node with singular data")
	}

	// insert left data into left node
	left := &node{parent: n.parent}
	left.insert(*n.data[0])

	// insert right data into right node
	right := &node{parent: n.parent}
	right.insert(*n.data[1])

	if childLen(n.children) > 0 {
		// if node has children, split an available
		// 3 node child and assign the children
		if dataLen(n.children[0].data) == 2 {
			ll, rr := n.children[0].split()
			left.children = [3]*node{ll, rr}
			right.children = [3]*node{n.children[1], n.children[2]}
		} else if dataLen(n.children[1].data) == 2 {
			lr, rl := n.children[1].split()
			left.children = [3]*node{n.children[0], lr}
			right.children = [3]*node{rl, n.children[2]}
		} else if dataLen(n.children[2].data) == 2 {
			rl, rr := n.children[2].split()
			left.children = [3]*node{n.children[0], n.children[1]}
			right.children = [3]*node{rl, rr}
		}

		// assign the new parent of the children
		left.children[0].parent = left
		left.children[1].parent = left
		right.children[0].parent = right
		right.children[1].parent = right
	}

	return left, right
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

	rands := randomNumbers(10, 20)

	t.InsertMany(rands)

	for i := 0; i < len(rands); i++ {
		if t.Contains(rands[i]) == false {
			fmt.Printf("Missing number in tree: %d\n", rands[i])
		}
	}

	fmt.Println("the root after a bunch of random inserts:")
	fmt.Println(t.root.toString(true))
}
