package tree

import (
	"fmt"
	"io"
	"log"
	"os"
)

//PRINTING

//Print for a Mermaid markdown file
func (t Tree) PrintMarkdownTree(filename string) {
	file, err := os.Create(filename + ".md")
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	fmt.Fprintf(mw, "```mermaid\ngraph TD\n")
	printRecursive(t.Root, mw, "")
	fmt.Fprintf(mw, "```")
}

func printRecursive(n IBranch, mw io.Writer, bff string) {

	node := n.(Node)
	nbff := bff + node.Want[0].toString()

	fmt.Fprintf(mw, "%s-->|Match|", nbff)
	if m, ok := node.Match.(Leaf); ok {
		fmt.Fprintf(mw, "%s;\n", m.Result)
	}
	if m, ok := node.Match.(Node); ok {
		fmt.Fprintf(mw, "%s;\n", nbff+m.Want[0].toString())
		printRecursive(node.Match, mw, nbff)
	}

	fmt.Fprintf(mw, "%s-->|Fail|", nbff)
	if m, ok := node.Fail.(Leaf); ok {
		fmt.Fprintf(mw, "%s;\n", m.Result)
	}
	if m, ok := node.Fail.(Node); ok {
		fmt.Fprintf(mw, "%s;\n", bff+m.Want[0].toString())
		printRecursive(node.Fail, mw, bff)
	}
}

//Print a Table markdown file
func (t Tree) PrintMarkdownTable(filename string) {
	file, err := os.Create(filename + ".md")
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, file)

	css := t.Conditions.combinations()
	check := "&#10004;"
	cross := "&#10007;"
	hdiv := ":-:"
	vdiv := " "

	fmt.Fprintf(mw, "|")
	for _, c := range t.Conditions {
		fmt.Fprintf(mw, " %s |", c)
	}
	fmt.Fprintf(mw, "   |")
	for _, o := range t.Outcomes {
		fmt.Fprintf(mw, " %s |", o)
	}
	fmt.Fprintf(mw, "\n")
	fmt.Fprintf(mw, "|")
	for range t.Conditions {
		fmt.Fprintf(mw, "%s|", hdiv)
	}
	fmt.Fprintf(mw, "%s|", hdiv)
	for range t.Outcomes {
		fmt.Fprintf(mw, "%s|", hdiv)
	}
	fmt.Fprintf(mw, "\n")

	for _, cmb := range css {
		fmt.Fprintf(mw, "|")

		r := t.Search(cmb)

		i := 0
		for i = 0; i < len(t.Conditions); i++ {
			if cmb.contains(t.Conditions[i]) {
				fmt.Fprintf(mw, " %s |", check)
				continue
			}
			fmt.Fprintf(mw, "   |")
		}
		fmt.Fprintf(mw, "%s|", vdiv)
		for i = 0; i < len(t.Outcomes); i++ {
			if r == t.Outcomes[i] {
				fmt.Fprintf(mw, " %s |", cross)
				continue
			}
			fmt.Fprintf(mw, "   |")
		}

		fmt.Fprintf(mw, "\n")
	}

}

func (cs Conditions) combinations() []Conditions {
	l := len(cs)
	comb := []Conditions{}
	r := Conditions{}

	combinationsRecursive(0, l, &comb, cs, r)
	return comb
}

func combinationsRecursive(i int, l int, comb *[]Conditions, cs Conditions, r Conditions) {
	if i == l {
		(*comb) = append((*comb), r)
		return
	}

	right := append(r, cs[i])
	combinationsRecursive(i+1, l, comb, cs, right)

	left := r
	combinationsRecursive(i+1, l, comb, cs, left)
}
