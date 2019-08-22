package instree

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"log"
	"os"
	"synthesis/pkg/config"
	"synthesis/gui/tree"
	"synthesis/gui/uiutil"
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
	treeWidget.ConnectItemDoubleClicked(treeWidget.customItemDoubleClicked)
	treeWidget.ConnectCurrentItemChanged(treeWidget.customCurrentItemChanged)
	//treeWidget.ConnectKeyPressEvent(treeWidget.customKeyPressEvent)

	treeWidget.ImportPreviousFile()

	treeWidget.SetAcceptDrops(true)
	treeWidget.SetDragEnabled(true)
	treeWidget.SetExpandsOnDoubleClick(false)
	treeWidget.ConnectDragEnterEvent(treeWidget.customDragEnterEvent)
	treeWidget.ConnectDropEvent(treeWidget.customDropEvent)
	treeWidget.ExpandAll()

	return treeWidget
}

func (t *InstructionTree) customKeyPressEvent(e *gui.QKeyEvent) {
	switch int32(e.Key()) {
	case int32(core.Qt__Key_Return), int32(core.Qt__Key_Enter):
		//e.Accept()
		t.ItemDoubleClicked(t.CurrentItem(), 0)
	default:
		//fmt.Println("key pressed", e.Key())
		e.Accept()
	}
}

func (t *InstructionTree) RefreshDetail() {
	t.detail.Refresh(t.CurrentItem())
	//fmt.Println("refresh instruction detail", t.detail.devInput.CurrentData(int(core.Qt__UserRole)).ToString())
	//t.detail.onDeviceChanged("")
}

func (t *InstructionTree) customDragEnterEvent(e *gui.QDragEnterEvent) {
	t.CurrentItem().SetExpanded(false)
	e.AcceptProposedAction()
}

func (t *InstructionTree) customDropEvent(e *gui.QDropEvent) {

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

func (t *InstructionTree) customCurrentItemChanged(
	curr *widgets.QTreeWidgetItem,
	prev *widgets.QTreeWidgetItem,
) {
	t.detail.Refresh(t.CurrentItem())
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
		core.NewQVariant15(line),
	)
	return treeItem
}

func (t *InstructionTree) customContextMenuRequested(p *core.QPoint) {
	if t.ContextMenu == nil {
		t.ContextMenu = widgets.NewQMenu(t)
		menuAdd := t.ContextMenu.AddAction("Add child")
		menuAdd.ConnectTriggered(func(checked bool) {
			newItem := NewInstructionItem("print instruction", "PRINT")
			t.AddItem(p, newItem)
			newItem.Parent().SetSelected(false)
			newItem.SetSelected(true)
			t.SetCurrentItem(newItem)
			//t.detail.Refresh(newItem)
		})
		menuRemove := t.ContextMenu.AddAction("Remove node")
		menuRemove.ConnectTriggered(func(checked bool) { t.RemoveItem(p) })
		menuImport := t.ContextMenu.AddAction("Import as sub node")
		menuImport.ConnectTriggered(func(checked bool) { t.importSubItem(p) })
		menuExport := t.ContextMenu.AddAction("Export selected node")
		menuExport.ConnectTriggered(func(checked bool) { t.exportCurItem(p) })
		menuRun := t.ContextMenu.AddAction("Execute single step")
		menuRun.ConnectTriggered(func(checked bool) { t.executeItem(p) })
		menuCopy := t.ContextMenu.AddAction("Copy")
		menuCopy.ConnectTriggered(func(checked bool) { t.copyItem(p) })
		menuPaste := t.ContextMenu.AddAction("Paste")
		menuPaste.ConnectTriggered(func(checked bool) { t.pasteItem(p) })
	}
	t.ContextMenu.Exec2(t.MapToGlobal(p), nil)
}

func (t *InstructionTree) importSubItem(p *core.QPoint) {
	item := t.ItemAt(p)
	if item.Pointer() == nil {
		log.Println("invalid tree item")
		return
	}
	filePath, err := uiutil.FilePath()
	if err != nil {
		uiutil.MessageBoxError(err.Error())
		return
	}
	node := new(Node)
	if err := tree.ImportNode(node, filePath); err != nil {
		uiutil.MessageBoxError(err.Error())
		return
	}
	for i := 0; i < len(node.Children); i++ {
		item.AddChild(t.ImportNode(node.Children[i]))
	}
	item.SetExpanded(true)
	uiutil.MessageBoxInfo("imported")
}

