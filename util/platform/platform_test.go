package platform_test

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"posam/util/platform"
	"testing"
)

func TestWrite(t *testing.T) {
	p := platform.NewPlatform(39, 41)
	myImage := image.NewRGBA(image.Rect(0, 0, p.Width, p.Height))
	block1 := &platform.Block{}
	block1.PositionX = 11
	block1.PositionY = 22
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

	// read
	platform.ParsePlatform("test.png")
	//fmt.Println(p)
	//existingImageFile, err := os.Open("test.png")
	//if err != nil {
	//// Handle error
	//}
	//defer existingImageFile.Close()
	//img, err := png.Decode(existingImageFile)
	//if err != nil {
	//// Handle error
	//}
	//fmt.Println("----")
	//width := img.Bounds().Max.X
	//height := img.Bounds().Max.Y
	//p = platform.NewPlatform(width, height)
	//for y := 0; y < height; y++ {
	//for x := 0; x < width; x++ {
	//c := img.At(x, y).(color.NRGBA)
	//if c == platform.BaseN.Color {
	//continue
	//}
	//p.AddBase(x, y, platform.ColorToBase(&c))
	//}
	//}
}
