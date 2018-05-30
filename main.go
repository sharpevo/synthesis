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
		`PRINT 10
PRINT 11
IMPORT testscripts/script2
PRINT 12
PRINT 13
RETRY -1 5
ASYNC testscripts/script3
PRINT 14
PRINT 15`)

	widget.Layout().AddWidget(input)

	button := widgets.NewQPushButton2("and click me!", nil)
	button.ConnectClicked(func(bool) {
		command.Init(CommandMap)
		statementGroup := command.StatementGroup{Execution: command.SYNC}
		command.ParseReader(
			strings.NewReader(input.ToPlainText()),
			&statementGroup,
		)
		resultList, _ := statementGroup.Execute()
		widgets.QMessageBox_Information(nil, "OK", strings.Join(resultList, "\n"), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})
	widget.Layout().AddWidget(button)

	window.Show()
	app.Exec()
}
