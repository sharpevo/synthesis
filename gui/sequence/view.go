package sequence

import (
	"bytes"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"image"
	"image/png"
	"time"
)

var scene *Scene
var paintedc chan struct{}

func NewViewGroup() *widgets.QGroupBox {
	group := widgets.NewQGroupBox2("Preview", nil)
	layout := widgets.NewQGridLayout2()
	group.SetLayout(layout)

	scene = NewScene(nil)
	view := widgets.NewQGraphicsView(nil)
	scene.ConnectKeyPressEvent(func(e *gui.QKeyEvent) {
		if e.Modifiers() == core.Qt__ControlModifier {
			switch int32(e.Key()) {
			case int32(core.Qt__Key_Plus):
				view.Scale(1.2, 1.2)
			case int32(core.Qt__Key_Minus):
				view.Scale(0.8, 0.8)
			}
		}
	})
	scene.ConnectWheelEvent(func(e *widgets.QGraphicsSceneWheelEvent) {
		if e.Modifiers() == core.Qt__ControlModifier {
			if e.Delta() > 0 {
				view.Scale(1.2, 1.2)
			} else {
				view.Scale(0.8, 0.8)
			}
		}
	})
	paintedc = make(chan struct{})
	scene.ConnectUpdatePixmap(func() {
		var imagebuff bytes.Buffer
		png.Encode(&imagebuff, scene.image)
		imagebytes := imagebuff.Bytes()
		pixmap := gui.NewQPixmap()
		pixmap.LoadFromData2(
			core.NewQByteArray2(string(imagebytes), len(imagebytes)),
			"png",
			core.Qt__AutoColor,
		)
		imageItem := widgets.NewQGraphicsPixmapItem2(pixmap.Scaled2(
			5*pixmap.Width(),
			5*pixmap.Height(),
			core.Qt__IgnoreAspectRatio,
			core.Qt__FastTransformation,
		), nil)
		scene.Clear()
		scene.AddItem(imageItem)
		go func() {
			select {
			case <-time.After(1 * time.Second):
				return
			case paintedc <- struct{}{}:
			}
		}()
	})
	view.SetScene(scene)
	view.Show()

	layout.AddWidget(view, 0, 0, 0)

	return group
}

type Scene struct {
	widgets.QGraphicsScene
	image *image.RGBA
	_     func() `signal:updatePixmap`
}
