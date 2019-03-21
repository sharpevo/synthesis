package substrate_test

import (
	"fmt"
	"io/ioutil"
	"posam/util/substrate"
	"reflect"
	"testing"
)

func TestParseLines(t *testing.T) { // {{{
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
} // }}}

func TestAddReagents(t *testing.T) { // {{{
	cases := []struct {
		names       []string
		activatable bool
		expect      []string
	}{
		{
			[]string{"A", "C", "T", "G"},
			false,
			[]string{"A", "C", "T", "G"},
		},
		{
			[]string{"A", "C", "T", "G"},
			true,
			[]string{"A", "Z", "C", "Z", "T", "Z", "G", "Z"},
		},
		{
			[]string{"A", "-", "T", "G"},
			false,
			[]string{"A", "-", "T", "G"},
		},
		{
			[]string{"A", "-", "T", "G"},
			true,
			[]string{"A", "Z", "-", "-", "T", "Z", "G", "Z"}, // never activate nil
		},
		{
			[]string{"A", "B", "T", "G"},
			false,
			[]string{"A", "B", "T", "G"},
		},
		{
			[]string{"A", "B", "T", "G"},
			true,
			[]string{"A", "Z", "B", "Z", "T", "Z", "G", "Z"},
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.names), func(t *testing.T) {
			s := substrate.NewSpot()
			s.AddReagents(c.names, c.activatable)
			names := []string{}
			for _, r := range s.Reagents {
				names = append(names, r.Name)
			}
			if !reflect.DeepEqual(names, c.expect) {
				t.Errorf(
					"\nEXPECT: %#v\n GET: %#v\n\n",
					c.expect,
					names,
				)
			}
		})
	}
} // }}}

func TestParseReagentName(t *testing.T) { // {{{
	cases := []struct {
		input  string
		expect []string
	}{
		{
			"ACTG",
			[]string{"A", "C", "T", "G"},
		},
		{
			"A CTG",
			[]string{"A", " ", "C", "T", "G"}, // allow space in sequences ?
		},
		{
			" ACTG ",
			[]string{"A", "C", "T", "G"},
		},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			actual := substrate.ParseReagentNames(c.input)
			if !reflect.DeepEqual(actual, c.expect) {
				t.Errorf(
					"\nEXPECT: %#v\n GET: %#v\n\n",
					c.expect,
					actual,
				)
			}
		})
	}
} // }}}

func TestParseSpots(t *testing.T) { // {{{
	cases := []struct {
		input             string
		activatable       bool
		expectSpotsLength int
		expectCycleCount  int
	}{
		{
			"ACGT",
			false,
			1,
			4,
		},
		{
			"ACGT\nCCCC",
			false,
			2,
			4,
		},
		{
			"A\nCCCCC",
			false,
			2,
			5,
		},
		{
			"A\nCCCCC",
			true,
			2,
			10,
		},
		{
			"A\n\n",
			false,
			1,
			1,
		},
		{
			"",
			true,
			0,
			0,
		},
		{
			"\n",
			false,
			0,
			0,
		},
	}
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			spots, cycleCount := substrate.ParseSpots(c.input, c.activatable)
			if len(spots) != c.expectSpotsLength ||
				cycleCount != c.expectCycleCount {
				t.Error(
					"\nEXPECT\n",
					c.expectSpotsLength,
					c.expectCycleCount,
					"\nGET\n",
					len(spots),
					cycleCount,
				)
			}
		})
	}
} // }}}
