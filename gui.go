package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"posam/dao/alientek"
	"posam/protocol/serialport"
	"runtime"
	"strconv"
	"strings"

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
	CMD_LED_SERIAL = `LED on
SLEEP 3
SENDSERIAL 010200010001E80A 55 018202c161`
)

var InstructionMap = map[string]instruction.Instructioner{
	"PRINT":      &Print,
	"SLEEP":      &instruction.Sleep,
	"IMPORT":     &instruction.Import,
	"ASYNC":      &instruction.Async,
	"RETRY":      &instruction.Retry,
	"LED":        &instruction.Led,
	"SENDSERIAL": &instruction.SendSerial,
}

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
	input.SetPlainText(CMD_LED_SERIAL)

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

	terminatecc := make(chan chan interface{}, 1)
	defer close(terminatecc)

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
			return
		}

		result.SetText("RUNNING")

		terminatec := make(chan interface{})
		terminatecc <- terminatec

		interpreter.InitParser(InstructionMap)
		statementGroup := interpreter.StatementGroup{Execution: interpreter.SYNC}
		interpreter.ParseReader(
			strings.NewReader(input.ToPlainText()),
			&statementGroup,
		)
		resultList := []string{}

		go func() {
			for resp := range statementGroup.Execute(terminatec, &suspend, nil) {
				if resp.Error != nil {
					go suspendExecution(&suspend, suspButton, resuButton)
					msgBox.ShowMessageBoxSlot(resp.Error.Error())
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
		go suspendExecution(&suspend, suspButton, resuButton)
	})

	resuButton.ConnectClicked(func(bool) {
		go func() {
			suspend = false
			suspButton.SetEnabled(true)
			resuButton.SetEnabled(false)
		}()
	})

	inputGroup := widgets.NewQGroupBox2("Instructions", nil)
	inputLayout := widgets.NewQGridLayout2()
	inputLayout.AddWidget(input, 0, 0, 0)
	inputLayout.AddWidget(serialGroup, 1, 0, 0)
	inputLayout.AddWidget(runButton, 2, 0, 0)
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

	alientek.AddInstance(&alientek.Dao{
		DeviceAddress: deviceAddress,
		SerialPort: &serialport.SerialPort{
			Name:     device,
			BaudRate: baudInt,
			DataBits: characterInt,
			StopBits: stopInt,
			Parity:   -1,
		},
	})

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
