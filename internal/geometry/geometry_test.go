package geometry_test

import (
	"fmt"
	"synthesis/internal/geometry"
	"testing"
)

const tmpl = "\nEXPECT: %v\nGET: %v\n"

func TestMillimeter2Dot(t *testing.T) {
	cases := []struct {
		input  float64
		output int
	}{
		{
			input:  1.1, // 25.9 + 0.5
			output: 26,
		},
		{
			input:  1.5, // 35.4 + 0.5
			output: 35,
		},
	}
	for index, c := range cases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			output := geometry.Millimeter2Dot(c.input)
			if output != c.output {
				t.Errorf(tmpl, c.output, output)
			}
		})
	}
}

func TestRoundedDot(t *testing.T) {
	cases := []struct {
		input  float64
		dpi    int
		output int
	}{
		{
			input:  1.1, // 25.9 + 0.5
			dpi:    150,
			output: 24, // 26
		},
		{
			input:  1.5, // 35.4 + 0.5
			dpi:    150,
			output: 32, // 35
		},
	}
	for index, c := range cases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			output := geometry.RoundedDot(c.input, c.dpi)
			if output != c.output {
				t.Errorf(tmpl, c.output, output)
			}
		})
	}
}

func TestRoundDot(t *testing.T) {
	cases := []struct {
		dividend int
		dpi      int
		output   int
	}{
		{
			dividend: 26,
			dpi:      150,
			output:   24,
		},
		{
			dividend: 31,
			dpi:      150,
			output:   28,
		},
		{
			dividend: 26,
			dpi:      300,
			output:   26,
		},
		{
			dividend: 31,
			dpi:      300,
			output:   30,
		},
		{
			dividend: 26,
			dpi:      600,
			output:   26,
		},
		{
			dividend: 31,
			dpi:      600,
			output:   31,
		},
	}
	for index, c := range cases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			output := geometry.RoundDot(c.dividend, c.dpi)
			if output != c.output {
				t.Errorf(tmpl, c.output, output)
			}
		})
	}
}

func TestDot2Millimeter(t *testing.T) {
	cases := []struct {
		input  int
		output float64
	}{
		{
			input:  26,
			output: 1.1006666666666667,
		},
		{
			input:  32,
			output: 1.3546666666666667,
		},
	}
	for index, c := range cases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			output := geometry.Dot2Millimeter(c.input)
			if output != c.output {
				t.Errorf(tmpl, c.output, output)
			}
		})
	}
}
