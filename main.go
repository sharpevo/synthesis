package main

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
	"os"
	"os/exec"
	"posam/ui/config"
	"runtime"
	"strconv"
	"strings"

	"github.com/therecipe/qt/widgets"
	"posam/interpreter"
	cmd "posam/ui/command"
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
)

var CommandMap = map[string]interpreter.Commander{
	"PRINT":      &Print,
	"SLEEP":      &interpreter.Sleep,
	"IMPORT":     &interpreter.Import,
	"ASYNC":      &interpreter.Async,
	"RETRY":      &interpreter.Retry,
	"LED":        &cmd.Led,
	"SENDSERIAL": &cmd.SendSerial,
}

type CommandPrint struct {
	interpreter.Command
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
	input.SetPlainText(CMD_SERIAL)

	// serial group

	serialGroup := widgets.NewQGroupBox2("Serial port", nil)

	serialDeviceLabel := widgets.NewQLabel2("Device name:", nil, 0)
	serialBaudLabel := widgets.NewQLabel2("Baud rate:", nil, 0)
	serialCharacterLabel := widgets.NewQLabel2("Character bits:", nil, 0)
	serialStopLabel := widgets.NewQLabel2("Stop bits:", nil, 0)
	serialParityLabel := widgets.NewQLabel2("Parity:", nil, 0)

	serialDeviceInput := widgets.NewQLineEdit(nil)
	serialDeviceInput.SetPlaceholderText("COM1, /dev/ttyUSB0...")
	//serialDeviceInput.SetText("/dev/ttyUSB0")
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

		interpreter.InitParser(CommandMap)
		statementGroup := interpreter.StatementGroup{Execution: interpreter.SYNC}
		interpreter.ParseReader(
			strings.NewReader(input.ToPlainText()),
			&statementGroup,
		)
		resultList := []string{}

		go func() {
			for resp := range statementGroup.Execute(terminatec, &suspend, nil) {
				if resp.Error != nil {
					suspendExecution(&suspend, suspButton, resuButton)

					//widgets.QMessageBox_Information(nil, "Error", resp.Error.Error(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
					//widgets.QMessageBox_Information(nil, "Error", "sus", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)

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
		go suspendExecution(&suspend, suspButton, resuButton)
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

	if config.SerialPortInstance != nil {
		return
	}

	config.Config["serialport"] = config.SerialPort{
		Device:    device,
		Baud:      baud,
		Character: character,
		Stop:      stop,
		Parity:    parity,
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

	config.SerialPortInstance, err = openSerialPort()
	if err != nil {
		return err
	}

	return

}

func openSerialPort() (serialPort *serial.Port, err error) {
	log.Println("Opening serial port...")
	sp := config.Config["serialport"].(config.SerialPort)
	baud, err := strconv.Atoi(sp.Baud)
	if err != nil {
		return
	}

	c := &serial.Config{
		Name: sp.Device,
		Baud: baud,
	}

	serialPort, err = serial.OpenPort(c)
	if err != nil {
		return
	}
	return
}
