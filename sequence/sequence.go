package sequence

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"image"
	"image/png"
	"os"
	"posam/util/platform"
	"strings"
)

var imageItem *widgets.QGraphicsPixmapItem

func NewSequence() *widgets.QGroupBox {

	group := widgets.NewQGroupBox2("sequence", nil)
	layout := widgets.NewQGridLayout2()

	img := gui.NewQImage()
	img.Load("test.png", "png")
	img = img.Scaled2(1000, 1000, core.Qt__IgnoreAspectRatio, core.Qt__FastTransformation)

	scene := widgets.NewQGraphicsScene(nil)
	view := widgets.NewQGraphicsView(nil)

	imageItem = widgets.NewQGraphicsPixmapItem2(gui.NewQPixmap().FromImage(img, 0), nil)
	scene.AddItem(imageItem)
	view.SetScene(scene)
	view.Show()

	layout.AddWidget(view, 0, 0, 0)
	layout.AddWidget(NewSequenceDetail(), 0, 1, 0)

	layout.SetColumnStretch(0, 1)
	layout.SetColumnStretch(1, 1)

	group.SetLayout(layout)
	return group
}

func NewSequenceDetail() *widgets.QGroupBox {
	group := widgets.NewQGroupBox2("Config", nil)
	layout := widgets.NewQGridLayout2()

	startxLabel := widgets.NewQLabel2("Starts at x:", nil, 0)
	startxInput := widgets.NewQLineEdit(nil)
	startyLabel := widgets.NewQLabel2("Starts at y:", nil, 0)
	startyInput := widgets.NewQLineEdit(nil)

	spacexLabel := widgets.NewQLabel2("Space x:", nil, 0)
	spacexInput := widgets.NewQLineEdit(nil)
	spaceyLabel := widgets.NewQLabel2("Space y:", nil, 0)
	spaceyInput := widgets.NewQLineEdit(nil)

	spaceBlockxLabel := widgets.NewQLabel2("Space slide x:", nil, 0)
	spaceBlockxInput := widgets.NewQLineEdit(nil)
	spaceBlockyLabel := widgets.NewQLabel2("Space slide y:", nil, 0)
	spaceBlockyInput := widgets.NewQLineEdit(nil)

	sequenceInput := widgets.NewQTextEdit(nil)

	generateButton := widgets.NewQPushButton2("GENERATE", nil)
	generateButton.ConnectClicked(func(bool) {
		generateImage(
			2, 5, 10, 20,
			sequenceInput.ToPlainText(),
		)
	})

	layout.AddWidget(startxLabel, 0, 0, 0)
	layout.AddWidget(startxInput, 0, 1, 0)
	layout.AddWidget(startyLabel, 1, 0, 0)
	layout.AddWidget(startyInput, 1, 1, 0)

	layout.AddWidget(spacexLabel, 2, 0, 0)
	layout.AddWidget(spacexInput, 2, 1, 0)
	layout.AddWidget(spaceyLabel, 3, 0, 0)
	layout.AddWidget(spaceyInput, 3, 1, 0)

	layout.AddWidget(spaceBlockxLabel, 4, 0, 0)
	layout.AddWidget(spaceBlockxInput, 4, 1, 0)
	layout.AddWidget(spaceBlockyLabel, 5, 0, 0)
	layout.AddWidget(spaceBlockyInput, 5, 1, 0)

	layout.AddWidget3(sequenceInput, 6, 0, 1, 2, 0)
	layout.AddWidget3(generateButton, 7, 0, 1, 2, 0)
	group.SetLayout(layout)
	return group
}

func generateImage(
	spaceX int,
	spaceY int,
	spaceBlockx int,
	spaceBlocky int,
	sequences string,
) {
	p := platform.NewPlatform(100, 100)
	img := image.NewRGBA(image.Rect(0, 0, p.Width, p.Height))
	lines := strings.Split(sequences, "\n")
	blocks := make([]*platform.Block, len(strings.Split(lines[0], ",")))
	xoffset := 0
	yoffset := 0
	for _, line := range strings.Split(sequences, "\n") {
		for x, seq := range strings.Split(line, ",") {
			if blocks[x] == nil {
				blocks[x] = &platform.Block{}
				blocks[x].PositionX = xoffset + spaceBlockx
				blocks[x].PositionY = yoffset + spaceBlocky
				blocks[x].SpaceX = spaceX
				blocks[x].SpaceY = spaceY
			}
			blocks[x].AddRow(seq)
			xoffset += len(blocks[x].Sequence[0])*spaceX + spaceBlockx
		}
		yoffset += spaceBlocky
	}
	fmt.Println(strings.Split(sequences, "\n"))

	for _, block := range blocks {
		p.AddBlock(block)
	}

	for posy, row := range p.Dots {
		for posx, dot := range row {
			if dot == nil {
				continue
			}
			fmt.Println(posx, posy, dot.Base.Name)
			img.Set(posx, posy, dot.Base.Color)
		}
	}
	outputFile, _ := os.Create("test.png")
	png.Encode(outputFile, img)
	outputFile.Close()

	qimg := gui.NewQImage()
	qimg.Load("test.png", "png")
	qimg = qimg.Scaled2(1000, 1000, core.Qt__IgnoreAspectRatio, core.Qt__FastTransformation)
	imageItem.SetPixmap(gui.NewQPixmap().FromImage(qimg, 0))
}
