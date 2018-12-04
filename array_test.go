package printhead_test

import (
	"fmt"
	"posam/util/geometry"
	"posam/util/printhead"
	"posam/util/reagent"
	"testing"
)

const tmpl = "\nEXPECT: %v\nGET: %v\n"

func RoundByStep(
	p0 float64,
	p1 float64,
	step int,
) int {
	delta := geometry.Unit(p1 - p0)
	rem := delta % step
	fmt.Println(">", delta)
	if rem > step/2 { // not include half of step
		delta += step - rem
	} else {
		delta -= rem
	}
	return delta
}

func TestRoundByStep(t *testing.T) {
	delta1 := RoundByStep(20.0, 65.0, 4)
	fmt.Println(delta1)
	delta2 := RoundByStep(20.0, 65.0, 2)
	fmt.Println(delta2)
}

func TestNewArray(t *testing.T) {
	p0 := printhead.NewPrinthead(
		"p0 path",
		[]*reagent.Reagent{
			reagent.NewReagent("A"),
			reagent.NewReagent("C"),
			reagent.NewReagent("G"),
			reagent.NewReagent("T"),
		},
	)
	p0x, p0y := 35.0, 20.0
	p0xUnit := geometry.Unit(p0x)
	p0yUnit := geometry.Unit(p0y)
	nozzles0 := p0.MakeNozzles(p0xUnit, p0yUnit)

	p1 := printhead.NewPrinthead(
		"p1 path",
		[]*reagent.Reagent{
			reagent.NewReagent("Z"),
			reagent.NewReagent("-"),
			reagent.NewReagent("-"),
			reagent.NewReagent("-"),
		},
	)
	p1x, p1y := 35.0, 65.0
	p1xUnit := geometry.Unit(p1x)
	deltay := RoundByStep(p0y, p1y, 4)
	nozzles1 := p1.MakeNozzles(p1xUnit, p0yUnit+deltay)
	array := printhead.NewArray(
		append(nozzles0, nozzles1...),
	)
	expectedTop := geometry.NewPosition(830, 1828)
	if !array.SightTop.Pos.Equal(expectedTop) {
		t.Errorf(
			tmpl,
			*expectedTop,
			*array.SightTop.Pos,
		)
	}
	expectedBottom := geometry.NewPosition(827, 472)
	if !array.SightBottom.Pos.Equal(expectedBottom) {
		t.Errorf(
			tmpl,
			expectedBottom,
			*array.SightBottom.Pos,
		)
	}
}
