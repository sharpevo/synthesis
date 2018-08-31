package instree

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"posam/gui/tree"
	"posam/interpreter"
	"sort"
	"strings"
)

const (
	TYPE_INS      = "Instruction"
	TYPE_SET      = "Instruction Set"
	INST_SET_SYNC = "IMPORT"
	INST_SET_ASYN = "ASYNC"

	TYPE_SET_ONCE = "ONCE"
	TYPE_SET_COND = "COND"
	TYPE_SET_LOOP = "LOOP"
)

type InstructionDetail struct {
	treeItem        *widgets.QTreeWidgetItem
	instructionList []string
	GroupBox        *widgets.QGroupBox
	titleInput      *widgets.QLineEdit
	lineInput       *widgets.QLineEdit
	typeInput       *widgets.QComboBox
	instInput       *widgets.QComboBox
	devInput        *widgets.QComboBox
	argsInput       *widgets.QLineEdit
	waveformGroup   *widgets.QGroupBox
	saveButton      *widgets.QPushButton
}

func NewInstructionDetail(instructionMap interpreter.InstructionMapt) *InstructionDetail {

	typeLabel := widgets.NewQLabel2("Type", nil, 0)
	titleLabel := widgets.NewQLabel2("Title", nil, 0)
	instLabel := widgets.NewQLabel2("Instruction", nil, 0)
	devLabel := widgets.NewQLabel2("Device", nil, 0)
	argsLabel := widgets.NewQLabel2("Arguments", nil, 0)

	d := InstructionDetail{}
	d.typeInput = widgets.NewQComboBox(nil)
	d.typeInput.AddItems([]string{TYPE_INS, TYPE_SET})
	d.typeInput.ConnectCurrentTextChanged(d.onInstructionTypeChanged)

	d.titleInput = widgets.NewQLineEdit(nil)
	d.lineInput = widgets.NewQLineEdit(nil)
	d.lineInput.SetVisible(false)

	d.instInput = widgets.NewQComboBox(nil)
	for k, _ := range instructionMap {
		if k != INST_SET_SYNC &&
			k != INST_SET_ASYN {
			d.instructionList = append(d.instructionList, k)
		}
	}
	sort.Sort(sort.StringSlice(d.instructionList))
	d.instInput.AddItems(d.instructionList)
	d.instInput.ConnectCurrentTextChanged(d.onInstructionChanged)

	d.devInput = widgets.NewQComboBox(nil)
	d.argsInput = widgets.NewQLineEdit(nil)

	// waveform group

	d.waveformGroup = widgets.NewQGroupBox2("WaveForm Builder", nil)

	waveformLineTimeLabel := widgets.NewQLabel2("Time", nil, 1)
	waveformLineTimeLabel.SetAlignment(core.Qt__AlignCenter)
	waveformLineVoltageLabel := widgets.NewQLabel2("Percentage", nil, 1)
	waveformLineVoltageLabel.SetAlignment(core.Qt__AlignCenter)
	waveformFallLabel := widgets.NewQLabel2("Fall:", nil, 0)
	waveformHoldLabel := widgets.NewQLabel2("Hold:", nil, 0)
	waveformRisingLabel := widgets.NewQLabel2("Rising:", nil, 0)
	waveformWaitLabel := widgets.NewQLabel2("Wait:", nil, 0)
	waveformMnLabel := widgets.NewQLabel2("Mn:", nil, 0)
	waveformVoltageLabel := widgets.NewQLabel2("Voltage:", nil, 0)

	waveformFallTimeInput := widgets.NewQDoubleSpinBox(nil)
	waveformFallTimeInput.SetMaximum(100)
	waveformFallPercentageInput := widgets.NewQDoubleSpinBox(nil)
	waveformFallPercentageInput.SetMaximum(100)
	waveformHoldTimeInput := widgets.NewQDoubleSpinBox(nil)
	waveformHoldTimeInput.SetMaximum(100)
	waveformHoldPercentageInput := widgets.NewQDoubleSpinBox(nil)
	waveformHoldPercentageInput.SetMaximum(100)
	waveformRisingTimeInput := widgets.NewQDoubleSpinBox(nil)
	waveformRisingTimeInput.SetMaximum(100)
	waveformRisingPercentageInput := widgets.NewQDoubleSpinBox(nil)
	waveformRisingPercentageInput.SetMaximum(100)
	waveformWaitTimeInput := widgets.NewQDoubleSpinBox(nil)
	waveformWaitTimeInput.SetMaximum(100)
	waveformWaitPercentageInput := widgets.NewQDoubleSpinBox(nil)
	waveformWaitPercentageInput.SetMaximum(100)
	waveformMnInput := widgets.NewQSpinBox(nil)
	waveformMnInput.SetMaximum(100)
	waveformVoltageInput := widgets.NewQDoubleSpinBox(nil)
	waveformVoltageInput.SetMaximum(100)

	waveformGenerateButton := widgets.NewQPushButton2("INSERT", nil)
	waveformGenerateButton.ConnectClicked(func(bool) {
		argumentList := []string{
			"VAR",
			"HEADBOARD", // headboard
			"ROW",       // row
			fmt.Sprintf("%.2f", waveformVoltageInput.Value()),
			"COUNT", // segment count

			fmt.Sprintf("%.2f", waveformFallTimeInput.Value()),
			fmt.Sprintf("%.2f", waveformFallPercentageInput.Value()),
			fmt.Sprintf("%.2f", waveformHoldPercentageInput.Value()),

			fmt.Sprintf("%.2f", waveformHoldTimeInput.Value()),
			fmt.Sprintf("%.2f", waveformHoldPercentageInput.Value()),
			fmt.Sprintf("%.2f", waveformRisingPercentageInput.Value()),

			fmt.Sprintf("%.2f", waveformRisingTimeInput.Value()),
			fmt.Sprintf("%.2f", waveformRisingPercentageInput.Value()),
			fmt.Sprintf("%.2f", waveformWaitPercentageInput.Value()),

			fmt.Sprintf("%.2f", waveformWaitTimeInput.Value()),
			fmt.Sprintf("%.2f", waveformWaitPercentageInput.Value()),
			fmt.Sprintf("%.2f", waveformWaitPercentageInput.Value()),

			fmt.Sprintf("%d", waveformMnInput.Value()),
		}

		d.argsInput.SetText(strings.Join(argumentList, " "))
	})

	waveformLayout := widgets.NewQGridLayout2()

	waveformLayout.AddWidget(waveformLineTimeLabel, 0, 1, 0)
	waveformLayout.AddWidget(waveformLineVoltageLabel, 0, 2, 0)

	waveformLayout.AddWidget(waveformFallLabel, 1, 0, 0)
	waveformLayout.AddWidget(waveformFallTimeInput, 1, 1, 0)
	waveformLayout.AddWidget(waveformFallPercentageInput, 1, 2, 0)
	waveformLayout.AddWidget(waveformHoldLabel, 2, 0, 0)
	waveformLayout.AddWidget(waveformHoldTimeInput, 2, 1, 0)
	waveformLayout.AddWidget(waveformHoldPercentageInput, 2, 2, 0)
	waveformLayout.AddWidget(waveformRisingLabel, 3, 0, 0)
	waveformLayout.AddWidget(waveformRisingTimeInput, 3, 1, 0)
	waveformLayout.AddWidget(waveformRisingPercentageInput, 3, 2, 0)
	waveformLayout.AddWidget(waveformWaitLabel, 4, 0, 0)
	waveformLayout.AddWidget(waveformWaitTimeInput, 4, 1, 0)
	waveformLayout.AddWidget(waveformWaitPercentageInput, 4, 2, 0)

	waveformLayout.AddWidget3(waveformVoltageLabel, 5, 0, 1, 1, 0)
	waveformLayout.AddWidget3(waveformVoltageInput, 5, 1, 1, 2, 0)
	waveformLayout.AddWidget3(waveformMnLabel, 6, 0, 1, 1, 0)
	waveformLayout.AddWidget3(waveformMnInput, 6, 1, 1, 2, 0)
	waveformLayout.AddWidget3(waveformGenerateButton, 7, 0, 1, 3, 0)

	d.waveformGroup.SetLayout(waveformLayout)
	d.waveformGroup.SetVisible(false)

	d.saveButton = widgets.NewQPushButton2("SAVE", nil)
	d.saveButton.ConnectClicked(func(bool) { d.saveInstruction() })
	d.saveButton.SetEnabled(false)

	d.GroupBox = widgets.NewQGroupBox2("Instruction", nil)
	layout := widgets.NewQGridLayout2()
	layout.AddWidget(titleLabel, 0, 0, 0)
	layout.AddWidget(d.titleInput, 0, 1, 0)
	layout.AddWidget(typeLabel, 1, 0, 0)
	layout.AddWidget(d.typeInput, 1, 1, 0)
	layout.AddWidget(instLabel, 2, 0, 0)
	layout.AddWidget(d.instInput, 2, 1, 0)
	layout.AddWidget(devLabel, 3, 0, 0)
	layout.AddWidget(d.devInput, 3, 1, 0)
	layout.AddWidget(argsLabel, 4, 0, 0)
	layout.AddWidget(d.argsInput, 4, 1, 0)
	layout.AddWidget3(d.lineInput, 5, 0, 1, 2, 0)
	layout.AddWidget3(d.waveformGroup, 6, 0, 1, 2, 0)

	layout.AddWidget3(d.saveButton, 7, 0, 1, 2, 0)
	d.GroupBox.SetLayout(layout)
	return &d
}

