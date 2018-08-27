package tree

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"log"
)

type InstructionTree struct {
	widgets.QTreeWidget
	contextMenu *widgets.QMenu
	detail      *InstructionDetail
	runButton   *widgets.QPushButton
	inputEdit   *widgets.QTextEdit
}

func NewTree(
	detail *InstructionDetail,
	runButton *widgets.QPushButton,
	inputEdit *widgets.QTextEdit,
) *InstructionTree {

	treeWidget := &InstructionTree{
		*widgets.NewQTreeWidget(nil),
		nil,
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
		menuAdd := t.contextMenu.AddAction("Add child")
		menuAdd.ConnectTriggered(func(checked bool) { t.addItem(p) })
		menuRemove := t.contextMenu.AddAction("Remove node")
		menuRemove.ConnectTriggered(func(checked bool) { t.removeItem(p) })
		menuRun := t.contextMenu.AddAction("Execute single step")
		menuRun.ConnectTriggered(func(checked bool) { t.executeItem(p) })
	}
	t.contextMenu.Exec2(t.MapToGlobal(p), nil)
}

func (t *InstructionTree) addItem(p *core.QPoint) {
	root := t.ItemAt(p)
	if root.Pointer() == nil {
		root = t.InvisibleRootItem()
	}
	item := NewInstructionItem("print instruction", "PRINT")
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
		MessageBox(err.Error())
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

func (t *InstructionTree) Import() {
	node := new(Node)
	err := node.Read()
	if err != nil {
		if err.Error() == "nothing selected" {
			return
		}
		log.Println(err)
	}
	t.Clear()
	for i := 0; i < len(node.Children); i++ {
		t.InvisibleRootItem().AddChild(t.ImportNode(node.Children[i]))
	}
	t.ExpandAll()

	MessageBox("Imported")
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
	MessageBox("Exported")
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

func MessageBox(message string) {
	widgets.QMessageBox_Information(
		nil,
		"OK",
		message,
		widgets.QMessageBox__Ok,
		widgets.QMessageBox__Close,
	)
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
		MessageBox(err.Error())
		return err
	}
	t.WriteInputEdit(filePath)
	t.runButton.Click()
	return nil
}
