package substrate

import (
	"fmt"
	"strings"
	"synthesis/internal/geometry"
)

const (
	RESOLUTION_X = 600
	RESOLUTION_Y = 600
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
		SpotSpaceu:   geometry.DPI / RESOLUTION_X,
		SlideNumh:    slideNumh,
		SlideNumv:    slideNumv,
		SlideWidth:   slideWidth,
		SlideHeight:  slideHeight,
		SlideWidthu:  geometry.RoundedDot(slideWidth, RESOLUTION_X),
		SlideHeightu: geometry.Millimeter2Dot(slideHeight),
		//SlideHeightu: geometry.RoundedDot(slideHeight, resolutionY),
		SlideSpacehu: geometry.RoundedDot(slideSpaceh, RESOLUTION_X),
		SlideSpacevu: geometry.RoundedDot(slideSpacev, RESOLUTION_Y),
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
	s.LeftMostu = geometry.RoundDot(leftmostu, RESOLUTION_X)
	fmt.Printf("%+v\n", s)
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
	spotsCount := 0
	x, y, err = s.allocate(x, y, s.marginRightu(), s.marginBottomu())
	if err != nil {
		return err
	}
	if s.Spots[y][x] != nil {
		return fmt.Errorf("invalid location")
	}
	placeholder := strings.Repeat("-", len(spots[0].Reagents))
	realSpotsPerRow := (s.SlideWidthu-1)*s.ResolutionX/geometry.DPI + 1
	prevS := 1
	for spotIndex, spot := range spots {
		fmt.Printf(
			"real spot(%d, %d): spotIndex %d; spotsCount %d; spotSum %d\n",
			x, y,
			spotIndex,
			spotsCount,
			len(spots),
		)
		spot.Pos = geometry.NewPosition(x, y)
		s.Spots[y][x] = spot
		x += s.SpotSpaceu
		spotsCount++

		prevS = s.curSlide

		if spotIndex == len(spots)-1 {
			break
		}
		x, y, err = s.allocate(x, y, s.marginRightu(), s.marginBottomu())
		if err != nil {
			return err
		}
		if s.Spots[y][x] != nil {
			return fmt.Errorf("invalid location")
		}
		fmt.Printf(
			"next: (%d, %d)\n",
			x, y,
		)

		// y of next available spot
		if (spotIndex+1)%realSpotsPerRow == 0 && // last spot of line
			s.curSlide == prevS {
			fmt.Println("new line", spotIndex, realSpotsPerRow, s.curSlide)
			spotsNumberPerLine := s.SlideWidthu
			for k := 0; k < spotsNumberPerLine*(geometry.DPI/s.ResolutionY-1); k++ {
				fmt.Printf("y: (%d, %d) %d %d\n", x, y, spotsCount, spotIndex)
				s.addPlaceholder(x, y, placeholder)
				x += s.SpotSpaceu
				spotsCount++
				x, y, err = s.allocate(x, y, s.marginRightu(), s.marginBottomu())
				if err != nil {
					return err
				}
				if s.Spots[y][x] != nil {
					return fmt.Errorf("invalid location")
				}
			}
			continue
		}

		// x
		if x <= s.marginRightu() && // add holders except the last spot
			s.curSlide == prevS { // next available spot is in the same slide
			for i := 1; i < geometry.DPI/s.ResolutionX; i++ {
				fmt.Printf(
					"x: (%d, %d) spotIndex %d; spotsCount %d; spotSum %d\n",
					x, y,
					spotIndex,
					spotsCount,
					len(spots),
				)
				// TODO: activatable
				s.addPlaceholder(x, y, placeholder)
				x += s.SpotSpaceu
				spotsCount++
				x, y, err = s.allocate(x, y, s.marginRightu(), s.marginBottomu())
				if err != nil {
					return err
				}
				if s.Spots[y][x] != nil {
					return fmt.Errorf("invalid location")
				}
			}
		}
	}
	return nil
}

func (s *Substrate) addPlaceholder(x int, y int, placeholder string) {
	// TODO: activatable
	spots, _ := ParseSpots(placeholder, false)
	spot := spots[0]
	spot.Pos = geometry.NewPosition(x, y)
	s.Spots[y][x] = spot
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
	// next spot in the same row
	if x < marginRightu { // not include last element since x starts from zero
		return x, y, nil
	}
	// next line in the same slide
	if y > marginBottomu+s.SpotSpaceu { //
		x = marginRightu - s.SlideWidthu
		y -= s.SpotSpaceu
		return x, y, nil
	}
	// next slide
	if s.curSlide <= s.SlideNumv-1 {
		fmt.Println("new slide", s.curSlide+1)
		x = marginRightu - s.SlideWidthu
		y -= s.SlideSpacevu + s.SpotSpaceu
		s.curSlide++
		return x, y, nil
	}
	x = (s.SlideWidthu + s.SlideSpacehu) * s.curColumn
	y = s.MaxSpotsv
	s.curColumn++
	s.curSlide = 1
	if s.curColumn > s.SlideNumh {
		fmt.Println("max", s.MaxSpotsh, s.MaxSpotsv)
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
	return s.MaxSpotsv % (geometry.DPI / RESOLUTION_Y)
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
	spotsHorizon := (geometry.Millimeter2Dot(s.SlideWidth)/(geometry.DPI/RESOLUTION_X) + 1)
	spotsVertical := (geometry.Millimeter2Dot(s.SlideHeight)/(geometry.DPI/RESOLUTION_Y) + 1)
	return spotsHorizon * spotsVertical
}
