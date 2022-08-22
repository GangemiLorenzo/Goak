package tree

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	mOpen  string = "```mermaid\ngraph TD\n"
	mClose string = "```"
	domain string = "### Domain: "
	match  string = "-->|Match|"
	fail   string = "-->|Fail|"
	check  string = "&#10004;"
	cross  string = "&#10007;"
	hdiv   string = ":-:"
	vdiv   string = "|"
	nl     string = "\n"
)

//Print the full tree in Mermaid markdown format
func (t Tree) PrintMarkdownTree(filename string) {
	file, err := os.Create(filename + ".md")
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	fmt.Fprintf(mw, "%s", mOpen)
	printRecursive(t.Root, mw, "")
	fmt.Fprintf(mw, "%s", mClose)
}

func printRecursive(n IBranch, mw io.Writer, bff string) {

	node := n.(Node)
	nbff := bff + node.Want.toPlainString()

	fmt.Fprintf(mw, "%s%s", nbff, match)
	if m, ok := node.Match.(Leaf); ok {
		fmt.Fprintf(mw, "%s;%s", m.Result, nl)
	}
	if m, ok := node.Match.(Node); ok {
		fmt.Fprintf(mw, "%s;%s", nbff+m.Want.toPlainString(), nl)
		printRecursive(node.Match, mw, nbff)
	}

	fmt.Fprintf(mw, "%s%s", nbff, fail)
	if m, ok := node.Fail.(Leaf); ok {
		fmt.Fprintf(mw, "%s;%s", m.Result, nl)
	}
	if m, ok := node.Fail.(Node); ok {
		fmt.Fprintf(mw, "%s;%s", bff+m.Want.toPlainString(), nl)
		printRecursive(node.Fail, mw, bff)
	}
}

//Print the truth table, divided by domains, in markdown format
func (t Tree) PrintMarkdownTable(filename string) {
	file, err := os.Create(filename + ".md")
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, file)

	for index, cs := range t.Conditions {
		os := t.Outcomes[index]
		css := cs.combinations()

		fmt.Fprintf(mw, "%s%d%s%s", domain, index, nl, vdiv)
		for _, c := range cs {
			fmt.Fprintf(mw, "%s%s", c, vdiv)
		}
		fmt.Fprintf(mw, "%s", vdiv)
		for _, o := range os {
			fmt.Fprintf(mw, "%s%s", o, vdiv)
		}
		fmt.Fprintf(mw, "%s%s", nl, vdiv)
		for range cs {
			fmt.Fprintf(mw, "%s%s", hdiv, vdiv)
		}
		fmt.Fprintf(mw, "%s%s", hdiv, vdiv)
		for range os {
			fmt.Fprintf(mw, "%s%s", hdiv, vdiv)
		}
		fmt.Fprintf(mw, "%s", nl)

		for _, cmb := range css {
			fmt.Fprintf(mw, "%s", vdiv)
			r := t.Search(cmb)
			if r != "DEFAULT" {
				i := 0
				for i = 0; i < len(cs); i++ {
					if ok, _ := cmb.contains(cs[i]); ok {
						fmt.Fprintf(mw, "%s%s", check, vdiv)
						continue
					}
					fmt.Fprintf(mw, "%s", vdiv)
				}
				fmt.Fprintf(mw, "%s", vdiv)
				for i = 0; i < len(os); i++ {
					if r == os[i] {
						fmt.Fprintf(mw, "%s%s", cross, vdiv)
						continue
					}
					fmt.Fprintf(mw, "%s", vdiv)
				}
				fmt.Fprintf(mw, "%s", nl)
			}
		}
		fmt.Fprintf(mw, "%s", nl)
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
