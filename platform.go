package platform

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"strings"
)

const (
	// units
	NM = 1
	UM = 1e3
	MM = 1e6
)

type Base struct {
	Name  string
	Color *color.NRGBA
}

var data uint8 = 0xff

var BaseA = &Base{
	"A",
	&color.NRGBA{0xff, 0x00, 0x00, data},
}

var BaseC = &Base{
	"C",
	&color.NRGBA{0x00, 0xff, 0x00, data},
}
var BaseG = &Base{
	"G",
	&color.NRGBA{0x00, 0x00, 0xff, data},
}

var BaseT = &Base{
	"T",
	&color.NRGBA{0xff, 0x00, 0xff, data},
}

var BaseN = &Base{
	"N",
	&color.NRGBA{0x00, 0x00, 0x00, 0x00},
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
	Minx   int
	Maxx   int
	Miny   int
	Maxy   int
	Dots   [][]*Dot
}

//func NewPlatform(width int, height int) *Platform {
//platform := &Platform{}
//platform.Width = width
//platform.Height = height
//platform.Dots = make([][]*Dot, width)
//for i := range platform.Dots {
//platform.Dots[i] = make([]*Dot, height)
//}
//return platform
//}

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

func (p *Platform) Top() int {
	p.rectangle()
	return p.Maxy
}

func (p *Platform) Bottom() int {
	p.rectangle()
	return p.Miny
}

func (p *Platform) Left() int {
	p.rectangle()
	return p.Minx
}

func (p *Platform) Right() int {
	p.rectangle()
	return p.Maxx
}

func (p *Platform) rectangle() {
	if p.Minx == 0 && p.Maxx == 0 && p.Miny == 0 && p.Maxy == 0 {
		inited := false
		for _, row := range p.Dots {
			for _, dot := range row {
				if dot == nil {
					continue
				}
				if !inited {
					p.Minx = dot.PositionX
					p.Maxx = p.Minx
					p.Miny = dot.PositionY
					p.Maxy = p.Miny
					inited = true
				}
				if dot.PositionX < p.Minx {
					p.Minx = dot.PositionX
				}
				if dot.PositionX > p.Maxx {
					p.Maxx = dot.PositionX
				}
				if dot.PositionY < p.Miny {
					p.Miny = dot.PositionY
				}
				if dot.PositionY > p.Maxy {
					p.Maxy = dot.PositionY
				}
			}
		}
	}
}

