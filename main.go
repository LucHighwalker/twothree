package main

import (
	"fmt"
	"log"
)

//---------------------------------------------Helpers-------------------------//

// dataLen counts the amount of non-nil elements in a node's data array
func dataLen(data [2]*int) int {
	count := 0
	for i := 0; i < 2; i++ {
		if data[i] != nil {
			count++
		}
	}
	return count
}

// childLen counts the amount of non-nil elements in a node's children array
func childLen(children [3]*node) int {
	count := 0
	for i := 0; i < 3; i++ {
		if children[i] != nil {
			count++
		}
	}
	return count
}

//---------------------------------------------Tree----------------------------//

// Tree struct for a two-three tree
type Tree struct {
	root *node
}

// Tree.FindNode(data) initiatiates the recursive findNode method on the root node.
func (t *Tree) FindNode(data int) *node {
	if t.root != nil {
		return t.root.findNode(data)
	} else {
		return nil
	}
}

// Tree.Insert(data) inserts the given data into the tree.
// TODO: Utilize FindNode's return of a nil value when tree is empty
func (t *Tree) Insert(data int) {
	n := t.FindNode(data)
	if n != nil {
		n.insert(data)
		t.refreshRoot()
	} else {
		t.root = &node{}
		t.root.insert(data)
	}
}

// Tree.refreshRoot() makes sure that the root is the top level node.
// If it isn't the top level node, the root is assigned to the top level node
// TODO: make into recursive function to get the true root.
// TODO: handle empty tree
func (t *Tree) refreshRoot() {
	if t.root.parent == nil {
		return
	} else {
		t.root = t.root.parent
	}
	if t.root.parent != nil {
		log.Fatal("wrong root stupid")
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

// node.findNode(data) recursive method to find the node that data belongs to or should
func (n *node) findNode(data int) *node {
	if data < *n.data[0] {
		if n.children[0] != nil {
			return n.children[0].findNode(data)
		}
	}

	// switch to handle two and three nodes slighlty differently
	switch dataLen(n.data) {
	case 1:
		if n.children[1] != nil {
			return n.children[1].findNode(data)
		}
		break

	case 2:
		if data < *n.data[1] {
			if n.children[1] != nil {
				return n.children[1].findNode(data)
			}
		} else {
			if n.children[2] != nil {
				return n.children[2].findNode(data)
			}
		}
		break

	}
	// node has been found
	return n
}

// node.toParent(data) pushes data to the node's parent
func (n *node) toParent(data int) {
	if n.parent == nil {
		n.parent = &node{data: [2]*int{&data}, children: [3]*node{n}}
	} else {
		n.parent.insert(data)
	}
	n.parent.validate()
}

// node.validate() ensures that the node is a valid two-three node.
func (n *node) validate() {
	childLength := childLen(n.children)
	if childLength == 0 {
		return
	}

	switch dataLen(n.data) {
	case 1:
		if childLength == 2 {
			return
		} else {
			if dataLen(n.children[0].data) == 2 {
				left, right := n.children[0].split()
				n.children = [3]*node{left, right}
				return
			}
		}
		break

	case 2:
		switch childLength {
		case 3:
			return

		case 2:
			if dataLen(n.children[0].data) == 2 {
				left, middle := n.children[0].split()
				n.children = [3]*node{left, middle, n.children[1]}
				return
			} else if dataLen(n.children[1].data) == 2 {
				middle, right := n.children[1].split()
				n.children = [3]*node{n.children[0], middle, right}
				return
			} else {
				data := n.data[1]
				n.data = [2]*int{n.data[0]}
				n.toParent(*data)
				return
			}
		}
		break
	}

	log.Fatal("Could not validate node. Something seems to be wrong....")
}

// node.split() splits the node into 2 nodes and returns them
func (n *node) split() (*node, *node) {
	if n.parent == nil {
		log.Fatal("Cannot split a node without a parent")
	}
	if dataLen(n.data) < 2 {
		log.Fatal("Cannot split a node with singular data")
	}

	left := &node{parent: n.parent}
	left.insert(*n.data[0])

	right := &node{parent: n.parent}
	right.insert(*n.data[1])

	if childLen(n.children) > 0 {

		if dataLen(n.children[0].data) == 2 {
			l, r := n.children[0].split()
			left.children = [3]*node{l, r}
			right.children = [3]*node{n.children[1], n.children[2]}
		} else if dataLen(n.children[1].data) == 2 {
			l, r := n.children[1].split()
			left.children = [3]*node{n.children[0], l}
			right.children = [3]*node{r, n.children[2]}
		} else if dataLen(n.children[2].data) == 2 {
			l, r := n.children[2].split()
			left.children = [3]*node{n.children[0], n.children[1]}
			right.children = [3]*node{l, r}
		}

		left.children[0].parent = left
		left.children[1].parent = left
		right.children[0].parent = right
		right.children[1].parent = right
	}

	return left, right
}

// node.toString() returns a string representation of the node
func (n *node) toString() string {
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

	switch childLen(n.children) {
	case 3:
		s += "\n/   |   \\"
		break
	case 2:
		s += "\n/     \\"
		break
	case 1:
		s += "\n   |   "
	}

	return s
}

func main() {
	t := &Tree{}

	for i := 1; i < 11111; i++ {
		t.Insert(i)
	}

	fmt.Println("new root after a bunch of inserts:")
	fmt.Println(t.root.toString())
}
