package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"posam/dao/alientek"
	"posam/dao/ricoh_g5"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"posam/instruction"
	"posam/interpreter"
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
	CMD_TMP = `ASYNC testscripts/serialconcurrency`
)

var InstructionMap = make(interpreter.InstructionMapt)

type InstructionPrint struct {
	instruction.Instruction
}

var Print InstructionPrint

func (c *InstructionPrint) Execute(args ...string) (interface{}, error) {
	return "Print: " + args[0], nil
}

type QMessageBoxWithCustomSlot struct {
	widgets.QMessageBox
	_ func(message string) `slot:showMessageBoxSlot`
}

func main() {

	InstructionMap.Set("PRINT", InstructionPrint{})
	InstructionMap.Set("SLEEP", instruction.InstructionSleep{})
	InstructionMap.Set("IMPORT", instruction.InstructionImport{})
	InstructionMap.Set("ASYNC", instruction.InstructionAsync{})
	InstructionMap.Set("RETRY", instruction.InstructionRetry{})
	InstructionMap.Set("LED", instruction.InstructionLed{})
	InstructionMap.Set("SENDSERIAL", instruction.InstructionSendSerial{})
	InstructionMap.Set("GETVAR", instruction.InstructionVariableGet{})
	InstructionMap.Set("ERRORCODE", instruction.InstructionPrinterHeadErrorCode{})
	InstructionMap.Set("PRINTERSTATUS", instruction.InstructionPrinterHeadPrinterStatus{})
	InstructionMap.Set("PRINTDATA", instruction.InstructionPrinterHeadPrintData{})
	InstructionMap.Set("WAVEFORM", instruction.InstructionPrinterHeadWaveform{})

	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(500, 400)
	window.SetWindowTitle("POSaM Control Software by iGeneTech")

	widget := widgets.NewQWidget(nil, 0)
	window.SetCentralWidget(widget)

	msgBox := NewQMessageBoxWithCustomSlot(nil)

	msgBox.ConnectShowMessageBoxSlot(func(message string) {
		msgBox.SetIcon(widgets.QMessageBox__Warning)
		msgBox.SetWindowTitle("Error")
		msgBox.SetText(message)
		msgBox.SetStandardButtons(widgets.QMessageBox__Ok)
		msgBox.Exec()
	})

	input := widgets.NewQTextEdit(nil)
	input.SetPlainText(CMD_WAVEFORM)

	// waveform group

	waveformGroup := widgets.NewQGroupBox2("WaveForm Builder", nil)

	waveformLineTimeLabel := widgets.NewQLabel2("Time", nil, 1)
	waveformLineTimeLabel.SetAlignment(core.Qt__AlignCenter)
	waveformLineVoltageLabel := widgets.NewQLabel2("Percentage", nil, 1)
	waveformLineVoltageLabel.SetAlignment(core.Qt__AlignCenter)
	waveformFallLabel := widgets.NewQLabel2("Fall:", nil, 0)
	waveformHoldLabel := widgets.NewQLabel2("Hold:", nil, 0)
	waveformRisingLabel := widgets.NewQLabel2("Rising:", nil, 0)
	waveformWaitLabel := widgets.NewQLabel2("Wait:", nil, 0)
	waveformMnLabel := widgets.NewQLabel2("Mn:", nil, 0)
	waveformVoltageLabel := widgets.NewQLabel2("Voltage:", nil, 0)

	waveformFallTimeInput := widgets.NewQDoubleSpinBox(nil)
	waveformFallTimeInput.SetMaximum(100)
	waveformFallPercentageInput := widgets.NewQDoubleSpinBox(nil)
	waveformFallPercentageInput.SetMaximum(100)
	waveformHoldTimeInput := widgets.NewQDoubleSpinBox(nil)
	waveformHoldTimeInput.SetMaximum(100)
	waveformHoldPercentageInput := widgets.NewQDoubleSpinBox(nil)
	waveformHoldPercentageInput.SetMaximum(100)
	waveformRisingTimeInput := widgets.NewQDoubleSpinBox(nil)
	waveformRisingTimeInput.SetMaximum(100)
	waveformRisingPercentageInput := widgets.NewQDoubleSpinBox(nil)
	waveformRisingPercentageInput.SetMaximum(100)
	waveformWaitTimeInput := widgets.NewQDoubleSpinBox(nil)
	waveformWaitTimeInput.SetMaximum(100)
	waveformWaitPercentageInput := widgets.NewQDoubleSpinBox(nil)
	waveformWaitPercentageInput.SetMaximum(100)
	waveformMnInput := widgets.NewQSpinBox(nil)
	waveformMnInput.SetMaximum(100)
	waveformVoltageInput := widgets.NewQDoubleSpinBox(nil)
	waveformVoltageInput.SetMaximum(100)

	waveformGenerateButton := widgets.NewQPushButton2("INSERT", nil)
	waveformGenerateButton.ConnectClicked(func(bool) {
		argumentList := []string{
			"WAVEFORM",
			"VAR",
			"HEADBOARD", // headboard
			"ROW",       // row
			fmt.Sprintf("%.2f", waveformVoltageInput.Value()),
			"COUNT", // segment count

			fmt.Sprintf("%.2f", waveformFallTimeInput.Value()),
			fmt.Sprintf("%.2f", waveformFallPercentageInput.Value()),
			fmt.Sprintf("%.2f", waveformHoldPercentageInput.Value()),

			fmt.Sprintf("%.2f", waveformHoldTimeInput.Value()),
			fmt.Sprintf("%.2f", waveformHoldPercentageInput.Value()),
			fmt.Sprintf("%.2f", waveformRisingPercentageInput.Value()),

			fmt.Sprintf("%.2f", waveformRisingTimeInput.Value()),
			fmt.Sprintf("%.2f", waveformRisingPercentageInput.Value()),
			fmt.Sprintf("%.2f", waveformWaitPercentageInput.Value()),

			fmt.Sprintf("%.2f", waveformWaitTimeInput.Value()),
			fmt.Sprintf("%.2f", waveformWaitPercentageInput.Value()),
			fmt.Sprintf("%.2f", waveformWaitPercentageInput.Value()),

			fmt.Sprintf("%d", waveformMnInput.Value()),
		}

		input.InsertPlainText(strings.Join(argumentList, " "))
	})

	waveformLayout := widgets.NewQGridLayout2()

	waveformLayout.AddWidget(waveformLineTimeLabel, 0, 1, 0)
	waveformLayout.AddWidget(waveformLineVoltageLabel, 0, 2, 0)

	waveformLayout.AddWidget(waveformFallLabel, 1, 0, 0)
	waveformLayout.AddWidget(waveformFallTimeInput, 1, 1, 0)
	waveformLayout.AddWidget(waveformFallPercentageInput, 1, 2, 0)
	waveformLayout.AddWidget(waveformHoldLabel, 2, 0, 0)
	waveformLayout.AddWidget(waveformHoldTimeInput, 2, 1, 0)
	waveformLayout.AddWidget(waveformHoldPercentageInput, 2, 2, 0)
	waveformLayout.AddWidget(waveformRisingLabel, 3, 0, 0)
	waveformLayout.AddWidget(waveformRisingTimeInput, 3, 1, 0)
	waveformLayout.AddWidget(waveformRisingPercentageInput, 3, 2, 0)
	waveformLayout.AddWidget(waveformWaitLabel, 4, 0, 0)
	waveformLayout.AddWidget(waveformWaitTimeInput, 4, 1, 0)
	waveformLayout.AddWidget(waveformWaitPercentageInput, 4, 2, 0)

	waveformLayout.AddWidget3(waveformVoltageLabel, 5, 0, 1, 1, 0)
	waveformLayout.AddWidget3(waveformVoltageInput, 5, 1, 1, 2, 0)
	waveformLayout.AddWidget3(waveformMnLabel, 6, 0, 1, 1, 0)
	waveformLayout.AddWidget3(waveformMnInput, 6, 1, 1, 2, 0)
	waveformLayout.AddWidget3(waveformGenerateButton, 7, 0, 1, 3, 0)

	waveformGroup.SetLayout(waveformLayout)

	// tcp group

	printerGroup := widgets.NewQGroupBox2("Printer", nil)

	printerNetworkLabel := widgets.NewQLabel2("Network:", nil, 0)
	printerAddressLabel := widgets.NewQLabel2("Address:", nil, 0)
	printerTimeoutLabel := widgets.NewQLabel2("Timeout:", nil, 0)

	printerNetworkInput := widgets.NewQLineEdit(nil)
	printerNetworkInput.SetPlaceholderText("tcp, tcp4, tcp6")
	printerNetworkInput.SetText("tcp")
	printerAddressInput := widgets.NewQLineEdit(nil)
	printerAddressInput.SetPlaceholderText("localhost:3000")
	printerAddressInput.SetText("192.168.100.215:21005")
	printerTimeoutInput := widgets.NewQLineEdit(nil)
	printerTimeoutInput.SetPlaceholderText("10")
	printerTimeoutInput.SetText("10")

	printerLayout := widgets.NewQGridLayout2()
	printerLayout.AddWidget(printerNetworkLabel, 0, 0, 0)
	printerLayout.AddWidget(printerNetworkInput, 0, 1, 0)
	printerLayout.AddWidget(printerAddressLabel, 1, 0, 0)
	printerLayout.AddWidget(printerAddressInput, 1, 1, 0)
	printerLayout.AddWidget(printerTimeoutLabel, 2, 0, 0)
	printerLayout.AddWidget(printerTimeoutInput, 2, 1, 0)

	printerGroup.SetLayout(printerLayout)

	// serial group

	serialGroup := widgets.NewQGroupBox2("Serial port", nil)

	serialDeviceLabel := widgets.NewQLabel2("Device name:", nil, 0)
	serialBaudLabel := widgets.NewQLabel2("Baud rate:", nil, 0)
	serialCharacterLabel := widgets.NewQLabel2("Character bits:", nil, 0)
	serialStopLabel := widgets.NewQLabel2("Stop bits:", nil, 0)
	serialParityLabel := widgets.NewQLabel2("Parity:", nil, 0)

	serialDeviceInput := widgets.NewQLineEdit(nil)
	serialDeviceInput.SetPlaceholderText("COM1, /dev/ttyUSB0...")
	serialDeviceInput.SetText("/dev/ttyUSB0")
	serialBaudInput := widgets.NewQLineEdit(nil)
	serialBaudInput.SetText("9600")
	serialCharacterInput := widgets.NewQLineEdit(nil)
	serialCharacterInput.SetText("8")
	serialStopInput := widgets.NewQLineEdit(nil)
	serialStopInput.SetText("1")
	serialParityInput := widgets.NewQLineEdit(nil)
	serialParityInput.SetText("n")
	serialParityInput.SetPlaceholderText("'n' means 'disable'")

	serialLayout := widgets.NewQGridLayout2()
	serialLayout.AddWidget(serialDeviceLabel, 0, 0, 0)
	serialLayout.AddWidget(serialDeviceInput, 0, 1, 0)
	serialLayout.AddWidget(serialBaudLabel, 1, 0, 0)
	serialLayout.AddWidget(serialBaudInput, 1, 1, 0)
	serialLayout.AddWidget(serialCharacterLabel, 2, 0, 0)
	serialLayout.AddWidget(serialCharacterInput, 2, 1, 0)
	serialLayout.AddWidget(serialStopLabel, 3, 0, 0)
	serialLayout.AddWidget(serialStopInput, 3, 1, 0)
	serialLayout.AddWidget(serialParityLabel, 4, 0, 0)
	serialLayout.AddWidget(serialParityInput, 4, 1, 0)

	serialGroup.SetLayout(serialLayout)

	// result group

	result := widgets.NewQTextEdit(nil)
	//result := widgets.NewQTextBrowser(nil)
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

	stack := interpreter.NewStack()

	runButton := widgets.NewQPushButton2("RUN", nil)
	runButton.ConnectClicked(func(bool) {

		if len(terminatecc) != 0 {
			return
		}
		suspButton.SetEnabled(true)

		err := initSerialDevice(
			serialDeviceInput.Text(),
			serialBaudInput.Text(),
			serialCharacterInput.Text(),
			serialStopInput.Text(),
			serialParityInput.Text(),
		)
		if err != nil {
			widgets.QMessageBox_Information(nil, "Error", err.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}

		err = initPrinter(
			printerNetworkInput.Text(),
			printerAddressInput.Text(),
			printerTimeoutInput.Text(),
		)
		if err != nil {
			widgets.QMessageBox_Information(nil, "Error", err.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}

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
	inputLayout.AddWidget(input, 0, 0, 0)
	inputLayout.AddWidget(waveformGroup, 1, 0, 0)
	inputLayout.AddWidget(printerGroup, 2, 0, 0)
	inputLayout.AddWidget(serialGroup, 3, 0, 0)
	inputLayout.AddWidget(runButton, 4, 0, 0)
	inputGroup.SetLayout(inputLayout)

	outputGroup := widgets.NewQGroupBox2("Results", nil)
	outputLayout := widgets.NewQGridLayout2()
	outputLayout.AddWidget3(result, 0, 0, 1, 2, 0)
	outputLayout.AddWidget3(termButton, 2, 0, 1, 2, 0)
	outputLayout.AddWidget(suspButton, 1, 0, 0)
	outputLayout.AddWidget(resuButton, 1, 1, 0)
	outputGroup.SetLayout(outputLayout)

	layout := widgets.NewQGridLayout2()
	layout.AddWidget(inputGroup, 0, 0, 0)
	layout.AddWidget(outputGroup, 0, 1, 0)
	layout.SetColumnStretch(0, 1)
	layout.SetColumnStretch(1, 1)
	widget.SetLayout(layout)

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
	parity string) (err error) {

	var deviceAddress byte
	deviceAddress = 0x01
	if alientek.Instance(string(deviceAddress)) != nil {
		return
	}

	baudInt, err := strconv.Atoi(baud)
	if err != nil {
		return
	}
	characterInt, err := strconv.Atoi(character)
	if err != nil {
		return
	}
	stopInt, err := strconv.Atoi(stop)
	if err != nil {
		return
	}
	//parity, err := strconv.Atoi(character)
	//if err != nil {
	//return
	//}

	alientek.NewDao(
		device,
		baudInt,
		characterInt,
		stopInt,
		-1,
		deviceAddress,
	)

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

func initPrinter(network string, address string, secondString string) (err error) {
	if ricoh_g5.Instance("") != nil {
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

	i := instruction.InstructionPrinterHeadPrinterStatus{}
	_, err = i.Execute()
	if err != nil {
		log.Println(err)
		return fmt.Errorf("printer is not ready")
	}
	return nil
}
