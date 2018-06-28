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
PRINT 0-2
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
	result.SetStyleSheet("QTextEdit { background-color: #e6e6e6}")

	suspButton := widgets.NewQPushButton2("SUSPEND", nil)
	resuButton := widgets.NewQPushButton2("RESUME", nil)

	suspButton.SetEnabled(false)
	resuButton.SetEnabled(false)

	suspend := false

	terminatecc := make(chan chan interface{}, 1)
	defer close(terminatecc)

	runButton := widgets.NewQPushButton2("RUN", nil)
	runButton.ConnectClicked(func(bool) {

		if len(terminatecc) != 0 {
			return
		}
		suspButton.SetEnabled(true)

		result.SetText("RUNNING")

		terminatec := make(chan interface{})
		terminatecc <- terminatec

		command.InitParser(CommandMap)
		statementGroup := command.StatementGroup{Execution: command.SYNC}
		command.ParseReader(
			strings.NewReader(input.ToPlainText()),
			&statementGroup,
		)
		resultList := []string{}

		go func() {
			for resp := range statementGroup.Execute(terminatec, &suspend, nil) {
				if resp.Error != nil {
					//fmt.Println(resp.Error)
					resultList = append(resultList, fmt.Sprintf("%s", resp.Error))
				}
				resultList = append(resultList, fmt.Sprintf("%s", resp.Output))
				result.SetText(strings.Join(resultList, "\n"))
			}
			result.SetText(strings.Join(resultList, "\n") + "\n\nDONE")
			if len(terminatecc) == 1 {
				t := <-terminatecc
				close(t)
			}

			suspButton.SetEnabled(false)
			resuButton.SetEnabled(false)
		}()
	})

	termButton := widgets.NewQPushButton2("TERMINATE", nil)
	termButton.ConnectClicked(func(bool) {
		go func() {
			if len(terminatecc) != 1 {
				return
			}
			//terminatec := <-terminatecc
			//close(terminatec)
			close(<-terminatecc)
			suspend = false
		}()
	})

	suspButton.ConnectClicked(func(bool) {
		go func() {
			suspend = true
			suspButton.SetEnabled(false)
			resuButton.SetEnabled(true)
		}()
	})

	resuButton.ConnectClicked(func(bool) {
		go func() {
			suspend = false
			suspButton.SetEnabled(true)
			resuButton.SetEnabled(false)
		}()
	})

	inputGroup := widgets.NewQGroupBox2("Commands", nil)
	inputLayout := widgets.NewQGridLayout2()
	inputLayout.AddWidget3(input, 0, 0, 1, 2, 0)
	inputLayout.AddWidget3(runButton, 1, 0, 1, 2, 0)
	inputGroup.SetLayout(inputLayout)

	outputGroup := widgets.NewQGroupBox2("Results", nil)
	outputLayout := widgets.NewQGridLayout2()
	outputLayout.AddWidget3(result, 0, 0, 1, 2, 0)
	outputLayout.AddWidget3(termButton, 1, 0, 1, 2, 0)
	outputLayout.AddWidget(suspButton, 2, 0, 0)
	outputLayout.AddWidget(resuButton, 2, 1, 0)
	outputGroup.SetLayout(outputLayout)

	layout := widgets.NewQGridLayout2()
	layout.AddWidget(inputGroup, 0, 0, 0)
	layout.AddWidget(outputGroup, 0, 1, 0)
	widget.SetLayout(layout)

	window.Show()
	app.Exec()
}
