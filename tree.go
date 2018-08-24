package main

import (
	"encoding/gob"
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
	anItem := NewInstructionItem("Hello World!", "PRINT Hello World!")
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
	item := NewInstructionItem("print instruction", "PRINT instruction")
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
	t.ExpandAll()

	widgets.QMessageBox_Information(
		nil,
		"OK",
		"Imported",
		widgets.QMessageBox__Ok,
		widgets.QMessageBox__Ok,
	)
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
	widgets.QMessageBox_Information(
		nil,
		"OK",
		"Exported",
		widgets.QMessageBox__Ok,
		widgets.QMessageBox__Ok,
	)
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
	filePath, err := FilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(filePath)
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
	filePath, err := FilePath()
	if err != nil {
		return err
	}
	file, err := os.Open(filePath)
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

func FilePath() (string, error) {
	dialog := widgets.NewQFileDialog2(nil, "Select file...", "", "")
	if dialog.Exec() != int(widgets.QDialog__Accepted) {
		return "", fmt.Errorf("nothing selected")
	}
	filePath := dialog.SelectedFiles()[0]
	return filePath, nil
}

func MessageBox(message string) {
	widgets.QMessageBox_Information(
		nil,
		"OK",
		message,
		widgets.QMessageBox__Ok,
		widgets.QMessageBox__Ok,
	)
}

func (t *InstructionTree) Generate() (string, error) {
	root := t.InvisibleRootItem()
	node := t.ExportNode(root)
	filePath, err := node.Generate()
	if err != nil {
		MessageBox(err.Error())
		return filePath, err
	}
	return filePath, nil
}

func (n *Node) Generate() (string, error) {
	dir, err := ioutil.TempDir("", "igenetech")
	if err != nil {
		return "", err
	}
	file, err := ioutil.TempFile(dir, "posam")
	if err != nil {
		return "", err
	}
	filePath := file.Name()

	defer file.Close()
	if err != nil {
		return filePath, err
	}
	offset := 0
	for _, child := range n.Children {
		offset += 1
		nodeType, err := child.Type()
		if err != nil {
			return filePath, err
		}
		switch nodeType {
		case TYPE_INS:
			file.WriteString(fmt.Sprintf("%s\n", child.Data))
			break
		case TYPE_SET_ONCE:
			setPath, err := child.Generate()
			if err != nil {
				return filePath, err
			}
			file.WriteString(fmt.Sprintf("%s %s\n", child.Instruction(), setPath))
			break
		case TYPE_SET_LOOP:
			setPath, err := child.Generate()
			if err != nil {
				return filePath, err
			}
			file.WriteString(fmt.Sprintf("%s %s\n", child.Instruction(), setPath))
			file.WriteString(fmt.Sprintf("LOOP %d %s\n", offset, child.Arguments()[0]))
			file.WriteString(fmt.Sprintf("PRINT loop done\n"))
			offset += 3
			break
		case TYPE_SET_COND:
			setPath, err := child.Generate()
			if err != nil {
				return filePath, err
			}
			var1 := child.Arguments()[0]
			opsb := child.Arguments()[1]
			var2 := child.Arguments()[2]
			var opst string
			switch opsb {
			case ">":
				opst = "GTGOTO"
				break
			case "<":
				opst = "LTGOTO"
				break
			case "!=":
				opst = "NEGOTO"
				break
			case "==":
				opst = "EQGOTO"
				break
			default:
				return filePath, fmt.Errorf(
					"invalid operator in %q",
					n.Title,
				)
			}

			file.WriteString(fmt.Sprintf("CMPVAR %s %s\n", var1, var2))
			file.WriteString(fmt.Sprintf("%s %d\n", opst, offset+3))
			file.WriteString(fmt.Sprintf("GOTO %d\n", offset+4))
			file.WriteString(fmt.Sprintf("%s %s\n", child.Instruction(), setPath))
			file.WriteString(fmt.Sprintf("PRINT condition done\n"))
			offset += 4
			break
		}
	}
	file.Sync()
	return filePath, nil
}

func shouldBeInstructionSet(instruction string) bool {
	return instruction == INST_SET_SYNC ||
		instruction == INST_SET_ASYN
}

func (n *Node) Instruction() string {
	return strings.Split(n.Data, " ")[0]
}

func (n *Node) Arguments() []string {
	return strings.Split(n.Data, " ")[1:]
}

func (n *Node) Type() (string, error) {
	dataList := strings.Split(n.Data, " ")
	instruction := dataList[0]
	argumentList := dataList[1:]

	if shouldBeInstructionSet(instruction) {
		if len(n.Children) == 0 {
			return "", fmt.Errorf(
				"instruction set %q has no instructions",
				n.Title,
			)
		}
		switch len(argumentList) {
		case 0:
			return TYPE_SET_ONCE, nil
		case 1:
			return TYPE_SET_LOOP, nil
		case 3:
			switch argumentList[1] {
			case ">":
				argumentList[1] = "GTGOTO"
				break
			case "<":
				argumentList[1] = "LTGOTO"
				break
			case "!=":
				argumentList[1] = "NEGOTO"
				break
			case "==":
				argumentList[1] = "EQGOTO"
				break
			default:
				return "", fmt.Errorf(
					"invalid operator in %q",
					n.Title,
				)
			}
			return TYPE_SET_COND, nil
		default:
			return "", fmt.Errorf(
				"instruction %q is not valid instruction set",
				n.Title,
			)
		}
	} else {
		if len(n.Children) > 0 {
			return "", fmt.Errorf(
				"instruction %q should be instruction set",
				n.Title,
			)
		}
		return TYPE_INS, nil
	}
}
