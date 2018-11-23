package sequence

import (
	"encoding/gob"
	"fmt"
	"github.com/therecipe/qt/widgets"
	"os"
	"posam/gui/uiutil"
	"posam/util/formation"
	"posam/util/geometry"
	"posam/util/printhead"
	"posam/util/reagent"
	"posam/util/slide"
	"strconv"
	"strings"
	//"time"
)

const (
	DPI_150 = "169.3"
	DPI_300 = "84.65"
	DPI_600 = "42.325"
)

var (
	offsetX = 0
	offsetY = 0
)

// const{{{

const (
	SEQUENCE_EXAMPLE = `GGGTCGGATGATCGGACACT
CATCATCTGGGTACAGCGGG
ATTATACAGTTTTGTCCAAT
`

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

func NewInputGroup() *widgets.QGroupBox {
	group := widgets.NewQGroupBox2("Parameters", nil)
	layout := widgets.NewQGridLayout2()
	group.SetLayout(layout)

	sequenceInput := widgets.NewQTextEdit(nil)
	sequenceInput.SetText(SEQUENCE_EXAMPLE)

	layout.AddWidget(sequenceInput, 0, 0, 0)

	// device group{{{

	deviceGroup := widgets.NewQGroupBox2("Device", nil)
	deviceLayout := widgets.NewQGridLayout2()
	deviceGroup.SetLayout(deviceLayout)
	layout.AddWidget(deviceGroup, 1, 0, 0)

	motorPathLabel := widgets.NewQLabel2("Motor path", nil, 0)
	motorPathInput := widgets.NewQLineEdit(nil)
	motorPathInput.SetText("/AOZTECH/Motor")
	motorSpeedLabel := widgets.NewQLabel2("Motor speed", nil, 0)
	motorSpeedInput := widgets.NewQLineEdit(nil)
	motorSpeedInput.SetText("10")
	motorAccelLabel := widgets.NewQLabel2("Motor acceleration", nil, 0)
	motorAccelInput := widgets.NewQLineEdit(nil)
	motorAccelInput.SetText("100")
	printhead0PathLabel := widgets.NewQLabel2("Printhead #1 path", nil, 0)
	printhead0PathInput := widgets.NewQLineEdit(nil)
	printhead0PathInput.SetText("/Ricoh-G5/Printer#1")
	printhead1PathLabel := widgets.NewQLabel2("Printhead #2 path", nil, 0)
	printhead1PathInput := widgets.NewQLineEdit(nil)
	printhead1PathInput.SetText("/Ricoh-G5/Printer#2")

	deviceLayout.AddWidget(motorPathLabel, 0, 0, 0)
	deviceLayout.AddWidget(motorPathInput, 0, 1, 0)
	deviceLayout.AddWidget(motorSpeedLabel, 1, 0, 0)
	deviceLayout.AddWidget(motorSpeedInput, 1, 1, 0)
	deviceLayout.AddWidget(motorAccelLabel, 2, 0, 0)
	deviceLayout.AddWidget(motorAccelInput, 2, 1, 0)
	deviceLayout.AddWidget(printhead0PathLabel, 3, 0, 0)
	deviceLayout.AddWidget(printhead0PathInput, 3, 1, 0)
	deviceLayout.AddWidget(printhead1PathLabel, 4, 0, 0)
	deviceLayout.AddWidget(printhead1PathInput, 4, 1, 0)

	// }}}

	// position gorup{{{

	positionGroup := widgets.NewQGroupBox2("Position (unit: mm)", nil)
	positionLayout := widgets.NewQGridLayout2()
	positionGroup.SetLayout(positionLayout)
	layout.AddWidget(positionGroup, 2, 0, 0)

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

	positionLayout.AddWidget(printhead0PositionLabel, 0, 0, 0)
	positionLayout.AddWidget(printhead0PositionXInput, 0, 1, 0)
	positionLayout.AddWidget(printhead0PositionYInput, 0, 2, 0)

	positionLayout.AddWidget(printhead1PositionLabel, 1, 0, 0)
	positionLayout.AddWidget(printhead1PositionXInput, 1, 1, 0)
	positionLayout.AddWidget(printhead1PositionYInput, 1, 2, 0)

	positionLayout.AddWidget(slide0PositionLabel, 2, 0, 0)
	positionLayout.AddWidget(slide0PositionXInput, 2, 1, 0)
	positionLayout.AddWidget(slide0PositionYInput, 2, 2, 0)

	positionLayout.AddWidget(slide1PositionLabel, 3, 0, 0)
	positionLayout.AddWidget(slide1PositionXInput, 3, 1, 0)
	positionLayout.AddWidget(slide1PositionYInput, 3, 2, 0)

	positionLayout.AddWidget(slide2PositionLabel, 4, 0, 0)
	positionLayout.AddWidget(slide2PositionXInput, 4, 1, 0)
	positionLayout.AddWidget(slide2PositionYInput, 4, 2, 0)

	// }}}

	// space group{{{

	spaceGroup := widgets.NewQGroupBox2("Space (unit: um)", nil)
	spaceLayout := widgets.NewQGridLayout2()
	spaceGroup.SetLayout(spaceLayout)
	layout.AddWidget(spaceGroup, 3, 0, 0)

	spaceLabel := widgets.NewQLabel2("Spot space", nil, 0)
	spacexInput := widgets.NewQLineEdit(nil)
	spaceyInput := widgets.NewQLineEdit(nil)

	spacexInput.SetText("169.3")
	spaceyInput.SetText("550.3")

	spaceLayout.AddWidget(spaceLabel, 0, 0, 0)
	spaceLayout.AddWidget(spacexInput, 0, 1, 0)
	spaceLayout.AddWidget(spaceyInput, 0, 2, 0)

	// }}}

	// reagent group{{{

	reagentGroup := widgets.NewQGroupBox2("Reagent", nil)
	reagentLayout := widgets.NewQGridLayout2()
	reagentGroup.SetLayout(reagentLayout)
	layout.AddWidget(reagentGroup, 4, 0, 0)

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
	printhead1Line0Input.SetText("Z")
	printhead1Line1Input.SetText("Z")
	printhead1Line2Input.SetText("-")
	printhead1Line3Input.SetText("-")

	reagentLayout.AddWidget(printhead0Line0Label, 0, 0, 0)
	reagentLayout.AddWidget(printhead0Line0Input, 0, 1, 0)

	reagentLayout.AddWidget(printhead0Line1Label, 1, 0, 0)
	reagentLayout.AddWidget(printhead0Line1Input, 1, 1, 0)

	reagentLayout.AddWidget(printhead0Line2Label, 2, 0, 0)
	reagentLayout.AddWidget(printhead0Line2Input, 2, 1, 0)

	reagentLayout.AddWidget(printhead0Line3Label, 3, 0, 0)
	reagentLayout.AddWidget(printhead0Line3Input, 3, 1, 0)

	reagentLayout.AddWidget(printhead1Line0Label, 4, 0, 0)
	reagentLayout.AddWidget(printhead1Line0Input, 4, 1, 0)

	reagentLayout.AddWidget(printhead1Line1Label, 5, 0, 0)
	reagentLayout.AddWidget(printhead1Line1Input, 5, 1, 0)

	reagentLayout.AddWidget(printhead1Line2Label, 6, 0, 0)
	reagentLayout.AddWidget(printhead1Line2Input, 6, 1, 0)

	reagentLayout.AddWidget(printhead1Line3Label, 7, 0, 0)
	reagentLayout.AddWidget(printhead1Line3Input, 7, 1, 0)

	// }}}

	// misc group{{{

	miscGroup := widgets.NewQGroupBox2("Misc", nil)
	miscLayout := widgets.NewQGridLayout2()
	miscGroup.SetLayout(miscLayout)
	layout.AddWidget(miscGroup, 5, 0, 0)

	toleranceLabel := widgets.NewQLabel2("Tolerance (um)", nil, 0)
	toleranceInput := widgets.NewQLineEdit(nil)
	toleranceInput.SetText("30")

	dpiLabel := widgets.NewQLabel2("Resolution", nil, 0)
	dpiInput := widgets.NewQComboBox(nil)
	dpiInput.AddItems([]string{
		DPI_150,
	})

	printhead0OffsetLabel := widgets.NewQLabel2("offset #0 (mm)", nil, 0)
	printhead0OffsetXInput := widgets.NewQLineEdit(nil)
	printhead0OffsetYInput := widgets.NewQLineEdit(nil)
	printhead0OffsetXInput.SetText("-30")
	printhead0OffsetYInput.SetText("-20")

	printhead1OffsetLabel := widgets.NewQLabel2("offset #1 (mm)", nil, 0)
	printhead1OffsetXInput := widgets.NewQLineEdit(nil)
	printhead1OffsetYInput := widgets.NewQLineEdit(nil)
	printhead1OffsetXInput.SetText("-90")
	printhead1OffsetYInput.SetText("-20")

	miscLayout.AddWidget(toleranceLabel, 0, 0, 0)
	miscLayout.AddWidget(toleranceInput, 0, 1, 0)
	miscLayout.AddWidget(dpiLabel, 1, 0, 0)
	miscLayout.AddWidget3(dpiInput, 1, 1, 1, 2, 0)
	miscLayout.AddWidget(printhead0OffsetLabel, 2, 0, 0)
	miscLayout.AddWidget(printhead0OffsetXInput, 2, 1, 0)
	miscLayout.AddWidget(printhead0OffsetYInput, 2, 2, 0)
	miscLayout.AddWidget(printhead1OffsetLabel, 3, 0, 0)
	miscLayout.AddWidget(printhead1OffsetXInput, 3, 1, 0)
	miscLayout.AddWidget(printhead1OffsetYInput, 3, 2, 0)

	// }}}

	// build button

	buildButton := widgets.NewQPushButton2("BUILD", nil)
	layout.AddWidget(buildButton, 6, 0, 0)

	buildProgressbar := widgets.NewQProgressBar(nil)
	layout.AddWidget(buildProgressbar, 7, 0, 0)

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
			slide0PositionX,
			slide0PositionY,
			slide1PositionX,
			slide1PositionY,
			slide2PositionX,
			slide2PositionY,
			spacex,
			spacey,
		)

		seqText := sequenceInput.ToPlainText()
		offsetX = printhead1OffsetX
		offsetY = printhead1OffsetY

		var step int
		var space int
		switch dpiInput.CurrentText() {
		case DPI_300:
			step = 126.975 * geometry.UM
			space = 84.65 * geometry.UM
		case DPI_600:
			step = 169.3 * geometry.UM
			space = 42.325 * geometry.UM
		default:
			step = 42.325 * geometry.UM
			space = 169.3 * geometry.UM
		}

		// create printhead{{{

		// rows of printhead 0

		nozzleSpace := int(169.3 * geometry.UM)

		row00 := printhead.NewRow(
			nozzleSpace,
			reagent.NewReagent(printhead0Line0Input.Text()),
		)
		row01 := printhead.NewRow(
			nozzleSpace,
			reagent.NewReagent(printhead0Line1Input.Text()),
		)
		row02 := printhead.NewRow(
			nozzleSpace,
			reagent.NewReagent(printhead0Line2Input.Text()),
		)
		row03 := printhead.NewRow(
			nozzleSpace,
			reagent.NewReagent(printhead0Line3Input.Text()),
		)
		rows0 := []*printhead.Row{row00, row01, row02, row03}

		// rows of printhead 1

		row10 := printhead.NewRow(
			nozzleSpace,
			reagent.NewReagent(printhead1Line0Input.Text()),
		)
		row11 := printhead.NewRow(
			nozzleSpace,
			reagent.NewReagent(printhead1Line1Input.Text()),
		)
		row12 := printhead.NewRow(
			nozzleSpace,
			reagent.NewReagent(printhead1Line2Input.Text()),
		)
		row13 := printhead.NewRow(
			nozzleSpace,
			reagent.NewReagent(printhead1Line3Input.Text()),
		)
		rows1 := []*printhead.Row{row10, row11, row12, row13}

		printhead0, err := printhead.NewPrinthead(
			printhead0OffsetX,
			printhead0OffsetY,
			//printhead0PositionX,
			//printhead0PositionY,
			rows0,
			nozzleSpace,
			550.3*geometry.UM,
			11.811*geometry.UM,
		)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}

		printhead1, err := printhead.NewPrinthead(
			//printhead1PositionX,
			//printhead1PositionY,
			printhead1OffsetX,
			printhead1OffsetY,
			rows1,
			nozzleSpace,
			550.3*geometry.UM,
			11.811*geometry.UM,
		)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		printheadArray := printhead.NewArray(printhead0, printhead1)

		// }}}
		fmt.Println("create printhead array", *printheadArray)

		// create slide array{{{

		//slide0 := slide.NewSlide(
		//slide0PositionX,
		//slide0PositionY,
		//spacex,
		//spacey,
		//)
		//slide1 := slide.NewSlide(
		//slide1PositionX,
		//slide1PositionY,
		//spacex,
		//spacey,
		//)
		//slide2 := slide.NewSlide(
		//slide2PositionX,
		//slide2PositionY,
		//spacex,
		//spacey,
		//)
		//slideArray := slide.NewArray(slide0, slide1, slide2)
		//slideArray := slide.NewDefaultArray(spacex, spacey, 3)
		slideArray := slide.NewDefaultArray(space, space, 3)

		cycleCount := 0
		for _, lineRaw := range strings.Split(seqText, "\n") {
			if lineRaw == "" {
				continue
			}
			spot := slide.NewSpot()
			line := strings.Trim(lineRaw, " ")
			bases := strings.Split(line, "")
			length := len(bases)
			if length > cycleCount {
				cycleCount = length
			}
			for _, base := range bases {
				r := reagent.NewReagent(base)
				spot.AddReagent(r)
				//if r.Name != reagent.Nil.Name {
				//spot.AddReagent(reagent.Activator)
				//}
			}
			slideArray.AddSpot(spot)
			//fmt.Println(spot.Pos.X, spot.Pos.Y, spot.Reagents)
			//fmt.Println(">>", spot.Reagents[0].Reagent.Name)
		}
		//fmt.Println(slideArray.ReagentMap)
		for c, spotmap := range slideArray.ReagentMap {
			for _, spot := range spotmap["A"] {
				fmt.Println(">>", c, spot.Pos.X, spot.Pos.Y)
			}
		}

		// }}}
		fmt.Println("create slide array", slideArray.AvailableSpots(), *slideArray.AvailableSpots()[0].Reagents[0])

		buildButton.SetVisible(false)
		buildProgressbar.SetVisible(true)

		build(
			step,
			cycleCount,
			printheadArray,
			slideArray,
			tolerance,
			motorPathInput.Text(),
			motorSpeedInput.Text(),
			motorAccelInput.Text(),
			printhead0PathInput.Text(),
			printhead1PathInput.Text(),
			buildProgressbar,
		)

	})

	return group
}

