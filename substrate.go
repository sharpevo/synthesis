package substrate

import (
	"fmt"
	"posam/util/geometry"
	"posam/util/log"
)

const (
	SLIDE_WIDTH  = 20 // mm
	SLIDE_HEIGHT = 50 // mm
	SLIDE_COUNT  = 3
)

type Substrate struct {
	Spots      [][]*Spot
	SpotCount  int
	SpotSpaceu int
	MaxSpotsh  int
	MaxSpotsv  int
	Width      int
	Height     int
	LeftMostu  int

	SlideHeight  float64
	SlideWidth   float64
	SlideHeightu int
	SlideWidthu  int
	SlideNumh    int
	SlideNumv    int
	SlideSpacehu int
	SlideSpacevu int
}

func NewSubstrate(
	slideNumh int,
	slideNumv int,
	slideWidth float64,
	slideHeight float64,
	slideSpaceHori float64,
	slideSpaceVert float64,
	spots []*Spot,
	spotSpaceu int,
	leftmostu int,
) (*Substrate, error) {
	s := &Substrate{
		SpotCount:    len(spots),
		SpotSpaceu:   spotSpaceu,
		SlideNumh:    slideNumh,
		SlideNumv:    slideNumv,
		SlideWidth:   slideWidth,
		SlideHeight:  slideHeight,
		SlideWidthu:  geometry.RoundedUnit(slideWidth),
		SlideHeightu: geometry.Unit(slideHeight),
		SlideSpacehu: geometry.RoundedUnit(slideSpaceHori),
		SlideSpacevu: geometry.RoundedUnit(slideSpaceVert),
	}
	if err := s.isOverloaded(); err != nil {
		return nil, err
	}
	s.MaxSpotsh = geometry.Unit(
		s.SlideWidth*float64(s.SlideNumh)+slideSpaceHori*float64(s.SlideNumh-1)) + 1
	s.MaxSpotsv = geometry.Unit(
		s.SlideHeight*float64(s.SlideNumv)+slideSpaceVert*float64(s.SlideNumv-1)) + 1
	s.Width = s.MaxSpotsh + 1
	s.Height = s.MaxSpotsv + 1
	s.LeftMostu = leftmostu
	if rem := s.LeftMostu % 4; rem != 0 {
		s.LeftMostu -= rem
	}
	log.Vs(log.M{
		"slideTop":    s.Top(),
		"slideBottom": s.Bottom(),
	}).Infof("substrate created %#v\n", s)
	if err := s.LoadSpots(spots); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Substrate) LoadSpots(spots []*Spot) error {
	s.Spots = make([][]*Spot, s.Height)
	for y := range s.Spots {
		s.Spots[y] = make([]*Spot, s.Width)
	}
	slideCount := 1
	columnCount := 1
	x := 0
	y := s.MaxSpotsv
	for spotCount, spot := range spots {
		right := (s.SlideWidthu+s.SlideSpacehu)*(columnCount-1) + s.SlideWidthu
		bottom := s.MaxSpotsv -
			(s.SlideHeightu+s.SlideSpacevu)*(slideCount-1) - s.SlideHeightu
		log.Vs(log.M{
			"spot":         spot,
			"right":        right,
			"bottom":       bottom,
			"x":            x,
			"y":            y,
			"SlideSpacehu": s.SlideSpacehu,
			"SlideSpacevu": s.SlideSpacevu,
			"spotCount":    spotCount,
			"slideCount":   slideCount,
			"columnCount":  columnCount,
		}).Debug()
		if x > right {
			log.Vs(log.M{
				"y":            y,
				"SlideSpacevu": s.SlideSpacevu,
				"maxspotsvert": s.MaxSpotsv,
				"maxspotshori": s.MaxSpotsh,
			}).Info("new line")
			if y-s.SpotSpaceu < bottom {
				log.D("new slide")
				x = right - s.SlideWidthu
				y -= s.SlideSpacevu
				slideCount += 1
				if slideCount > s.SlideNumv {
					log.D("new column")
					columnCount += 1
					x = (s.SlideWidthu + s.SlideSpacehu) * (columnCount - 1)
					y = s.MaxSpotsv
					slideCount = 1
					if columnCount > s.SlideNumh {
						return fmt.Errorf(
							"not enough space for spots in horizon: %v > %v",
							columnCount, s.SlideNumh)
					}
				}
			} else {
				log.D("same slide")
				x = right - s.SlideWidthu
				y -= s.SpotSpaceu
			}
		}
		log.Vs(log.M{
			"x":            x,
			"y":            y,
			"maxspotshori": s.MaxSpotsh,
			"maxspotsvert": s.MaxSpotsv,
			"spot":         spot.Reagents[0].Name,
		}).Debug("update spot")
		if s.Spots[y][x] != nil {
			return fmt.Errorf("invalid location")
		}
		spot.Pos = geometry.NewPosition(x, y)
		s.Spots[y][x] = spot
		x += s.SpotSpaceu
	}
	return nil
}

func (s *Substrate) Top() int {
	return s.MaxSpotsv
}

func (s *Substrate) Bottom() int {
	return s.MaxSpotsv % 4
}

func (s *Substrate) Strip() (count int) {
	extension := s.MaxSpotsh + s.LeftMostu
	quo, rem := extension/1280, extension%1280
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

func (s *Substrate) isOverloaded() error {
	capacity := s.SlideNumh * s.SlideNumv
	required := int(float64(s.SpotCount)/float64(s.spotsPerSlide()) + 0.5)
	log.Vs(log.M{
		"capacity":      capacity,
		"required":      required,
		"spotsPerSlide": s.spotsPerSlide(),
		"SlideWidthu":   s.SlideWidthu,
		"SlideHeightu":  s.SlideHeightu,
		"slideNumHori":  s.SlideNumh,
		"slideNumVert":  s.SlideNumv,
		"spotSpaceUnit": s.SpotSpaceu,
		"spots":         s.SpotCount,
		"leftmost":      s.LeftMostu,
	}).Info("isOverloaded()")
	if required > capacity {
		return fmt.Errorf("not enough slide: %v > %v", required, capacity)
	}
	return nil
}

func (s *Substrate) spotsPerSlide() int {
	return (geometry.Unit(s.SlideWidth)/4 + 1) *
		(geometry.Unit(s.SlideHeight)/4 + 1)
}