func (d *InstructionDetail) saveInstruction() {
	if d.treeItem == nil {
		return
	}

	d.SetLineInput()

	d.treeItem.SetText(0, d.titleInput.Text())
	tree.SetTreeItemData(d.treeItem, d.lineInput.Text())

}

func (d *InstructionDetail) Refresh(item *widgets.QTreeWidgetItem) {
	line := tree.GetTreeItemData(item)
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

func (d *InstructionDetail) onInstructionChanged(selected string) {
	switch selected {
	case "WAVEFORM":
		d.waveformGroup.SetVisible(true)
	default:
		d.waveformGroup.SetVisible(false)
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
			return
		}
	}
	if instruction == INST_SET_SYNC {
		d.instInput.SetCurrentText(INST_SET_SYNC)
		return
	}
	if instruction == INST_SET_ASYN {
		d.instInput.SetCurrentText(INST_SET_ASYN)
		return
	}
}

func (d *InstructionDetail) GetArgumentsFromLine() string {
	instruction := d.GetInstructionFromLine()
	return strings.Trim(d.Line(), fmt.Sprintf("%s ", instruction))
}

func (d *InstructionDetail) SetDevInput(items []string) {
	d.devInput.Clear()
	d.devInput.AddItems(items)
}

func (d *InstructionDetail) SetArgsInput() {
	d.argsInput.SetText(d.GetArgumentsFromLine())
}

func (d *InstructionDetail) SetLineInput() {
	d.lineInput.SetText(fmt.Sprintf("%s %s", d.instInput.CurrentText(), d.argsInput.Text()))
}
