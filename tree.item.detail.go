package main

import (
	"github.com/therecipe/qt/widgets"
)

type InstructionDetail struct {
	treeItem   *widgets.QTreeWidgetItem
	GroupBox   *widgets.QGroupBox
	titleInput *widgets.QLineEdit
	lineInput  *widgets.QLineEdit
	saveButton *widgets.QPushButton
}

func NewInstructionDetail() *InstructionDetail {
	d := InstructionDetail{}
	d.titleInput = widgets.NewQLineEdit(nil)
	d.lineInput = widgets.NewQLineEdit(nil)
	d.saveButton = widgets.NewQPushButton2("SAVE", nil)
	d.saveButton.ConnectClicked(func(bool) { d.saveInstruction() })
	d.saveButton.SetEnabled(false)

	d.GroupBox = widgets.NewQGroupBox2("Instruction", nil)
	layout := widgets.NewQGridLayout2()
	layout.AddWidget(d.titleInput, 0, 0, 0)
	layout.AddWidget(d.lineInput, 1, 0, 0)
	layout.AddWidget(d.saveButton, 2, 0, 0)
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
	d.saveButton.SetEnabled(true)
}
