package sequence

import (
	"github.com/therecipe/qt/widgets"
)

func NewSequence() *widgets.QGroupBox {
	group := widgets.NewQGroupBox2("Sequence", nil)
	layout := widgets.NewQGridLayout2()
	group.SetLayout(layout)

	viewGroup := NewViewGroup()
	layout.AddWidget2(viewGroup, 0, 0, 0)

	inputGroup, previewGroup := NewInputGroup()
	layout.AddWidget2(inputGroup, 0, 1, 0)
	viewGroup.Layout().AddWidget(previewGroup)

	layout.SetColumnStretch(0, 1)
	layout.SetColumnStretch(1, 1)

	return group
}
