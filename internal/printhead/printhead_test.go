package printhead_test

import (
	"fmt"
	"synthesis/internal/printhead"
	"synthesis/internal/reagent"
	"testing"
)

func TestMakeNozzles(t *testing.T) {
	reagents := []*reagent.Reagent{
		reagent.NewReagent("A"),
		reagent.NewReagent("C"),
		reagent.NewReagent("G"),
		reagent.NewReagent("T"),
	}
	p := printhead.NewPrinthead("A", reagents)
	nozzles := p.MakeNozzles(0, 0)
	for _, n := range nozzles[:8] {
		//fmt.Printf("%#v\n", n.Pos.X, n.Pos.Y)
		fmt.Println(n.Pos.X, n.Pos.Y)
	}

	p2 := printhead.NewPrinthead("B", reagents)
	nozzles2 := p2.MakeNozzles(20, 30)
	for _, n := range nozzles2[:30] {
		//fmt.Printf("%#v\n", n.Pos.X, n.Pos.Y)
		fmt.Println(n.Pos.X, n.Pos.Y)
	}
}
