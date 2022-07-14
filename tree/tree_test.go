package tree

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type tests struct {
	TestOptions []testOption `json:"tests"`
}

type testOption struct {
	Exp Outcome    `json:"expected"`
	C   Conditions `json:"conditions"`
}

func Test(t *testing.T) {
	tree := BuildTree("./assets/test_tree.json")
	jsonFile, _ := os.Open("./assets/test_options.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var ts tests
	json.Unmarshal(byteValue, &ts)
	defer jsonFile.Close()

	for _, to := range ts.TestOptions {

		d := "Given the conditions ["
		for _, c := range to.C {
			d = d + c.toString() + ", "
		}
		d = d + "] i expect the result: " + string(to.Exp)

		t.Run(d, func(t *testing.T) {
			want := to.Exp
			got := tree.Search(to.C)

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
