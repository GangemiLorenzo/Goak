package tree

//Contains

func (c Conditions) contains(k Condition) (bool, int) {
	return contains(c.toString(), k.toString())
}

func (o Outcomes) contains(k Outcome) (bool, int) {
	return contains(o.toString(), k.toString())
}

func contains(s []string, str string) (bool, int) {
	for i, v := range s {
		if v == str {
			return true, i
		}
	}

	return false, -1
}

//ToString

func (k Condition) toString() string {
	return string(k)
}

func (c Conditions) toString() []string {
	r := make([]string, len(c))
	for i, k := range c {
		r[i] = k.toString()
	}
	return r
}

func (cs Conditions) toPlainString() (r string) {
	r = ""
	for _, c := range cs {
		r = r + c.toString()
	}
	return
}

func (k Outcome) toString() string {
	return string(k)
}

func (o Outcomes) toString() []string {
	r := make([]string, len(o))
	for i, k := range o {
		r[i] = k.toString()
	}
	return r
}
