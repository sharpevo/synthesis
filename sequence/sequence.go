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
	"math"
	"os"
	"posam/dao/printheads"
	"posam/gui/uiutil"
	"posam/util/platform"
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

	viewLayout.AddWidget(view, 0, 0, 0)
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
	exportButton := widgets.NewQPushButton2("EXPORT", nil)
	exportButton.ConnectClicked(func(bool) {
		_,
			startxFloat,
			startyFloat,
			spacexFloat,
			spaceyFloat,
			spaceBlockxFloat,
			spaceBlockyFloat,
			spaceSlideyFloat,
			err := parseFloatArg(
			resolutionInput.Text(),
			startxInput.Text(),
			startyInput.Text(),
			spacexInput.Text(),
			spaceyInput.Text(),
			spaceBlockxInput.Text(),
			spaceBlockyInput.Text(),
			spaceSlideyInput.Text(),
		)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		export(
			int(startxFloat*platform.UM+50*platform.MM),
			int(50*platform.MM-startyFloat*platform.UM),
			int(spacexFloat*platform.UM),
			int(spaceyFloat*platform.UM),
			int(spaceBlockxFloat*platform.UM),
			int(spaceBlockyFloat*platform.UM),
			int(spaceSlideyFloat*platform.UM),
			sequenceInput.ToPlainText(),
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
	layout.AddWidget3(generateButton, 9, 0, 1, 1, 0)
	layout.AddWidget3(exportButton, 9, 1, 1, 1, 0)
	layout.SetColumnStretch(0, 1)
	layout.SetColumnStretch(1, 1)

	group.SetLayout(layout)
	return group
}

func parseFloatArg(
	resolution string,
	startx string,
	starty string,
	spacex string,
	spacey string,
	spaceBlockx string,
	spaceBlocky string,
	spaceSlidey string,
) (
	resolutionFloat float64,
	startxFloat float64,
	startyFloat float64,
	spacexFloat float64,
	spaceyFloat float64,
	spaceBlockxFloat float64,
	spaceBlockyFloat float64,
	spaceSlideyFloat float64,
	err error,
) {
	resolutionFloat, err = strconv.ParseFloat(resolution, 32)
	if err != nil {
		return resolutionFloat,
			startxFloat,
			startyFloat,
			spacexFloat,
			spaceyFloat,
			spaceBlockxFloat,
			spaceBlockyFloat,
			spaceSlideyFloat,
			err
	}
	startxFloat, err = strconv.ParseFloat(startx, 32)
	if err != nil {
		return resolutionFloat,
			startxFloat,
			startyFloat,
			spacexFloat,
			spaceyFloat,
			spaceBlockxFloat,
			spaceBlockyFloat,
			spaceSlideyFloat,
			err
	}
	startyFloat, err = strconv.ParseFloat(starty, 32)
	if err != nil {
		return resolutionFloat,
			startxFloat,
			startyFloat,
			spacexFloat,
			spaceyFloat,
			spaceBlockxFloat,
			spaceBlockyFloat,
			spaceSlideyFloat,
			err
	}
	spacexFloat, err = strconv.ParseFloat(spacex, 32)
	if err != nil {
		return resolutionFloat,
			startxFloat,
			startyFloat,
			spacexFloat,
			spaceyFloat,
			spaceBlockxFloat,
			spaceBlockyFloat,
			spaceSlideyFloat,
			err
	}
	spaceyFloat, err = strconv.ParseFloat(spacey, 32)
	if err != nil {
		return resolutionFloat,
			startxFloat,
			startyFloat,
			spacexFloat,
			spaceyFloat,
			spaceBlockxFloat,
			spaceBlockyFloat,
			spaceSlideyFloat,
			err
	}
	spaceBlockxFloat, err = strconv.ParseFloat(spaceBlockx, 32)
	if err != nil {
		return resolutionFloat,
			startxFloat,
			startyFloat,
			spacexFloat,
			spaceyFloat,
			spaceBlockxFloat,
			spaceBlockyFloat,
			spaceSlideyFloat,
			err
	}
	spaceBlockyFloat, err = strconv.ParseFloat(spaceBlocky, 32)
	if err != nil {
		return resolutionFloat,
			startxFloat,
			startyFloat,
			spacexFloat,
			spaceyFloat,
			spaceBlockxFloat,
			spaceBlockyFloat,
			spaceSlideyFloat,
			err
	}
	spaceSlideyFloat, err = strconv.ParseFloat(spaceSlidey, 32)
	if err != nil {
		return resolutionFloat,
			startxFloat,
			startyFloat,
			spacexFloat,
			spaceyFloat,
			spaceBlockxFloat,
			spaceBlockyFloat,
			spaceSlideyFloat,
			err
	}
	return resolutionFloat,
		startxFloat,
		startyFloat,
		spacexFloat,
		spaceyFloat,
		spaceBlockxFloat,
		spaceBlockyFloat,
		spaceSlideyFloat,
		nil
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
	fmt.Println("start", startX, startY)
	fmt.Println("space spot", spaceX, spaceY)
	fmt.Println("space block", spaceBlockx, spaceBlocky)
	fmt.Println("space slide", spaceSlidey)
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

func ToBase(base string) *platform.Base {
	switch base {
	case "A":
		return platform.BaseA
	case "C":
		return platform.BaseC
	case "G":
		return platform.BaseG
	case "T":
		return platform.BaseT
	default:
		return platform.BaseN
	}
}

func export(
	startX int,
	startY int,
	spaceX int,
	spaceY int,
	spaceBlockx int,
	spaceBlocky int,
	spaceSlidey int,
	sequences string,
) {
	//filePath, err := uiutil.FilePath()
	//if err != nil {
	//uiutil.MessageBoxError(err.Error())
	//return
	//}

	dots := [][]*platform.Dot{}

	rowCount := 0
	columnCount := 0
	xoffset := startX
	yoffset := startY
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
		dots = append(dots, []*platform.Dot{})
		xoffset = startX
		baseCount := 0
		for _, seq := range strings.Split(line, ",") {
			for _, base := range strings.Split(strings.Trim(seq, " "), "") {
				fmt.Println("location", xoffset, yoffset)
				baseCount += 1
				dots[rowCount] = append(dots[rowCount], &platform.Dot{
					platform.NewBase(base),
					false,
					xoffset,
					yoffset,
				})
				xoffset += spaceX
			}
			xoffset += spaceBlockx - spaceX
		}
		yoffset += spaceY
		if baseCount > columnCount {
			columnCount = baseCount
		}
		rowCount += 1
		count = 0
	}

	pf := platform.NewPlatform(columnCount+1, rowCount+1)
	for y, row := range dots {
		for x, dot := range row {
			pf.Dots[y][x] = dot
		}
	}
	fmt.Println(pf.Dots[0])

	h, _ := printheads.NewPrintHeadLineD(
		4,
		1280,
		169.3*platform.UM,
		84.65*platform.UM,
		550.3*platform.UM,
		11.811*platform.MM,
		0,
		0,
	)

	_, py, dot, err := pf.NextDot()
	if err != nil {
		fmt.Println(err)
	}
	h.UpdatePositionStar(dot.PositionX, dot.PositionX)
	for h.Rows[3].Nozzles[0].X < 50*printheads.MM {

		// loop vertically, from printhead bottom to printhead top
		// downward
		fmt.Println(">>>downward")
		for dposy := dot.PositionY; h.Rows[3].Nozzles[0].Y > -50*printheads.MM; dposy -= h.RowOffset {
			genData(h, pf, py)
			h.UpdatePositionStar(dot.PositionX, dposy)
		}

		dposx := dot.PositionX + h.RowOffset
		h.UpdatePositionStar(dposx, h.Rows[0].Nozzles[0].Y)

		// upward
		fmt.Println(">>>upword")
		for dposy := h.Rows[0].Nozzles[0].Y; h.Rows[0].Nozzles[0].Y < 50*printheads.MM; dposy += h.RowOffset {
			genData(h, pf, py)
			h.UpdatePositionStar(dposx, dposy)
		}

		_, py, dot, err = pf.NextDot()
		if err != nil {
			break
		}
		h.UpdatePositionStar(dot.PositionX, dot.PositionY)
	}
}

func genData(h *printheads.PrintHead, pf *platform.Platform, py int) []int {
	data := make([]int, 1280)

	//printable := false
	// traverse nozzles
	for _, row := range h.Rows {
		for _, nozzle := range row.Nozzles {

			// check available nozzles
			for _, dot := range pf.DotsInRow(py) {
				dotx, doty := dot.PositionX, dot.PositionY
				if math.Abs(float64(nozzle.X-dotx)) < float64(h.RowOffset) &&
					math.Abs(float64(nozzle.Y-doty)) < float64(h.RowOffset) {
					if (dot.Base.Name == "A" && row.Index == 0) ||
						(dot.Base.Name == "C" && row.Index == 1) ||
						(dot.Base.Name == "G" && row.Index == 2) ||
						(dot.Base.Name == "T" && row.Index == 3) {
						dot.Printed = true
						//img.Set(dot.PositionX, dot.PositionY, dot.Base.Color)
						fmt.Println(dot.Base.Name, nozzle, " || ", dot, " >> ", dotx, doty)
						data[nozzle.Index] = int(dot.Base.Color.A)
						//printable = true
					}
				}
			}

		}
	}
	//if printable {
	//fileName := fmt.Sprintf("output/%02d.png", *imageIndex)
	//outputFile, _ := os.Create(fileName)
	//png.Encode(outputFile, img)
	//outputFile.Close()
	//*imageIndex = *imageIndex + 1
	//}
	return data
}
