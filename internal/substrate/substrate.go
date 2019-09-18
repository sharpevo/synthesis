package substrate

import (
	"fmt"
	"synthesis/internal/geometry"
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

	curSlide  int
	curColumn int

	ResolutionX int
	ResolutionY int
}

func NewSubstrate(
	slideNumh int,
	slideNumv int,
	slideWidth float64,
	slideHeight float64,
	slideSpaceh float64,
	slideSpacev float64,
	spots []*Spot,
	leftmostu int,
	resolutionX int,
	resolutionY int,
) (*Substrate, error) {
	s := &Substrate{
		SpotCount:    len(spots),
		SpotSpaceu:   geometry.DPI / resolutionX,
		SlideNumh:    slideNumh,
		SlideNumv:    slideNumv,
		SlideWidth:   slideWidth,
		SlideHeight:  slideHeight,
		SlideWidthu:  geometry.RoundedDot(slideWidth, resolutionX),
		SlideHeightu: geometry.Millimeter2Dot(slideHeight),
		SlideSpacehu: geometry.RoundedDot(slideSpaceh, resolutionX),
		SlideSpacevu: geometry.RoundedDot(slideSpacev, resolutionY),
		ResolutionX:  resolutionX,
		ResolutionY:  resolutionY,
	}
	if err := s.isOverloaded(); err != nil {
		return nil, err
	}
	s.MaxSpotsh = s.MaxSpotshu(slideSpaceh)
	s.MaxSpotsv = s.MaxSpotsvu(slideSpacev)
	s.Width = s.MaxSpotsh + 1
	s.Height = s.MaxSpotsv + 1
	s.LeftMostu = geometry.RoundDot(leftmostu, resolutionX)
	if err := s.loadSpots(spots); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Substrate) MaxSpotshu(slideSpaceh float64) int {
	return geometry.Millimeter2Dot(
		s.SlideWidth*float64(s.SlideNumh)+slideSpaceh*float64(s.SlideNumh-1)) + 1
}

func (s *Substrate) MaxSpotsvu(slideSpacev float64) int {
	return geometry.Millimeter2Dot(
		s.SlideHeight*float64(s.SlideNumv)+slideSpacev*float64(s.SlideNumv-1)) + 1
}

func (s *Substrate) loadSpots(spots []*Spot) (err error) {
	s.Spots = make([][]*Spot, s.Height)
	for y := range s.Spots {
		s.Spots[y] = make([]*Spot, s.Width)
	}
	s.curSlide = 1
	s.curColumn = 1
	x := 0
	y := s.MaxSpotsv
	for _, spot := range spots {
		x, y, err = s.allocate(x, y, s.marginRightu(), s.marginBottomu())
		if err != nil {
			return err
		}
		if s.Spots[y][x] != nil {
			return fmt.Errorf("invalid location")
		}
		spot.Pos = geometry.NewPosition(x, y)
		s.Spots[y][x] = spot
		x += s.SpotSpaceu
	}
	return nil
}

func (s *Substrate) marginRightu() int {
	return (s.SlideWidthu+s.SlideSpacehu)*(s.curColumn-1) + s.SlideWidthu
}

func (s *Substrate) marginBottomu() int {
	return s.MaxSpotsv -
		(s.SlideHeightu+s.SlideSpacevu)*(s.curSlide-1) - s.SlideHeightu
}

func (s *Substrate) allocate(
	x, y, marginRightu, marginBottomu int) (int, int, error) {
	if x <= marginRightu {
		return x, y, nil
	}
	if y >= marginBottomu+s.SpotSpaceu {
		x = marginRightu - s.SlideWidthu
		y -= s.SpotSpaceu
		return x, y, nil
	}
	if s.curSlide <= s.SlideNumv-1 {
		x = marginRightu - s.SlideWidthu
		y -= s.SlideSpacevu
		s.curSlide++
		return x, y, nil
	}
	x = (s.SlideWidthu + s.SlideSpacehu) * s.curColumn
	y = s.MaxSpotsv
	s.curColumn++
	s.curSlide = 1
	if s.curColumn > s.SlideNumh {
		return x, y, fmt.Errorf(
			"not enough space for spots in horizon: %v > %v",
			s.curColumn, s.SlideNumh)
	}
	return x, y, nil
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
	if required > capacity {
		return fmt.Errorf("not enough slide: %v > %v", required, capacity)
	}
	return nil
}

func (s *Substrate) spotsPerSlide() int {
	return (geometry.Unit(s.SlideWidth)/4 + 1) *
		(geometry.Unit(s.SlideHeight)/4 + 1)
}
