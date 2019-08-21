package instree

import (
	"github.com/therecipe/qt/widgets"
	"io/ioutil"
	"log"
	"synthesis/gui/uiutil"
)

func NewInsTree(
	detail *InstructionDetail,
	runButton *widgets.QPushButton,
	input *widgets.QTextEdit,
) *widgets.QGroupBox {
	treeGroup := widgets.NewQGroupBox2(
		"Stepwise instruction programming interface", nil)
	layout := widgets.NewQGridLayout2()
	treeWidget := NewTree(detail, runButton, input)
	layout.AddWidget3(treeWidget, 0, 0, 1, 2, 0)

	treeExportButton := widgets.NewQPushButton2("EXPORT", nil)
	treeExportButton.ConnectClicked(func(bool) {
		filePath, err := uiutil.FilePath()
		if err != nil {
			uiutil.MessageBoxError(err.Error())
			return
		}
		treeWidget.Export(filePath)
	})
	layout.AddWidget2(treeExportButton, 1, 0, 0)

	treeImportButton := widgets.NewQPushButton2("IMPORT", nil)
	treeImportButton.ConnectClicked(func(bool) {
		filePath, err := uiutil.FilePath()
		if err != nil {
			if err.Error() == "nothing selected" {
				return
			}
			uiutil.MessageBoxError(err.Error())
			return
		}
		err = treeWidget.Import(filePath)
		if err != nil {
			uiutil.MessageBoxError(err.Error())
		}
		uiutil.MessageBoxInfo("Imported")
	})
	layout.AddWidget2(treeImportButton, 1, 1, 0)

	treeGenerateButton := widgets.NewQPushButton2("RUN", nil)
	treeGenerateButton.ConnectClicked(func(bool) {
		filePath, err := treeWidget.Generate()
		if err != nil {
			log.Println(err)
			return
		}
		instBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Println(err)
			return
		}
		input.SetPlainText(string(instBytes))
		runButton.Click()
	})
	layout.AddWidget3(treeGenerateButton, 2, 0, 1, 2, 0)

	treeGroup.SetLayout(layout)
	return treeGroup
}
