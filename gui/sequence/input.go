package sequence

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"synthesis/gui/uiutil"
	"synthesis/internal/formation"
	"synthesis/internal/geometry"
	"synthesis/internal/printhead"
	"synthesis/internal/reagent"
	"synthesis/internal/slide"
	"synthesis/internal/substrate"
	//"time"
)

const (
	DPI_150 = "169.3"
	DPI_300 = "84.65"
	DPI_600 = "42.325"
)

var offsetX, offsetY float64 // mm
//var (
//offsetX = 0.0 //mm
//offsetY = 0.0 //mm
//)
var (
	IMAGABLE  = os.Getenv("IMG") == "true"
	DEBUGABLE = os.Getenv("DBG") == "true"
)

// const{{{

const (
	SEQUENCE_EXAMPLE = `GGTCATC
CATTGAT
ATCCCGG
`
	//SEQUENCE_EXAMPLE = `GGGTCGGATGATCGGACACT
	//CATCATCTGGGTACAGCGGG
	//ATTATACAGTTTTGTCCAAT
	//`

	//SEQUENCE_EXAMPLE = `GGGTCGGATGATCGGACACT
//CATCATCTGGGTACAGCGGG
//ATTATACAGTTTTGTCCAAT
//CTATCTTGGAGGGGTAGGCG
//AGGCTGGCCATGTTGTCTTA
//ACTTTATGATGCGTAAGCAC
//CAGCCTCAACCGCTCTGCAA
//CATGCTCCTATCGTAGGAAG
//CAGGAGTCCATTCCGTGCTG
//ATTGCCGTTAATCGGCAGGA
//AGAGTGCCGGAACACTGTTG
//TCACGAGGGGGCAAAGAAAG
//ATTTGCCGGGGCGTGTCCTG
//GGATGCTGACACGTCGTGTT
//GTATCTACTTGACTACGGCC
//GGTTTGAAGTAAGACCCCCC
//CGTCTCGGCCCGTAATCTCC
//TGATCCAAATCGATTAATGT
//AAGATCCCAGTTTTTTAGAG
//AATCACTGCATTGCGAAAAA
//CTGCACGATTATGGGGTGAG
//GTCCGACCAGGGGTCTATCT
//CTGGAAATGCCTGGGCGGTG
//TTCCAAGTGATAGCTACGAA
//GTTCCGTTATGCCGAGGAAG
//AGATCCACGGCTCGTCAGAC
//GATGAATTAGCGGAGGATCC
//GGCACGGTAAGTTCCCACGC
//GCGCTCGAGACGAACACTAA
//CGATAGATGAATGGGCACCT
//CCAGACCGGAGTTGGAGGAG
//GTTTGCTCCTCTTCACTCCG
//TCAAGGCTGATATCACCAAT
//CAGCATCTTAACTCCAGGAC
//GTATCTCTCGTAACATGCTA
//ATCACGAGATGAAAGTCTGG
//TTCTCGTTCCACCCAGTCGT
//GAAGCTCAACACATAGCAAC
//GACCGGACGAGAAAACTCCG
//TACTCCCTCAAGTAAGTCTA
//CTAGACCGCAGCAAAATCGT
//TCACTTTCGCGCGCACAGGG
//AGGGTCGGACTTCTAGGTAG
//GATCAGACACCTCATCACGA
//GTGCCTCCTGCCCTAGTCGA
//CCAACATGTGCCAACGATTA
//ATGAGCTGAAGACAGAGGGC
//ATAGCCGCTGGCGTTCGTGG
//GCGTAGCAAAGGGGCGGAGT
//CAGTTATTTCAGAGGTACCG
//`
)

// }}}

