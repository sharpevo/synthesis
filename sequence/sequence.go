package sequence

import (
	//"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var imageItem *widgets.QGraphicsPixmapItem

func NewSequence() *widgets.QGroupBox {
	group := widgets.NewQGroupBox2("Sequence", nil)
	layout := widgets.NewQGridLayout2()
	group.SetLayout(layout)

	qimg := gui.NewQImage()
	imageItem = widgets.NewQGraphicsPixmapItem2(
		gui.NewQPixmap().FromImage(qimg, 0),
		nil,
	)

	viewGroup := NewViewGroup(imageItem)
	layout.AddWidget(viewGroup, 0, 0, 0)

	inputGroup := NewInputGroup()
	layout.AddWidget(inputGroup, 0, 1, 0)

	layout.SetColumnStretch(0, 1)
	layout.SetColumnStretch(1, 1)

	return group
}
