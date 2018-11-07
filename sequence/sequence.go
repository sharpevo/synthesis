package sequence

import (
	"bytes"
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"image"
	"image/color"
	"image/png"
	"posam/gui/uiutil"
	"posam/util/platform"
	"strconv"
	"strings"
)

var imageItem *widgets.QGraphicsPixmapItem

func NewSequence() *widgets.QGroupBox {

	group := widgets.NewQGroupBox2("sequence", nil)
	layout := widgets.NewQGridLayout2()

	scene := widgets.NewQGraphicsScene(nil)
	view := widgets.NewQGraphicsView(nil)
	img := gui.NewQImage()
	imageItem = widgets.NewQGraphicsPixmapItem2(gui.NewQPixmap().FromImage(img, 0), nil)
	scene.AddItem(imageItem)
	view.SetScene(scene)
	view.Show()

	layout.AddWidget(view, 0, 0, 0)
	layout.AddWidget(NewSequenceDetail(), 0, 1, 0)

	layout.SetColumnStretch(0, 1)
	layout.SetColumnStretch(1, 1)

	group.SetLayout(layout)
	return group
}

func NewSequenceDetail() *widgets.QGroupBox {
	group := widgets.NewQGroupBox2("Config", nil)
	layout := widgets.NewQGridLayout2()

	startxLabel := widgets.NewQLabel2("Starts at x:", nil, 0)
	startxInput := widgets.NewQLineEdit(nil)
	startyLabel := widgets.NewQLabel2("Starts at y:", nil, 0)
	startyInput := widgets.NewQLineEdit(nil)

	spacexLabel := widgets.NewQLabel2("Drop space x:", nil, 0)
	spacexInput := widgets.NewQLineEdit(nil)
	spaceyLabel := widgets.NewQLabel2("Drop space y:", nil, 0)
	spaceyInput := widgets.NewQLineEdit(nil)

	spaceBlockxLabel := widgets.NewQLabel2("Block space x:", nil, 0)
	spaceBlockxInput := widgets.NewQLineEdit(nil)
	spaceBlockyLabel := widgets.NewQLabel2("Block space y:", nil, 0)
	spaceBlockyInput := widgets.NewQLineEdit(nil)
	spaceSlideyLabel := widgets.NewQLabel2("Slide space y:", nil, 0)
	spaceSlideyInput := widgets.NewQLineEdit(nil)

	sequenceInput := widgets.NewQTextEdit(nil)

	startxInput.SetText("25")
	startyInput.SetText("25")
	spacexInput.SetText("2")
	spaceyInput.SetText("5")
	spaceBlockxInput.SetText("20")
	spaceBlockyInput.SetText("30")
	spaceSlideyInput.SetText("50")
	sequenceInput.SetText(`AGGTGCGTGT,TGAATCATTG,AGGTGCGTGT
GGAGGGAATG,CTAGTACTTT,GGAGGGAATG
CTGTGCGTGA,ACACCCTTGG,CTGTGCGTGA

TATAGCCTAC,GTACTCGTAG,TATAGCCTAC
ACACATACGG,ACTCGACTGA,ACACATACGG
GTCAGCATAC,AAGCTTGTTC,GTCAGCATAC
CATACGCAGC,TGTACATGAC,CATACGCAGC

GTACTCGTAG,TATAGCCTAC,TATAGCCTAC
ACTCGACTGA,ACACATACGG,ACACATACGG
AAGCTTGTTC,GTCAGCATAC,GTCAGCATAC
TGTACATGAC,CATACGCAGC,CATACGCAGC


GGAGGGAATG,CTAGTACTTT,GGAGGGAATG

AAGCTTGTTC,GTCAGCATAC,GTCAGCATAC
GTCAGCATAC,AAGCTTGTTC,GTCAGCATAC
`)
	generateButton := widgets.NewQPushButton2("GENERATE", nil)
	generateButton.ConnectClicked(func(bool) {
		startxInt, err := strconv.Atoi(startxInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		startyInt, err := strconv.Atoi(startyInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		spacexInt, err := strconv.Atoi(spacexInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		spaceyInt, err := strconv.Atoi(spaceyInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		spaceBlockxInt, err := strconv.Atoi(spaceBlockxInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		spaceBlockyInt, err := strconv.Atoi(spaceBlockyInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		spaceSlideyInt, err := strconv.Atoi(spaceSlideyInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		generateImage(
			startxInt,
			startyInt,
			spacexInt,
			spaceyInt,
			spaceBlockxInt,
			spaceBlockyInt,
			spaceSlideyInt,
			sequenceInput.ToPlainText(),
		)
	})

	layout.AddWidget(startxLabel, 0, 0, 0)
	layout.AddWidget(startxInput, 0, 1, 0)
	layout.AddWidget(startyLabel, 1, 0, 0)
	layout.AddWidget(startyInput, 1, 1, 0)

	layout.AddWidget(spacexLabel, 2, 0, 0)
	layout.AddWidget(spacexInput, 2, 1, 0)
	layout.AddWidget(spaceyLabel, 3, 0, 0)
	layout.AddWidget(spaceyInput, 3, 1, 0)

	layout.AddWidget(spaceBlockxLabel, 4, 0, 0)
	layout.AddWidget(spaceBlockxInput, 4, 1, 0)
	layout.AddWidget(spaceBlockyLabel, 5, 0, 0)
	layout.AddWidget(spaceBlockyInput, 5, 1, 0)
	layout.AddWidget(spaceSlideyLabel, 6, 0, 0)
	layout.AddWidget(spaceSlideyInput, 6, 1, 0)

	layout.AddWidget3(sequenceInput, 7, 0, 1, 2, 0)
	layout.AddWidget3(generateButton, 8, 0, 1, 2, 0)
	group.SetLayout(layout)
	return group
}

func generateImage(
	startX int,
	startY int,
	spaceX int,
	spaceY int,
	spaceBlockx int,
	spaceBlocky int,
	//spaceSlidex int, // use spaceBlockx
	spaceSlidey int,
	sequences string,
) {
	yoffset := startY
	xoffset := startX
	width := 0
	height := 0
	pixels := make(map[int]map[int]*color.NRGBA)

	count := 0
	for _, line := range strings.Split(sequences, "\n") {
		if line == "" {
			count += 1
			continue
		}
		switch count {
		case 0:
		case 1:
			if yoffset != startY {
				yoffset += spaceBlocky - spaceY
			} else {
				yoffset += spaceBlocky
			}
		case 2:
			if yoffset != startY {
				yoffset += spaceSlidey - spaceY
			} else {
				yoffset += spaceSlidey
			}
		default:
			uiutil.MessageBoxError("invalid sequences")
			return
		}
		if _, ok := pixels[yoffset]; !ok {
			pixels[yoffset] = make(map[int]*color.NRGBA)
		}
		xoffset = startX
		for _, seq := range strings.Split(line, ",") {
			for _, base := range strings.Split(strings.Trim(seq, " "), "") {
				pixels[yoffset][xoffset] = ToColor(base)
				xoffset += spaceX
				if xoffset > width {
					width = xoffset
				}
			}
			xoffset += spaceBlockx - spaceX
		}
		yoffset += spaceY
		if yoffset > height {
			height = yoffset
		}
		count = 0
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y, pixel := range pixels {
		for x, c := range pixel {
			img.Set(x, y, *c)
		}
	}
	var imagebuff bytes.Buffer
	png.Encode(&imagebuff, img)
	imagebyte := imagebuff.Bytes()
	qimg := gui.NewQImage()
	qimg.LoadFromData2(core.NewQByteArray2(string(imagebyte), len(imagebyte)), "png")
	qimg = qimg.Scaled2(5*width, 5*height, core.Qt__IgnoreAspectRatio, core.Qt__FastTransformation)
	imageItem.SetPixmap(gui.NewQPixmap().FromImage(qimg, 0))
}

func ToColor(base string) *color.NRGBA {
	switch base {
	case "A":
		return platform.BaseA.Color
	case "C":
		return platform.BaseC.Color
	case "G":
		return platform.BaseG.Color
	case "T":
		return platform.BaseT.Color
	default:
		return platform.BaseN.Color
	}
}

func generateImage2(
	startX int,
	startY int,
	spaceX int,
	spaceY int,
	spaceBlockx int,
	spaceBlocky int,
	//spaceSlidex int, // use spaceBlockx
	spaceSlidey int,
	sequences string,
) {

	width := 0
	height := 0
	blocks := map[int]map[int]map[int]*platform.Block{}
	yoffset := startY
	for z, slide := range strings.Split(sequences, "\n\n\n") { // slide
		if _, ok := blocks[z]; !ok {
			blocks[z] = make(map[int]map[int]*platform.Block)
		}

		for y, block := range strings.Split(slide, "\n\n") { // block
			if _, ok := blocks[z][y]; !ok {
				blocks[z][y] = make(map[int]*platform.Block)
			}

			for _, line := range strings.Split(block, "\n") { // line

				xoffset := startX - spaceBlockx
				for x, seq := range strings.Split(line, ",") { // seq
					if _, ok := blocks[z][y][x]; !ok {
						xoffset += spaceBlockx
						b := &platform.Block{}
						b.PositionX = xoffset
						b.PositionY = yoffset
						b.SpaceX = spaceX
						b.SpaceY = spaceY
						blocks[z][y][x] = b
					}
					blocks[z][y][x].AddRow(seq)
					xoffset += len(blocks[z][y][x].Sequence[0]) * (spaceX + 1)
				}
				if xoffset > width {
					width = xoffset
				}
			}
			yoffset += len(blocks[z][y][0].Sequence)*(spaceY+1) + spaceBlocky
		}
		yoffset += spaceSlidey
	}
	height = yoffset

	fmt.Println(width, height)
	p := platform.NewPlatform(width, height)
	for _, slide := range blocks {
		for _, row := range slide {
			for _, block := range row {
				fmt.Println("block", block.PositionX, block.PositionY)
				p.AddBlock(block)
			}
		}
	}
	img := image.NewRGBA(image.Rect(0, 0, p.Width, p.Height))
	for posy, row := range p.Dots {
		for posx, dot := range row {
			if dot == nil {
				continue
			}
			img.Set(posx, posy, dot.Base.Color)
		}
	}
	var imagebuff bytes.Buffer
	png.Encode(&imagebuff, img)
	imagebyte := imagebuff.Bytes()
	qimg := gui.NewQImage()
	qimg.LoadFromData2(core.NewQByteArray2(string(imagebyte), len(imagebyte)), "png")
	qimg = qimg.Scaled2(5*p.Width, 5*p.Height, core.Qt__IgnoreAspectRatio, core.Qt__FastTransformation)
	imageItem.SetPixmap(gui.NewQPixmap().FromImage(qimg, 0))
	fmt.Println(startX, startY, spaceX, spaceY, spaceBlockx, spaceBlocky)
}
