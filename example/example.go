package main

import (
	"fmt"

	"github.com/GangemiLorenzo/decisional-tree/tree"
)

const (
	A tree.Condition = "A"
	B tree.Condition = "B"
	C tree.Condition = "C"
	D tree.Condition = "D"
	E tree.Condition = "E"
)

func main() {
	t := tree.BuildTree("./assets/test_tree.json")
	t.PrintMarkdown("result")

	c := tree.Conditions{
		B,
		C,
		D,
		E,
	}
	res := t.Search(c)
	fmt.Println("RESULT --> " + res)
}
