package uiutil

import (
	"fmt"
	"github.com/therecipe/qt/widgets"
)

func FilePath() (string, error) {
	dialog := widgets.NewQFileDialog2(nil, "Select file...", "", "")
	if dialog.Exec() != int(widgets.QDialog__Accepted) {
		return "", fmt.Errorf("nothing selected")
	}
	filePath := dialog.SelectedFiles()[0]
	return filePath, nil
}

func MessageBoxInfo(message string) {
	widgets.QMessageBox_Information(
		nil,
		"Information",
		message,
		widgets.QMessageBox__Ok,
		widgets.QMessageBox__Close,
	)
}

func MessageBoxError(message string) {
	widgets.QMessageBox_Critical(
		nil,
		"Error",
		message,
		widgets.QMessageBox__Ok,
		widgets.QMessageBox__Close,
	)
}

func ShowDialog(title string, content interface{}) {
	dialog := widgets.NewQDialog(nil, 0)
	dialog.SetAttribute(core.Qt__WA_DeleteOnClose, true)
	dialog.SetWindowTitle(title)
	dialogLayout := widgets.NewQGridLayout2()
	dialogText := widgets.NewQLabel2(fmt.Sprintf("%v", content), nil, 0)
	dialogLayout.AddWidget3(dialogText, 0, 0, 1, 2, 0)
	dialog.SetLayout(dialogLayout)
	dialog.Show()
}
