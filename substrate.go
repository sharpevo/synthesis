package substrate

import (
	"fmt"
	"math"
	"posam/util/geometry"
	"posam/util/log"
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
	LeftMost  int
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
	leftmost int, // unit
) (*Substrate, error) {
	s := &Substrate{
		SpotCount: len(spots),
	}
	slideWidthUnit := geometry.Unit(slideWidth)
	if rem := slideWidthUnit % 4; rem != 0 {
		slideWidthUnit -= rem
	} // new column
	slideHeightUnit := geometry.Unit(slideHeight)
	slideSpaceHoriUnit := geometry.Unit(slideSpaceHori)
	if rem := slideSpaceHoriUnit % 4; rem != 0 {
		slideSpaceHoriUnit -= rem
	} // counted in xoffset
	slideSpaceVertUnit := geometry.Unit(slideSpaceVert)
	if rem := slideSpaceVertUnit % 4; rem != 0 {
		slideSpaceVertUnit -= rem
	} // counted in yoffset
	s.LeftMost = leftmost
	if rem := s.LeftMost % 4; rem != 0 {
		s.LeftMost -= rem
	}

	capacity := slideNumHori * slideNumVert
	//spotsPerSlide := (slideWidthUnit / (1 + spotSpaceUnit)) * (slideHeightUnit / (1 + spotSpaceUnit)) // rem has been removed
	spotsPerSlide := (geometry.Unit(slideWidth)/4 + 1) * (geometry.Unit(slideHeight)/4 + 1)
	required := int(math.Ceil(float64(len(spots)) / float64(spotsPerSlide)))
	log.Vs(log.M{
		"capacity": capacity,
		"required": required,

		"spotsPerSlide": spotsPerSlide,

		"slideWidthUnit":  slideWidthUnit,
		"slideHeightUnit": slideHeightUnit,
		"slideNumHori":    slideNumHori,
		"slideNumVert":    slideNumVert,
		"spotSpaceUnit":   spotSpaceUnit,
		"spots":           len(spots),
		"leftmost":        s.LeftMost,
	}).Info()
	if required > capacity {
		return nil, fmt.Errorf("not enough slide: %v > %v", required, capacity)
	}

	maxSpotsHori := geometry.Unit(slideWidth*float64(slideNumHori)+slideSpaceHori*float64(slideNumHori-1)) + 1 // include 0
	maxSpotsVert := geometry.Unit(slideHeight*float64(slideNumVert)+slideSpaceVert*float64(slideNumVert-1)) + 1
	s.Width = maxSpotsHori + 1  // most right spot missed
	s.Height = maxSpotsVert + 1 // top line, but C & T will not printed since they moves downward
	log.Vs(log.M{
		"slideTop":    s.Top(),
		"slideBottom": s.Bottom(),
		"slideLeft":   s.Left(),
		"slideRight":  s.Right(),
	}).Info("slide info")

	s.Spots = make([][]*Spot, maxSpotsVert+1) // last spot
	for y := range s.Spots {
		s.Spots[y] = make([]*Spot, maxSpotsHori+1)
	}

	slideCount := 1
	columnCount := 1
	xOffset := 0
	yOffset := maxSpotsVert

	fmt.Println(slideWidthUnit, slideHeightUnit, slideSpaceHoriUnit, maxSpotsHori)
	for spotCount, spot := range spots {
		right := (slideWidthUnit+slideSpaceHoriUnit)*(columnCount-1) + slideWidthUnit
		bottom := maxSpotsVert - (slideHeightUnit+slideSpaceVertUnit)*(slideCount-1) - slideHeightUnit
		//bottom := maxSpotsVert - (slideHeightUnit+geometry.Unit(slideSpaceVert))*(slideCount-1) - slideHeightUnit
		log.Vs(log.M{
			"spot":    spot,
			"right":   right,
			"bottom":  bottom,
			"xoffset": xOffset,
			"yoffset": yOffset,

			"slideSpaceHoriUnit": slideSpaceHoriUnit,
			"slideSpaceVertUnit": slideSpaceVertUnit,
			"spotCount":          spotCount,
			"slideCount":         slideCount,
			"columnCount":        columnCount,
		}).Debug()
		//if xOffset >= right {
		if xOffset > right { // xOffset starts from 0 and == should be in the same line, so that 3000 counts is done
			//xOffset >= maxSpotsHori { // xOffset starts from 0
			log.Vs(log.M{
				"yoffset":            yOffset,
				"slideSpaceVertUnit": slideSpaceVertUnit,
				"maxspotsvert":       maxSpotsVert,
				"maxspotshori":       maxSpotsHori,
			}).Info("new line")
			if yOffset-spotSpaceUnit < bottom { // == is the last dots
				//if yOffset <= bottom {
				log.D("new slide")
				xOffset = right - slideWidthUnit
				yOffset -= slideSpaceVertUnit
				slideCount += 1
				if slideCount > slideNumVert {
					// new column
					log.D("new column")
					columnCount += 1
					xOffset = (slideWidthUnit + slideSpaceHoriUnit) * (columnCount - 1)
					yOffset = maxSpotsVert
					slideCount = 1
					if columnCount > slideNumHori {
						return nil, fmt.Errorf(
							"not enough space for spots in horizon: %v > %v", columnCount, slideNumHori)
					}
				}
			} else {
				log.D("same slide")
				xOffset = right - slideWidthUnit
				yOffset -= spotSpaceUnit
				// yOffset may below the bottom
			}
		}

		log.Vs(log.M{
			"x":            xOffset,
			"y":            yOffset,
			"maxspotshori": maxSpotsHori,
			"maxspotsvert": maxSpotsVert,
			"spot":         spot.Reagents[0].Name,
		}).Debug("update spot")
		x := xOffset
		//y := yOffset - 1
		y := yOffset
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
	return s.Height - 1 // include 0
	//return s.Height
}

func (s *Substrate) Right() int {
	//return s.Width - 1
	return s.Width
}

func (s *Substrate) Bottom() int {
	return (s.Height - 1) % 4 // bottom = 6, then nothing in 4x21
}

func (s *Substrate) Left() int {
	return 0
}

func (s *Substrate) Strip() (count int) {
	quo, rem := (s.Width+s.LeftMost)/1280, (s.Width+s.LeftMost)%1280
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
