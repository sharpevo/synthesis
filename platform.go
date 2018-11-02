package platform

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"strings"
)

type Base struct {
	Name  string
	Color color.NRGBA
}

var BaseA = &Base{
	"A",
	color.NRGBA{0xff, 0x00, 0x00, 0xff},
}

var BaseC = &Base{
	"C",
	color.NRGBA{0x00, 0xff, 0x00, 0xff},
}
var BaseG = &Base{
	"G",
	color.NRGBA{0x00, 0x00, 0xff, 0xff},
}

var BaseT = &Base{
	"T",
	color.NRGBA{0xff, 0x00, 0xff, 0xff},
}

var BaseN = &Base{
	"N",
	color.NRGBA{0x00, 0x00, 0x00, 0x00},
}

func NewBase(base string) *Base {
	switch base {
	case "A":
		return BaseA
	case "C":
		return BaseC
	case "G":
		return BaseG
	case "T":
		return BaseT
	default:
		return BaseN
	}
}

type Block struct {
	Sequence  [][]*Base
	PositionX int
	PositionY int
	SpaceX    int
	SpaceY    int
}

func (b *Block) AddRow(row string) {
	bases := []*Base{}
	for _, base := range strings.Split(row, "") {
		bases = append(bases, NewBase(base))
	}
	b.Sequence = append(b.Sequence, bases)
}

type Dot struct {
	Base      *Base
	Printed   bool
	PositionX int
	PositionY int
}

type Platform struct {
	Width  int
	Height int
	Dots   [][]*Dot
}

func NewPlatform(width int, height int) *Platform {
	platform := &Platform{}
	platform.Width = width
	platform.Height = height
	platform.Dots = make([][]*Dot, height)
	for i := range platform.Dots {
		platform.Dots[i] = make([]*Dot, width)
	}
	return platform
}

func (p *Platform) AddBase(x int, y int, base *Base) {
	fmt.Println(x, y, base)
	dot := &Dot{
		Base:      base,
		Printed:   false,
		PositionX: x,
		PositionY: y,
	}
	p.Dots[y][x] = dot
}

func (p *Platform) AddBlock(block *Block) {
	for rowIndex, row := range block.Sequence {
		posy := block.PositionY + rowIndex*(block.SpaceY+1)
		for baseIndex, base := range row {
			posx := block.PositionX + baseIndex*(block.SpaceX+1)
			p.AddBase(posx, posy, base)
		}
	}
}

func (p *Platform) NextPosition() (int, int) {
	for posy, row := range p.Dots {
		for posx, dot := range row {
			if dot == nil {
				continue
			}
			if !dot.Printed {
				return posx, posy
			}
		}
	}
	return 0, 0
}

func (p *Platform) DotsInRow(y int) []*Dot {
	output := []*Dot{}
	for _, dot := range p.Dots[y] {
		if dot == nil {
			continue
		}
		if !dot.Printed {
			output = append(output, dot)
		}
	}
	return output
}

func ParsePlatform(pngPath string) (*Platform, error) {
	existingImageFile, err := os.Open(pngPath)
	if err != nil {
		return nil, err
	}
	defer existingImageFile.Close()
	img, err := png.Decode(existingImageFile)
	if err != nil {
		return nil, err
	}
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	platform := NewPlatform(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y).(color.NRGBA)
			if c == BaseN.Color {
				continue
			}
			platform.AddBase(x, y, ColorToBase(&c))
		}
	}
	return platform, nil
}

func ColorToBase(c *color.NRGBA) *Base {
	if c.R == BaseA.Color.R && c.G == BaseA.Color.G && c.B == BaseA.Color.B {
		return BaseA
	}
	if c.R == BaseC.Color.R && c.G == BaseC.Color.G && c.B == BaseC.Color.B {
		return BaseC
	}
	if c.R == BaseG.Color.R && c.G == BaseG.Color.G && c.B == BaseG.Color.B {
		return BaseG
	}
	if c.R == BaseT.Color.R && c.G == BaseT.Color.G && c.B == BaseT.Color.B {
		return BaseT
	}
	return BaseN
}
