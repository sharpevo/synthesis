package sequence

import (
	"github.com/therecipe/qt/widgets"
)

func NewViewGroup(imageItem *widgets.QGraphicsPixmapItem) *widgets.QGroupBox {
	group := widgets.NewQGroupBox2("Preview", nil)
	layout := widgets.NewQGridLayout2()
	group.SetLayout(layout)

	scene := widgets.NewQGraphicsScene(nil)
	view := widgets.NewQGraphicsView(nil)
	scene.AddItem(imageItem)
	view.SetScene(scene)
	view.Show()

	layout.AddWidget(view, 0, 0, 0)

	return group
}
