package main

import (
	"fmt"
	"log"
)

//---------------------------------------------Helpers--------------------------//

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
func childLen(data [3]*node) int {
	count := 0
	for i := 0; i < 3; i++ {
		if data[i] != nil {
			count++
		}
	}
	return count
}

//---------------------------------------------Tree----------------------------//
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
		node := t.FindNode(data)
		node.insert(data)
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
type node struct {
	data     [2]*int
	children [3]*node
	parent   *node
}

// Node.insert(data) inserts the data into the node.
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
		n.split()
		break
	}
}

// Node.findNode(data) recursive method to find the node that data belongs to or should
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

// Node.toParent(data) pushes data to the node's parent
func (n *node) toParent(data int) {
	if n.parent == nil {
		n.parent = &node{data: [2]*int{&data}, children: [3]*node{n}}
	} else {
		n.parent.insert(data)
	}
}

// Node.split() splits the node and reassigns the parent's children.
// TODO: Finish other cases
// TODO: Might be better to split from the parent down?
func (n *node) split() {
	if n.parent == nil {
		log.Fatal("Cannot split a node without a parent")
	}

	dataLength := dataLen(n.data)

	if dataLength < 2 {
		log.Fatal("Cannot split a node with singular data")
	}

	parentLength := dataLen(n.parent.data)
	// parentChildren := childLen(n.parent.children)

	leftNode := &node{parent: n.parent}
	leftNode.insert(*n.data[0])

	rightNode := &node{parent: n.parent}
	rightNode.insert(*n.data[1])

	switch parentLength {
	case 1:
		n.parent.children = [3]*node{leftNode, rightNode}
	}
}

func main() {
	t := &Tree{}
	t.Insert(1)
	fmt.Println("Root")
	fmt.Println(t.root.data)

	t.Insert(2)
	fmt.Println("inserted 2")
	fmt.Println(t.root.data)

	fmt.Println("inserting 3...")
	t.Insert(3)
	fmt.Println(t.root.data)
	fmt.Println(t.root.children)
	fmt.Println(t.root.children[0].data)

	fmt.Println("finding node for 4")
	node := t.FindNode(4)
	fmt.Println(node.data)
	fmt.Println(*node.data[0])
	// fmt.Println(t.root.parent.data)
}
