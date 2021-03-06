package devtree

import (
	"github.com/therecipe/qt/widgets"
	"synthesis/gui/tree/instree"
)

func NewDevTree(instDetail *instree.InstructionDetail) *widgets.QGroupBox {

	detail := NewDeviceDetail()

	treeGroup := widgets.NewQGroupBox2("Referenced configurations of devices", nil)
	treeLayout := widgets.NewQGridLayout2()
	treeWidget := NewTree(detail, instDetail)
	treeLayout.AddWidget2(treeWidget, 0, 0, 0)
	treeLayout.AddWidget3(detail.GroupBox, 0, 1, 2, 1, 0)

	treeSaveButton := widgets.NewQPushButton2("SAVE", nil)
	treeSaveButton.ConnectClicked(func(bool) { treeWidget.Save() })
	treeSaveButton.SetVisible(false)
	treeLayout.AddWidget2(treeSaveButton, 1, 0, 0)

	treeGroup.SetLayout(treeLayout)
	return treeGroup
}
