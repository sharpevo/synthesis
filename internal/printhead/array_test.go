package printhead_test

import (
	//"fmt"
	"synthesis/internal/geometry"
	"synthesis/internal/printhead"
	"synthesis/internal/reagent"
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
	if rem > step/2 { // not include half of step
		delta += step - rem
	} else {
		delta -= rem
	}
	return delta
}

func newArray() *printhead.Array {
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
	deltay := RoundByStep(p0y, p1y, 4) // 1064
	nozzles1 := p1.MakeNozzles(p1xUnit, p0yUnit+deltay)
	array := printhead.NewArray(
		append(nozzles0, nozzles1...),
	)
	return array
}

func TestRoundByStep(t *testing.T) {
	var delta, expected int
	delta = RoundByStep(20.0, 65.0, 4) // 1063
	expected = 1064
	if delta != expected {
		t.Errorf(tmpl, expected, delta)
	}
	delta = RoundByStep(20.0, 65.0, 2)
	expected = 1062
	if delta != expected {
		t.Errorf(tmpl, expected, delta)
	}
}

func TestNewArray(t *testing.T) {
	array := newArray()
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

func TestMoveRowUpward(t *testing.T) {
	array := newArray()
	array.MoveTopRow(0, 10, 20)
	expectedTop := geometry.NewPosition(13, 312) // 10 + 3, 20 + 292
	if !array.SightTop.Pos.Equal(expectedTop) {
		t.Errorf(
			tmpl,
			*expectedTop,
			*array.SightTop.Pos,
		)
	}
	expectedBottom := geometry.NewPosition(10, -1044) // 10, 312 - (1064 + 292)
	if !array.SightBottom.Pos.Equal(expectedBottom) {
		t.Errorf(
			tmpl,
			expectedBottom,
			*array.SightBottom.Pos,
		)
	}
}

func TestMoveRowDownward(t *testing.T) {
	array := newArray()
	array.MoveBottomRow(0, 10, 20)
	expectedTop := geometry.NewPosition(13, 1376) // 10 + 3, 20 + 1064 + 292
	if !array.SightTop.Pos.Equal(expectedTop) {
		t.Errorf(
			tmpl,
			*expectedTop,
			*array.SightTop.Pos,
		)
	}
	expectedBottom := geometry.NewPosition(10, 20)
	if !array.SightBottom.Pos.Equal(expectedBottom) {
		t.Errorf(
			tmpl,
			expectedBottom,
			*array.SightBottom.Pos,
		)
	}
}
