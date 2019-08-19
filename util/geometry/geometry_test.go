package geometry

import (
	//"fmt"
	"posam/util/geometry"
	"testing"
)

const tmpl = "\nEXPECT: %v\nGET: %v\n"

func TestUnit(t *testing.T) {
	inputs := []float64{
		50.0,
		25.0,
		25.4,
	}
	outputs := []int{
		1181,
		591,
		600,
	}
	for index, input := range inputs {
		actual := geometry.Unit(input)
		expected := outputs[index]
		if actual != expected {
			t.Errorf(tmpl, expected, actual)
		}
	}
}
