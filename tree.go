package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"log"
)

type InstructionTree struct {
	widgets.QTreeWidget
	contextMenu *widgets.QMenu
	detail      *InstructionDetail
}

func NewInstructionTree(detail *InstructionDetail) *InstructionTree {

	treeWidget := &InstructionTree{*widgets.NewQTreeWidget(nil), nil, detail}
	treeWidget.SetWindowTitle("Graphical Programming")
	treeWidget.SetContextMenuPolicy(core.Qt__CustomContextMenu)
	treeWidget.ConnectCustomContextMenuRequested(treeWidget.customContextMenuRequested)
	treeWidget.ConnectItemClicked(treeWidget.customItemClicked)
	rootNode := treeWidget.InvisibleRootItem()
	anItem := NewInstructionItem("first instruction", "preset")
	rootNode.AddChild(anItem)
	treeWidget.ExpandAll()

	return treeWidget
}

func DataRole() int {
	return int(core.Qt__UserRole) + 1
}

func GetTreeItemData(item *widgets.QTreeWidgetItem) string {
	return item.Data(0, DataRole()).ToString()
}

func SetTreeItemData(item *widgets.QTreeWidgetItem, data string) {
	item.SetData(
		0,
		DataRole(),
		core.NewQVariant17(data),
	)
}

func NewInstructionItem(title string, line string) *widgets.QTreeWidgetItem {
	treeItem := widgets.NewQTreeWidgetItem2([]string{title}, 0)
	treeItem.SetData(
		0,
		DataRole(),
		core.NewQVariant17(line),
	)
	return treeItem
}

func (t *InstructionTree) customContextMenuRequested(p *core.QPoint) {
	if t.contextMenu == nil {
		t.contextMenu = widgets.NewQMenu(t)
		menuAdd := t.contextMenu.AddAction("Add")
		menuAdd.ConnectTriggered(func(checked bool) { t.addItem(p) })
		menuRemove := t.contextMenu.AddAction("Remove")
		menuRemove.ConnectTriggered(func(checked bool) { t.removeItem(p) })
	}
	t.contextMenu.Exec2(t.MapToGlobal(p), nil)
}

func (t *InstructionTree) addItem(p *core.QPoint) {
	root := t.ItemAt(p)
	if root.Pointer() == nil {
		root = t.InvisibleRootItem()
	}
	item := NewInstructionItem("instruction", "new")
	root.AddChild(item)
	root.SetExpanded(true)
}

func (t *InstructionTree) removeItem(p *core.QPoint) {
	item := t.ItemAt(p)
	if item.Pointer() == nil {
		log.Println("invalid tree item")
		return
	}
	parent := item.Parent()
	if parent.Pointer() == nil {
		parent = t.InvisibleRootItem()
	}
	parent.RemoveChild(item)
}

func (t *InstructionTree) customItemClicked(item *widgets.QTreeWidgetItem, column int) {
	t.detail.Refresh(item)
}
