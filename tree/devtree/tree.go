package devtree

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"posam/gui/tree"
	"posam/gui/uiutil"
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
	rootNode := treeWidget.InvisibleRootItem()
	anItem := NewDeviceItem("Root Device", "value")
	rootNode.AddChild(anItem)

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
	err := node.Write("/tmp/exportnode.gob")
	if err != nil {
		return err
	}
	uiutil.MessageBoxInfo("Exported")
	return nil
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
