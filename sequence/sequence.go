package sequence

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"posam/dao/printheads"
	"posam/gui/uiutil"
	"posam/util/platform"
	"strconv"
	"strings"
)

const DEBUG = false

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

	resolutionLabel := widgets.NewQLabel2("Tolerance:", nil, 0)
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
	if !DEBUG {
		spaceSlideyLabel.SetVisible(false)
		spaceSlideyInput.SetVisible(false)
	}

	sequenceInput := widgets.NewQTextEdit(nil)

	resolutionInput.SetText("30")
	startxInput.SetText("-50000")
	startyInput.SetText("50000")
	spacexInput.SetText("169.3")
	spaceyInput.SetText("550.3")
	spaceBlockxInput.SetText("3000")
	spaceBlockyInput.SetText("3000")
	spaceSlideyInput.SetText("5000")
	sequenceInput.SetText(`CTGGTTCCTCATATAAGCTT, CGTTAAAACATCGACTGACT
CAACATTTAGACAATAAACG, CACCAGGTGAACATTTTTGA
TGAGCGTCCGACGCGGTCCT, ATATAGAAAGTTATTTGATG
GAAATATCACTTCGTGAACA, GTAGCTCATGAGCTGCAGTG

GTAACTCTAACGTATAGGCA, TATGATCTTCTAACCATATG
TTGGGAAGACAGCACCTGAC, TGCCTACAGCGCTACGCGCA
CCCAGTACCCTGGGCCACGA, AAAACCGGTAAGGTGCGAAG
CCCGGCAATGGATACCGTAG, TGCTCGCCAGGATCACCTAT`)
	if DEBUG {
		sequenceInput.SetText(`CTGGTTCCTCATATAAGCTT, CGTTAAAACATCGACTGACT, GTTGTGAGAGTCAGTTATAG, TTGTCGTAACTTTCTGCCCT, CATAGGTTTAATATTGGATC, GTGAATACTTCGGCGGGTTG, AGGGTCTGAACGCTCATAGT, GCGTTATCGCTAGTGCGCAA
CAACATTTAGACAATAAACG, CACCAGGTGAACATTTTTGA, GTTTCAAAATACGCCAAAGT, CCGGTATTTCTACCGAATGT, GATCCACGTCAGTGTGCTAG, GTTCGACCATTTCCAGCAGT, GGGATCGTTCGCGGTCTGTT, AAACCATCGCACCCCCAAGC
TGAGCGTCCGACGCGGTCCT, ATATAGAAAGTTATTTGATG, CGATTGAGGAGCGCAAAGGA, TATAGTTACTGAGGATGTGC, GAAGAATACGAAACACGCTC, GTGAATTGTACGCCAACGGG, ACTCGCAGTCAGAGCATTCA, ACAGAGGAGCCTCGAATCCT
GAAATATCACTTCGTGAACA, GTAGCTCATGAGCTGCAGTG, GGCTCTATATTTCAACGAAC, TAAAGGTGCAAGCGGATACG, ACGAGAGACGCCAATGACTC, TGCTGTAGCGGTAGTGATCA, GTTGGGTTGGTCTGCACAGC, AATCCTGAGGAATGTGTTTT
GTAACTCTAACGTATAGGCA, TATGATCTTCTAACCATATG, CAACCGAGGGGCTCTGAAGA, CGTTGCGGTATGGCGTCATG, ACTTCAAGGGCACTTCGCCT, AACATACGCGCCTGGGGTAG, AGTCACTCTACCACGAAGAC, AGGCGGCAGCTTACGTTGCC
TTGGGAAGACAGCACCTGAC, TGCCTACAGCGCTACGCGCA, CTTCCTAGGAGACCCCATGA, ACCTTAAGGACACACAGGGA, CGACGCGTGTCCTTAGTACC, TGAATTGAGACCTGAGGCCT, CGCCGCCGTTAGGTGAAAAT, CGTTCGCGAAAACGCCAACC
CCCAGTACCCTGGGCCACGA, AAAACCGGTAAGGTGCGAAG, ATGTAACGGTTGAAATTTAG, TGTATGTCGATACTAGGCTA, GTCGAGATCTATTACCGGCC, GTGTGGCGCAACTCAAATTG, AACAATTTGTTTCGTCGTTG, GCAGCAAAGCCTTTTGAGCC
CCCGGCAATGGATACCGTAG, TGCTCGCCAGGATCACCTAT, CCACCACAAGATCGGCGACT, CGTGGGTACTCTTCATACCG, CTGCGAACAAGTCCCCCCTA, ACCCATTGCTGGCAGATGAT, AATTGGCACAGTAGACAGTC, GTCGCAACGGGACGTCAAAC
TATTGATGTACCAAACGCAG, TGTTAAAGACCGCCGCGTCG, CATCTGTCTGGTCTGCCATT, GCGTTACACGATGGTTGCAG, CACCGCCATCCTACTTTGCA, GCGCCTCTCGTAATCAACAC, CCGGGACACCCTGTACACCG, CCACGTGAATAATTCGTTCC
AGACGGGAACGGGCCGGCTC, GGACGCCAAGCAGGGACAAC, CGAACTGTAGGGGTCACGCG, CCGGGGCGGTTTGGCGATAC, GGTACCGCTGAGTTCACCGG, ACTTAGCATATCACGAATAC, ACCCGGGGTAGGCTACCAGG, CCTGAAGTCTTGTTTGGTCC

CATCCCAACTTGACGATGGT, TCGGCCTTCTTAAAGTAAAC, AGACTGGCGTTAGTGTACCA, GACAAGCGTCATCGCGGAGA, GGTTATGCATGCGAGGACAG, ACTATTCCTGGGATGCACTG, CACCGGCAACGGCTCCAATA, CCTGACAGAGCCAGTCGTTT
TAAAGATTGCACGACTGGTC, GGGGCGGGCGAAATTACGTT, ACCCAGACGAGACAGTCAGC, AACACGGACAGAGCCCCCGA, GTCTGACTCGTGTAGTGTTC, CGAGAAACATCCCCTATCAG, CCCACGTACTGTCCCAGCTT, CCGATATGGACGAAACCCGA
CGGACAGCTCGTATTAATGT, AAATAGCTGCCTTTGTGGGG, ACCAATTTGTGTGAAACCAA, CTGAGAACAGAGCCGGTTCC, TCTTCATGGTTAGGGCCTCG, CGTCGCTCTACTGTCCGATG, TTATGGCGCGTCACCTGCCA, TCGAAAGCGTTCTACACTCG
CATAACCGCTTCCCTGGGAG, AGGCGACATGGCCCACATGC, TGTGTGAGATTTGGTGTGGA, GAAACTATAATGATCAATCA, TCGATCCGGAGAGGCCGAGA, ATAGACTCGAGAGTCTCCCA, CACATGACCACAAACGGACG, CACCGGGGAACTATAGATAA
ACGTAGCATGGAGTACTCCT, ATGCAGCCTCACACAGCGTC, CTGGACTCAAGCGTTACCTG, GCTCCGCCTCCTCAATGACG, CAATTTATTATTGAGTCTGA, AAAACTATGCCGTGGTGCTC, CCACCCAGTGGTGAGCTGGA, CACATCGTAGCTAGAAATGC
CAACCGGAAAATGTCTCTTG, ATTATGCCACATGCCACGCT, ATCTGCGACTTAAGTCCATG, GCAAACCTGCTGCCCGAAAC, GAACCCGTATAGAACGCACA, AGCACTTCAGATATTCAGAG, AATTCCAGACCTATACGGCG, CGGCAGAATGGAAGTATGCT
AGACAAGTCGGTGGCCATGT, GTGGAGCGGCATTGCACTCC, AGGGTTGCACTTAAGGTCTT, CCGTCCAACCTGTATTCGAT, CAACTGGGATGAGTTGTTTT, CCGACGGTCGTTCACTACTG, CGCCGACGTCTTACCCAGAT, TAAGCTAGTATTACACGATA
CGCAATTTTCGTATAAATAT, CGAGCAGAGCCCCGAACATG, GCACCGAAATCCCCATAAAG, ACAGATTTCACAAACGTAGG, CCCGAGTCCATCCCTACGGC, AAAAACAACATTCGACGCAG, TCCTTTGCGGTTGAACAGGA, GTGTCACCTTACTGTGAACG
GCACGTGAGGTTAAGTTATG, ACCATTTCCCTCGTGTTAGT, TAGCATCATATGGCCGAAGC, TTGCTTAGTAGACGCATATC, CCCGCTTGTATCAAGAAAGT, GCCGGGGGGTATCACGGGGA, CACGTTACACAGCTGCTCTC, TTATCACCAAATCATCTCCT
GGATTTAATCTGTTCGGATA, GACGCTGAATCGTGATAAAC, TGGACCTCCCTTGTTAACTC, AGTAATTCTTCGGGTCGATG, ACTTCGGCCTGAGGCTCGAC, CACGACAAAGCCATCTTATG, TCTGCCTGCACATGCTTGGC, ATGACCTATGTTCAGCTCTA
`)
	}

	generateButton := widgets.NewQPushButton2("PREVIEW", nil)

	generateProgressbar := widgets.NewQProgressBar(nil)
	generateProgressbar.SetWindowTitle("generating...")
	generateProgressbar.SetMinimum(0)
	generateProgressbar.SetMaximum(1000)
	generateProgressbar.SetValue(0)
	generateProgressbar.SetVisible(false)

	generateProgressbar.ConnectValueChanged(func(value int) {
		if value == generateProgressbar.Maximum() {
			generateProgressbar.SetValue(generateProgressbar.Minimum())
			generateButton.SetVisible(true)
			generateProgressbar.SetVisible(false)
		}
	})

	generateButton.ConnectClicked(func(bool) {
		resolutionFloat, err := strconv.ParseFloat(resolutionInput.Text(), 32)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		if resolutionFloat < 10.0 {
			uiutil.MessageBoxError("tolerance should be greater than the resolution of motor")
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

		generateButton.SetVisible(false)
		generateProgressbar.SetVisible(true)
		go func() {
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
				generateProgressbar,
			)
		}()
	})

	motorPathLabel := widgets.NewQLabel2("Motor path:", nil, 0)
	motorPathInput := widgets.NewQLineEdit(nil)
	motorPathInput.SetText("/AOZTECH/Motor")
	motorSpeedLabel := widgets.NewQLabel2("Motor speed:", nil, 0)
	motorSpeedInput := widgets.NewQLineEdit(nil)
	motorSpeedInput.SetText("10")
	motorAccelLabel := widgets.NewQLabel2("Motor acceleration:", nil, 0)
	motorAccelInput := widgets.NewQLineEdit(nil)
	motorAccelInput.SetText("100")
	printheadPathLabel := widgets.NewQLabel2("Printhead path:", nil, 0)
	printheadPathInput := widgets.NewQLineEdit(nil)
	printheadPathInput.SetText("/Ricoh-G5/Printer#1")
	printheadxLabel := widgets.NewQLabel2("Printhead x:", nil, 0)
	printheadxInput := widgets.NewQLineEdit(nil)
	printheadyLabel := widgets.NewQLabel2("Printhead y:", nil, 0)
	printheadyInput := widgets.NewQLineEdit(nil)

	exportButton := widgets.NewQPushButton2("BUILD", nil)

	exportProgressbar := widgets.NewQProgressBar(nil)
	exportProgressbar.SetWindowTitle("exporting...")
	exportProgressbar.SetMinimum(0)
	exportProgressbar.SetMaximum(1000)
	exportProgressbar.SetValue(0)
	exportProgressbar.SetVisible(false)

	exportProgressbar.ConnectValueChanged(func(value int) {
		if value == exportProgressbar.Maximum() {
			exportProgressbar.SetValue(exportProgressbar.Minimum())
			exportButton.SetVisible(true)
			exportProgressbar.SetVisible(false)
		}
	})

	exportButton.ConnectClicked(func(bool) {
		resolutionFloat,
			startxFloat,
			startyFloat,
			spacexFloat,
			spaceyFloat,
			spaceBlockxFloat,
			spaceBlockyFloat,
			spaceSlideyFloat,
			printheadxFloat,
			printheadyFloat,
			err := parseFloatArg(
			resolutionInput.Text(),
			startxInput.Text(),
			startyInput.Text(),
			spacexInput.Text(),
			spaceyInput.Text(),
			spaceBlockxInput.Text(),
			spaceBlockyInput.Text(),
			spaceSlideyInput.Text(),
			printheadxInput.Text(),
			printheadyInput.Text(),
		)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		exportButton.SetVisible(false)
		exportProgressbar.SetVisible(true)
		export(
			int(resolutionFloat*platform.UM),
			int(startxFloat*platform.UM),
			int(startyFloat*platform.UM),
			int(spacexFloat*platform.UM),
			int(spaceyFloat*platform.UM),
			int(spaceBlockxFloat*platform.UM),
			int(spaceBlockyFloat*platform.UM),
			int(spaceSlideyFloat*platform.UM),
			sequenceInput.ToPlainText(),
			motorPathInput.Text(),
			motorSpeedInput.Text(),
			motorAccelInput.Text(),
			printheadPathInput.Text(),
			int(printheadxFloat*platform.UM),
			int(printheadyFloat*platform.UM),
			exportProgressbar,
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
	layout.AddWidget3(generateProgressbar, 9, 0, 1, 2, 0)
	layout.AddWidget3(generateButton, 10, 0, 1, 2, 0)
	layout.AddWidget(motorPathLabel, 11, 0, 0)
	layout.AddWidget(motorPathInput, 11, 1, 0)
	layout.AddWidget(motorSpeedLabel, 12, 0, 0)
	layout.AddWidget(motorSpeedInput, 12, 1, 0)
	layout.AddWidget(motorAccelLabel, 13, 0, 0)
	layout.AddWidget(motorAccelInput, 13, 1, 0)
	layout.AddWidget(printheadPathLabel, 14, 0, 0)
	layout.AddWidget(printheadPathInput, 14, 1, 0)
	layout.AddWidget(printheadxLabel, 15, 0, 0)
	layout.AddWidget(printheadxInput, 15, 1, 0)
	layout.AddWidget(printheadyLabel, 16, 0, 0)
	layout.AddWidget(printheadyInput, 16, 1, 0)
	layout.AddWidget3(exportProgressbar, 17, 0, 1, 2, 0)
	layout.AddWidget3(exportButton, 18, 0, 1, 2, 0)

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
	printheadx string,
	printheady string,
) (
	resolutionFloat float64,
	startxFloat float64,
	startyFloat float64,
	spacexFloat float64,
	spaceyFloat float64,
	spaceBlockxFloat float64,
	spaceBlockyFloat float64,
	spaceSlideyFloat float64,
	printheadxFloat float64,
	printheadyFloat float64,
	err error,
) {
	resolutionFloat, err = strconv.ParseFloat(resolution, 32)
	if err != nil {
		return
	}
	if resolutionFloat < 10.0 {
		err = fmt.Errorf("tolerance should be greater than the resolution of motor")
		return
	}
	startxFloat, err = strconv.ParseFloat(startx, 32)
	if err != nil {
		return
	}
	startyFloat, err = strconv.ParseFloat(starty, 32)
	if err != nil {
		return
	}
	spacexFloat, err = strconv.ParseFloat(spacex, 32)
	if err != nil {
		return
	}
	spaceyFloat, err = strconv.ParseFloat(spacey, 32)
	if err != nil {
		return
	}
	spaceBlockxFloat, err = strconv.ParseFloat(spaceBlockx, 32)
	if err != nil {
		return
	}
	spaceBlockyFloat, err = strconv.ParseFloat(spaceBlocky, 32)
	if err != nil {
		return
	}
	spaceSlideyFloat, err = strconv.ParseFloat(spaceSlidey, 32)
	if err != nil {
		return
	}
	printheadxFloat, err = strconv.ParseFloat(printheadx, 32)
	if err != nil {
		return
	}
	printheadyFloat, err = strconv.ParseFloat(printheady, 32)
	if err != nil {
		return
	}
	return
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
	generateProgressbar *widgets.QProgressBar,
) {
	xoffset := startX
	yoffset := startY
	width := 0
	height := 0
	pixels := make(map[int]map[int]*color.NRGBA)

	pixelSum := 0
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
				pixelSum += 1
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
	pixelCount := -1
	for y, pixel := range pixels {
		for x, c := range pixel {
			img.Set(x, y, *c)
			pixelCount += 1
			generateProgressbar.SetValue(pixelCount * generateProgressbar.Maximum() / pixelSum)
		}
	}
	if DEBUG {
		// with file
		outputFile, _ := os.Create("test.png")
		png.Encode(outputFile, img)
		outputFile.Close()
	}
	// nofile
	var imagebuff bytes.Buffer
	png.Encode(&imagebuff, img)
	imagebyte := imagebuff.Bytes()
	qimg := gui.NewQImage()
	qimg.LoadFromData2(core.NewQByteArray2(string(imagebyte), len(imagebyte)), "png")
	qimg = qimg.Scaled2(5*width, 5*height, core.Qt__IgnoreAspectRatio, core.Qt__FastTransformation)
	imageItem.SetPixmap(gui.NewQPixmap().FromImage(qimg, 0))
	generateProgressbar.SetValue(generateProgressbar.Maximum())
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
	resolution int,
	startX int,
	startY int,
	spaceX int,
	spaceY int,
	spaceBlockx int,
	spaceBlocky int,
	spaceSlidey int,
	sequences string,
	motorPath string,
	motorSpeed string,
	motorAccel string,
	printheadPath string,
	printheadX int,
	printheadY int,
	exportProgressbar *widgets.QProgressBar,
) {
	filePath, err := uiutil.FilePath()
	if err != nil {
		uiutil.MessageBoxError(err.Error())
		return
	}

	bin := NewBin(motorPath, motorSpeed, motorAccel, printheadPath)

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
				yoffset -= spaceBlocky + spaceY
			} else {
				yoffset -= spaceBlocky
			}
		case 2:
			if yoffset != startY {
				yoffset -= spaceSlidey + spaceY
			} else {
				yoffset -= spaceSlidey
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
				log.Println("location", xoffset, yoffset)
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
		yoffset -= spaceY
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
			//pf.Dots[x][y] = dot
		}
	}
	log.Println("platform", columnCount+1, rowCount+1)

	h, _ := printheads.NewPrintHeadLineD(
		4,
		1280,
		169.3*platform.UM,
		84.65*platform.UM,
		550.3*platform.UM,
		11.811*platform.MM,
		printheadX,
		printheadY,
	)

	sum, px, py, dot, err := pf.NextDotVertical()
	if err != nil {
		fmt.Println(err)
	}
	//dotIndexX := px
	//dotIndexY := py
	h.UpdatePositionStar(dot.PositionX, dot.PositionY)
	// loop horizontally, from left to right
	// 1. offset, upward
	// 2. next, downward
	imageIndex := 0
	img := image.NewRGBA(image.Rect(0, 0, 100*platform.MM/resolution, 100*platform.MM/resolution))
	//img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	//fmt.Println("IMAGE", 100*platform.MM/resolution*platform.UM)

	var direction string

	go func() {
		count := 0
		log.Println("rect", pf.Top(), pf.Right(), pf.Bottom(), pf.Left())
		for h.Rows[3].Nozzles[0].X <= pf.Right() {
			direction = "downward"
			log.Println("moving downward", dot.PositionY, pf.Bottom())
			for dposy := dot.PositionY; h.Rows[3].Nozzles[0].Y >= pf.Bottom(); dposy -= resolution {
				data := genData(exportProgressbar, &count, sum, h, pf, py, &imageIndex, img, resolution, direction)
				if data != "" {
					bin.AddMotion(dot.PositionX, dposy)
					bin.AddPrint(data)
				}
				h.UpdatePositionStar(dot.PositionX, dposy)
			}

			dposx := dot.PositionX + h.RowOffset
			//dposx := dot.PositionX + resolution
			h.UpdatePositionStar(dposx, h.Rows[0].Nozzles[0].Y)

			direction = "upward"
			log.Println("moving upward", h.Rows[0].Nozzles[0].Y, pf.Top())
			for dposy := h.Rows[0].Nozzles[0].Y; h.Rows[0].Nozzles[0].Y <= pf.Top(); dposy += resolution {
				data := genData(exportProgressbar, &count, sum, h, pf, py, &imageIndex, img, resolution, direction)
				if data != "" {
					bin.AddMotion(dposx, dposy)
					bin.AddPrint(data)
				}
				h.UpdatePositionStar(dposx, dposy)
			}

			sum, px, py, dot, err = pf.NextDotVertical()
			if err != nil {
				log.Println("DONE", err, count, sum)
				break
			}
			log.Println("moving to the next spot", count, sum, "||", px, py, dot.Base.Name, dot.PositionX, dot.PositionY)
			h.UpdatePositionStar(dot.PositionX, dot.PositionY)
		}

		file, err := os.Create(filePath)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		}
		encoder := gob.NewEncoder(file)
		encoder.Encode(bin.Node)
		//uiutil.MessageBoxInfo(fmt.Sprintf("Sequences have been built into %q", filePath))
	}()

}

