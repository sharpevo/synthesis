package uiutil

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type Application struct {
	widgets.QApplication
	_ func(message string) `slot:showMessageSlot`
	_ func(status string)  `slot:updateMotorStatusSlot`

	MotorStatusLabel *widgets.QLabel
}

var App *Application

func NewApp(args []string) *Application {
	App = NewApplication(len(args), args)
	App.MotorStatusLabel = widgets.NewQLabel2("Motor:", nil, 0)
	//App = &Application{}
	//App.QApplication = widgets.NewQApplication(len(args), args)
	App.ConnectShowMessageSlot(func(message string) {
		App.showMessage(message)
	})
	App.ConnectUpdateMotorStatusSlot(func(status string) {
		App.updateMotorStatus(status)
	})
	return App
}

func (a *Application) showMessage(msg interface{}) {
	msgBox := widgets.NewQMessageBox(nil)
	msgBox.SetIcon(widgets.QMessageBox__Warning)
	msgBox.SetWindowTitle("Warning")
	msgBox.SetStandardButtons(widgets.QMessageBox__Ok)
	msgBox.SetModal(false)
	msgBox.SetAttribute(core.Qt__WA_DeleteOnClose, true)
	msgBox.SetText(fmt.Sprintf("%v", msg))
	msgBox.Show()
}

func (a *Application) updateMotorStatus(status interface{}) {
	a.MotorStatusLabel.SetText(fmt.Sprintf("%v", status))
}

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
