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

func main() {
	fmt.Println("Hello froma two three world")
}
