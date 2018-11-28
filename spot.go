package substrate

import (
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

func ParseSpots(input string) ([]*Spot, int) {
	spots := []*Spot{}
	cycleCount := 0
	for _, line := range strings.Split(input, "\n") {
		if line == "" {
			continue
		}
		spot := NewSpot()
		reagents := strings.Split(strings.Trim(line, " "), "")
		length := len(reagents)
		if length > cycleCount {
			cycleCount = length
		}
		for _, name := range reagents {
			r := reagent.NewReagent(name)
			spot.AddReagent(r)
			//if r.Name != reagent.Nil.Name{
			//spot.AddReagent(reagent.Activator)
			//}
		}
		spots = append(spots, spot)
	}
	return spots, cycleCount
}
