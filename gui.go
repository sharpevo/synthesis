package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"posam/dao"
	"posam/dao/alientek"
	"posam/dao/ricoh_g5"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/therecipe/qt/widgets"
	"posam/gui/tree/devtree"
	"posam/gui/tree/instree"
	"posam/gui/uiutil"
	"posam/instruction"
	"posam/interpreter"
	"posam/interpreter/vrb"
)

const (
	CMD_ASYNC = `PRINT 0-1
PRINT 0-2
IMPORT testscripts/script1
PRINT 0-3
PRINT 0-4
RETRY -2 5
SLEEP 5
ASYNC testscripts/script2
PRINT 0-5
PRINT 0-6`
	CMD_LED = `LED on
SLEEP 1
LED off`
	CMD_SERIAL = `SENDSERIAL 010300010001D5CA 55 018302c0f1
SLEEP 3
SENDSERIAL 010200010001E80A 55 018202c161`
	CMD_LED_SERIAL_SIMPLE = `LED on
SLEEP 3
SENDSERIAL 010200010001E80A 55 018202c161`
	CMD_LED_SERIAL = `LED on
PRINT 1
SLEEP 1
PRINT 2
SENDSERIAL 010200010001E80A 525 018202c161
RETRY -2 3
PRINT 3
SLEEP 1
PRINT 4`
	CMD_PRINTER = `ERRORCODE var1
GETVAR var1`
	CMD_WAVEFORM = `WAVEFORM var1 1 2 11.22 1 1.1 2.2 3.3 4.4 5.5 6.6 7.7 8.8 9.9 10.10 11.11 12.12 1
ASYNC testscripts/tcpconcurrency`
	CMD_VARIABLE_SETTER_GETTER = `SETVAR var1 This is a string variable
GETVAR var1
SETVAR var2 2
GETVAR var2
SETVAR var3 3.0
GETVAR var3`
	CMD_VARIABLE_GLOBAL = `SETVAR globalvar1 This is a global string variable
GETVAR globalvar1
IMPORT testscripts/variable/modification
GETVAR localvar1
GETVAR globalvar1`
	CMD_CF = `PRINT ---- start ----
SETVAR var1 11.11
GETVAR var1
SETVAR var2 11.11
GETVAR var2
CMPVAR var1 var2
EQGOTO 10
PRINT not here
PRINT not here
PRINT equal redirected
SETVAR var3 33.33
GETVAR var3
CMPVAR var1 var3
LTGOTO 17
PRINT not here
PRINT not here
PRINT less than redirected
CMPVAR var3 var1
GTGOTO 22
PRINT not here
PRINT not here
PRINT greater than redirected
SETVAR var4 string1
GETVAR var4
SETVAR var5 string2
GETVAR var5
CMPVAR var4 var5
NEGOTO 31
PRINT not here
PRINT not here
PRINT not equal redirected
CMPVAR var1 var5
ERRGOTO 36
PRINT not here
PRINT not here
PRINT error redirected
SETVAR loopcount 3
PRINT loop start
PRINT loop body
PRINT loop end
SLEEP 2
LOOP 38 loopcount
PRINT last command
PRINT ---- end ----
RETURN nil
PRINT never executed
PRINT never executed
PRINT never executed`
)

var InstructionMap = make(interpreter.InstructionMapt)
var InstructionDaoMap = make(map[string]interpreter.InstructionMapt)

type QMessageBoxWithCustomSlot struct {
	widgets.QMessageBox
	_ func(message string) `slot:showMessageBoxSlot`
}

func init() {

	InstructionDaoMap[devtree.DEV_TYPE_UNK] = dao.InstructionMap
	InstructionDaoMap[devtree.DEV_TYPE_ALT] = alientek.InstructionMap
	InstructionDaoMap[devtree.DEV_TYPE_RCG] = ricoh_g5.InstructionMap
	// TODO: CAN
	buildInstructionMap()
}

