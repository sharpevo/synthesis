package instruction

import (
	"encoding/hex"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"posam/dao/alientek"
)

type InstructionSendSerial struct {
	Instruction
}

func (c *InstructionSendSerial) Execute(args ...string) (resp interface{}, err error) {
	if len(args) != 3 {
		return resp, fmt.Errorf("not enough arguments")
	}

	instruction := args[0]
	sentResp := args[1]
	doneResp := args[2]

	output, err := send(instruction, sentResp, doneResp)

	resp = output
	//resp = interpreter.Response{
	//Error:  err,
	//Output: output,
	//}
	return
}

func initSerialPort() (serialPort *serial.Port, err error) {
	dao := alientek.Instance(string(0x01))
	return dao.SerialPort.Instance(), nil
}

func send(
	instruction string,
	sentResp string,
	doneResp string) (resp string, err error) {

	serialp, err := initSerialPort()
	if err != nil {
		return resp, err
	}

	serialp.Flush()

	data, err := hex.DecodeString(instruction)
	if err != nil {
		return resp, err
	}

	if _, err = serialp.Write(data); err != nil {
		return resp, err
	}

	// sent

	log.Printf("%s: check sent response %q", instruction, sentResp)
	sentResult, err := collect(serialp, sentResp)
	if err != nil {
		serialp.Flush()
		return toHexString(sentResult), fmt.Errorf("failed to send instruction %q: %s", instruction, err)
	}

	// done

	log.Printf("%s: check complete response %q", instruction, doneResp)
	doneResult, err := collect(serialp, doneResp)
	if err != nil {
		return toHexString(doneResult), fmt.Errorf("failed to run instruction %q: %s", instruction, err)
	}

	return toHexString(doneResult), nil

}

func collect(serialp *serial.Port, resp string) (result []byte, err error) {
	max := len(resp) / 2
	buf := make([]byte, max)
	cnt := 0

	for {
		n, err := serialp.Read(buf)
		if err != nil {
			return result, err
		}
		cnt += n
		result = append(result, buf[:n]...)
		if cnt >= max || n == 0 {
			break
		}
	}

	//log.Printf("collect responses %q", toHexString(result))

	if resp != toHexString(result) {
		msg := fmt.Sprintf("invalid response cocde %s (%s)", resp, toHexString(result))
		log.Printf(msg)
		return result, fmt.Errorf(msg)
	}
	return

}

func toHexString(input []byte) (output string) {
	return fmt.Sprintf("%x", input)
}
