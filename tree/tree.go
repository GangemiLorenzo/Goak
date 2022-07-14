package tree

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

//JSON marshalling
type (
	TreeConf struct {
		Domains []Domain `json:"domains"`
	}
	Domain struct {
		Options []Option `json:"options"`
	}
	Option struct {
		Result Outcome    `json:"result"`
		Want   Conditions `json:"want"`
		Or     Conditions `json:"or"`
	}
)

//Basic types
type (
	Outcome    string
	Condition  string
	Conditions []Condition
)

//DAFAULT outcome for not handled path
const Default Outcome = "DEFAULT"

//Tree branches interfaces and structs

type (
	IBranch interface {
		Search(c Conditions) Outcome
	}
	Node struct {
		Want  Conditions
		Match IBranch
		Fail  IBranch
	}
	Leaf struct {
		Result Outcome
	}
)

//Tree navigation with recursion
func (n Node) Search(c Conditions) Outcome {
	for _, k := range n.Want {
		if ok := c.containsCondition(k); !ok {
			return n.Fail.Search(c)
		}
	}
	return n.Match.Search(c)
}

func (l Leaf) Search(c Conditions) Outcome {
	return l.Result
}

//Given a Tree json, builds the Tree with recursion
func BuildTree(filename string) (t Node) {
	jsonFile, _ := os.Open(filename)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var tc TreeConf
	json.Unmarshal(byteValue, &tc)
	defer jsonFile.Close()

	return buildDomains(
		tc.Domains,
		Leaf{Result: Default},
	).(Node)
}

func buildDomains(ds []Domain, def IBranch) (b IBranch) {

	if len(ds) == 0 {
		return def
	}

	i := len(ds) - 1
	return buildDomains(ds[:i], buildOptions(ds[i].expand().Options, def))
}

func buildOptions(os []Option, def IBranch) (b IBranch) {
	occ := countOccurrences(os)

	//Terminal condition
	if len(occ) == 0 {
		if len(os) == 0 {
			return def
		}
		if len(os) == 1 {
			return Leaf{
				Result: os[0].Result,
			}
		}
		if len(os) != 1 {
			return Leaf{
				Result: Default,
			}
		}
	}

	k := higherOccurrence(occ)

	//Split in 2 slices
	los := []Option{}
	ros := []Option{}

	for _, o := range os {
		if o.Want.containsCondition(*k) {
			o.Want = o.Want.removeCondition(*k)
			ros = append(ros, o)
		} else {
			los = append(los, o)
		}
	}

	return Node{
		Want:  Conditions{(*k)},
		Match: buildOptions(ros, def),
		Fail:  buildOptions(los, def),
	}
}

// Tree building utils

func countOccurrences(os []Option) (occ map[Condition]int) {
	occ = make(map[Condition]int)
	for i := 0; i < len(os); i++ {
		for j := 0; j < len(os[i].Want); j++ {
			w := os[i].Want[j]
			occ[w] = occ[w] + 1
		}
	}
	return
}

func higherOccurrence(occ map[Condition]int) *Condition {
	max := 0
	var cnd Condition

	for k, n := range occ {
		if n > max {
			max = n
			cnd = k
		}
	}

	return &cnd
}

func (d Domain) expand() Domain {
	no := []Option{}
	for _, o := range d.Options {
		if len(o.Or) == 0 {
			no = append(no, o)
			continue
		}

		os := Option{
			Result: o.Result,
			Want:   o.Want,
			Or:     o.Or,
		}.expand()
		no = append(no, os...)
	}

	return Domain{Options: no}
}

func (o Option) expand() []Option {
	res := []Option{}

	for _, or := range o.Or {
		supp := make(Conditions, len(o.Want)+1)
		copy(supp, o.Want)
		supp[len(o.Want)] = or

		res = append(res, Option{
			Result: o.Result,
			Want:   supp,
		})
	}
	return res
}

func (c Conditions) removeCondition(r Condition) Conditions {
	res := Conditions{}
	for _, k := range c {
		if k != r {
			res = append(res, k)
		}
	}
	return res
}

//Basic utils

func (k Condition) toString() string {
	return string(k)
}

func (c Conditions) containsCondition(k Condition) bool {
	return contains(c.toString(), string(k))
}

func (c Conditions) toString() []string {
	r := make([]string, len(c))
	for i, k := range c {
		r[i] = k.toString()
	}
	return r
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

//PRINTING

//Print for a Mermaid markdown file
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

func (n Node) PrintMarkdown(filename string) {
	file, err := os.Create(filename + ".md")
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	fmt.Fprintf(mw, "```mermaid\ngraph TD\n")
	printRecursive(n, mw, "")
	fmt.Fprintf(mw, "```")
}
