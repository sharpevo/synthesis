package printheads_test

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"reflect"
	"synthesis/internal/dao/printheads"
	"synthesis/internal/platform"
	"testing"
	"time"
)

const (
	ROW_COUNT    = 4
	NOZZLE_COUNT = 1280
	NOZZLE_SPACE = 169.3 * printheads.UM
	ROW_OFFSET   = 84.65 * printheads.UM
	ROW_SPACE_A  = 550.3 * printheads.UM
	ROW_SPACE_B  = 11.811 * printheads.MM
)

var NozzleMap0 = map[int]printheads.Nozzle{
	1: printheads.Nozzle{
		Index: 0,
		Row:   0,
		Position: printheads.Position{
			X: -84650,
			Y: -12361300,
		},
	},
	3: printheads.Nozzle{
		Index: 2,
		Row:   1,
		Position: printheads.Position{
			X: 0,
			Y: -11811000,
		},
	},
	2: printheads.Nozzle{
		Index: 1,
		Row:   2,
		Position: printheads.Position{
			X: -84650,
			Y: -550300,
		},
	},
	4: printheads.Nozzle{
		Index: 3,
		Row:   3,
		Position: printheads.Position{
			X: 0,
			Y: 0,
		},
	},
	1277: printheads.Nozzle{
		Index: 1276,
		Row:   0,
		Position: printheads.Position{
			X: 53922050,
			Y: -12361300,
		},
	},
	1279: printheads.Nozzle{
		Index: 1278,
		Row:   1,
		Position: printheads.Position{
			X: 54006700,
			Y: -11811000,
		},
	},
	1278: printheads.Nozzle{
		Index: 1277,
		Row:   2,
		Position: printheads.Position{
			X: 53922050,
			Y: -550300,
		},
	},
	1280: printheads.Nozzle{
		Index: 1279,
		Row:   3,
		Position: printheads.Position{
			X: 54006700,
			Y: 0,
		},
	},
}

var NozzleMap1 = map[int]printheads.Nozzle{
	1: printheads.Nozzle{
		Index: 0,
		Row:   0,
		Position: printheads.Position{
			X: 27088000,
			Y: -6180650,
		},
	},
	3: printheads.Nozzle{
		Index: 2,
		Row:   1,
		Position: printheads.Position{
			X: 27172650,
			Y: -5630350,
		},
	},
	2: printheads.Nozzle{
		Index: 1,
		Row:   2,
		Position: printheads.Position{
			X: 27088000,
			Y: 5630350,
		},
	},
	4: printheads.Nozzle{
		Index: 3,
		Row:   3,
		Position: printheads.Position{
			X: 27172650,
			Y: 6180650,
		},
	},
	1277: printheads.Nozzle{
		Index: 1276,
		Row:   0,
		Position: printheads.Position{
			X: 81094700,
			Y: -6180650,
		},
	},
	1279: printheads.Nozzle{
		Index: 1278,
		Row:   1,
		Position: printheads.Position{
			X: 81179350,
			Y: -5630350,
		},
	},
	1278: printheads.Nozzle{
		Index: 1277,
		Row:   2,
		Position: printheads.Position{
			X: 81094700,
			Y: 5630350,
		},
	},
	1280: printheads.Nozzle{
		Index: 1279,
		Row:   3,
		Position: printheads.Position{
			X: 81179350,
			Y: 6180650,
		},
	},
}

func TestOne(t *testing.T) {
	fmt.Println(printheads.SPACE_ROW_B)
	fmt.Println(printheads.SPACE_NOZZLE_PITCH)
}

func TestPosition(t *testing.T) {
	index := []int{
		0, 1, 2, 3,
		4, 5, 6, 7,
		8, 9, 10, 11,
		12, 13, 14, 15,
		16, 17, 18, 19,
		1276, 1277, 1278, 1279,
	}
	row := []int{
		0, 2, 1, 3,
		0, 2, 1, 3,
		0, 2, 1, 3,
		0, 2, 1, 3,
		0, 2, 1, 3,
		0, 2, 1, 3,
	}

	for k, v := range index {
		n, err := printheads.NewNozzle(v)
		if err != nil || n.Row != row[k] {
			t.Errorf(
				"\n%d: EXPECT: %v\nGET: %v\n",
				k,
				row[k],
				n.Row,
			)
		}
	}
}

