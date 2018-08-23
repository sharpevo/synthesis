package main

import (
	"github.com/therecipe/qt/widgets"
)

type InstructionDetail struct {
	GroupBox  *widgets.QGroupBox
	lineInput *widgets.QLineEdit
	lineLabel *widgets.QLabel
}

func NewInstructionDetail() *InstructionDetail {
	detail := InstructionDetail{}
	detail.lineInput = widgets.NewQLineEdit(nil)
	detail.lineLabel = widgets.NewQLabel2("-", nil, 0)

	detail.GroupBox = widgets.NewQGroupBox2("Instruction", nil)
	detailLayout := widgets.NewQGridLayout2()
	detailLayout.AddWidget(detail.lineInput, 0, 0, 0)
	detailLayout.AddWidget(detail.lineLabel, 1, 0, 0)
	detail.GroupBox.SetLayout(detailLayout)
	return &detail
}

func (d *InstructionDetail) Refresh(line string) {
	d.lineLabel.SetText(line)
}
