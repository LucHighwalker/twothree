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
	if t.root == nil {
		t.root = &node{}
		t.root.insert(data)
	} else {
		n := t.FindNode(data)
		n.insert(data)
		t.refreshRoot()
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
		n.parent.validate()
		break
	}
}

// node.findNode(data) recursive method to find the node that data belongs to or should
func (n *node) findNode(data int) *node {
	// switch to handle two and three nodes slighlty differently
	switch dataLen(n.data) {
	case 1:
		if data < *n.data[0] {
			if n.children[0] != nil {
				return n.children[0].findNode(data)
			}
		} else {
			if n.children[1] != nil {
				return n.children[1].findNode(data)
			}
		}
		break

	case 2:
		if data < *n.data[0] {
			if n.children[0] != nil {
				return n.children[0].findNode(data)
			}
		} else if data < *n.data[1] {
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
				n.parent.validate()
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

	return left, right
}

// node.toString() returns a string representation of the node
func (n *node) toString() string {
	var s string

	if dataLen(n.data) == 2 {
		s = fmt.Sprintf("[ %d | %d ]", *n.data[0], *n.data[1])
	} else {
		s = fmt.Sprintf("[ %d | <nil> ]", *n.data[0])
	}

	childLength := childLen(n.children)

	if childLength == 3 {
		s += "\n/   |   \\"
	} else if childLength == 2 {
		s += "\n/     \\"
	}

	return s
}

func main() {
	t := &Tree{}
	t.Insert(1)
	fmt.Println("Root")
	fmt.Println(t.root.toString())

	t.Insert(2)
	fmt.Println("inserted 2")
	fmt.Println(t.root.toString())

	fmt.Println("inserting 3...")
	t.Insert(3)
	fmt.Println(t.root.toString())
	fmt.Println(t.root.children)
	fmt.Println(t.root.children[0].toString())
	fmt.Println(t.root.children[1].toString())

	fmt.Println("finding node for 4")
	n := t.FindNode(4)
	fmt.Println(n.data)
	fmt.Println(*n.data[0])

	t.Insert(4)
	fmt.Println(t.root.children[0].toString())
	fmt.Println(t.root.children[1].toString())
	t.Insert(5)
	t.Insert(6)
	t.Insert(7)
	t.Insert(8)
	t.Insert(9)
	t.Insert(10)

	fmt.Println("new root after a bunch of inserts:")
	fmt.Println(t.root.toString())
}
