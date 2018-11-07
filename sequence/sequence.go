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
	"os"
	"posam/gui/uiutil"
	"posam/util/platform"
	//"posam/util/printheads"
	"strconv"
	"strings"
)

var imageItem *widgets.QGraphicsPixmapItem

func NewSequence() *widgets.QGroupBox {

	group := widgets.NewQGroupBox2("sequence", nil)
	layout := widgets.NewQGridLayout2()

	viewGroup := widgets.NewQGroupBox2("sequence", nil)
	viewLayout := widgets.NewQGridLayout2()
	scene := widgets.NewQGraphicsScene(nil)
	view := widgets.NewQGraphicsView(nil)
	qimg := gui.NewQImage()
	imageItem = widgets.NewQGraphicsPixmapItem2(gui.NewQPixmap().FromImage(qimg, 0), nil)
	scene.AddItem(imageItem)
	view.SetScene(scene)
	view.Show()

	exportButton := widgets.NewQPushButton2("EXPORT", nil)
	exportButton.ConnectClicked(func(bool) {
	})

	viewLayout.AddWidget(view, 0, 0, 0)
	viewLayout.AddWidget(exportButton, 1, 0, 0)
	viewGroup.SetLayout(viewLayout)

	layout.AddWidget(viewGroup, 0, 0, 0)
	layout.AddWidget(NewSequenceDetail(), 0, 1, 0)

	layout.SetColumnStretch(0, 1)
	layout.SetColumnStretch(1, 1)

	group.SetLayout(layout)
	return group
}

func NewSequenceDetail() *widgets.QGroupBox {
	group := widgets.NewQGroupBox2("Config", nil)
	layout := widgets.NewQGridLayout2()

	resolutionLabel := widgets.NewQLabel2("Resolution:", nil, 0)
	resolutionInput := widgets.NewQLineEdit(nil)
	startxLabel := widgets.NewQLabel2("Starts at x:", nil, 0)
	startxInput := widgets.NewQLineEdit(nil)
	startyLabel := widgets.NewQLabel2("Starts at y:", nil, 0)
	startyInput := widgets.NewQLineEdit(nil)

	spacexLabel := widgets.NewQLabel2("Spot space x:", nil, 0)
	spacexInput := widgets.NewQLineEdit(nil)
	spaceyLabel := widgets.NewQLabel2("Spot space y:", nil, 0)
	spaceyInput := widgets.NewQLineEdit(nil)

	spaceBlockxLabel := widgets.NewQLabel2("Block space x:", nil, 0)
	spaceBlockxInput := widgets.NewQLineEdit(nil)
	spaceBlockyLabel := widgets.NewQLabel2("Block space y:", nil, 0)
	spaceBlockyInput := widgets.NewQLineEdit(nil)
	spaceSlideyLabel := widgets.NewQLabel2("Slide space y:", nil, 0)
	spaceSlideyInput := widgets.NewQLineEdit(nil)

	sequenceInput := widgets.NewQTextEdit(nil)

	resolutionInput.SetText("84.65")
	startxInput.SetText("-50000")
	startyInput.SetText("50000")
	spacexInput.SetText("169.3")
	spaceyInput.SetText("84.65")
	spaceBlockxInput.SetText("1000")
	spaceBlockyInput.SetText("42.65")
	spaceSlideyInput.SetText("100")
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
		resolutionFloat, err := strconv.ParseFloat(resolutionInput.Text(), 32)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		if resolutionFloat < 12.0 {
			uiutil.MessageBoxError("invaild resolution")
			return
		}
		resolutionInt := int(resolutionFloat * platform.UM)
		maxWidth := 100 * platform.MM / resolutionInt
		maxHeight := 100 * platform.MM / resolutionInt
		startxIntRaw, err := parseArg(startxInput.Text(), resolutionInt)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		startxInt := startxIntRaw + 50*platform.MM/resolutionInt
		startyIntRaw, err := parseArg(startyInput.Text(), resolutionInt)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		startyInt := 50*platform.MM/resolutionInt - startyIntRaw
		spacexInt, err := parseArg(spacexInput.Text(), resolutionInt)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		spaceyInt, err := parseArg(spaceyInput.Text(), resolutionInt)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		spaceBlockxInt, err := parseArg(spaceBlockxInput.Text(), resolutionInt)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		spaceBlockyInt, err := parseArg(spaceBlockyInput.Text(), resolutionInt)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		spaceSlideyInt, err := parseArg(spaceSlideyInput.Text(), resolutionInt)
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
			maxWidth,
			maxHeight,
		)
	})

	layout.AddWidget(resolutionLabel, 0, 0, 0)
	layout.AddWidget(resolutionInput, 0, 1, 0)
	layout.AddWidget(startxLabel, 1, 0, 0)
	layout.AddWidget(startxInput, 1, 1, 0)
	layout.AddWidget(startyLabel, 2, 0, 0)
	layout.AddWidget(startyInput, 2, 1, 0)

	layout.AddWidget(spacexLabel, 3, 0, 0)
	layout.AddWidget(spacexInput, 3, 1, 0)
	layout.AddWidget(spaceyLabel, 4, 0, 0)
	layout.AddWidget(spaceyInput, 4, 1, 0)

	layout.AddWidget(spaceBlockxLabel, 5, 0, 0)
	layout.AddWidget(spaceBlockxInput, 5, 1, 0)
	layout.AddWidget(spaceBlockyLabel, 6, 0, 0)
	layout.AddWidget(spaceBlockyInput, 6, 1, 0)
	layout.AddWidget(spaceSlideyLabel, 7, 0, 0)
	layout.AddWidget(spaceSlideyInput, 7, 1, 0)

	layout.AddWidget3(sequenceInput, 8, 0, 1, 2, 0)
	layout.AddWidget3(generateButton, 9, 0, 1, 2, 0)
	group.SetLayout(layout)
	return group
}

func parseArg(argString string, resolution int) (int, error) {
	argFloat, err := strconv.ParseFloat(argString, 32)
	if err != nil {
		return 0, err
	}
	return int(argFloat*platform.UM) / resolution, nil
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
	maxWidth int,
	maxHeight int,
) {
	xoffset := startX
	yoffset := startY
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
	if width > maxWidth || height > maxHeight {
		uiutil.MessageBoxError(fmt.Sprintf("invalid size: %d x %d (%d x %d)", width, height, maxWidth, maxHeight))
		return
	}
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y, pixel := range pixels {
		for x, c := range pixel {
			img.Set(x, y, *c)
		}
	}
	// with file
	outputFile, _ := os.Create("test.png")
	png.Encode(outputFile, img)
	outputFile.Close()
	// nofile
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
