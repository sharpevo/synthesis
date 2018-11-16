package sequence

import (
	"github.com/therecipe/qt/widgets"
)

// const{{{

const (
	SEQUENCE_EXAMPLE = `GGGTCGGATGATCGGACACT
CATCATCTGGGTACAGCGGG
ATTATACAGTTTTGTCCAAT
CTATCTTGGAGGGGTAGGCG
AGGCTGGCCATGTTGTCTTA
ACTTTATGATGCGTAAGCAC
CAGCCTCAACCGCTCTGCAA
CATGCTCCTATCGTAGGAAG
CAGGAGTCCATTCCGTGCTG
ATTGCCGTTAATCGGCAGGA
AGAGTGCCGGAACACTGTTG
TCACGAGGGGGCAAAGAAAG
ATTTGCCGGGGCGTGTCCTG
GGATGCTGACACGTCGTGTT
GTATCTACTTGACTACGGCC
GGTTTGAAGTAAGACCCCCC
CGTCTCGGCCCGTAATCTCC
TGATCCAAATCGATTAATGT
AAGATCCCAGTTTTTTAGAG
AATCACTGCATTGCGAAAAA
CTGCACGATTATGGGGTGAG
GTCCGACCAGGGGTCTATCT
CTGGAAATGCCTGGGCGGTG
TTCCAAGTGATAGCTACGAA
GTTCCGTTATGCCGAGGAAG
AGATCCACGGCTCGTCAGAC
GATGAATTAGCGGAGGATCC
GGCACGGTAAGTTCCCACGC
GCGCTCGAGACGAACACTAA
CGATAGATGAATGGGCACCT
CCAGACCGGAGTTGGAGGAG
GTTTGCTCCTCTTCACTCCG
TCAAGGCTGATATCACCAAT
CAGCATCTTAACTCCAGGAC
GTATCTCTCGTAACATGCTA
ATCACGAGATGAAAGTCTGG
TTCTCGTTCCACCCAGTCGT
GAAGCTCAACACATAGCAAC
GACCGGACGAGAAAACTCCG
TACTCCCTCAAGTAAGTCTA
CTAGACCGCAGCAAAATCGT
TCACTTTCGCGCGCACAGGG
AGGGTCGGACTTCTAGGTAG
GATCAGACACCTCATCACGA
GTGCCTCCTGCCCTAGTCGA
CCAACATGTGCCAACGATTA
ATGAGCTGAAGACAGAGGGC
ATAGCCGCTGGCGTTCGTGG
GCGTAGCAAAGGGGCGGAGT
CAGTTATTTCAGAGGTACCG
`
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
	printheadPathLabel := widgets.NewQLabel2("Printhead path:", nil, 0)
	printheadPathInput := widgets.NewQLineEdit(nil)
	printheadPathInput.SetText("/Ricoh-G5/Printer#1")

	deviceLayout.AddWidget(motorPathLabel, 0, 0, 0)
	deviceLayout.AddWidget(motorPathInput, 0, 1, 0)
	deviceLayout.AddWidget(motorSpeedLabel, 1, 0, 0)
	deviceLayout.AddWidget(motorSpeedInput, 1, 1, 0)
	deviceLayout.AddWidget(motorAccelLabel, 2, 0, 0)
	deviceLayout.AddWidget(motorAccelInput, 2, 1, 0)
	deviceLayout.AddWidget(printheadPathLabel, 3, 0, 0)
	deviceLayout.AddWidget(printheadPathInput, 3, 1, 0)

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
	slide1PositionXInput.SetText("-26")
	slide1PositionYInput.SetText("0")
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

	printhead0line0Label := widgets.NewQLabel2("Row A of Printhead #0", nil, 0)
	printhead0line1Label := widgets.NewQLabel2("Row B of Printhead #0", nil, 0)
	printhead0line2Label := widgets.NewQLabel2("Row C of Printhead #0", nil, 0)
	printhead0line3Label := widgets.NewQLabel2("Row D of Printhead #0", nil, 0)

	printhead1line0Label := widgets.NewQLabel2("Row A of Printhead #1", nil, 0)
	printhead1line1Label := widgets.NewQLabel2("Row B of Printhead #1", nil, 0)
	printhead1line2Label := widgets.NewQLabel2("Row C of Printhead #1", nil, 0)
	printhead1line3Label := widgets.NewQLabel2("Row D of Printhead #1", nil, 0)

	printhead0line0Input := widgets.NewQLineEdit(nil)
	printhead0line1Input := widgets.NewQLineEdit(nil)
	printhead0line2Input := widgets.NewQLineEdit(nil)
	printhead0line3Input := widgets.NewQLineEdit(nil)

	printhead1line0Input := widgets.NewQLineEdit(nil)
	printhead1line1Input := widgets.NewQLineEdit(nil)
	printhead1line2Input := widgets.NewQLineEdit(nil)
	printhead1line3Input := widgets.NewQLineEdit(nil)

	printhead0line0Input.SetText("T")
	printhead0line1Input.SetText("Z")
	printhead0line2Input.SetText("Z")
	printhead0line3Input.SetText("G")
	printhead1line0Input.SetText("A")
	printhead1line1Input.SetText("Z")
	printhead1line2Input.SetText("Z")
	printhead1line3Input.SetText("C")

	reagentLayout.AddWidget(printhead0line0Label, 0, 0, 0)
	reagentLayout.AddWidget(printhead0line0Input, 0, 1, 0)

	reagentLayout.AddWidget(printhead0line1Label, 1, 0, 0)
	reagentLayout.AddWidget(printhead0line1Input, 1, 1, 0)

	reagentLayout.AddWidget(printhead0line2Label, 2, 0, 0)
	reagentLayout.AddWidget(printhead0line2Input, 2, 1, 0)

	reagentLayout.AddWidget(printhead0line3Label, 3, 0, 0)
	reagentLayout.AddWidget(printhead0line3Input, 3, 1, 0)

	reagentLayout.AddWidget(printhead1line0Label, 4, 0, 0)
	reagentLayout.AddWidget(printhead1line0Input, 4, 1, 0)

	reagentLayout.AddWidget(printhead1line1Label, 5, 0, 0)
	reagentLayout.AddWidget(printhead1line1Input, 5, 1, 0)

	reagentLayout.AddWidget(printhead1line2Label, 6, 0, 0)
	reagentLayout.AddWidget(printhead1line2Input, 6, 1, 0)

	reagentLayout.AddWidget(printhead1line3Label, 7, 0, 0)
	reagentLayout.AddWidget(printhead1line3Input, 7, 1, 0)

	// }}}

	// build button

	buildButton := widgets.NewQPushButton2("BUILD", nil)
	layout.AddWidget(buildButton, 5, 0, 0)

	buildProgressbar := widgets.NewQProgressBar(nil)
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
	})

	return group
}