func genData(exportProgressbar *widgets.QProgressBar, count *int, sum int, h *printheads.PrintHead, pf *platform.Platform, py int, imageIndex *int, img *image.RGBA, resolution int, direction string) (output string) {
	data := make([]string, 1280)

	printable := false
	// traverse nozzles
	for _, row := range h.Rows {
		for _, nozzle := range row.Nozzles {

			data[nozzle.Index] = "0"
			for _, dot := range pf.AvailableDots() {
				dotx, doty := dot.PositionX, dot.PositionY
				if nozzle.IsAvailable(dotx, doty, resolution) {
					if dot.Base.Name == row.Reagent {
						*count = *count + 1
						exportProgressbar.SetValue(*count * exportProgressbar.Maximum() / sum)
						dot.Printed = true
						if DEBUG {
							img.Set((dotx+50*platform.MM)/resolution, (50*platform.MM-doty)/resolution, dot.Base.Color)
							log.Println(dot.Base.Name, nozzle, " || ", dot, " >> ", dotx, doty) //, "..", (dotx+50*platform.MM)/resolution, (50*platform.MM-doty)/resolution)
						}
						data[nozzle.Index] = "1"
						printable = true
					} else {
						// should not happen
					}
				}
			}
		}
	}
	if printable {
		if DEBUG {
			fileName := fmt.Sprintf("output/%03d.%s.png", *imageIndex, direction)
			outputFile, _ := os.Create(fileName)
			png.Encode(outputFile, img)
			outputFile.Close()
			*imageIndex = *imageIndex + 1
		}
		outputSlice := make([]string, 160)
		for i := 0; i < len(data); i += 8 {
			value, _ := strconv.ParseInt(strings.Join(data[i:i+8], ""), 2, 64)
			outputSlice = append(outputSlice, fmt.Sprintf("%02x", value))
		}
		output = strings.Join(outputSlice, "")
		//fmt.Println(len(outputSlice), outputSlice)
		//fmt.Println(len(output), output)
		//fmt.Println(data)
	}
	return output
}
