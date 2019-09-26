package substrate_test

import (
	"fmt"
	"reflect"
	"synthesis/internal/substrate"
	"testing"
)

const tmpl = "\nEXPECT: %#v\n GET: %#v\n\n"

func TestSplitByLine(t *testing.T) { // {{{
	cases := []struct {
		input  string
		output []string
	}{
		{
			input:  "one\ntwo\nthree",
			output: []string{"one", "two", "three"},
		},
		{
			input:  "one\r\ntwo\r\nthree\r\n",
			output: []string{"one", "two", "three", ""},
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			output := substrate.SplitByLine(c.input)
			if !reflect.DeepEqual(output, c.output) {
				t.Errorf(tmpl, c.output, output)
			}
		})
	}
} // }}}

func TestSplitByChar(t *testing.T) { // {{{
	cases := []struct {
		input  string
		output []string
	}{
		{
			input:  "abcd",
			output: []string{"a", "b", "c", "d"},
		},
		{
			input:  "  abcd   ",
			output: []string{"a", "b", "c", "d"},
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			output := substrate.SplitByChar(c.input)
			if !reflect.DeepEqual(output, c.output) {
				t.Errorf(tmpl, c.output, output)
			}
		})
	}
} // }}}

func TestAddReagent(t *testing.T) {
	cases := []struct {
		input       string
		spotsLength int
		cycleCount  int
	}{
		{
			input:       "ACGT",
			spotsLength: 1,
			cycleCount:  4,
		},
		{
			input:       "A\nCGT",
			spotsLength: 2,
			cycleCount:  3,
		},
		{
			input:       "",
			spotsLength: 0,
			cycleCount:  0,
		},
		{
			input:       "\n",
			spotsLength: 0,
			cycleCount:  0,
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			spots, cycleCount := substrate.ParseSpots(c.input)
			if cycleCount != c.cycleCount || len(spots) != c.spotsLength {
				t.Errorf(
					tmpl,
					fmt.Sprintf("%+v %+v", c.cycleCount, c.spotsLength),
					fmt.Sprintf("%+v %+v", cycleCount, len(spots)),
				)
			}
		})
	}

}
