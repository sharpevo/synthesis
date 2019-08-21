package devtree

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"synthesis/dao"
	"synthesis/dao/alientek"
	"synthesis/dao/aoztech"
	"synthesis/dao/canalystii"
	"synthesis/dao/ricoh_g5"
	"synthesis/gui/tree"
)

const (
	DEV_TYPE_UNK = dao.NAME
	DEV_TYPE_ALT = alientek.NAME
	DEV_TYPE_RCG = ricoh_g5.NAME
	DEV_TYPE_AOZ = aoztech.NAME
	DEV_TYPE_CAN = canalystii.NAME

	PRT_CONN = "CONN"
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

	descriptionInput *widgets.QLineEdit
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
		DEV_TYPE_ALT,
		DEV_TYPE_RCG,
		DEV_TYPE_AOZ,
		DEV_TYPE_CAN,
	})
	d.typeInput.ConnectCurrentTextChanged(d.onDeviceTypeChanged)
	descriptionLabel := widgets.NewQLabel2("Description:", nil, 0)
	d.descriptionInput = widgets.NewQLineEdit(nil)
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
	layout.AddWidget(descriptionLabel, 3, 0, 0)
	layout.AddWidget(d.descriptionInput, 3, 1, 0)
	layout.AddWidget(d.enabledInput, 4, 0, 0)
	layout.AddWidget3(d.saveButton, 5, 0, 1, 2, 0)
	layout.AddWidget3(logGroup, 6, 0, 1, 2, 0)
	d.GroupBox.SetLayout(layout)
	return &d
}

func (d *DeviceDetail) saveDeviceDetail() {
	if d.treeItem == nil {
		return
	}
	d.treeItem.SetText(0, d.titleInput.Text())
	variantMap := MakeVariantMap(
		d.lineInput.Text(),
		d.typeInput.CurrentText(),
		d.descriptionInput.Text(),
		d.enabledInput.CheckState() == core.Qt__Checked,
	)
	d.treeItem.SetData(
		0,
		tree.DataRole(),
		core.NewQVariant25(variantMap),
	)
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
	d.descriptionInput.SetText(variantMap.Description())
	d.enabledInput.SetCheckState(core.Qt__Unchecked)
	if variantMap.Enabled() {
		d.enabledInput.SetCheckState(core.Qt__Checked)
	}
}

func (d *DeviceDetail) onDeviceTypeChanged(selected string) {

	if selected == DEV_TYPE_UNK {
		return
	}

	var connItem *widgets.QTreeWidgetItem
	for i := 0; i < d.treeItem.ChildCount(); i++ {
		if item := d.treeItem.Child(i); PRT_CONN == item.Text(0) {
			connItem = item
			break
		}
	}
	if connItem.Pointer() == nil {
		connItem = NewDeviceConnItem(PRT_CONN)
		d.treeItem.InsertChild(0, connItem)
	}

	switch selected {
	case DEV_TYPE_RCG:
		for _, title := range ricoh_g5.CONN_ATTRIBUTES {
			seen := false
			for i := 0; i < connItem.ChildCount(); i++ {
				if item := connItem.Child(i); title == item.Text(0) {
					seen = true
					break
				}
			}
			if !seen {
				item := NewDeviceConnItem(title)
				connItem.InsertChild(0, item)
			}
		}
	case DEV_TYPE_ALT:
		for _, title := range alientek.CONN_ATTRIBUTES {
			seen := false
			for i := 0; i < connItem.ChildCount(); i++ {
				if item := connItem.Child(i); title == item.Text(0) {
					seen = true
					break
				}
			}
			if !seen {
				item := NewDeviceConnItem(title)
				connItem.InsertChild(0, item)
			}
		}
	case DEV_TYPE_AOZ:
		for _, title := range aoztech.CONN_ATTRIBUTES {
			seen := false
			for i := 0; i < connItem.ChildCount(); i++ {
				if item := connItem.Child(i); title == item.Text(0) {
					seen = true
					break
				}
			}
			if !seen {
				item := NewDeviceConnItem(title)
				connItem.InsertChild(0, item)
			}
		}
	case DEV_TYPE_CAN:
		for _, title := range canalystii.CONN_ATTRIBUTES {
			seen := false
			for i := 0; i < connItem.ChildCount(); i++ {
				if item := connItem.Child(i); title == item.Text(0) {
					seen = true
					break
				}
			}
			if !seen {
				item := NewDeviceConnItem(title)
				connItem.InsertChild(0, item)
			}
		}
	}
}

type VariantMap map[string]*core.QVariant

func MakeVariantMap(
	lineText string,
	typeText string,
	descriptionText string,
	enabledState bool,
) VariantMap {
	variantMap := make(VariantMap)
	variantMap["data"] = core.NewQVariant17(lineText)
	variantMap["type"] = core.NewQVariant17(typeText)
	variantMap["description"] = core.NewQVariant17(descriptionText)
	variantMap["enabled"] = core.NewQVariant11(enabledState)
	return variantMap
}

func (v VariantMap) Data() string {
	return v["data"].ToString()
}

func (v VariantMap) Type() string {
	return v["type"].ToString()
}

func (v VariantMap) Description() string {
	return v["description"].ToString()
}

func (v VariantMap) Enabled() bool {
	return v["enabled"].ToBool()
}
