package tree

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"log"
	"synthesis/gui/uiutil"
)

type Tree struct {
	widgets.QTreeWidget
	ContextMenu *widgets.QMenu
}

func (t *Tree) CustomDragEnterEvent(e *gui.QDragEnterEvent) {
	t.CurrentItem().SetExpanded(false)
	e.AcceptProposedAction()
}

func (t *Tree) CustomDropEvent(e *gui.QDropEvent) {
	index := t.IndexAt(e.Pos())
	item := t.CurrentItem()
	target := t.ItemFromIndex(index)
	if target.IsSelected() {
		return
	}
	if !index.IsValid() ||
		item.Pointer() == nil {
		return
	}

	parent := item.Parent()
	if parent.Pointer() == nil {
		fmt.Println("nil pointer")
		parent = t.InvisibleRootItem()
	}
	parent.RemoveChild(item)

	tparent := target.Parent()
	if tparent.Pointer() == nil {
		fmt.Println("nil pointer")
		tparent = t.InvisibleRootItem()
	}

	indic := t.DropIndicatorPosition()

	switch indic {
	case widgets.QAbstractItemView__OnItem:
		target.AddChild(item)
		target.SetExpanded(true)
		break
	case widgets.QAbstractItemView__AboveItem:
		tparent.InsertChild(tparent.IndexOfChild(target), item)
		break
	case widgets.QAbstractItemView__BelowItem:
		tparent.InsertChild(tparent.IndexOfChild(target)+1, item)
		break
	case widgets.QAbstractItemView__OnViewport:
		fmt.Println("viewport")
		break
	}
	t.SetCurrentItem(item)
	//item.SetExpanded(true)
}

func (t *Tree) AddItem(p *core.QPoint, item *widgets.QTreeWidgetItem) {
	root := t.ItemAt(p)
	if root.Pointer() == nil {
		root = t.InvisibleRootItem()
	}
	root.AddChild(item)
	root.SetExpanded(true)
}

func (t *Tree) RemoveItem(p *core.QPoint) {
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

func (t *Tree) Import(filePath string) error {
	node := new(Node)
	err := ImportNode(node, filePath)
	if err != nil {
		if err.Error() == "nothing selected" {
			return nil
		}
		log.Println(err)
		return err
	}
	t.Clear()
	for i := 0; i < len(node.Children); i++ {
		t.InvisibleRootItem().AddChild(t.ImportNode(node.Children[i]))
	}
	t.ExpandAll()
	return nil
}

func (t *Tree) ImportNode(node Node) *widgets.QTreeWidgetItem {
	item := widgets.NewQTreeWidgetItem2([]string{node.Title}, 0)
	item.SetData(
		0,
		DataRole(),
		core.NewQVariant15(node.Data),
	)
	for i := 0; i < len(node.Children); i++ {
		item.AddChild(t.ImportNode(node.Children[i]))
	}
	return item
}

func (t *Tree) Export(filePath string) error {
	item := t.InvisibleRootItem()
	node := t.ExportNode(item)
	err := ExportNode(node, filePath)
	if err != nil {
		return err
	}
	uiutil.MessageBoxInfo("Exported")
	return nil
}

func (t *Tree) ExportNode(root *widgets.QTreeWidgetItem) Node {
	node := Node{
		Title: root.Text(0),
		Data:  GetTreeItemData(root),
	}
	for i := 0; i < root.ChildCount(); i++ {
		node.Children = append(node.Children, t.ExportNode(root.Child(i)))
	}
	return node
}
