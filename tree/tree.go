package tree

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//JSON marshalling
type (
	TreeConf struct {
		Domains []Domain `json:"domains"`
	}
	Domain struct {
		Options  []Option   `json:"options"`
		Priority Conditions `json:"priority"`
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
	Outcomes   []Outcome
)

//DAFAULT outcome for not handled path
const Default Outcome = "DEFAULT"

//Tree branches interfaces and structs

type Tree struct {
	Root Node
	//Conditions and Outcomes are divided by domain
	Conditions []Conditions
	Outcomes   []Outcomes
}

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
func (t Tree) Search(c Conditions) Outcome {
	return t.Root.Search(c)
}

func (n Node) Search(c Conditions) Outcome {
	for _, k := range n.Want {
		if ok, _ := c.contains(k); !ok {
			return n.Fail.Search(c)
		}
	}
	return n.Match.Search(c)
}

func (l Leaf) Search(c Conditions) Outcome {
	return l.Result
}

//Given a Tree json, builds the Tree with recursion
func BuildTree(filename string) (t Tree) {
	jsonFile, _ := os.Open(filename)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var tc TreeConf
	json.Unmarshal(byteValue, &tc)
	defer jsonFile.Close()

	t.Root = buildDomains(
		tc.Domains,
		Leaf{Result: Default},
	).(Node)

	t.Conditions, t.Outcomes = tc.extractConditionsAndOutcomes()
	return
}

func buildDomains(ds []Domain, def IBranch) (b IBranch) {

	if len(ds) == 0 {
		return def
	}

	i := len(ds) - 1
	return buildDomains(ds[:i], buildOptions(ds[i].expand().Options, def, ds[i].Priority))
}

func buildOptions(os []Option, def IBranch, pr Conditions) (b IBranch) {
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
			//TODO: if this condition happens, then something is wrong (probably duplicated option with different result)
			return Leaf{
				Result: Default,
			}
		}
	}

	k := nextCondition(occ, pr)

	//Split in 2 slices
	los := []Option{}
	ros := []Option{}

	for _, o := range os {
		if ok, _ := o.Want.contains(*k); ok {
			o.Want = o.Want.removeCondition(*k)
			ros = append(ros, o)
		} else {
			los = append(los, o)
		}
	}

	return shrink(
		buildOptions(ros, def, pr),
		buildOptions(los, def, pr),
		k,
	)
}

func shrink(match IBranch, fail IBranch, k *Condition) Node {
	if m, ok := match.(Node); ok {
		if l, ok := m.Fail.(Leaf); ok {
			if l.Result == Default {
				return Node{
					Want:  append(m.Want, (*k)),
					Match: m.Match,
					Fail:  fail,
				}
			}
		}
	}
	return Node{
		Want:  Conditions{(*k)},
		Match: match,
		Fail:  fail,
	}
}

//Get all conditions and outcomes

func (tf TreeConf) extractConditionsAndOutcomes() ([]Conditions, []Outcomes) {
	conditions := []Conditions{}
	outcomes := []Outcomes{}

	for _, d := range tf.Domains {
		c, o := d.extractConditionsAndOutcomes()
		conditions = append(conditions, c)
		outcomes = append(outcomes, o)
	}

	return conditions, outcomes
}

func (d Domain) extractConditionsAndOutcomes() (Conditions, Outcomes) {
	conditions := Conditions{}
	outcomes := Outcomes{}

	for _, o := range d.Options {
		if ok, _ := outcomes.contains(o.Result); !ok {
			outcomes = append(outcomes, o.Result)
		}
		for _, c := range append(o.Want, o.Or...) {
			if ok, _ := conditions.contains(c); !ok {
				conditions = append(conditions, c)
			}
		}
	}

	return conditions, outcomes
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

func nextCondition(occ map[Condition]int, pr Conditions) *Condition {
	min := len(pr)
	var cnd Condition

	for k := range occ {
		if ok, i := pr.contains(k); ok {
			if i < min {
				min = i
				cnd = k
			}
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
