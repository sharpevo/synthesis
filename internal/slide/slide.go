package slide

import (
	"synthesis/internal/geometry"
)

const (
	WIDTH  = 20 * geometry.MM
	HEIGHT = 50 * geometry.MM
)

type Slide struct {
	Pos        *geometry.Position
	SpotSpaceh int
	SpotSpacev int
	MaxSpotsh  int
	MaxSpotsv  int
	MaxSpots   int
	Spots      [][]*Spot
}

func NewSlide(
	posx int,
	posy int,
	spaceh int,
	spacev int,
) *Slide {
	s := &Slide{
		Pos:        geometry.NewPosition(posx, posy),
		SpotSpaceh: spaceh,
		SpotSpacev: spacev,
		MaxSpotsh:  WIDTH/spaceh + 1,
		MaxSpotsv:  HEIGHT/spacev + 1,
	}
	s.MaxSpots = s.MaxSpotsh * s.MaxSpotsv
	s.Spots = make([][]*Spot, s.MaxSpotsv)
	for k, _ := range s.Spots {
		s.Spots[k] = make([]*Spot, s.MaxSpotsh)
	}
	return s
}

func (s *Slide) Top() int {
	return s.Pos.Y + HEIGHT
}

func (s *Slide) Bottom() int {
	return s.Pos.Y
}

func (s *Slide) Left() int {
	return s.Pos.X
}

func (s *Slide) AddSpot(spot *Spot) bool {
	if s.IsFull() {
		return false
	}
	for y, _ := range s.Spots {
		for x, _ := range s.Spots[y] {
			if s.Spots[y][x] != nil {
				continue
			}
			spot.Pos = geometry.NewPosition(
				s.Left()+x*s.SpotSpaceh,
				s.Top()-y*s.SpotSpacev,
			)
			s.Spots[y][x] = spot
			return true
		}
	}
	return false
}

func (s *Slide) IsFull() bool {
	return s.SpotCount() == s.MaxSpots
}

func (s *Slide) SpotCount() int {
	count := 0
	for y, _ := range s.Spots {
		for x, _ := range s.Spots[y] {
			if s.Spots[y][x] != nil {
				count++
			}
		}

	}
	return count
}

func (s *Slide) ReagentCount() int {
	count := 0
	for _, row := range s.Spots {
		for _, spot := range row {
			if spot != nil {
				count += len(spot.Reagents)
			}
		}
	}
	return count
}

func (s *Slide) AvailableSpots() []*Spot {
	spots := []*Spot{}
	for _, row := range s.Spots {
		for _, spot := range row {
			if spot == nil {
				continue
			}
			spots = append(spots, spot)
		}
	}
	return spots
}

func (s *Slide) SpotsIn(top int, right int, bottom int, left int) []*Spot {
	spots := []*Spot{}
	for _, row := range s.Spots {
		for _, spot := range row {
			if spot == nil {
				continue
			}
			if spot.Pos.X < left || spot.Pos.X > right ||
				spot.Pos.Y < bottom || spot.Pos.Y > top {
				continue
			}
			spots = append(spots, spot)
		}
	}
	return spots
}
