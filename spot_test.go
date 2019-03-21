package substrate_test

import (
	//"fmt"
	"io/ioutil"
	"posam/util/substrate"
	"reflect"
	"testing"
)

func TestParseLines(t *testing.T) {
	windata, err := ioutil.ReadFile("testfiles/seq.with.windows.new.line")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		input  string
		output []string
	}{
		{
			"1ACTG\nACGT",
			[]string{"1ACTG", "ACGT"},
		},
		{
			"2ACTG\r\nACGT",
			[]string{"2ACTG", "ACGT"},
		},
		{
			string(windata),
			[]string{"1", "2", "3", "4", "5"},
		},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			actual := substrate.ParseLines(c.input)
			if !reflect.DeepEqual(actual, c.output) {
				t.Errorf(
					"\nEXPECT: %#v\n GET: %#v\n\n",
					c.output,
					actual,
				)
			}
		})
	}
}
