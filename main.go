package main

import (
	//"fmt"
	"os"
	"strings"

	"github.com/therecipe/qt/widgets"
	command "posam/commandparser"
)

var CommandMap = map[string]command.FunctionType{
	"PRINT":  CmdPrint,
	"IMPORT": command.CmdImport,
	"ASYNC":  command.CmdAsync,
	"RETRY":  command.CmdRetry,
}

func CmdPrint(args ...string) (string, error) {
	return "Print: " + args[0], nil
}

func main() {

	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(250, 200)
	window.SetWindowTitle("Hello Widgets Example")

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(widget)

	input := widgets.NewQTextEdit(nil)
	input.SetPlainText(
		`PRINT 0-1
PRINT 0-2
IMPORT testscripts/script1
PRINT 0-3
PRINT 0-4
RETRY -2 5
ASYNC testscripts/script2
PRINT 0-5
PRINT 0-6`)

	widget.Layout().AddWidget(input)

	button := widgets.NewQPushButton2("RUN", nil)
	button.ConnectClicked(func(bool) {
		command.Init(CommandMap)
		statementGroup := command.StatementGroup{Execution: command.SYNC}
		command.ParseReader(
			strings.NewReader(input.ToPlainText()),
			&statementGroup,
		)
		resultList, _ := statementGroup.Execute(nil)
		widgets.QMessageBox_Information(nil, "OK", strings.Join(resultList, "\n"), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})
	widget.Layout().AddWidget(button)

	window.Show()
	app.Exec()
}
