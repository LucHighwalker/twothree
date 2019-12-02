package main

import "fmt"

type Tree struct {
	root *Node
}
type Node struct {
	data     [2]*int
	children [3]*Node
	parent   *Node
	tree     *Tree
}

func dataLen(data [2]*int) int {
	count := 0
	for i := 0; i < 2; i++ {
		if data[i] != nil {
			count++
		}
	}
	return count
}

func childLen(data [3]*Node) int {
	count := 0
	for i := 0; i < 3; i++ {
		if data[i] != nil {
			count++
		}
	}
	return count
}

func (t *Tree) Insert(data int) {
	if t.root == nil {
		t.root = &Node{}
	}
	t.root.Insert(data)
}

func (n *Node) Insert(data int) {
	length := dataLen(n.data)
	fmt.Println("Inserting data to root node")
	fmt.Println(data, length)
	switch length {
	case 0:
		n.data[0] = &data
		break

	case 1:
		stored := *n.data[0]
		if data < stored {
			n.data[0] = &data
			n.data[1] = &stored
		} else {
			n.data[1] = &data
		}
		break

	case 2:
		// shift up to parent and split
		break
	}
}

// if len(n.data) == 0 {
// 	n.data[0] = &data
// }

func main() {
	t := &Tree{}
	t.Insert(1)
	fmt.Println("Root")
	fmt.Println(t.root.data)
}
