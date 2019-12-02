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
		fmt.Println("inserting a third data")
		left := *n.data[0]
		right := *n.data[1]
		if left > data {
			// data becomes left and left gets pushed up
			n.data[0] = &data
			n.ToParent(left)
		} else if right < data {
			// data becomes right and right gets pushed up
			n.data[1] = &data
			n.ToParent(right)
		} else {
			// data gets pushed up
			n.ToParent(data)
		}
		break
	}
}

func (n *Node) ToParent(data int) *Node {
	if n.parent == nil {
		n.parent = &Node{data: [2]*int{&data, nil}}
	} else {
		n.parent.Insert(data)
	}
	return n.parent
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
	fmt.Println(t.root.parent.data)
	fmt.Println(*t.root.parent.data[0])
}
