package substrate

import (
	"fmt"
	"posam/util/geometry"
	"posam/util/reagent"
	"strings"
)

type Spot struct {
	Pos      *geometry.Position
	Reagents []*reagent.Reagent
}

func NewSpot() *Spot {
	return &Spot{}
}

func (s *Spot) AddReagent(r *reagent.Reagent) {
	s.Reagents = append(s.Reagents, r)
}

func ParseSpots(input string, activatable bool) ([]*Spot, int) {
	fmt.Println("activatable", activatable)
	spots := []*Spot{}
	cycleCount := 0
	sep := "\n"
	for _, line := range strings.Split(input, sep) {
		if line == "" {
			continue
		}
		spot := NewSpot()
		reagents := strings.Split(strings.Trim(line, " "), "")
		length := len(reagents)
		if length > cycleCount {
			cycleCount = length
		}
		//fmt.Println(line)
		for _, name := range reagents {
			r := reagent.NewReagent(name)
			spot.AddReagent(r)
			if activatable {
				if r.Name != reagent.Nil.Name {
					spot.AddReagent(reagent.Activator)
				} else {
					spot.AddReagent(reagent.Nil)
				}
			}
			//fmt.Printf("'%#v'\n", name)
		}
		spots = append(spots, spot)
	}
	if activatable {
		return spots, cycleCount * 2
	}
	return spots, cycleCount
}
