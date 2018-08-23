package main

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"log"
)

type InstructionTree struct {
	widgets.QTreeWidget
	contextMenu *widgets.QMenu
}

func NewInstructionTree() *InstructionTree {

	treeWidget := &InstructionTree{*widgets.NewQTreeWidget(nil), nil}
	treeWidget.SetWindowTitle("Graphical Programming")
	treeWidget.SetContextMenuPolicy(core.Qt__CustomContextMenu)
	treeWidget.ConnectCustomContextMenuRequested(treeWidget.customContextMenuRequested)
	rootNode := treeWidget.InvisibleRootItem()
	anItem := widgets.NewQTreeWidgetItem2([]string{"first instruction"}, 0)
	rootNode.AddChild(anItem)
	treeWidget.ExpandAll()

	return treeWidget
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
	item := widgets.NewQTreeWidgetItem2([]string{"instruction"}, 0)
	//item.SetFlags(^core.Qt__ItemIsEditable)
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
