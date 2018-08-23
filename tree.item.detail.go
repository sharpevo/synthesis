package main

import (
	"github.com/therecipe/qt/widgets"
	"strings"
)

const (
	TYPE_INS = "Instruction"
	TYPE_SET = "Instruction Set"
)

type InstructionDetail struct {
	treeItem   *widgets.QTreeWidgetItem
	GroupBox   *widgets.QGroupBox
	titleInput *widgets.QLineEdit
	lineInput  *widgets.QLineEdit
	typeInput  *widgets.QComboBox
	saveButton *widgets.QPushButton
}

func NewInstructionDetail() *InstructionDetail {

	typeLabel := widgets.NewQLabel2("Type", nil, 0)
	titleLabel := widgets.NewQLabel2("Title", nil, 0)
	lineLabel := widgets.NewQLabel2("Instruction", nil, 0)

	d := InstructionDetail{}
	d.typeInput = widgets.NewQComboBox(nil)
	d.typeInput.AddItems([]string{TYPE_INS, TYPE_SET})
	d.titleInput = widgets.NewQLineEdit(nil)
	d.lineInput = widgets.NewQLineEdit(nil)

	d.saveButton = widgets.NewQPushButton2("SAVE", nil)
	d.saveButton.ConnectClicked(func(bool) { d.saveInstruction() })
	d.saveButton.SetEnabled(false)

	d.GroupBox = widgets.NewQGroupBox2("Instruction", nil)
	layout := widgets.NewQGridLayout2()
	layout.AddWidget(typeLabel, 0, 0, 0)
	layout.AddWidget(d.typeInput, 0, 1, 0)
	layout.AddWidget(titleLabel, 1, 0, 0)
	layout.AddWidget(d.titleInput, 1, 1, 0)
	layout.AddWidget(lineLabel, 2, 0, 0)
	layout.AddWidget(d.lineInput, 2, 1, 0)
	layout.AddWidget3(d.saveButton, 3, 0, 1, 2, 0)
	d.GroupBox.SetLayout(layout)
	return &d
}

func (d *InstructionDetail) saveInstruction() {
	if d.treeItem == nil {
		return
	}
	d.treeItem.SetText(0, d.titleInput.Text())
	SetTreeItemData(d.treeItem, d.lineInput.Text())
}

func (d *InstructionDetail) Refresh(item *widgets.QTreeWidgetItem) {
	line := GetTreeItemData(item)
	d.treeItem = item
	d.titleInput.SetText(item.Text(0))
	d.lineInput.SetText(line)
	d.SetTypeInput()
	d.saveButton.SetEnabled(true)
}

func (d *InstructionDetail) TypeInput() string {
	return d.typeInput.CurrentText()
}

func (d *InstructionDetail) SetTypeInput() {
	if strings.HasPrefix(d.lineInput.Text(), "ASYNC") ||
		strings.HasPrefix(d.lineInput.Text(), "IMPORT") {
		d.typeInput.SetCurrentText(TYPE_SET)
	} else {
		d.typeInput.SetCurrentText(TYPE_INS)
	}

}
