package substrate

import (
	"fmt"
	"posam/util/geometry"
)

const (
	SLIDE_WIDTH  = 20 // mm
	SLIDE_HEIGHT = 50 // mm
	SLIDE_COUNT  = 3
)

type Substrate struct {
	Spots     [][]*Spot
	Width     int
	Height    int
	SpotCount int
}

func NewSubstrate(
	spotSpaceUnit int, // unit
	slideNumHori int,
	slideNumVert int,
	slideWidth float64, // mm
	slideHeight float64, // mm
	slideSpaceHori float64, // mm
	slideSpaceVert float64, // mm
	spots []*Spot,
) (*Substrate, error) {
	s := &Substrate{
		SpotCount: len(spots),
	}
	slideWidthUnit := geometry.Unit(slideWidth)
	slideHeightUnit := geometry.Unit(slideHeight)
	slideSpaceHoriUnit := geometry.Unit(slideSpaceHori)
	slideSpaceVertUnit := geometry.Unit(slideSpaceVert)
	rem := slideSpaceHoriUnit % 4
	if rem != 0 {
		slideSpaceHoriUnit -= rem
	}

	maxSpotsHori := geometry.Unit(slideWidth*float64(slideNumHori) + slideSpaceHori*float64(slideNumHori-1))
	maxSpotsVert := geometry.Unit(slideHeight*float64(slideNumVert) + slideSpaceVert*float64(slideNumVert-1))
	s.Width = maxSpotsHori
	s.Height = maxSpotsVert
	s.Spots = make([][]*Spot, maxSpotsVert)
	for y := range s.Spots {
		s.Spots[y] = make([]*Spot, maxSpotsHori)
	}

	slideCount := 1
	columnCount := 1
	xOffset := 0
	yOffset := slideHeightUnit // 1181

	fmt.Println(slideWidthUnit, slideHeightUnit, slideSpaceHoriUnit, maxSpotsHori)
	for _, spot := range spots {
		right := (slideWidthUnit+slideSpaceHoriUnit)*(columnCount-1) + slideWidthUnit
		bottom := (slideHeightUnit+slideSpaceVertUnit)*(slideCount-1) + slideHeightUnit
		// new line
		if xOffset > right {
			xOffset = right - slideWidthUnit
			yOffset -= spotSpaceUnit
		}
		// new column
		if yOffset <= 0 {
			columnCount += 1
			xOffset = (slideWidthUnit + slideSpaceHoriUnit) * (columnCount - 1)
			yOffset = slideHeightUnit
		}
		if columnCount > slideNumHori {
			return nil, fmt.Errorf("not enough space for spots")
		}

		x := xOffset
		y := yOffset - 1
		if s.Spots[y][x] != nil {
			return nil, fmt.Errorf("invalid location")
		}

		spot.Pos = geometry.NewPosition(x, y)
		s.Spots[y][x] = spot

		xOffset += spotSpaceUnit
	}
	return s, nil
}

func (s *Substrate) Top() int {
	return s.Height - 1
}

func (s *Substrate) Right() int {
	return s.Width - 1
}

func (s *Substrate) Bottom() int {
	return (s.Height - 1) % 4
}

func (s *Substrate) Left() int {
	return 0
}

func (s *Substrate) Strip() (count int) {
	quo, rem := s.Width/1280, s.Width%1280
	if rem != 0 {
		count = quo + 1
	} else {
		count = quo
	}

	// TODO: segs
	//for i := 0; i < s.Width; i += 1280{
	//}

	return count
}