func TestPrintHeadLineD(t *testing.T) {
	h, _ := printheads.NewPrintHeadLineD(
		ROW_COUNT,
		NOZZLE_COUNT,
		NOZZLE_SPACE,
		ROW_OFFSET,
		ROW_SPACE_A,
		ROW_SPACE_B,
		0,
		0,
	)
	for _, row := range h.Rows {
		for _, nozzle := range row.Nozzles {
			//fmt.Printf("%#v\n", nozzle)
			if nozzle.Row != row.Index {
				t.Errorf(
					"\nEXPECT: %v\nGET: %v\n",
					row.Index,
					nozzle.Row,
				)
			}
			if n, ok := NozzleMap0[nozzle.Index+1]; ok {
				if !reflect.DeepEqual(n, *nozzle) {
					t.Errorf(
						"\n%d EXPECT: %#v\nGET: %#v\n",
						nozzle.Index+1,
						n,
						*nozzle,
					)
				}
			}
		}
	}
}

func TestPrintHead(t *testing.T) {
	h, _ := printheads.NewPrintHead(
		ROW_COUNT,
		NOZZLE_COUNT,
		NOZZLE_SPACE,
		ROW_OFFSET,
		ROW_SPACE_A,
		ROW_SPACE_B,
	)
	for _, v := range NozzleMap1 {
		fmt.Printf("%s\n", v.String())
	}
	for _, row := range h.Rows {
		for _, nozzle := range row.Nozzles {
			//fmt.Printf("%#v\n", nozzle)
			if nozzle.Row != row.Index {
				t.Errorf(
					"\nEXPECT: %v\nGET: %v\n",
					row.Index,
					nozzle.Row,
				)
			}
			if n, ok := NozzleMap1[nozzle.Index+1]; ok {
				if !reflect.DeepEqual(n, *nozzle) {
					t.Errorf(
						"\n%d EXPECT: %#v\nGET: %#v\n",
						nozzle.Index+1,
						n,
						*nozzle,
					)
				}
			}
		}
	}
}

func TestUpdatePosition(t *testing.T) {
	zero := printheads.Position{
		X: 0,
		Y: 0,
	}
	test := printheads.Position{
		X: 27172650,
		Y: 6180650,
	}
	h, _ := printheads.NewPrintHeadLineD(
		ROW_COUNT,
		NOZZLE_COUNT,
		NOZZLE_SPACE,
		ROW_OFFSET,
		ROW_SPACE_A,
		ROW_SPACE_B,
		zero.X,
		zero.Y,
	)
	expected := NozzleMap0[1]
	actual := *h.Rows[0].Nozzles[0]
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"\nEXPECT: %#v\nGET: %#v\n",
			expected,
			actual,
		)
	}
	h.UpdatePosition(test.X, test.Y)
	expected = NozzleMap1[1]
	actual = *h.Rows[0].Nozzles[0]
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"\nEXPECT: %#v\nGET: %#v\n",
			expected,
			actual,
		)
	}
	h.UpdatePosition(zero.X, zero.Y)
	expected = NozzleMap0[1]
	actual = *h.Rows[0].Nozzles[0]
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"\nEXPECT: %#v\nGET: %#v\n",
			expected,
			actual,
		)
	}
}

