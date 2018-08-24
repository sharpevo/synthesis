package main

import (
	"encoding/gob"
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"log"
	"os"
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

	treeWidget.SetAcceptDrops(true)
	treeWidget.SetDragEnabled(true)
	//treeWidget.SetDragDropMode()
	treeWidget.ConnectDropEvent(treeWidget.customDropEvent)
	treeWidget.ExpandAll()

	return treeWidget
}

func (t *InstructionTree) customDropEvent(e *gui.QDropEvent) {

	//e.SetDropAction(core.Qt__CopyAction)
	//e.AcceptProposedAction()
	//e.SetAccepted(true)
	index := t.IndexAt(e.Pos())
	item := t.CurrentItem()
	target := t.ItemFromIndex(index)
	if !index.IsValid() ||
		item.Pointer() == nil {
		return
	}
	indic := t.DropIndicatorPosition()
	fmt.Println("---")
	fmt.Printf("mime: %#v\n", e.MimeData())
	fmt.Printf("selected: %#v\n", t.CurrentItem().Text(0))
	fmt.Printf("parent: %#v\n", index.Parent())
	fmt.Printf("row: %#v\n", index.Row())
	fmt.Printf("indic: %#v\n", indic)
	fmt.Printf("text: %#v\n", t.ItemFromIndex(index).Text(0))
	isAbove := false
	switch indic {
	case widgets.QAbstractItemView__OnItem:
		fmt.Println("on")
		isAbove = true
		break
	case widgets.QAbstractItemView__AboveItem:
		fmt.Println("above")
		break
	case widgets.QAbstractItemView__BelowItem:
		fmt.Println("below")
		break
	case widgets.QAbstractItemView__OnViewport:
		fmt.Println("viewport")
		break
	}
	parent := item.Parent()
	if parent.Pointer() == nil {
		fmt.Println("nil pointer")
		parent = t.InvisibleRootItem()
	}
	parent.RemoveChild(item)

	if isAbove {
		target.AddChild(item)
		target.SetExpanded(true)
	} else {

		tparent := target.Parent()
		if tparent.Pointer() == nil {
			fmt.Println("nil pointer")
			tparent = t.InvisibleRootItem()
		}
		tparent.InsertChild(index.Row(), item)
	}
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
	treeItem.SetFlags(core.Qt__ItemIsEnabled |
		core.Qt__ItemIsSelectable |
		core.Qt__ItemIsDragEnabled |
		core.Qt__ItemIsDropEnabled)
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

func (t *InstructionTree) Import() {
	node := new(Node)
	node.Read()
	t.Clear()
	for i := 0; i < len(node.Children); i++ {
		t.InvisibleRootItem().AddChild(t.ImportNode(node.Children[i]))
	}
}

func (t *InstructionTree) ImportNode(node Node) *widgets.QTreeWidgetItem {
	item := NewInstructionItem(node.Title, node.Data)
	for i := 0; i < len(node.Children); i++ {
		item.AddChild(t.ImportNode(node.Children[i]))
	}
	return item
}

func (t *InstructionTree) ExportAll() error {
	item := t.InvisibleRootItem()
	node := t.ExportNode(item)
	err := node.Write()
	if err != nil {
		return err
	}
	return nil
}

func (t *InstructionTree) ExportNode(root *widgets.QTreeWidgetItem) Node {
	node := Node{
		Title: root.Text(0),
		Data:  GetTreeItemData(root),
	}
	for i := 0; i < root.ChildCount(); i++ {
		node.Children = append(node.Children, t.ExportNode(root.Child(i)))
	}
	return node
}

type Node struct {
	Title    string
	Data     string
	Children []Node
}

func (n *Node) Write() error {
	file, err := os.Create("/tmp/gob/test")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	encoder := gob.NewEncoder(file)
	encoder.Encode(n)
	return nil
}

func (n *Node) Read() error {
	file, err := os.Open("/tmp/gob/test")
	defer file.Close()
	if err != nil {
		return err
	}
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(n)
	if err != nil {
		return err
	}
	return nil
}