func ToFloat(inputString string) (float64, error) {
	if inputString == "" {
		return 0.0, fmt.Errorf("should not be null")
	}
	inputFloat, err := strconv.ParseFloat(inputString, 64)
	if err != nil {
		return 0.0, err
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
	slideArray *slide.Array,
	tolerance int,
	motorPath string,
	motorSpeed string,
	motorAccel string,
	printhead0Path string,
	printhead1Path string,
	buildProgressbar *widgets.QProgressBar,
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
			"1280",
			"160",
		),
		formation.NewPrintheadConf(
			printhead1Path,
			"1",
			"1280",
			"160",
		),
	)
	fmt.Println("create bin", bin)

	count := 0
	sum := slideArray.ReagentCount()
	fmt.Println("reagents sum:", sum)
	go func() {
		for cycleIndex := 0; cycleIndex < cycleCount; cycleIndex++ {
			fmt.Println("loop cycle", cycleIndex)

			//var spot *slide.Spot
			//spot = slideArray.NextSpotInVert(cycleIndex)
			//if spot == nil {
			//break
			//}
			//fmt.Println("spot", spot.Pos.X, spot.Pos.Y)
			//x, y := RawPos(spot.Pos.X, spot.Pos.Y)
			//bin.AddFormation( // move sight printhead to the top left of slide array
			//cycleIndex,
			//x, y,
			//"", "",
			//)
			//printheadArray.UpdatePos(spot.Pos.X, spot.Pos.Y)

			for pi, p := range printheadArray.Printheads {
				dataMap := map[int]string{}
				for _, row := range p.Rows {
					spots := slideArray.ReagentMap[cycleIndex][row.Reagent.Name]
					if len(spots) == 0 {
						continue
					}
					for _, v := range spots {
						fmt.Println("s", v.Reagents[cycleIndex].Reagent.Name)
					}
					target := MostLeftSpot(cycleIndex, spots)
					fmt.Println("row", pi, row.Index, row.Reagent.Name)
					for target != nil {
						//time.Sleep(1 * time.Second)
						fmt.Println("next spot", pi, row.Reagent.Name, target.Pos.X, target.Pos.Y)
						//fmt.Println("before", row.Index, p.Pos.X, p.Pos.Y)
						p.UpdatePos(target.Pos.X, target.Pos.Y, row)
						//fmt.Println("after", row.Index, p.Pos.X, p.Pos.Y)

						dataBinSlice := make([]string, 1280)
						printable := false
						// try print
						// spot over nozzel is more effective
						// but the dataBinSlice is overwrite every time
						for _, nozzle := range row.Nozzles {
							dataBinSlice[nozzle.Index] = "0"
							for _, spot := range spots {
								if spot.Reagents[cycleIndex].Printed {
									continue
								}
								if nozzle.Pos.Equal(spot.Pos) &&
									row.Reagent.Equal(spot.Reagents[cycleIndex].Reagent) {
									count += 1
									printable = true
									dataBinSlice[nozzle.Index] = "1"
									spot.Reagents[cycleIndex].Printed = true
									buildProgressbar.SetValue(count * buildProgressbar.Maximum() / sum)
									fmt.Println("spot printed", spot.Reagents[cycleIndex].Reagent.Name, spot.Pos.X, spot.Pos.Y)
									// TODO: check the nozzles in other row for high resolution printing
								}
							}
						}
						if printable {
							dataHexSlice := make([]string, 160)
							for i := 0; i < len(dataBinSlice); i += 8 {
								value, _ := strconv.ParseInt(strings.Join(dataBinSlice[i:i+8], ""), 2, 64)
								dataHexSlice = append(dataHexSlice, fmt.Sprintf("%02x", value))
							}
							dataMap[pi] = strings.Join(dataHexSlice, "")
							x, y := RawPos(target.Pos.X, target.Pos.Y)
							bin.AddFormation(
								cycleIndex, x, y, dataMap[0], dataMap[1],
							)
						}
						target = MostLeftSpot(cycleIndex, spots)
					}
				}
			}

			//for printheadArray.Top() > slideArray.Bottom() {
			//fmt.Println("downward")
			//for posy := spot.Pos.Y; printheadArray.Top() > slideArray.Bottom(); posy -= tolerance {
			//data, c := genData(cycleIndex, slideArray, printheadArray, tolerance)
			//if len(data) != 0 {
			//x, y := RawPos(spot.Pos.X, posy)
			//bin.AddFormation(
			//cycleIndex, x, y, data[0], data[1])
			//fmt.Println("spots detected", c)
			//count += c
			//buildProgressbar.SetValue(count * buildProgressbar.Maximum() / sum)
			//}
			//printheadArray.UpdatePos(spot.Pos.X, posy)
			////fmt.Println("next", spot.Pos.X, posy)
			//}

			//sight := printheadArray.SightPrinthead()
			//posx := spot.Pos.X + sight.RowOffset
			//printheadArray.UpdatePos(posx, sight.Pos.Y)

			//fmt.Println("upward", posx, sight.Pos.Y)
			//for posy := sight.Bottom(); sight.Bottom() <= slideArray.Top(); posy += tolerance {
			//data, c := genData(cycleIndex, slideArray, printheadArray, tolerance)
			//if len(data) != 0 {
			//x, y := RawPos(posx, posy)
			//bin.AddFormation(cycleIndex, x, y, data[0], data[1])
			//fmt.Println("spots detected", c)
			//count += c
			//buildProgressbar.SetValue(count * buildProgressbar.Maximum() / sum)
			//}
			//printheadArray.UpdatePos(posx, posy)
			////fmt.Println("next", posx, posy)
			//}

			//spot = slideArray.NextSpotInVert(cycleIndex)
			//if spot == nil {
			//break
			//}
			//fmt.Println("spot", spot.Pos.X, spot.Pos.Y)
			//printheadArray.UpdatePos(spot.Pos.X, spot.Pos.Y)
			//}
		}
		buildProgressbar.SetValue(buildProgressbar.Maximum())

		filePath := "test.bin"
		fmt.Printf("%#v\n", bin)
		//go func() {
		file, err := os.Create(filePath)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		}
		encoder := gob.NewEncoder(file)
		encoder.Encode(bin)
		//}()
	}()

}