func (t *InstructionTree) exportCurItem(p *core.QPoint) {
	item := t.ItemAt(p)
	if item.Pointer() == nil {
		log.Println("invalid tree item")
		return
	}
	node := t.ExportNode(item)
	step := new(Node)
	step.Title = "step"
	step.Children = append(step.Children, node)
	filePath, err := uiutil.FilePath()
	if err != nil {
		uiutil.MessageBoxError(err.Error())
		return
	}
	if err := tree.ExportNode(step, filePath); err != nil {
		uiutil.MessageBoxError(err.Error())
		return
	}
	uiutil.MessageBoxInfo(fmt.Sprintf("exported to %q", filePath))
}

func (t *InstructionTree) pasteItem(p *core.QPoint) {
	item := t.ItemAt(p)
	if item.Pointer() == nil {
		log.Println("invalid tree item")
		return
	}
	filePath := "tmp.bin"
	node := new(Node)
	if err := tree.ImportNode(node, filePath); err != nil {
		uiutil.MessageBoxError(err.Error())
		return
	}
	for i := 0; i < len(node.Children); i++ {
		item.AddChild(t.ImportNode(node.Children[i]))
	}
	item.SetExpanded(true)
}

func (t *InstructionTree) copyItem(p *core.QPoint) {
	item := t.ItemAt(p)
	if item.Pointer() == nil {
		log.Println("invalid tree item")
		return
	}
	node := t.ExportNode(item)
	step := new(Node)
	step.Title = "step"
	step.Children = append(step.Children, node)
	filePath := "tmp.bin"
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Println("error", err)
		uiutil.MessageBoxError(err.Error())
		return
	}
	f.Close()
	if err := tree.ExportNode(step, filePath); err != nil {
		uiutil.MessageBoxError(err.Error())
		return
	}
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

func (t *InstructionTree) customItemDoubleClicked(item *widgets.QTreeWidgetItem, column int) {
	node := t.ExportNode(item)
	pseudop := Node{}
	pseudop.Children = []Node{node}
	filePath, err := pseudop.Generate()
	if err != nil {
		uiutil.MessageBoxError(err.Error())
	}
	t.WriteInputEdit(filePath)
	t.SetCurrentItem(nil)
	t.runButton.Click()
	t.SetCurrentItem(item)
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

func (t *InstructionTree) Import(filePath string) error {
	node := new(Node)
	err := tree.ImportNode(node, filePath)
	if err != nil {
		if err.Error() == "nothing selected" {
			return nil
		}
		return err
	}
	t.SaveImportedFilePath(filePath)
	t.Clear()
	for i := 0; i < len(node.Children); i++ {
		t.InvisibleRootItem().AddChild(t.ImportNode(node.Children[i]))
	}
	t.ExpandAll()
	return nil
}

func (t *InstructionTree) ImportNode(node Node) *widgets.QTreeWidgetItem {
	item := widgets.NewQTreeWidgetItem2([]string{node.Title}, 0)
	variantMap := MakeVariantMap(
		node.DevicePath,
		node.DeviceType,
		node.Instruction,
		node.Arguments,
	)
	item.SetData(
		0,
		tree.DataRole(),
		core.NewQVariant23(variantMap),
	)
	for i := 0; i < len(node.Children); i++ {
		item.AddChild(t.ImportNode(node.Children[i]))
	}
	return item
}

func (t *InstructionTree) Export(filePath string) error {
	item := t.InvisibleRootItem()
	node := t.ExportNode(item)
	err := tree.ExportNode(node, filePath)
	if err != nil {
		return err
	}
	uiutil.MessageBoxInfo("Exported")
	return nil
}

func (t *InstructionTree) ExportNode(root *widgets.QTreeWidgetItem) Node {
	node := Node{}
	node.Title = root.Text(0)
	variantMap := VariantMap(root.Data(0, tree.DataRole()).ToMap())
	node.DevicePath = variantMap.Device()
	node.DeviceType = variantMap.DeviceType()
	node.Instruction = variantMap.Instruction()
	node.Arguments = variantMap.Arguments()
	for i := 0; i < root.ChildCount(); i++ {
		node.Children = append(node.Children, t.ExportNode(root.Child(i)))
	}
	return node
}

const IMPORTED_FILE_KEY = "general.importedfile"

func (t *InstructionTree) SaveImportedFilePath(filePath string) {
	config.Set(IMPORTED_FILE_KEY, filePath)
}

func (t *InstructionTree) ImportPreviousFile() {
	filePath := config.GetString(IMPORTED_FILE_KEY)
	if err := t.Import(filePath); err != nil {
		uiutil.MessageBoxError(err.Error())
	} else {
		log.Println(
			fmt.Sprintf("latest file imported: %q", filePath))
	}
}
