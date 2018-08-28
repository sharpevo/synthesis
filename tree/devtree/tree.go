package devtree

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"posam/gui/tree"
	"posam/gui/uiutil"
)

const (
	DEVICE_CONF_FILEPATH = "devices.bin"
)

type DeviceTree struct {
	tree.Tree
	detail *DeviceDetail
}

func NewTree(
	detail *DeviceDetail,
) *DeviceTree {

	treeWidget := &DeviceTree{
		tree.Tree{
			*widgets.NewQTreeWidget(nil),
			nil,
		},
		detail,
	}
	treeWidget.SetWindowTitle("Devices")
	treeWidget.SetContextMenuPolicy(core.Qt__CustomContextMenu)
	treeWidget.ConnectCustomContextMenuRequested(treeWidget.customContextMenuRequested)
	treeWidget.ConnectItemClicked(treeWidget.customItemClicked)
	treeWidget.ConnectItemChanged(treeWidget.customItemChanged)

	err := treeWidget.Import()
	if err != nil {
		uiutil.MessageBoxError(err.Error())
	}

	treeWidget.SetAcceptDrops(true)
	treeWidget.SetDragEnabled(true)
	treeWidget.ConnectDropEvent(treeWidget.CustomDropEvent)
	treeWidget.ExpandAll()

	return treeWidget
}

func NewDeviceItem(title string, data string) *widgets.QTreeWidgetItem {
	treeItem := widgets.NewQTreeWidgetItem2([]string{title}, 0)
	treeItem.SetFlags(core.Qt__ItemIsEnabled |
		core.Qt__ItemIsSelectable |
		core.Qt__ItemIsDragEnabled |
		core.Qt__ItemIsDropEnabled)
	treeItem.SetData(
		0,
		tree.DataRole(),
		core.NewQVariant17(data),
	)
	return treeItem
}

func (t *DeviceTree) customItemClicked(item *widgets.QTreeWidgetItem, column int) {
	t.detail.Refresh(item)
}

func (t *DeviceTree) customItemChanged(item *widgets.QTreeWidgetItem, column int) {
	err := t.Save()
	if err != nil {
		uiutil.MessageBoxError(err.Error())
		return
	}
}

func (t *DeviceTree) customContextMenuRequested(p *core.QPoint) {
	if t.ContextMenu == nil {
		t.ContextMenu = widgets.NewQMenu(t)
		menuAdd := t.ContextMenu.AddAction("Add child")
		menuAdd.ConnectTriggered(func(checked bool) { t.AddItem(p, NewDeviceItem("New item", "value")) })
		menuRemove := t.ContextMenu.AddAction("Remove node")
		menuRemove.ConnectTriggered(func(checked bool) { t.RemoveItem(p) })
	}
	t.ContextMenu.Exec2(t.MapToGlobal(p), nil)
}

func (t *DeviceTree) Save() error {
	item := t.InvisibleRootItem()
	node := t.ExportNode(item)
	err := tree.ExportNode(node, DEVICE_CONF_FILEPATH)
	if err != nil {
		return err
	}
	uiutil.MessageBoxInfo(fmt.Sprintf(
		"Configuration is saved as %q", DEVICE_CONF_FILEPATH))
	return nil
}

func (t *DeviceTree) Import() error {
	node := new(Node)
	err := tree.ImportNode(node, DEVICE_CONF_FILEPATH)
	if err != nil {
		if err.Error() == "nothing selected" {
			return nil
		}
		return err
	}
	t.Clear()
	for i := 0; i < len(node.Children); i++ {
		t.InvisibleRootItem().AddChild(t.ImportNode(node.Children[i]))
	}
	t.ExpandAll()
	return nil
}

func (t *DeviceTree) ImportNode(node Node) *widgets.QTreeWidgetItem {
	item := widgets.NewQTreeWidgetItem2([]string{node.Title}, 0)
	item.SetData(
		0,
		tree.DataRole(),
		core.NewQVariant17(node.Data),
	)
	for i := 0; i < len(node.Children); i++ {
		item.AddChild(t.ImportNode(node.Children[i]))
	}
	return item
}

func (t *DeviceTree) ExportNode(root *widgets.QTreeWidgetItem) Node {
	node := Node{}
	node.Title = root.Text(0)
	node.Data = tree.GetTreeItemData(root)
	for i := 0; i < root.ChildCount(); i++ {
		node.Children = append(node.Children, t.ExportNode(root.Child(i)))
	}
	return node
}