func (p *Platform) AddBase(x int, y int, base *Base) {
	//fmt.Println(x-50, 50-y, base)

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

// deprecated
func (p *Platform) NextPosition() (int, int, error) {
	for posy, row := range p.Dots {
		for posx, dot := range row {
			if dot == nil {
				continue
			}
			if !dot.Printed {
				//fmt.Println(dot.Base.Name, dot.Printed, dot.PositionX, dot.PositionY)
				return posx, posy, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("all dots are printed")
}

func (p *Platform) NextDot() (int, int, *Dot, error) {
	for y, row := range p.Dots {
		for x, dot := range row {
			if dot == nil {
				continue
			}
			if !dot.Printed {
				return x, y, dot, nil
			}
		}
	}
	return 0, 0, nil, fmt.Errorf("all dots are printed")
}

func (p *Platform) AvailableDots() []*Dot {
	dots := []*Dot{}
	for _, row := range p.Dots {
		for _, dot := range row {
			if dot == nil {
				continue
			}
			if !dot.Printed {
				dots = append(dots, dot)
			}
		}
	}
	return dots
}

func (p *Platform) NextDotVertical() (int, int, int, *Dot, error) {
	var target *Dot
	var targetX int
	var targetY int
	sum := 0
	for y, row := range p.Dots {
		for x, dot := range row {
			if dot == nil {
				continue
			}
			sum += 1
			if !dot.Printed {
				if target == nil {
					target = dot
					targetX = x
					targetY = y
				} else {
					if dot.PositionX < target.PositionX {
						target = dot
						targetX = x
						targetY = y
					}
				}
			}
		}
	}
	if target == nil {
		return sum, targetX, targetY, nil, fmt.Errorf("all dots are printed")
	}
	return sum, targetX, targetY, target, nil
}

func (p *Platform) NextDotPositionY(direction int, posx int, posy int) int {
	var target *Dot
	for _, row := range p.Dots {
		for _, dot := range row {
			if dot == nil {
				continue
			}
			if !dot.Printed {
				//if dot.PositionY == position {
				//continue
				//}
				//if target == nil {
				//target = dot
				//} else {
				//if dot.PositionX < target.PositionX {
				//target = dot
				//}
				//}
				if direction < 0 { // downward
					if dot.PositionY >= posy {
						continue
					}
				}
				if direction > 0 { // upward
					if dot.PositionY <= posy || dot.PositionX < posx {
						continue
					}
				}
				if target == nil {
					target = dot
				} else {
					if dot.PositionX < target.PositionX {
						target = dot
					}
				}

				//if target == nil {
				//if direction < 0 { // downward
				//if dot.PositionY < position {
				//target = dot
				//} else {
				//continue
				//}
				//} else {
				//if dot.PositionY > position {
				//target = dot
				//} else {
				//continue
				//}
				//}
				//} else {
				//if dot.PositionX < target.PositionX {
				//if direction < 0 { // upward
				//if dot.PositionY < position {
				//target = dot
				//} else {
				//continue
				//}
				//} else {
				//if dot.PositionY > position {
				//target = dot
				//} else {
				//continue
				//}
				//}
				//}
				//}
			}
		}
	}
	if target == nil {
		return (50*MM + 1) * direction
	}
	fmt.Println("........NEXT STEP Y", posx, posy, target.Base.Name, target.PositionX, target.PositionY)
	return target.PositionY
}

func (p *Platform) NextDotPositionX() int {
	var target *Dot
	for _, row := range p.Dots {
		for _, dot := range row {
			if dot == nil {
				continue
			}
			if !dot.Printed {
				if target == nil {
					target = dot
				} else {
					if dot.PositionX < target.PositionX {
						target = dot
					}
				}
			}
		}
	}
	if target == nil {
		return 50*MM + 1
	}
	fmt.Println("........NEXT STEP X", target.Base.Name, target.PositionX, target.PositionY)
	return target.PositionX
}

func (p *Platform) NextDotInColumn(x int) (int, int, *Dot) {
	var target *Dot
	var targetX int
	var targetY int
	for y, row := range p.Dots {
		for x, dot := range row {
			if dot == nil {
				continue
			}
			if dot.PositionX == x {
				if !dot.Printed {
					if target == nil {
						target = dot
						targetX = x
						targetY = y
					} else {
						if dot.PositionY > target.PositionY { // most top
							target = dot
							targetX = x
							targetY = y
						}
					}
				}
			}
		}
	}
	return targetX, targetY, target
}

func (p *Platform) PreviousDot(x int, y int) (int, int, *Dot) {
	y--
	if y < 0 {
		return 0, 0, nil
	}
	return x, y, p.Dots[y][x]
}

func (p *Platform) PreviousDotInColumn(x int) *Dot {
	var target *Dot
	for _, row := range p.Dots {
		for _, dot := range row {
			if dot == nil {
				continue
			}
			if dot.PositionX == x {
				if !dot.Printed {
					if target == nil {
						target = dot
					} else {
						if dot.PositionY < target.PositionY { // most bottom
							target = dot
						}
					}
				}
			}
		}
	}
	return target
}

func (p *Platform) NextDotAfter(x int) *Dot {
	var target *Dot
	for _, row := range p.Dots {
		for _, dot := range row {
			if dot == nil {
				continue
			}
			if dot.PositionX > x {
				if target == nil {
					target = dot
				} else {
					if dot.PositionX < target.PositionX { // most bottom
						target = dot
					}
				}
			}
		}
	}
	return target
}

//func (p *Platform) DotsInRow(y int) []*Dot {
//output := []*Dot{}
//for _, column := range p.Dots {
//dot := column[y]
//if dot == nil {
//continue
//}
////fmt.Println(dot.Base.Name, dot.Printed, dot.PositionX, dot.PositionY)
//if !dot.Printed {
//output = append(output, dot)
//}
//}
////fmt.Printf("%#v\n", output)
////fmt.Println(len(output))
//return output
//}

func (p *Platform) DotsInRow(y int) []*Dot {
	output := []*Dot{}
	//fmt.Println(p.Dots[y][:5])
	for _, dot := range p.Dots[y] {
		if dot == nil {
			continue
		}
		//fmt.Println(dot.Base.Name, dot.Printed, dot.PositionX, dot.PositionY)
		if !dot.Printed {
			output = append(output, dot)
		}
	}
	//fmt.Printf("%#v\n", output)
	//fmt.Println(len(output))
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
			if c == *BaseN.Color {
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
