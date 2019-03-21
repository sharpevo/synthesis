package substrate

import (
	"fmt"
	"posam/util/geometry"
	"posam/util/log"
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

func (s *Spot) addReagent(r *reagent.Reagent) {
	s.Reagents = append(s.Reagents, r)
}

func (s *Spot) AddReagents(names []string, activatable bool) {
	for _, name := range names {
		r := reagent.NewReagent(name)
		s.addReagent(r)
		if activatable {
			if r.Name != reagent.Nil.Name {
				s.addReagent(reagent.Activator)
			} else {
				s.addReagent(reagent.Nil)
			}
		}
	}
}

func ParseSpots(input string, activatable bool) ([]*Spot, int) {
	log.Dv(log.M{"activatable": activatable})
	spots := []*Spot{}
	cycleCount := 0
	for _, line := range parseLines(input) {
		if line == "" {
			continue
		}
		spot := NewSpot()
		reagents := strings.Split(strings.Trim(line, " "), "")
		length := len(reagents)
		if length > cycleCount {
			cycleCount = length
		}
		spot.AddReagents(names, activatable)
		spots = append(spots, spot)
	}
	if activatable {
		return spots, cycleCount * 2
	}
	return spots, cycleCount
}

func parseLines(input string) []string {
	return strings.Split(strings.Replace(input, "\r\n", "\n", -1), "\n")
}
