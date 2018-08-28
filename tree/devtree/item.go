package devtree

import (
	"github.com/therecipe/qt/widgets"
	"posam/gui/tree"
)

type DeviceDetail struct {
	treeItem   *widgets.QTreeWidgetItem
	GroupBox   *widgets.QGroupBox
	titleInput *widgets.QLineEdit
	lineInput  *widgets.QLineEdit
	logInput   *widgets.QTextEdit
	saveButton *widgets.QPushButton
}

func NewDeviceDetail() *DeviceDetail {

	titleLabel := widgets.NewQLabel2("Name:", nil, 0)
	lineLabel := widgets.NewQLabel2("Value:", nil, 0)

	d := DeviceDetail{}
	d.titleInput = widgets.NewQLineEdit(nil)
	d.lineInput = widgets.NewQLineEdit(nil)

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
	layout.AddWidget3(d.saveButton, 2, 0, 1, 2, 0)
	layout.AddWidget3(logGroup, 3, 0, 1, 2, 0)
	d.GroupBox.SetLayout(layout)
	return &d
}

func (d *DeviceDetail) saveDeviceDetail() {
	if d.treeItem == nil {
		return
	}
	d.treeItem.SetText(0, d.titleInput.Text())
	tree.SetTreeItemData(d.treeItem, d.lineInput.Text())
}

func (d *DeviceDetail) Refresh(item *widgets.QTreeWidgetItem) {
	line := tree.GetTreeItemData(item)
	d.treeItem = item
	d.titleInput.SetText(item.Text(0))
	d.lineInput.SetText(line)
}
