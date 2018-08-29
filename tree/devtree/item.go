package devtree

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"posam/gui/tree"
)

const (
	DEV_TYPE_UNK = "UNKNOWN"
	DEV_TYPE_TCP = "TCP"
	DEV_TYPE_CAN = "CAN"
	DEV_TYPE_SRL = "SERIAL"
)

type DeviceDetail struct {
	treeItem     *widgets.QTreeWidgetItem
	GroupBox     *widgets.QGroupBox
	titleInput   *widgets.QLineEdit
	lineInput    *widgets.QLineEdit
	typeInput    *widgets.QComboBox
	enabledInput *widgets.QCheckBox
	logInput     *widgets.QTextEdit
	saveButton   *widgets.QPushButton
}

func NewDeviceDetail() *DeviceDetail {

	titleLabel := widgets.NewQLabel2("Name:", nil, 0)
	lineLabel := widgets.NewQLabel2("Value:", nil, 0)
	typeLabel := widgets.NewQLabel2("Type:", nil, 0)
	//enabledLabel := widgets.NewQLabel2("Enabled:", nil, 0)

	d := DeviceDetail{}
	d.titleInput = widgets.NewQLineEdit(nil)
	d.lineInput = widgets.NewQLineEdit(nil)

	d.typeInput = widgets.NewQComboBox(nil)
	d.typeInput.AddItems([]string{
		DEV_TYPE_UNK,
		DEV_TYPE_SRL,
		DEV_TYPE_TCP,
		DEV_TYPE_CAN,
	})
	d.enabledInput = widgets.NewQCheckBox2("Enabled", nil)

	d.logInput = widgets.NewQTextEdit(nil)
	d.logInput.SetReadOnly(true)
	d.logInput.SetStyleSheet("QTextEdit { background-color: #e6e6e6}")
	logGroup := widgets.NewQGroupBox2("Logs", nil)
	logGroupLayout := widgets.NewQGridLayout2()
	logGroupLayout.AddWidget(d.logInput, 0, 0, 0)
	logGroup.SetLayout(logGroupLayout)

	d.saveButton = widgets.NewQPushButton2("SAVE", nil)
	d.saveButton.ConnectClicked(func(bool) { d.saveDeviceDetail() })

	d.GroupBox = widgets.NewQGroupBox2("Device", nil)
	layout := widgets.NewQGridLayout2()
	layout.AddWidget(titleLabel, 0, 0, 0)
	layout.AddWidget(d.titleInput, 0, 1, 0)
	layout.AddWidget(lineLabel, 1, 0, 0)
	layout.AddWidget(d.lineInput, 1, 1, 0)
	layout.AddWidget(typeLabel, 2, 0, 0)
	layout.AddWidget(d.typeInput, 2, 1, 0)
	layout.AddWidget(d.enabledInput, 3, 0, 0)
	layout.AddWidget3(d.saveButton, 4, 0, 1, 2, 0)
	layout.AddWidget3(logGroup, 5, 0, 1, 2, 0)
	d.GroupBox.SetLayout(layout)
	return &d
}

func (d *DeviceDetail) saveDeviceDetail() {
	if d.treeItem == nil {
		return
	}
	d.treeItem.SetText(0, d.titleInput.Text())
	tree.SetTreeItemData(d.treeItem, d.lineInput.Text())
	variantMap := MakeVariantMap(
		d.lineInput.Text(),
		d.typeInput.CurrentText(),
		d.enabledInput.CheckState() == core.Qt__Checked,
	)
	d.treeItem.SetData(
		0,
		tree.DataRole(),
		core.NewQVariant25(variantMap),
	)
	fmt.Println(d.treeItem)
}

func (d *DeviceDetail) Refresh(item *widgets.QTreeWidgetItem) {
	d.treeItem = item
	d.titleInput.SetText(item.Text(0))
	variantMap := VariantMap(item.Data(0, tree.DataRole()).ToMap())
	d.lineInput.SetText(variantMap.Data())
	if variantMap.Type() != "" {
		d.typeInput.SetCurrentText(variantMap.Type())
	} else {
		d.typeInput.SetCurrentText(DEV_TYPE_UNK)
	}
	d.enabledInput.SetCheckState(core.Qt__Unchecked)
	if variantMap.Enabled() {
		d.enabledInput.SetCheckState(core.Qt__Checked)
	}
}

type VariantMap map[string]*core.QVariant

func MakeVariantMap(
	lineText string,
	typeText string,
	enabledState bool,
) VariantMap {
	variantMap := make(VariantMap)
	variantMap["data"] = core.NewQVariant17(lineText)
	variantMap["type"] = core.NewQVariant17(typeText)
	variantMap["enabled"] = core.NewQVariant11(enabledState)
	return variantMap
}

func (v VariantMap) Data() string {
	return v["data"].ToString()
}

func (v VariantMap) Type() string {
	return v["type"].ToString()
}

func (v VariantMap) Enabled() bool {
	return v["enabled"].ToBool()
}
