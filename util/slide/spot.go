package slide

import (
	"synthesis/util/geometry"
	"synthesis/util/reagent"
)

type SpotReagent struct {
	Reagent *reagent.Reagent
	Printed bool
}

type Spot struct {
	Pos      *geometry.Position
	Reagents []*SpotReagent
}

func NewSpot() *Spot {
	return &Spot{}
}

func (s *Spot) AddReagent(r *reagent.Reagent) {
	s.Reagents = append(
		s.Reagents,
		&SpotReagent{
			Reagent: r,
			Printed: false,
		})
}

func (s *Spot) NextReagent(cycleIndex int) *reagent.Reagent {
	if cycleIndex > len(s.Reagents)-1 {
		return nil
	}
	reagent := s.Reagents[cycleIndex]
	if !reagent.Printed {
		return reagent.Reagent
	}
	return nil
}
