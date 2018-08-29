package uiutil

import (
	"fmt"
	"github.com/therecipe/qt/widgets"
)

const (
	DEVICE_CONF_FILE = "devices.bin"
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
