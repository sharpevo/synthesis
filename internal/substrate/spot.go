package substrate

import (
	"strings"
	"synthesis/internal/geometry"
	"synthesis/internal/reagent"
)

type Spot struct {
	Pos      *geometry.Position
	Reagents []*reagent.Reagent
}

func (s *Spot) addReagents(names []string) {
	for _, name := range names {
		r := reagent.NewReagent(name)
		s.Reagents = append(s.Reagents, r)
	}
}

func ParseSpots(input string) ([]*Spot, int) {
	spots := []*Spot{}
	cycleCount := 0
	for _, line := range splitByLine(input) {
		if line == "" {
			continue
		}
		spot := &Spot{}
		names := splitByChar(line)
		length := len(names)
		if length > cycleCount {
			cycleCount = length
		}
		spot.addReagents(names)
		spots = append(spots, spot)
	}
	return spots, cycleCount
}

func splitByLine(input string) []string {
	return strings.Split(strings.Replace(input, "\r\n", "\n", -1), "\n")
}

func splitByChar(input string) []string {
	return strings.Split(strings.Trim(input, " "), "")
}