func main() {

	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(500, 400)
	window.SetWindowTitle("POSaM Control Software by iGeneTech")

	tabWidget := widgets.NewQTabWidget(nil)
	window.SetCentralWidget(tabWidget)

	msgBox := NewQMessageBoxWithCustomSlot(nil)

	msgBox.ConnectShowMessageBoxSlot(func(message string) {
		msgBox.SetIcon(widgets.QMessageBox__Warning)
		msgBox.SetWindowTitle("Error")
		msgBox.SetText(message)
		msgBox.SetStandardButtons(widgets.QMessageBox__Ok)
		msgBox.Exec()
	})

	input := widgets.NewQTextEdit(nil)
	input.SetPlainText(CMD_CF)
	input.SetVisible(false)

	instDetail := instree.NewInstructionDetail(InstructionDaoMap)
	instDetail.InitDevInput(devtree.ParseConnList())

	// result group

	result := widgets.NewQTextEdit(nil)
	result.SetReadOnly(true)
	result.SetStyleSheet("QTextEdit { background-color: #e6e6e6}")

	suspButton := widgets.NewQPushButton2("SUSPEND", nil)
	resuButton := widgets.NewQPushButton2("RESUME", nil)

	suspButton.SetEnabled(false)
	resuButton.SetEnabled(false)

	suspend := false
	resumec := make(chan<- interface{})

	terminatecc := make(chan chan interface{}, 1)
	defer close(terminatecc)

	runButton := widgets.NewQPushButton2("RUN", nil)
	runButton.SetVisible(false)
	runButton.ConnectClicked(func(bool) {

		if len(terminatecc) != 0 {
			return
		}
		suspButton.SetEnabled(true)

		stack := interpreter.NewStack()

		devtree.ParseDeviceConf()
		for k, v := range devtree.ConfMap {
			variable, err := vrb.NewVariable(k, v)
			if err != nil {
				uiutil.MessageBoxError(err.Error())
				return
			}
			stack.Set(variable)
		}

		for _, s := range devtree.ConnMap[devtree.DEV_TYPE_ALT] {
			base := devtree.ComposeVarName(s, devtree.PRT_CONN)
			name, _ := stack.Get(
				devtree.ComposeVarName(base, devtree.PRT_SRL_NAME))
			baud, _ := stack.Get(
				devtree.ComposeVarName(base, devtree.PRT_SRL_BAUD))
			character, _ := stack.Get(
				devtree.ComposeVarName(base, devtree.PRT_SRL_CHARACTER))
			stop, _ := stack.Get(
				devtree.ComposeVarName(base, devtree.PRT_SRL_STOP))
			parity, _ := stack.Get(
				devtree.ComposeVarName(base, devtree.PRT_SRL_PARITY))
			deviceCode, _ := stack.Get(
				devtree.ComposeVarName(base, devtree.PRT_SRL_CODE))
			err := initSerialDevice(
				fmt.Sprintf("%v", name.Value),
				fmt.Sprintf("%v", baud.Value),
				fmt.Sprintf("%v", character.Value),
				fmt.Sprintf("%v", stop.Value),
				fmt.Sprintf("%v", parity.Value),
				fmt.Sprintf("%v", deviceCode.Value),
			)
			if err != nil {
				uiutil.MessageBoxError(err.Error())
			}
		}

		for _, s := range devtree.ConnMap[devtree.DEV_TYPE_RCG] {
			base := devtree.ComposeVarName(s, devtree.PRT_CONN)
			network, _ := stack.Get(
				devtree.ComposeVarName(base, devtree.PRT_TCP_NETWORK))
			address, _ := stack.Get(
				devtree.ComposeVarName(base, devtree.PRT_TCP_ADDRESS))
			timeout, _ := stack.Get(
				devtree.ComposeVarName(base, devtree.PRT_TCP_TIMEOUT))
			err := initTCPDevice(
				fmt.Sprintf("%v", network.Value),
				fmt.Sprintf("%v", address.Value),
				fmt.Sprintf("%v", timeout.Value),
			)
			if err != nil {
				uiutil.MessageBoxError(err.Error())
			}
		}

		instDetail.InitDevInput(devtree.GetConnMap())

		// TODO: init CAN devices

		result.SetText("RUNNING")

		terminatec := make(chan interface{})
		terminatecc <- terminatec

		interpreter.InitParser(InstructionMap)
		statementGroup := interpreter.StatementGroup{
			Execution: interpreter.SYNC,
			Stack:     stack,
		}
		interpreter.ParseReader(
			strings.NewReader(input.ToPlainText()),
			&statementGroup,
		)
		resultList := []string{}

		go func() {
			completec := make(chan interface{})
			go func() {
				<-completec
			}()
			for resp := range statementGroup.Execute(terminatec, completec) {
				if resp.Error != nil {
					resumec = resp.Completec
					if resp.IgnoreError {
						resp.Completec <- true
					} else {
						msgBox.ShowMessageBoxSlot(resp.Error.Error())
						suspendExecution(&suspend, suspButton, resuButton)
					}
				} else {
					if suspend {
						for {
							if !suspend {
								break
							}
							time.Sleep(1 * time.Second)
						}
					}
					resp.Completec <- true
				}
				resultList = append(resultList, fmt.Sprintf("%v", resp.Output))
				result.SetText(strings.Join(resultList, "\n"))
			}
			result.SetText(strings.Join(resultList, "\n") + "\n\nDONE")
			if len(terminatecc) == 1 {
				t := <-terminatecc
				close(t)
			}

			msgBox.ShowMessageBoxSlot("Done")

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
		suspendExecution(&suspend, suspButton, resuButton)
	})

	resuButton.ConnectClicked(func(bool) {
		go func() {
			suspend = false
			suspButton.SetEnabled(true)
			resuButton.SetEnabled(false)
			resumec <- true
		}()
	})

	inputGroup := widgets.NewQGroupBox2("Instructions", nil)
	inputLayout := widgets.NewQGridLayout2()
	inputLayout.AddWidget(instDetail.GroupBox, 0, 0, 0)
	inputLayout.AddWidget(input, 1, 0, 0)
	inputLayout.AddWidget(runButton, 4, 0, 0)
	inputGroup.SetLayout(inputLayout)

	outputGroup := widgets.NewQGroupBox2("Results", nil)
	outputLayout := widgets.NewQGridLayout2()
	outputLayout.AddWidget3(result, 0, 0, 1, 2, 0)
	outputLayout.AddWidget3(termButton, 2, 0, 1, 2, 0)
	outputLayout.AddWidget(suspButton, 1, 0, 0)
	outputLayout.AddWidget(resuButton, 1, 1, 0)
	outputGroup.SetLayout(outputLayout)

	insTab := widgets.NewQWidget(nil, 0)
	insTabLayout := widgets.NewQGridLayout2()
	insTabLayout.AddWidget3(instree.NewInsTree(instDetail, runButton, input), 0, 0, 2, 1, 0)
	insTabLayout.AddWidget3(inputGroup, 0, 1, 1, 1, 0)
	insTabLayout.AddWidget3(outputGroup, 1, 1, 1, 1, 0)

	insTabLayout.SetColumnStretch(0, 1)
	insTabLayout.SetColumnStretch(1, 1)

	insTab.SetLayout(insTabLayout)
	tabWidget.AddTab(insTab, "Instructions")

	devTab := widgets.NewQWidget(nil, 0)
	devTabLayout := widgets.NewQGridLayout2()
	devTabLayout.AddWidget(devtree.NewDevTree(instDetail), 0, 0, 0)
	devTab.SetLayout(devTabLayout)

	tabWidget.AddTab(devTab, "Devices")
	window.Show()
	app.Exec()
}

func suspendExecution(
	suspend *bool,
	suspButton *widgets.QPushButton,
	resuButton *widgets.QPushButton) {
	*suspend = true
	suspButton.SetEnabled(false)
	resuButton.SetEnabled(true)
}

func initSerialDevice(
	device string,
	baud string,
	character string,
	stop string,
	parity string,
	deviceCode string,
) (err error) {
	if alientek.Instance(deviceCode) != nil {
		log.Printf("Device %q has been initialized\n", deviceCode)
		return
	}
	_, err = alientek.NewDao(
		device,
		baud,
		character,
		stop,
		parity,
		deviceCode,
	)
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command(
			"mode",
			device,
			"BAUD="+baud,
			"PARITY="+parity,
			"DATA="+character,
			"STOP="+stop,
		)
		log.Printf(
			"Initializing device %q with baud rate %q, character bits %q, %q stop bits per characeter, and parity bits %q\n",
			device,
			baud,
			character,
			stop,
			parity,
		)
		err = cmd.Run()
		if err != nil {
			msg := fmt.Sprintf("failed to init device %q: %s", device, err)
			log.Println(msg)
			return fmt.Errorf(msg)
		}
	case "linux":
		cmd := exec.Command(
			"stty",
			"-F",
			device,
			baud,
			"cs"+character,
			"-parenb",
			"-cstopb",
		)
		log.Printf(
			"Initializing device %q with baud rate %q, character bits %q, 1 stop bits per characeter, and none parity bits\n",
			device,
			baud,
			character,
		)
		err = cmd.Run()
		if err != nil {
			msg := fmt.Sprintf("failed to init device %q: %s", device, err)
			log.Println(msg)
			return fmt.Errorf(msg)
		}
	case "darwin":
		log.Println("TBD...")
	default:
		msg := "unknown os"
		log.Println(msg)
		return fmt.Errorf(msg)
	}

	if err != nil {
		return err
	}

	return
}

func initTCPDevice(network string, address string, secondString string) (err error) {
	if ricoh_g5.Instance(address) != nil {
		log.Printf("Device %q has been initialized\n", address)
		return
	}
	secondInt, err := strconv.Atoi(secondString)
	if err != nil {
		return err
	}
	_, err = ricoh_g5.NewDao(network, address, secondInt)
	if err != nil {
		return err
	}
	return nil
}

func buildInstructionMap() {
	for _, instructionMap := range InstructionDaoMap {
		for k, v := range instructionMap {
			fmt.Println(k)
			InstructionMap[k] = v
		}
	}
}