func TestPrint(t *testing.T) {

	fmt.Println("====")
	fmt.Println()
	fmt.Println()
	p := platform.NewPlatform(100, 100)
	myImage := image.NewRGBA(image.Rect(0, 0, p.Width, p.Height))
	block1 := &platform.Block{}
	block1.PositionX = 25
	block1.PositionY = 25
	block1.SpaceX = 2
	block1.SpaceY = 5
	block1.AddRow("TTTTTCTGGA")
	block1.AddRow("AGGTGCGTGT")
	block1.AddRow("GGAGGGAATG")
	block1.AddRow("CTGTGCGTGA")
	minWidth := 10 + block1.SpaceX*(10-1) + block1.PositionX
	minHeight := 4 + block1.SpaceY*(4-1) + block1.PositionY
	fmt.Println("min platform: ", minWidth, minHeight)
	p.AddBlock(block1)
	for posy, row := range p.Dots {
		for posx, dot := range row {
			if dot == nil {
				continue
			}
			myImage.Set(posx, posy, dot.Base.Color)
		}
	}
	outputFile, _ := os.Create("test.png")
	png.Encode(outputFile, myImage)
	outputFile.Close()

	// printhead
	//h, _ := printheads.NewPrintHead(
	//ROW_COUNT,
	//NOZZLE_COUNT,
	//NOZZLE_SPACE,
	//ROW_OFFSET,
	//ROW_SPACE_A,
	//ROW_SPACE_B,
	//)
	h, _ := printheads.NewPrintHeadLineD(
		ROW_COUNT,
		NOZZLE_COUNT,
		NOZZLE_SPACE,
		ROW_OFFSET,
		ROW_SPACE_A,
		ROW_SPACE_B,
		0,
		0,
	)
	for _, v := range NozzleMap1 {
		fmt.Printf("%s\n", v.String())
	}
	time.Sleep(time.Second)

	// read
	pf, _ := platform.ParsePlatform("test.png")

	// next dots
	px, py, err := pf.NextPosition()
	if err != nil {
		fmt.Println(err)
	}
	hx, hy := ToPrintheadPosition(px, py)
	h.UpdatePositionStar(hx*printheads.MM, hy*printheads.MM)

	// loop horizontally, from left to right
	// 1. offset, upward
	// 2. next, downward
	imageIndex := 0
	img := image.NewRGBA(image.Rect(0, 0, pf.Width, pf.Height))
	for h.Rows[3].Nozzles[0].X < 50*printheads.MM {

		// loop vertically, from printhead bottom to printhead top
		// downward
		fmt.Println(">>>downward")
		for dposy := hy * printheads.MM; h.Rows[3].Nozzles[0].Y > -50*printheads.MM; dposy -= h.RowOffset {
			genData(h, pf, py, &imageIndex, img)
			h.UpdatePositionStar(hx*printheads.MM, dposy)
		}

		dposx := hx*printheads.MM + h.RowOffset
		h.UpdatePositionStar(dposx, h.Rows[0].Nozzles[0].Y)

		// upward
		fmt.Println(">>>upword")
		for dposy := h.Rows[0].Nozzles[0].Y; h.Rows[0].Nozzles[0].Y < 50*printheads.MM; dposy += h.RowOffset {
			genData(h, pf, py, &imageIndex, img)
			h.UpdatePositionStar(dposx, dposy)
		}

		px, py, err = pf.NextPosition()
		if err != nil {
			break
		}
		hx, hy = ToPrintheadPosition(px, py)
		//fmt.Println(hx, hy, h.Rows[3].Nozzles[0].X)
		h.UpdatePositionStar(hx*printheads.MM, hy*printheads.MM)
	}
}

func printheadPosition(h *printheads.PrintHead) {
	fmt.Println("printhead left bottom: ", h.Rows[0].X, h.Rows[0].Y)
}

func moveAbs(x int, y int) {
	hx, hy := ToPrintheadPosition(x, y)
	fmt.Printf("abs move dot (%v, %v), printhead (%v, %v)\n", x, y, hx, hy)
}

func ToPrintheadPosition(x int, y int) (int, int) {
	return x - 50, 50 - y
}
func ToPlatformPosition(x int, y int) (int, int) {
	return x + 50, 50 - y
}

func genData(h *printheads.PrintHead, pf *platform.Platform, py int, imageIndex *int, img *image.RGBA) []int {
	data := make([]int, 1280)

	printable := false
	// traverse nozzles
	for _, row := range h.Rows {
		for _, nozzle := range row.Nozzles {

			// check available nozzles
			for _, dot := range pf.DotsInRow(py) {
				dotx, doty := ToPrintheadPosition(dot.PositionX, dot.PositionY)
				if math.Abs(float64(nozzle.X-dotx*printheads.MM)) < float64(h.RowOffset) &&
					math.Abs(float64(nozzle.Y-doty*printheads.MM)) < float64(h.RowOffset) {
					if (dot.Base.Name == "A" && row.Index == 0) ||
						(dot.Base.Name == "C" && row.Index == 1) ||
						(dot.Base.Name == "G" && row.Index == 2) ||
						(dot.Base.Name == "T" && row.Index == 3) {
						dot.Printed = true
						img.Set(dot.PositionX, dot.PositionY, dot.Base.Color)
						fmt.Println(dot.Base.Name, nozzle, " || ", dot, " >> ", dotx, doty)
						data[nozzle.Index] = int(dot.Base.Color.A)
						printable = true
					}
				}
			}

		}
	}
	if printable {
		fileName := fmt.Sprintf("output/%02d.png", *imageIndex)
		outputFile, _ := os.Create(fileName)
		png.Encode(outputFile, img)
		outputFile.Close()
		*imageIndex = *imageIndex + 1
	}
	return data
}
