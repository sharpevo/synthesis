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

func (w *InstructionTree) customContextMenuRequested(p *core.QPoint) {
	if w.contextMenu == nil {
		w.contextMenu = widgets.NewQMenu(w)
		menuAdd := w.contextMenu.AddAction("Add")
		menuAdd.ConnectTriggered(func(checked bool) { w.addItem(p) })
		menuRemove := w.contextMenu.AddAction("Remove")
		menuRemove.ConnectTriggered(func(checked bool) { w.removeItem(p) })
	}
	w.contextMenu.Exec2(w.MapToGlobal(p), nil)
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
