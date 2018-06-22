package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/therecipe/qt/widgets"
	command "posam/commandparser"
)

var CommandMap = map[string]command.Commander{
	"PRINT":  &Print,
	"SLEEP":  &command.Sleep,
	"IMPORT": &command.Import,
	"ASYNC":  &command.Async,
	"RETRY":  &command.Retry,
}

type CommandPrint struct {
	command.Command
}

var Print CommandPrint

func (c *CommandPrint) Execute(args ...string) (interface{}, error) {
	return "Print: " + args[0], nil
}

func main() {

	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(500, 400)
	window.SetWindowTitle("POSaM Control Software by iGeneTech")

	widget := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(widget)

	input := widgets.NewQTextEdit(nil)
	input.SetPlainText(
		`PRINT 0-1
PRIN2T 0-2
IMPORT testscripts/script1
PRINT 0-3
PRINT 0-4
RETRY -2 5
SLEEP 5
ASYNC testscripts/script2
PRINT 0-5
PRINT 0-6`)

	result := widgets.NewQTextEdit(nil)
	//result := widgets.NewQTextBrowser(nil)
	result.SetReadOnly(true)

	terminatec := make(chan interface{})
	defer close(terminatec)

	runButton := widgets.NewQPushButton2("RUN", nil)
	runButton.ConnectClicked(func(bool) {
		command.InitParser(CommandMap)
		statementGroup := command.StatementGroup{Execution: command.SYNC}
		command.ParseReader(
			strings.NewReader(input.ToPlainText()),
			&statementGroup,
		)
		resultList := []string{}

		go func() {
			for resp := range statementGroup.Execute(terminatec, nil) {
				if resp.Error != nil {
					//fmt.Println(resp.Error)
					resultList = append(resultList, fmt.Sprintf("%s", resp.Error))
				}
				resultList = append(resultList, fmt.Sprintf("%s", resp.Output))
				result.SetText(strings.Join(resultList, "\n"))
			}
		}()
	})

	termButton := widgets.NewQPushButton2("TERMINATE", nil)
	termButton.ConnectClicked(func(bool) {
		go func() {
			close(terminatec)
		}()
	})

	inputGroup := widgets.NewQGroupBox2("Input", nil)
	inputLayout := widgets.NewQGridLayout2()
	inputLayout.AddWidget(input, 0, 0, 0)
	inputLayout.AddWidget(runButton, 1, 0, 0)
	inputGroup.SetLayout(inputLayout)

	outputGroup := widgets.NewQGroupBox2("Output", nil)
	outputLayout := widgets.NewQGridLayout2()
	outputLayout.AddWidget(result, 0, 0, 0)
	outputLayout.AddWidget(termButton, 1, 0, 0)
	outputGroup.SetLayout(outputLayout)

	layout := widgets.NewQGridLayout2()
	layout.AddWidget(inputGroup, 0, 0, 0)
	layout.AddWidget(outputGroup, 0, 1, 0)
	widget.SetLayout(layout)

	window.Show()
	app.Exec()
}