func NewInputGroup() (
	group *widgets.QGroupBox,
	previewGroup *widgets.QWidget,
) {
	group = widgets.NewQGroupBox2("Parameters", nil)
	layout := widgets.NewQGridLayout2()
	group.SetLayout(layout)

	fileInput := widgets.NewQLineEdit(nil)
	fileInput.SetReadOnly(true)
	fileInput.ConnectMousePressEvent(func(e *gui.QMouseEvent) {
		filePath, err := uiutil.FilePath()
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		fileInput.SetText(filePath)
	})
	layout.AddWidget2(fileInput, 0, 0, 0)

	sequenceInput := widgets.NewQTextEdit(nil)
	sequenceInput.SetText(SEQUENCE_EXAMPLE)
	layout.AddWidget2(sequenceInput, 1, 0, 0)

	// device group{{{

	deviceGroup := widgets.NewQGroupBox2("Device", nil)
	deviceLayout := widgets.NewQGridLayout2()
	deviceGroup.SetLayout(deviceLayout)
	layout.AddWidget2(deviceGroup, 2, 0, 0)

	motorPathLabel := widgets.NewQLabel2("Motor path", nil, 0)
	motorPathInput := widgets.NewQLineEdit(nil)
	motorPathInput.SetText("/AOZTECH/Motor")
	motorSpeedLabel := widgets.NewQLabel2("Motor speed", nil, 0)
	motorSpeedInput := widgets.NewQLineEdit(nil)
	motorSpeedInput.SetText("10")
	motorAccelLabel := widgets.NewQLabel2("Motor acceleration", nil, 0)
	motorAccelInput := widgets.NewQLineEdit(nil)
	motorAccelInput.SetText("100")
	printhead0PathLabel := widgets.NewQLabel2("Printhead path", nil, 0)
	printhead0PathInput := widgets.NewQLineEdit(nil)
	printhead0PathInput.SetText("/Ricoh-G5/Printer#2")
	//printhead1PathLabel := widgets.NewQLabel2("Printhead #2 path", nil, 0)
	//printhead1PathInput := widgets.NewQLineEdit(nil)
	//printhead1PathInput.SetText("/Ricoh-G5/Printer#1")

	deviceLayout.AddWidget2(motorPathLabel, 0, 0, 0)
	deviceLayout.AddWidget2(motorPathInput, 0, 1, 0)
	deviceLayout.AddWidget2(motorSpeedLabel, 1, 0, 0)
	deviceLayout.AddWidget2(motorSpeedInput, 1, 1, 0)
	deviceLayout.AddWidget2(motorAccelLabel, 2, 0, 0)
	deviceLayout.AddWidget2(motorAccelInput, 2, 1, 0)
	deviceLayout.AddWidget2(printhead0PathLabel, 3, 0, 0)
	deviceLayout.AddWidget2(printhead0PathInput, 3, 1, 0)
	//deviceLayout.AddWidget2(printhead1PathLabel, 4, 0, 0)
	//deviceLayout.AddWidget2(printhead1PathInput, 4, 1, 0)

	// }}}

	// position gorup{{{

	positionGroup := widgets.NewQGroupBox2("Position (unit: mm)", nil)
	positionLayout := widgets.NewQGridLayout2()
	positionGroup.SetLayout(positionLayout)
	//layout.AddWidget2(positionGroup, 3, 0, 0)

	printhead0PositionLabel := widgets.NewQLabel2("Printhead #0", nil, 0)
	printhead1PositionLabel := widgets.NewQLabel2("Printhead #1", nil, 0)
	slide0PositionLabel := widgets.NewQLabel2("Slide #0", nil, 0)
	slide1PositionLabel := widgets.NewQLabel2("Slide #1", nil, 0)
	slide2PositionLabel := widgets.NewQLabel2("Slide #2", nil, 0)

	printhead0PositionXInput := widgets.NewQLineEdit(nil)
	printhead0PositionYInput := widgets.NewQLineEdit(nil)
	printhead1PositionXInput := widgets.NewQLineEdit(nil)
	printhead1PositionYInput := widgets.NewQLineEdit(nil)

	printhead0PositionXInput.SetText("0")
	printhead0PositionYInput.SetText("0")
	printhead1PositionXInput.SetText("50")
	printhead1PositionYInput.SetText("0")

	slide0PositionXInput := widgets.NewQLineEdit(nil)
	slide0PositionYInput := widgets.NewQLineEdit(nil)
	slide1PositionXInput := widgets.NewQLineEdit(nil)
	slide1PositionYInput := widgets.NewQLineEdit(nil)
	slide2PositionXInput := widgets.NewQLineEdit(nil)
	slide2PositionYInput := widgets.NewQLineEdit(nil)

	slide0PositionXInput.SetText("0")
	slide0PositionYInput.SetText("0")
	slide1PositionXInput.SetText("22")
	slide1PositionYInput.SetText("-5")
	slide2PositionXInput.SetText("26")
	slide2PositionYInput.SetText("0")

	positionLayout.AddWidget2(printhead0PositionLabel, 0, 0, 0)
	positionLayout.AddWidget2(printhead0PositionXInput, 0, 1, 0)
	positionLayout.AddWidget2(printhead0PositionYInput, 0, 2, 0)

	positionLayout.AddWidget2(printhead1PositionLabel, 1, 0, 0)
	positionLayout.AddWidget2(printhead1PositionXInput, 1, 1, 0)
	positionLayout.AddWidget2(printhead1PositionYInput, 1, 2, 0)

	positionLayout.AddWidget2(slide0PositionLabel, 2, 0, 0)
	positionLayout.AddWidget2(slide0PositionXInput, 2, 1, 0)
	positionLayout.AddWidget2(slide0PositionYInput, 2, 2, 0)

	positionLayout.AddWidget2(slide1PositionLabel, 3, 0, 0)
	positionLayout.AddWidget2(slide1PositionXInput, 3, 1, 0)
	positionLayout.AddWidget2(slide1PositionYInput, 3, 2, 0)

	positionLayout.AddWidget2(slide2PositionLabel, 4, 0, 0)
	positionLayout.AddWidget2(slide2PositionXInput, 4, 1, 0)
	positionLayout.AddWidget2(slide2PositionYInput, 4, 2, 0)

	// }}}

	// space group{{{

	spaceGroup := widgets.NewQGroupBox2("Space (unit: um)", nil)
	spaceLayout := widgets.NewQGridLayout2()
	spaceGroup.SetLayout(spaceLayout)
	layout.AddWidget2(spaceGroup, 3, 0, 0)

	spaceLabel := widgets.NewQLabel2("Spot space", nil, 0)
	spacexInput := widgets.NewQLineEdit(nil)
	spaceyInput := widgets.NewQLineEdit(nil)

	spacexInput.SetText("169.3")
	spaceyInput.SetText("550.3")

	spaceLayout.AddWidget2(spaceLabel, 0, 0, 0)
	spaceLayout.AddWidget2(spacexInput, 0, 1, 0)
	spaceLayout.AddWidget2(spaceyInput, 0, 2, 0)

	// }}}

	// reagent group{{{

	reagentGroup := widgets.NewQGroupBox2("Reagent", nil)
	reagentLayout := widgets.NewQGridLayout2()
	reagentGroup.SetLayout(reagentLayout)
	layout.AddWidget2(reagentGroup, 4, 0, 0)

	printhead0Line0Label := widgets.NewQLabel2("Row A of Printhead #0", nil, 0)
	printhead0Line1Label := widgets.NewQLabel2("Row B of Printhead #0", nil, 0)
	printhead0Line2Label := widgets.NewQLabel2("Row C of Printhead #0", nil, 0)
	printhead0Line3Label := widgets.NewQLabel2("Row D of Printhead #0", nil, 0)

	printhead1Line0Label := widgets.NewQLabel2("Row A of Printhead #1", nil, 0)
	printhead1Line1Label := widgets.NewQLabel2("Row B of Printhead #1", nil, 0)
	printhead1Line2Label := widgets.NewQLabel2("Row C of Printhead #1", nil, 0)
	printhead1Line3Label := widgets.NewQLabel2("Row D of Printhead #1", nil, 0)

	printhead0Line0Input := widgets.NewQLineEdit(nil)
	printhead0Line1Input := widgets.NewQLineEdit(nil)
	printhead0Line2Input := widgets.NewQLineEdit(nil)
	printhead0Line3Input := widgets.NewQLineEdit(nil)

	printhead1Line0Input := widgets.NewQLineEdit(nil)
	printhead1Line1Input := widgets.NewQLineEdit(nil)
	printhead1Line2Input := widgets.NewQLineEdit(nil)
	printhead1Line3Input := widgets.NewQLineEdit(nil)

	printhead0Line0Input.SetText("A")
	printhead0Line1Input.SetText("C")
	printhead0Line2Input.SetText("G")
	printhead0Line3Input.SetText("T")
	printhead1Line0Input.SetText("-")
	printhead1Line1Input.SetText("-")
	printhead1Line2Input.SetText("-")
	printhead1Line3Input.SetText("-")

	reagentLayout.AddWidget2(printhead0Line0Label, 0, 0, 0)
	reagentLayout.AddWidget2(printhead0Line0Input, 0, 1, 0)

	reagentLayout.AddWidget2(printhead0Line1Label, 1, 0, 0)
	reagentLayout.AddWidget2(printhead0Line1Input, 1, 1, 0)

	reagentLayout.AddWidget2(printhead0Line2Label, 2, 0, 0)
	reagentLayout.AddWidget2(printhead0Line2Input, 2, 1, 0)

	reagentLayout.AddWidget2(printhead0Line3Label, 3, 0, 0)
	reagentLayout.AddWidget2(printhead0Line3Input, 3, 1, 0)

	reagentLayout.AddWidget2(printhead1Line0Label, 4, 0, 0)
	reagentLayout.AddWidget2(printhead1Line0Input, 4, 1, 0)

	reagentLayout.AddWidget2(printhead1Line1Label, 5, 0, 0)
	reagentLayout.AddWidget2(printhead1Line1Input, 5, 1, 0)

	reagentLayout.AddWidget2(printhead1Line2Label, 6, 0, 0)
	reagentLayout.AddWidget2(printhead1Line2Input, 6, 1, 0)

	reagentLayout.AddWidget2(printhead1Line3Label, 7, 0, 0)
	reagentLayout.AddWidget2(printhead1Line3Input, 7, 1, 0)

	// }}}

	// misc group{{{

	miscGroup := widgets.NewQGroupBox2("Misc", nil)
	miscLayout := widgets.NewQGridLayout2()
	miscGroup.SetLayout(miscLayout)
	layout.AddWidget2(miscGroup, 5, 0, 0)

	toleranceLabel := widgets.NewQLabel2("Tolerance (um)", nil, 0)
	toleranceInput := widgets.NewQLineEdit(nil)
	toleranceInput.SetText("30")

	dpiLabel := widgets.NewQLabel2("Resolution (um)", nil, 0)
	dpiInput := widgets.NewQComboBox(nil)
	dpiInput.AddItems([]string{
		DPI_150,
	})

	slideCountLabel := widgets.NewQLabel2("Slide count", nil, 0)
	slideCountHoriInput := widgets.NewQLineEdit(nil)
	slideCountHoriInput.SetText("3")
	slideCountVertInput := widgets.NewQLineEdit(nil)
	slideCountVertInput.SetText("1")

	slideAreaSpaceLabel := widgets.NewQLabel2("Slide space (mm)", nil, 0)
	slideAreaSpaceHoriInput := widgets.NewQLineEdit(nil)
	slideAreaSpaceHoriInput.SetText("5")
	slideAreaSpaceVertInput := widgets.NewQLineEdit(nil)
	slideAreaSpaceVertInput.SetText("25")

	printhead1OffsetLabel := widgets.NewQLabel2("offset #1 (mm)", nil, 0)
	printhead1OffsetXInput := widgets.NewQLineEdit(nil)
	printhead1OffsetYInput := widgets.NewQLineEdit(nil)
	printhead1OffsetXInput.SetText("35") // 35
	printhead1OffsetYInput.SetText("65") // 65

	printhead0OffsetLabel := widgets.NewQLabel2("offset #0 (mm)", nil, 0)
	printhead0OffsetXInput := widgets.NewQLineEdit(nil)
	printhead0OffsetYInput := widgets.NewQLineEdit(nil)
	printhead0OffsetXInput.SetText("35") // 35
	printhead0OffsetYInput.SetText("20") // 20

	slideGeometryLabel := widgets.NewQLabel2("slide (mm)", nil, 0)
	slideGeometryWidthInput := widgets.NewQLineEdit(nil)
	slideGeometryHeightInput := widgets.NewQLineEdit(nil)
	slideGeometryWidthInput.SetText("20")
	slideGeometryHeightInput.SetText("29")

	activatorInput := widgets.NewQCheckBox2("Activator", nil)
	previewInput := widgets.NewQCheckBox2("Preview", nil)

	printModeLabel := widgets.NewQLabel2("Print mode", nil, 0)
	printModeInput := widgets.NewQComboBox(nil)
	printModeInput.AddItems([]string{
		"Drop on Demand",
		"Continuous Inkjet",
	})

	miscLayout.AddWidget2(toleranceLabel, 0, 0, 0)
	miscLayout.AddWidget3(toleranceInput, 0, 1, 1, 2, 0)
	miscLayout.AddWidget2(dpiLabel, 1, 0, 0)
	miscLayout.AddWidget3(dpiInput, 1, 1, 1, 2, 0)
	miscLayout.AddWidget2(slideCountLabel, 2, 0, 0)
	miscLayout.AddWidget2(slideCountHoriInput, 2, 1, 0)
	miscLayout.AddWidget2(slideCountVertInput, 2, 2, 0)
	miscLayout.AddWidget2(slideAreaSpaceLabel, 3, 0, 0)
	miscLayout.AddWidget2(slideAreaSpaceHoriInput, 3, 1, 0)
	miscLayout.AddWidget2(slideAreaSpaceVertInput, 3, 2, 0)
	miscLayout.AddWidget2(slideGeometryLabel, 4, 0, 0)
	miscLayout.AddWidget2(slideGeometryWidthInput, 4, 1, 0)
	miscLayout.AddWidget2(slideGeometryHeightInput, 4, 2, 0)
	miscLayout.AddWidget2(printhead1OffsetLabel, 5, 0, 0)
	miscLayout.AddWidget2(printhead1OffsetXInput, 5, 1, 0)
	miscLayout.AddWidget2(printhead1OffsetYInput, 5, 2, 0)
	miscLayout.AddWidget2(printhead0OffsetLabel, 6, 0, 0)
	miscLayout.AddWidget2(printhead0OffsetXInput, 6, 1, 0)
	miscLayout.AddWidget2(printhead0OffsetYInput, 6, 2, 0)
	miscLayout.AddWidget2(printModeLabel, 7, 0, 0)
	miscLayout.AddWidget3(printModeInput, 7, 1, 1, 2, 0)
	miscLayout.AddWidget2(activatorInput, 8, 0, 0)
	miscLayout.AddWidget2(previewInput, 8, 1, 0)

	// }}}

	// build button

	var lotc chan int
	var stepwise = false

	buildButton := widgets.NewQPushButton2("BUILD", nil)
	layout.AddWidget2(buildButton, 6, 0, 0)

	buildProgressbar := widgets.NewQProgressBar(nil)
	layout.AddWidget2(buildProgressbar, 7, 0, 0)

	buildProgressbar.SetWindowTitle("building...")
	buildProgressbar.SetMinimum(0)
	buildProgressbar.SetMaximum(1000)
	buildProgressbar.SetValue(0)
	buildProgressbar.SetVisible(false)

	buildProgressbar.ConnectValueChanged(func(value int) {
		if value == buildProgressbar.Maximum() {
			buildProgressbar.SetValue(buildProgressbar.Minimum())
			buildButton.SetVisible(true)
			buildProgressbar.SetVisible(false)
		}
	})

	// previews

	previewGroup = widgets.NewQWidget(nil, 0)
	previewLayout := widgets.NewQGridLayout2()
	previewLayout.SetContentsMargins(0, 0, 0, 0)
	previewGroup.SetLayout(previewLayout)

	buildButton.ConnectClicked(func(bool) {
		printhead0PositionX,
			printhead0PositionY,
			printhead1PositionX,
			printhead1PositionY,
			slide0PositionX,
			slide0PositionY,
			slide1PositionX,
			slide1PositionY,
			slide2PositionX,
			slide2PositionY,
			spacex,
			spacey,
			tolerance,
			//slideCountHori,
			//slideCountVert,
			slideAreaSpaceHori,
			slideAreaSpaceVert,
			printhead0OffsetX,
			printhead0OffsetY,
			printhead1OffsetX,
			printhead1OffsetY,
			err := ParseParameters(
			printhead0PositionXInput.Text(),
			printhead0PositionYInput.Text(),
			printhead1PositionXInput.Text(),
			printhead1PositionYInput.Text(),
			slide0PositionXInput.Text(),
			slide0PositionYInput.Text(),
			slide1PositionXInput.Text(),
			slide1PositionYInput.Text(),
			slide2PositionXInput.Text(),
			slide2PositionYInput.Text(),
			spacexInput.Text(),
			spaceyInput.Text(),
			toleranceInput.Text(),
			//slideCountHoriInput.Text(),
			//slideCountVertInput.Text(),
			slideAreaSpaceHoriInput.Text(),
			slideAreaSpaceVertInput.Text(),
			printhead0OffsetXInput.Text(),
			printhead0OffsetYInput.Text(),
			printhead1OffsetXInput.Text(),
			printhead1OffsetYInput.Text(),
		)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		fmt.Println(
			printhead0PositionX,
			printhead0PositionY,
			printhead1PositionX,
			printhead1PositionY,
			printhead0OffsetX,
			printhead0OffsetY,
			printhead1OffsetX,
			printhead1OffsetY,
			slide0PositionX,
			slide0PositionY,
			slide1PositionX,
			slide1PositionY,
			slide2PositionX,
			slide2PositionY,
			spacex,
			spacey,
			//slideCountHori,
			//slideCountVert,
			slideAreaSpaceHori,
			slideAreaSpaceVert,
		)

		seqText := sequenceInput.ToPlainText()
		filePath := fileInput.Text()
		if filePath != "" {
			seqBytes, err := ioutil.ReadFile(filePath)
			if err != nil {
				uiutil.MessageBoxError(err.Error())
				return
			}
			seqText = string(seqBytes)
		}

		var step int
		var space int
		switch dpiInput.CurrentText() {
		case DPI_300:
			step = 4
			space = 2
			//step = 126.975 * geometry.UM
			//space = 84.65 * geometry.UM
		case DPI_600:
			step = 1
			space = 1
			//step = 169.3 * geometry.UM
			//space = 42.325 * geometry.UM
		default:
			step = 4
			space = 4
			//step = 42.325 * geometry.UM
			//space = 169.3 * geometry.UM
		}

		// create printhead{{{

		printhead0OffsetXFloat, err := ToFloat(printhead0OffsetXInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}

		var maxOffsetX = 50 - 4*25.4/600
		if printhead0OffsetXFloat > maxOffsetX {
			uiutil.MessageBoxError(fmt.Sprintf(
				"invalid offset of printhead #0: %v > %v",
				printhead0OffsetXFloat,
				maxOffsetX,
			))
			return
		}

		printhead0OffsetYFloat, err := ToFloat(printhead0OffsetYInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		printhead1OffsetXFloat, err := ToFloat(printhead1OffsetXInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		printhead1OffsetYFloat, err := ToFloat(printhead1OffsetYInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		offsetX = printhead0OffsetXFloat
		offsetY = printhead0OffsetYFloat
		fmt.Println("offsets", offsetX, offsetY)

		p0 := printhead.NewPrinthead(
			0,
			[]*reagent.Reagent{
				reagent.NewReagent(printhead0Line0Input.Text()),
				reagent.NewReagent(printhead0Line1Input.Text()),
				reagent.NewReagent(printhead0Line2Input.Text()),
				reagent.NewReagent(printhead0Line3Input.Text()),
			},
			printhead0PathInput.Text(),
			false,
			printhead0OffsetXFloat,
			printhead0OffsetYFloat,
		)
		p0x := geometry.Unit(printhead0OffsetXFloat)
		p0y := geometry.Unit(printhead0OffsetYFloat)
		nozzles0 := p0.MakeNozzles(p0x, p0y)
		fmt.Println("printhead 0", p0x, p0y)

		p1 := printhead.NewPrinthead(
			1,
			[]*reagent.Reagent{
				reagent.NewReagent(printhead1Line0Input.Text()),
				reagent.NewReagent(printhead1Line1Input.Text()),
				reagent.NewReagent(printhead1Line2Input.Text()),
				reagent.NewReagent(printhead1Line3Input.Text()),
			},
			// TODO: printhead #1 path input
			printhead0PathInput.Text(),
			false,
			printhead1OffsetXFloat,
			printhead1OffsetYFloat,
		)

		// automatically adujstment for different row alignment
		//p1x := geometry.Unit(printhead1OffsetXFloat)
		deltay := geometry.Unit(printhead1OffsetYFloat - printhead0OffsetYFloat)
		yrem := deltay % step
		if yrem > step/2 {
			deltay += step - yrem
		} else {
			deltay -= yrem
		}
		deltax := geometry.Unit(printhead1OffsetXFloat - printhead0OffsetXFloat)
		xrem := deltax % step
		if xrem > step/2 {
			deltax += step - xrem
		} else {
			deltax -= xrem
		}
		fmt.Println("printhead 1", p0x-deltax, p0y+deltay)
		nozzles1 := p1.MakeNozzles(p0x-deltax, p0y+deltay)

		printheadArray := printhead.NewArray(
			append(nozzles0, nozzles1...),
			2,
			[]*printhead.Printhead{p0, p1},
		)
		fmt.Println(
			"sights",
			printheadArray.SightTop.Pos.X,
			printheadArray.SightTop.Pos.Y,
			printheadArray.SightBottom.Pos.X,
			printheadArray.SightBottom.Pos.Y,
		)
		fmt.Println(
			"nozzles",
			printheadArray.Nozzles[0].Reagent.Name,
			printheadArray.Nozzles[0].RowIndex,

			printheadArray.Nozzles[1].Reagent.Name,
			printheadArray.Nozzles[1].RowIndex,

			printheadArray.Nozzles[2].Reagent.Name,
			printheadArray.Nozzles[2].RowIndex,

			printheadArray.Nozzles[3].Reagent.Name,
			printheadArray.Nozzles[3].RowIndex,
		)

		// }}}

		// create substrate{{{
		slideCountHoriInt, err := strconv.Atoi(slideCountHoriInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		slideCountVertInt, err := strconv.Atoi(slideCountVertInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		slideAreaSpaceHoriFloat, err := ToFloat(slideAreaSpaceHoriInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		slideAreaSpaceVertFloat, err := ToFloat(slideAreaSpaceVertInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}

		slideGeometryWidthFloat, err := ToFloat(slideGeometryWidthInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		slideGeometryHeightFloat, err := ToFloat(slideGeometryHeightInput.Text())
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}

		maxY := printhead0OffsetYFloat - slideGeometryHeightFloat
		if printhead0OffsetYFloat < 0 {
			maxY = slideGeometryHeightFloat - printhead0OffsetYFloat
		}
		if maxY <= -50.0 {
			uiutil.MessageBoxError(
				fmt.Sprintf(
					"invalid config: slide height '%v', position y '%v'",
					slideGeometryHeightFloat,
					printhead0OffsetYFloat,
				))
			return
		}

		spots, cycleCount := substrate.ParseSpots(
			seqText,
			activatorInput.CheckState() == core.Qt__Checked,
		)
		subs, err := substrate.NewSubstrate(
			slideCountHoriInt,
			slideCountVertInt,
			slideGeometryWidthFloat,
			slideGeometryHeightFloat,
			slideAreaSpaceHoriFloat,
			slideAreaSpaceVertFloat,
			spots,
			space,
			deltax,
		)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		fmt.Println("substrate", subs.Top(), subs.Bottom(), subs.Width, subs.Height)
		if DEBUGABLE {
			for _, spot := range spots {
				fmt.Println(spot.Pos.X, spot.Pos.Y)
				for _, r := range spot.Reagents {
					fmt.Println(r.Name)
				}
			}
		}

		//return // for testing
		// }}}

		buildButton.SetVisible(false)
		buildProgressbar.SetVisible(true)

		mode := formation.MODE_DOD
		if printModeInput.CurrentText() == "Continuous Inkjet" {
			mode = formation.MODE_CIJ
		}

		fmt.Println("initialize lotc")
		lotc = make(chan int)
		if !stepwise {
			fmt.Println("close lotc")
			close(lotc)
		}
		fmt.Println("going to build")

		build(
			step,
			cycleCount,
			printheadArray,
			subs,
			tolerance,
			motorPathInput.Text(),
			motorSpeedInput.Text(),
			motorAccelInput.Text(),
			printhead0PathInput.Text(),
			buildProgressbar,
			byte(mode),
			lotc,
			&stepwise,
			previewInput.CheckState() == core.Qt__Checked,
		)

	})

	// preview group 2

	stepInput := widgets.NewQLineEdit(nil)
	stepInput.SetFixedWidth(50)
	stepInput.SetText("1")
	stepInput.SetAlignment(core.Qt__AlignCenter)

	nextButton := widgets.NewQPushButton2("NEXT", nil)
	prevButton := widgets.NewQPushButton2("PREV", nil)
	previewLayout.AddWidget2(prevButton, 0, 1, 0)
	previewLayout.AddWidget2(stepInput, 0, 2, 0)
	previewLayout.AddWidget2(nextButton, 0, 3, 0)

	nextButton.ConnectClicked(func(bool) {
		previewGroup.SetEnabled(false)
		step := getStep(stepInput)
		fmt.Println("next", step, stepwise)
		if !stepwise {
			stepwise = true
			buildButton.Clicked(false)
		}
		go func() {
			lotc <- step
			previewGroup.SetEnabled(true)
		}()
	})
	prevButton.ConnectClicked(func(bool) {
		previewGroup.SetEnabled(false)
		step := getStep(stepInput)
		fmt.Println("prev", step)
		if !stepwise {
			stepwise = true
			buildButton.Clicked(false)
		}
		go func() {
			lotc <- -step
			previewGroup.SetEnabled(true)
		}()
	})

	return group, previewGroup
}

func ToFloat(inputString string) (float64, error) {
	if inputString == "" {
		return 0.0, fmt.Errorf("should not be null")
	}
	inputFloat, err := strconv.ParseFloat(inputString, 64)
	if err != nil {
		return 0.0, fmt.Errorf(
			"failed to convert %q to float: %v",
			inputString,
			err.Error(),
		)
	}
	return inputFloat, nil
}

func ParseParameters( // {{{
	printhead0PositionXString string,
	printhead0PositionYString string,
	printhead1PositionXString string,
	printhead1PositionYString string,
	slide0PositionXString string,
	slide0PositionYString string,
	slide1PositionXString string,
	slide1PositionYString string,
	slide2PositionXString string,
	slide2PositionYString string,
	spacexString string,
	spaceyString string,
	toleranceString string,
	//slideCountHoriString string,
	//slideCountVertString string,
	slideAreaSpaceHoriString string,
	slideAreaSpaceVertString string,
	printhead0OffsetXString string,
	printhead0OffsetYString string,
	printhead1OffsetXString string,
	printhead1OffsetYString string,
) (
	printhead0PositionXInt int,
	printhead0PositionYInt int,
	printhead1PositionXInt int,
	printhead1PositionYInt int,
	slide0PositionXInt int,
	slide0PositionYInt int,
	slide1PositionXInt int,
	slide1PositionYInt int,
	slide2PositionXInt int,
	slide2PositionYInt int,
	spacexInt int,
	spaceyInt int,
	toleranceInt int,
	slideAreaSpaceHoriInt int,
	slideAreaSpaceVertInt int,
	printhead0OffsetXInt int,
	printhead0OffsetYInt int,
	printhead1OffsetXInt int,
	printhead1OffsetYInt int,
	err error,
) {
	var printhead0PositionXFloat,
		printhead0PositionYFloat,
		printhead1PositionXFloat,
		printhead1PositionYFloat,
		slide0PositionXFloat,
		slide0PositionYFloat,
		slide1PositionXFloat,
		slide1PositionYFloat,
		slide2PositionXFloat,
		slide2PositionYFloat,
		spacexFloat,
		spaceyFloat,
		toleranceFloat,
		slideAreaSpaceHoriFloat,
		slideAreaSpaceVertFloat,
		printhead0OffsetXFloat,
		printhead0OffsetYFloat,
		printhead1OffsetXFloat,
		printhead1OffsetYFloat float64

	printhead0PositionXFloat, err = ToFloat(printhead0PositionXString)
	if err != nil {
		return
	}
	printhead0PositionXInt = int(printhead0PositionXFloat * geometry.MM)

	printhead0PositionYFloat, err = ToFloat(printhead0PositionYString)
	if err != nil {
		return
	}
	printhead0PositionYInt = int(printhead0PositionYFloat * geometry.MM)

	printhead1PositionXFloat, err = ToFloat(printhead1PositionXString)
	if err != nil {
		return
	}
	printhead1PositionXInt = int(printhead1PositionXFloat * geometry.MM)

	printhead1PositionYFloat, err = ToFloat(printhead1PositionYString)
	if err != nil {
		return
	}
	printhead1PositionYInt = int(printhead1PositionYFloat * geometry.MM)

	slide0PositionXFloat, err = ToFloat(slide0PositionXString)
	if err != nil {
		return
	}
	slide0PositionXInt = int(slide0PositionXFloat * geometry.MM)

	slide0PositionYFloat, err = ToFloat(slide0PositionYString)
	if err != nil {
		return
	}
	slide0PositionYInt = int(slide0PositionYFloat * geometry.MM)

	slide1PositionXFloat, err = ToFloat(slide1PositionXString)
	if err != nil {
		return
	}
	slide1PositionXInt = int(slide1PositionXFloat * geometry.MM)

	slide1PositionYFloat, err = ToFloat(slide1PositionYString)
	if err != nil {
		return
	}
	slide1PositionYInt = int(slide1PositionYFloat * geometry.MM)

	slide2PositionXFloat, err = ToFloat(slide2PositionXString)
	if err != nil {
		return
	}
	slide2PositionXInt = int(slide2PositionXFloat * geometry.MM)

	slide2PositionYFloat, err = ToFloat(slide2PositionYString)
	if err != nil {
		return
	}
	slide2PositionYInt = int(slide2PositionYFloat * geometry.MM)

	spacexFloat, err = ToFloat(spacexString)
	if err != nil {
		return
	}
	spacexInt = int(spacexFloat * geometry.UM)

	spaceyFloat, err = ToFloat(spaceyString)
	if err != nil {
		return
	}
	spaceyInt = int(spaceyFloat * geometry.UM)

	toleranceFloat, err = ToFloat(toleranceString)
	if err != nil {
		return
	}
	toleranceInt = int(toleranceFloat * geometry.UM)

	slideAreaSpaceHoriFloat, err = ToFloat(slideAreaSpaceHoriString)
	if err != nil {
		return
	}
	slideAreaSpaceHoriInt = int(slideAreaSpaceHoriFloat * geometry.MM)

	slideAreaSpaceVertFloat, err = ToFloat(slideAreaSpaceVertString)
	if err != nil {
		return
	}
	slideAreaSpaceVertInt = int(slideAreaSpaceVertFloat * geometry.MM)

	printhead0OffsetXFloat, err = ToFloat(printhead0OffsetXString)
	if err != nil {
		return
	}
	printhead0OffsetXInt = int(printhead0OffsetXFloat * geometry.MM)

	printhead0OffsetYFloat, err = ToFloat(printhead0OffsetYString)
	if err != nil {
		return
	}
	printhead0OffsetYInt = int(printhead0OffsetYFloat * geometry.MM)

	printhead1OffsetXFloat, err = ToFloat(printhead1OffsetXString)
	if err != nil {
		return
	}
	printhead1OffsetXInt = int(printhead1OffsetXFloat * geometry.MM)

	printhead1OffsetYFloat, err = ToFloat(printhead1OffsetYString)
	if err != nil {
		return
	}
	printhead1OffsetYInt = int(printhead1OffsetYFloat * geometry.MM)

	return
}

// }}}

func build(
	step int,
	cycleCount int,
	printheadArray *printhead.Array,
	subs *substrate.Substrate,
	tolerance int,
	motorPath string,
	motorSpeed string,
	motorAccel string,
	printhead0Path string,
	buildProgressbar *widgets.QProgressBar,
	mode byte,
	lotc chan int,
	stepwise *bool,
	preview bool,
) {
	//filePath, err := uiutil.FilePath()
	//if err != nil {
	//uiutil.MessageBoxError(err.Error())
	//return
	//}
	//fmt.Println(len(slideArray.Slides[0].Spots))    // 91
	//fmt.Println(len(slideArray.Slides[0].Spots[0])) //119
	//fmt.Printf("%#v\n", slideArray.Slides[0].Spots[0][2].Reagents)
	//for _, r := range slideArray.Slides[0].Spots[0][2].Reagents {
	//fmt.Printf(r.Reagent.Name)
	//}

	bin := formation.NewBin(
		cycleCount,
		formation.NewMotorConf(
			motorPath,
			motorSpeed,
			motorAccel,
		),
		formation.NewPrintheadConf(
			printhead0Path,
			"1",
			"2560",
			"320",
		),
		mode,
		subs,
		printheadArray,
	)
	fmt.Println("create bin", bin)
	img := image.NewRGBA(image.Rect(0, 0, subs.Width, subs.Height+1))
	scene.image = img
	countc := bin.Build(step, lotc, img, paintedc, preview)
	go func() {
		for count := range countc {
			buildProgressbar.SetValue(count * buildProgressbar.Maximum() / cycleCount)
			//buildProgressbar.SetValue(count * buildProgressbar.Maximum())
			//if stepwise {
			if preview {
				scene.UpdatePixmap()
			} else {
				go func() {
					paintedc <- struct{}{}
					fmt.Println(">>> printedc")
				}()
			}
			//}
		}
		//close(lotc)
		//lotc = make(chan int)
		*stepwise = false
		fmt.Println("reinitialize paintedc")
		//close(paintedc)
		//paintedc = make(chan struct{})
	}()
	return

	// reagent first mode{{{

	//count := -1
	//sum := slideArray.ReagentCount()
	//fmt.Println("reagents sum:", sum)
	//go func() {
	//for cycleIndex := 0; cycleIndex < cycleCount; cycleIndex++ {
	//fmt.Println("loop cycle", cycleIndex)
	//for pi, p := range printheadArray.Printheads {
	//dataMap := map[int]string{}
	//for _, row := range p.Rows {
	//spots := slideArray.ReagentMap[cycleIndex][row.Reagent.Name]
	//if len(spots) == 0 {
	//continue
	//}
	//for _, v := range spots {
	//fmt.Println("s", v.Reagents[cycleIndex].Reagent.Name)
	//}
	//target := MostLeftSpot(cycleIndex, spots)
	//fmt.Println("row", pi, row.Index, row.Reagent.Name)
	//for target != nil {
	////time.Sleep(1 * time.Second)
	//fmt.Println("next spot", pi, row.Reagent.Name, target.Pos.X, target.Pos.Y)
	////fmt.Println("before", row.Index, p.Pos.X, p.Pos.Y)
	//p.UpdatePos(target.Pos.X, target.Pos.Y, row)
	////fmt.Println("after", row.Index, p.Pos.X, p.Pos.Y)

	//dataBinSlice := make([]string, 1280)
	//for index := range dataBinSlice {
	//dataBinSlice[index] = "0"
	//}
	//printable := false
	//// try print
	//// spot over nozzel is more effective
	//// but the dataBinSlice is overwrite every time
	//for _, nozzle := range row.Nozzles {
	//for _, spot := range spots {
	//if spot.Reagents[cycleIndex].Printed {
	//continue
	//}
	//if nozzle.Pos.Equal(spot.Pos) &&
	//row.Reagent.Equal(spot.Reagents[cycleIndex].Reagent) {
	//count += 1
	//printable = true
	//dataBinSlice[nozzle.Index] = "1"
	//spot.Reagents[cycleIndex].Printed = true
	//buildProgressbar.SetValue(count * buildProgressbar.Maximum() / sum)
	//fmt.Printf(
	//"spot printed: at(%d, %d), reagent %q, nozzle %v\n",
	//spot.Pos.X, spot.Pos.Y,
	//spot.Reagents[cycleIndex].Reagent.Name,
	//nozzle.Index,
	//)
	//// TODO: check the nozzles in other row for high resolution printing
	//}
	//}
	//}
	//if printable {
	//dataHexSlice := make([]string, 160)
	//for i := 0; i < len(dataBinSlice); i += 8 {
	//value, _ := strconv.ParseInt(strings.Join(dataBinSlice[i:i+8], ""), 2, 64)
	//dataHexSlice = append(dataHexSlice, fmt.Sprintf("%02x", value))
	//}
	//data := strings.Join(dataHexSlice, "")
	//dataMap[pi] = data
	//x, y := RawPos(target.Pos.X, target.Pos.Y)
	//bin.AddFormation(
	//cycleIndex, x, y, dataMap[0], dataMap[1],
	//)
	//fmt.Printf(">>> move to (%d, %d)\n", target.Pos.X, target.Pos.Y)
	//fmt.Printf(">>> print: row: %d, data: %#v\n", row.Index, dataBinSlice[:8])
	//fmt.Printf(">>> data: %#v\n", data[:5])
	//}
	//target = MostLeftSpot(cycleIndex, spots)
	//}
	//}
	//}
	//}
	//buildProgressbar.SetValue(buildProgressbar.Maximum())

	//filePath := "test.bin"
	//fmt.Printf("%#v\n", bin)
	////go func() {
	//file, err := os.Create(filePath)
	//defer file.Close()
	//if err != nil {
	//fmt.Println(err)
	//}
	//encoder := gob.NewEncoder(file)
	//encoder.Encode(bin)
	////}()
	//}()

	// }}}

	// normal mode

	imageIndex := 0

	go func() {
		for cycleIndex := 0; cycleIndex < cycleCount; cycleIndex++ {
			//for cycleIndex := 0; cycleIndex < 2; cycleIndex++ {
			img := image.NewRGBA(image.Rect(0, 0, subs.Width, subs.Height+1))
			fmt.Println("cycle ", cycleIndex)
			stripSum := subs.Strip()
			for stripCount := 0; stripCount < stripSum; stripCount++ {

				fmt.Println("strip ", stripCount)
				posx := stripCount * 1280
				posy := subs.Top()

				rowIndex := 3

				printheadArray.MoveBottomRow(rowIndex, posx, posy)
				fmt.Printf("move row %v to %v %v\n", rowIndex, posx, posy)

				for printheadArray.Top() >= subs.Bottom() {
					dataMap, count := genData(cycleIndex, printheadArray, subs, img, &imageIndex)
					if count > 0 {
						if DEBUGABLE {
							fmt.Println("data downward #1", count, dataMap)
						}
						x, y := RawPos(
							printheadArray.SightBottom.Pos.X,
							printheadArray.SightBottom.Pos.Y,
						)
						bin.AddFormation(
							cycleIndex, x, y, dataMap,
						)
					}
					posy -= step
					printheadArray.MoveBottomRow(rowIndex, posx, posy)
				}

				if step == 1 {
					continue
				}

				//posx -= 1
				posy = subs.Bottom()
				// distance is integer multiple of 4
				// so that move one row will match the others
				rowIndex = 2
				printheadArray.MoveTopRow(rowIndex, posx, posy)
				fmt.Printf("move row %v to %v %v\n", rowIndex, posx, posy)

				for printheadArray.Bottom() <= subs.Top() {
					dataMap, count := genData(cycleIndex, printheadArray, subs, img, &imageIndex)
					if count > 0 {
						if DEBUGABLE {
							fmt.Println("data upward #2", count, dataMap)
						}
						// use the bottom position
						// sinc the offset is bottomed
						x, y := RawPos(
							printheadArray.SightBottom.Pos.X,
							printheadArray.SightBottom.Pos.Y,
						)
						bin.AddFormation(
							cycleIndex, x, y, dataMap,
						)
					}
					posy += step
					printheadArray.MoveTopRow(rowIndex, posx, posy)
				}

				//posx -= 1
				posy = subs.Top()
				rowIndex = 1
				printheadArray.MoveBottomRow(rowIndex, posx, posy)
				fmt.Printf("move row %v to %v %v\n", rowIndex, posx, posy)

				for printheadArray.Top() >= subs.Bottom() {
					dataMap, count := genData(cycleIndex, printheadArray, subs, img, &imageIndex)
					if count > 0 {
						if DEBUGABLE {
							fmt.Println("data downward #3", count, dataMap)
						}
						x, y := RawPos(
							printheadArray.SightBottom.Pos.X,
							printheadArray.SightBottom.Pos.Y,
						)
						bin.AddFormation(
							cycleIndex, x, y, dataMap,
						)
					}
					posy -= step
					printheadArray.MoveBottomRow(rowIndex, posx, posy)
				}

				//posx -= 1
				posy = subs.Bottom()
				rowIndex = 0
				printheadArray.MoveTopRow(rowIndex, posx, posy)
				fmt.Printf("move row %v to %v %v\n", rowIndex, posx, posy)

				for printheadArray.Bottom() <= subs.Top() {
					dataMap, count := genData(cycleIndex, printheadArray, subs, img, &imageIndex)
					if count > 0 {
						if DEBUGABLE {
							fmt.Println("data upward #4", count, dataMap)
						}
						// use the bottom position
						// sinc the offset is bottomed
						x, y := RawPos(
							printheadArray.SightBottom.Pos.X,
							printheadArray.SightBottom.Pos.Y,
						)
						bin.AddFormation(
							cycleIndex, x, y, dataMap,
						)
					}
					posy += step
					printheadArray.MoveTopRow(rowIndex, posx, posy)
				}
			}
			buildProgressbar.SetValue((cycleIndex + 1) * buildProgressbar.Maximum() / cycleCount)
		}
		err := bin.SaveToFile("test.bin")
		if err != nil {
			fmt.Println(err)
			//uiutil.MessageBoxError(err.Error())
		}
	}()

}

func genData(
	cycleIndex int,
	printheadArray *printhead.Array,
	subs *substrate.Substrate,
	img *image.RGBA,
	imageIndex *int,
) ([]string, int) {
	count := 0
	dataSlice := make([][]string, printheadArray.PrintheadCount)
	for _, nozzle := range printheadArray.Nozzles {
		if dataSlice[nozzle.Printhead.Index] == nil {
			dataSlice[nozzle.Printhead.Index] = make([]string, 1280)
		}
		dataSlice[nozzle.Printhead.Index][nozzle.Index] = "0"
		if nozzle.Reagent.Equal(reagent.Nil) {
			continue
		}
		//fmt.Println(nozzle.Pos.X, nozzle.Pos.Y, subs.Width, subs.Height)
		if nozzle.Pos.Y >= subs.Height ||
			nozzle.Pos.Y < 0 ||
			nozzle.Pos.X >= subs.Width ||
			nozzle.Pos.X < 0 {
			continue
		}
		spot := subs.Spots[nozzle.Pos.Y][nozzle.Pos.X]
		if spot == nil || cycleIndex > len(spot.Reagents)-1 {
			//fmt.Println("not enough reagents")
			continue
		}
		if spot != nil &&
			nozzle.Reagent.Equal(spot.Reagents[cycleIndex]) {
			count += 1
			dataSlice[nozzle.Printhead.Index][nozzle.Index] = "1"

			if DEBUGABLE {
				//fmt.Printf(" | printing ", nozzle.Reagent.Name, nozzle.Pos.X, nozzle.Pos.Y)
			}
			if IMAGABLE {
				img.Set(nozzle.Pos.X, subs.Height-nozzle.Pos.Y, nozzle.Reagent.Color)
			}
		}
	}

	output := make([]string, printheadArray.PrintheadCount)
	if count > 0 {
		for deviceIndex, dataBinSlice := range dataSlice {
			dataHexSlice := make([]string, 160)
			for i := 0; i < len(dataBinSlice); i += 8 {
				value, _ := strconv.ParseInt(strings.Join(dataBinSlice[i:i+8], ""), 2, 64)
				dataHexSlice = append(dataHexSlice, fmt.Sprintf("%02x", value))
			}
			output[deviceIndex] = strings.Join(dataHexSlice, "")
			if DEBUGABLE {
				fmt.Println("print device", deviceIndex)
				fmt.Printf("data: %#v\n", dataBinSlice[:16])
				fmt.Printf("linebuffer: %#v\n", output[deviceIndex][:8])
			}
		}

		if IMAGABLE {
			outputFile, _ := os.Create(fmt.Sprintf("output/%06d.%03d.png", *imageIndex, cycleIndex))
			png.Encode(outputFile, img)
			outputFile.Close()
			*imageIndex = *imageIndex + 1
		}
	}
	return output, count
}

// genData{{{
//func genData2(
//cycleIndex int,
//slideArray *slide.Array,
//printheadArray *printhead.Array,
//tolerance int,
//) (map[int]string, int) {
//count := 0
//dataSlice := map[int]string{}
//for printheadIndex, printhead := range printheadArray.Printheads {
//printable := false
//data := make([]string, 1280)
//dataStringSlide := make([]string, 160)
//for _, row := range printhead.Rows {
//for _, nozzle := range row.Nozzles {
//data[nozzle.Index] = "0"
//spots := slideArray.SpotsIn(
//printheadArray.Top(),
//printheadArray.Right(),
//printheadArray.Bottom(),
//printheadArray.Left(),
//)
////if len(spots) > 0 {
////fmt.Println("spots count", len(spots))
////}
//for _, spot := range spots {
//if spot == nil {
//continue
//}
//if nozzle.IsAvailable(
//spot.Pos.X,
//spot.Pos.Y,
//tolerance,
//) {
//if spot.Reagents[cycleIndex].Reagent.Name == row.Reagent.Name {
//count += 1
//spot.Reagents[cycleIndex].Printed = true
//data[nozzle.Index] = "0"
//printable = true
//}
//}

//}

////for _, slide := range slideArray.Slides {
////for _, spots := range slide.Spots {
////for _, spot := range spots {
////if spot == nil {
////continue
////}
////if nozzle.IsAvailable(
////spot.Pos.X,
////spot.Pos.Y,
////tolerance,
////) {
////if spot.Reagents[cycleIndex].Reagent.Name == row.Reagent.Name {
////count += 1
////spot.Reagents[cycleIndex].Printed = true
////data[nozzle.Index] = "0"
////printable = true
////}
////}

////}
////}
////}

//}
//}
//if printable {
//for i := 0; i < len(data); i += 8 {
//value, _ := strconv.ParseInt(strings.Join(data[i:i+8], ""), 2, 64)
//dataStringSlide = append(dataStringSlide, fmt.Sprintf("%02x", value))
//}
//dataSlice[printheadIndex] = strings.Join(dataStringSlide, "")
//}
//}
//return dataSlice, count
//}

// }}}

func RawPos(
	posx int,
	posy int,
) (string, string) {
	x := offsetX - geometry.Mm(posx)
	y := offsetY - geometry.Mm(posy)
	if DEBUGABLE {
		fmt.Println("move to", x, y)
	}
	return fmt.Sprintf("%.6f", x), fmt.Sprintf("%.6f", y)
}

func MostLeftSpot(cycleIndex int, spots []*slide.Spot) *slide.Spot {
	var target *slide.Spot
	for _, spot := range spots {
		if spot.Reagents[cycleIndex].Printed {
			continue
		}
		if target == nil {
			target = spot
		} else {
			if spot.Pos.AtLeft(target.Pos) {
				target = spot
			}
		}
	}
	return target
}

func getStep(stepInput *widgets.QLineEdit) int {
	result, err := strconv.Atoi(stepInput.Text())
	if err != nil {
		stepInput.SetText("1")
		return 1
	}
	return result
}
