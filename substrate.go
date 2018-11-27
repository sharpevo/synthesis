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
	SpotCount int
}

func NewSubstrate(
	spotSpaceUnit int, // unit
	slideNum int,
	slideWidth int, // mm
	slideHeight int, // mm
	slideSpace int, // mm
	spots []*Spot,
) (*Substrate, error) {
	s := &Substrate{
		SpotCount: len(spots),
	}
	slideWidthUnit := geometry.Unit(slideWidth)
	slideHeightUnit := geometry.Unit(slideHeight)
	slideSpaceUnit := geometry.Unit(slideSpace)

	maxSpotsHori := slideWidthUnit*slideNum + slideSpaceUnit*(slideNum-1)
	s.Spots = make([][]*Spot, slideHeightUnit)
	for y := range s.Spots {
		s.Spots[y] = make([]*Spot, maxSpotsHori)
	}

	slideCount := 1
	xOffset := 0
	yOffset := slideHeightUnit // 1181

	fmt.Println(slideWidthUnit, slideHeightUnit, slideSpaceUnit, maxSpotsHori)
	for _, spot := range spots {
		right := (slideWidthUnit+slideSpaceUnit)*(slideCount-1) + slideWidthUnit
		// new line
		if xOffset > right {
			xOffset = right - slideWidthUnit
			yOffset -= spotSpaceUnit
		}
		// new slide
		if yOffset < 0 {
			slideCount += 1
			xOffset = (slideWidthUnit + slideSpaceUnit) * (slideCount - 1)
			yOffset = slideHeightUnit
		}
		if slideCount > slideNum {
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