func genData(
	cycleIndex int,
	slideArray *slide.Array,
	printheadArray *printhead.Array,
	tolerance int,
) (map[int]string, int) {
	count := 0
	dataSlice := map[int]string{}
	for printheadIndex, printhead := range printheadArray.Printheads {
		printable := false
		data := make([]string, 1280)
		dataStringSlide := make([]string, 160)
		for _, row := range printhead.Rows {
			for _, nozzle := range row.Nozzles {
				data[nozzle.Index] = "0"
				spots := slideArray.SpotsIn(
					printheadArray.Top(),
					printheadArray.Right(),
					printheadArray.Bottom(),
					printheadArray.Left(),
				)
				//if len(spots) > 0 {
				//fmt.Println("spots count", len(spots))
				//}
				for _, spot := range spots {
					if spot == nil {
						continue
					}
					if nozzle.IsAvailable(
						spot.Pos.X,
						spot.Pos.Y,
						tolerance,
					) {
						if spot.Reagents[cycleIndex].Reagent.Name == row.Reagent.Name {
							count += 1
							spot.Reagents[cycleIndex].Printed = true
							data[nozzle.Index] = "0"
							printable = true
						}
					}

				}

				//for _, slide := range slideArray.Slides {
				//for _, spots := range slide.Spots {
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
				//}
				//}

			}
		}
		if printable {
			for i := 0; i < len(data); i += 8 {
				value, _ := strconv.ParseInt(strings.Join(data[i:i+8], ""), 2, 64)
				dataStringSlide = append(dataStringSlide, fmt.Sprintf("%02x", value))
			}
			dataSlice[printheadIndex] = strings.Join(dataStringSlide, "")
		}
	}
	return dataSlice, count
}

func RawPos(
	posx int,
	posy int,
	//) (string, string) {
) (int, int) {
	return posx - offsetX, posy - offsetY
	//return fmt.Sprintf("%v", posx-offsetX), fmt.Sprintf("%v", posy-offsetY)
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
