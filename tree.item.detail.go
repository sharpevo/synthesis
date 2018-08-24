package main

import (
	"fmt"
	"github.com/therecipe/qt/widgets"
	"strings"
)

const (
	TYPE_INS      = "Instruction"
	TYPE_SET      = "Instruction Set"
	INST_SET_SYNC = "IMPORT"
	INST_SET_ASYN = "ASYNC"
)

type InstructionDetail struct {
	treeItem        *widgets.QTreeWidgetItem
	instructionList []string
	GroupBox        *widgets.QGroupBox
	titleInput      *widgets.QLineEdit
	lineInput       *widgets.QLineEdit
	typeInput       *widgets.QComboBox
	instInput       *widgets.QComboBox
	argsInput       *widgets.QLineEdit
	saveButton      *widgets.QPushButton
}

func NewInstructionDetail() *InstructionDetail {

	typeLabel := widgets.NewQLabel2("Type", nil, 0)
	titleLabel := widgets.NewQLabel2("Title", nil, 0)
	lineLabel := widgets.NewQLabel2("Line", nil, 0)
	instLabel := widgets.NewQLabel2("Instruction", nil, 0)
	argsLabel := widgets.NewQLabel2("Arguments", nil, 0)

	d := InstructionDetail{}
	d.typeInput = widgets.NewQComboBox(nil)
	d.typeInput.AddItems([]string{TYPE_INS, TYPE_SET})
	d.typeInput.ConnectCurrentTextChanged(d.onInstructionTypeChanged)

	d.titleInput = widgets.NewQLineEdit(nil)
	d.lineInput = widgets.NewQLineEdit(nil)

	d.instInput = widgets.NewQComboBox(nil)
	for k, _ := range InstructionMap {
		if k != INST_SET_SYNC &&
			k != INST_SET_ASYN {
			d.instructionList = append(d.instructionList, k)
		}
	}
	d.instInput.AddItems(d.instructionList)

	d.argsInput = widgets.NewQLineEdit(nil)

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
	layout.AddWidget(instLabel, 3, 0, 0)
	layout.AddWidget(d.instInput, 3, 1, 0)
	layout.AddWidget(argsLabel, 4, 0, 0)
	layout.AddWidget(d.argsInput, 4, 1, 0)

	layout.AddWidget3(d.saveButton, 5, 0, 1, 2, 0)
	d.GroupBox.SetLayout(layout)
	return &d
}

func (d *InstructionDetail) saveInstruction() {
	if d.treeItem == nil {
		return
	}

	d.SetLineInput()

	d.treeItem.SetText(0, d.titleInput.Text())
	SetTreeItemData(d.treeItem, d.lineInput.Text())

}

func (d *InstructionDetail) Refresh(item *widgets.QTreeWidgetItem) {
	line := GetTreeItemData(item)
	d.treeItem = item
	d.titleInput.SetText(item.Text(0))
	d.lineInput.SetText(line)
	d.SetTypeInput()
	d.SetInstInput()
	d.SetArgsInput()
	d.saveButton.SetEnabled(true)
}

func (d *InstructionDetail) onInstructionTypeChanged(selected string) {
	switch selected {
	case TYPE_SET:
		d.instInput.Clear()
		d.instInput.AddItems([]string{INST_SET_SYNC, INST_SET_ASYN})
	default:
		d.instInput.Clear()
		d.instInput.AddItems(d.instructionList)
	}
}

func (d *InstructionDetail) Line() string {
	return d.lineInput.Text()
}

func (d *InstructionDetail) SetTypeInput() {
	if strings.HasPrefix(d.lineInput.Text(), INST_SET_SYNC) ||
		strings.HasPrefix(d.lineInput.Text(), INST_SET_ASYN) {
		d.typeInput.SetCurrentText(TYPE_SET)
	} else {
		d.typeInput.SetCurrentText(TYPE_INS)
	}
}

func (d *InstructionDetail) GetInstructionFromLine() string {
	list := strings.Split(d.Line(), " ")
	return list[0]
}

func (d *InstructionDetail) SetInstInput() {
	instruction := d.GetInstructionFromLine()
	for _, v := range d.instructionList {
		if instruction == v {
			d.instInput.SetCurrentText(instruction)
			break
		}
	}
}

func (d *InstructionDetail) GetArgumentsFromLine() string {
	instruction := d.GetInstructionFromLine()
	return strings.Trim(d.Line(), fmt.Sprintf("%s ", instruction))
}

func (d *InstructionDetail) SetArgsInput() {
	d.argsInput.SetText(d.GetArgumentsFromLine())
}

func (d *InstructionDetail) SetLineInput() {
	d.lineInput.SetText(fmt.Sprintf("%s %s", d.instInput.CurrentText(), d.argsInput.Text()))
}
