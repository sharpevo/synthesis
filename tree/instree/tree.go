package instree

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"log"
	"posam/gui/tree"
	"posam/gui/uiutil"
)

type InstructionTree struct {
	tree.Tree
	detail    *InstructionDetail
	runButton *widgets.QPushButton
	inputEdit *widgets.QTextEdit
}

func NewTree(
	detail *InstructionDetail,
	runButton *widgets.QPushButton,
	inputEdit *widgets.QTextEdit,
) *InstructionTree {

	treeWidget := &InstructionTree{
		tree.Tree{
			*widgets.NewQTreeWidget(nil),
			nil,
		},
		detail,
		runButton,
		inputEdit,
	}
	treeWidget.SetWindowTitle("Graphical Programming")
	treeWidget.SetContextMenuPolicy(core.Qt__CustomContextMenu)
	treeWidget.ConnectCustomContextMenuRequested(treeWidget.customContextMenuRequested)
	treeWidget.ConnectItemClicked(treeWidget.customItemClicked)
	rootNode := treeWidget.InvisibleRootItem()
	anItem := NewInstructionItem("Hello World!", "PRINT Hello World!")
	rootNode.AddChild(anItem)

	treeWidget.SetAcceptDrops(true)
	treeWidget.SetDragEnabled(true)
	treeWidget.ConnectDropEvent(treeWidget.customDropEvent)
	treeWidget.ExpandAll()

	return treeWidget
}

func (t *InstructionTree) customDropEvent(e *gui.QDropEvent) {

	index := t.IndexAt(e.Pos())
	item := t.CurrentItem()
	target := t.ItemFromIndex(index)
	if !index.IsValid() ||
		item.Pointer() == nil {
		return
	}
	indic := t.DropIndicatorPosition()
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

func NewInstructionItem(title string, line string) *widgets.QTreeWidgetItem {
	treeItem := widgets.NewQTreeWidgetItem2([]string{title}, 0)
	treeItem.SetFlags(core.Qt__ItemIsEnabled |
		core.Qt__ItemIsSelectable |
		core.Qt__ItemIsDragEnabled |
		core.Qt__ItemIsDropEnabled)
	treeItem.SetData(
		0,
		tree.DataRole(),
		core.NewQVariant17(line),
	)
	return treeItem
}

func (t *InstructionTree) customContextMenuRequested(p *core.QPoint) {
	if t.ContextMenu == nil {
		t.ContextMenu = widgets.NewQMenu(t)
		menuAdd := t.ContextMenu.AddAction("Add child")
		menuAdd.ConnectTriggered(func(checked bool) { t.AddItem(p, NewInstructionItem("print instruction", "PRINT")) })
		menuRemove := t.ContextMenu.AddAction("Remove node")
		menuRemove.ConnectTriggered(func(checked bool) { t.RemoveItem(p) })
		menuRun := t.ContextMenu.AddAction("Execute single step")
		menuRun.ConnectTriggered(func(checked bool) { t.executeItem(p) })
	}
	t.ContextMenu.Exec2(t.MapToGlobal(p), nil)
}

func (t *InstructionTree) executeItem(p *core.QPoint) {
	item := t.ItemAt(p)
	if item.Pointer() == nil {
		log.Println("invalid tree item")
		return
	}
	node := t.ExportNode(item)
	pseudop := Node{}
	pseudop.Children = []Node{node}
	filePath, err := pseudop.Generate()
	if err != nil {
		uiutil.MessageBoxError(err.Error())
	}
	t.WriteInputEdit(filePath)
	t.runButton.Click()
}

func (t *InstructionTree) WriteInputEdit(filePath string) {
	instBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	t.inputEdit.SetPlainText(string(instBytes))
}

func (t *InstructionTree) customItemClicked(item *widgets.QTreeWidgetItem, column int) {
	t.detail.Refresh(item)
}

func (t *InstructionTree) Generate() (string, error) {
	root := t.InvisibleRootItem()
	node := t.ExportNode(root)
	filePath, err := node.Generate()
	if err != nil {
		return filePath, err
	}
	return filePath, nil
}

func (t *InstructionTree) Execute() error {
	filePath, err := t.Generate()
	if err != nil {
		uiutil.MessageBoxError(err.Error())
		return err
	}
	t.WriteInputEdit(filePath)
	t.runButton.Click()
	return nil
}

func (t *InstructionTree) ExportNode(root *widgets.QTreeWidgetItem) Node {
	node := Node{}
	node.Title = root.Text(0)
	node.Data = tree.GetTreeItemData(root)
	for i := 0; i < root.ChildCount(); i++ {
		node.Children = append(node.Children, t.ExportNode(root.Child(i)))
	}
	return node
}
