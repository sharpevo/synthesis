package instruction

import (
	"encoding/hex"
	"fmt"
	"synthesis/internal/dao/alientek"
)

func init() {
	alientek.InstructionMap.Set("SENDSERIAL", InstructionSendSerial{})
}

type InstructionSendSerial struct {
	Instruction
}

func (i *InstructionSendSerial) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 3 {
		return resp, fmt.Errorf("not enough arguments")
	}
	cm, found := i.Env.Get(args[0])
	if !found {
		resp = fmt.Sprintf("device %q is not defined", args[0])
		return
	}
	cm.Lock()
	variable, _ := i.GetVarLockless(cm, args[0])
	deviceCode := variable.GetValue().(string)
	cm.Unlock()
	instruction := args[1]
	doneResp := args[2]
	sentResp := ""
	if len(args) == 4 {
		sentResp = args[3]
	}
	output, err := send(deviceCode, instruction, sentResp, doneResp)
	resp = output
	return
}

func send(
	name string,
	instruction string,
	sentResp string,
	doneResp string) (resp string, err error) {

	data, err := hex.DecodeString(instruction)
	if err != nil {
		return resp, err
	}

	sentBytes := []byte{}
	if sentResp != "" {
		sentBytes, err = hex.DecodeString(sentResp)
		if err != nil {
			return resp, err
		}
	}

	doneBytes, err := hex.DecodeString(doneResp)
	if err != nil {
		return resp, err
	}

	devInstance := alientek.Instance(name)
	if devInstance == nil {
		return resp, fmt.Errorf("invalid device %q", name)
	}

	if _, err = devInstance.SerialClient.Send(data, sentBytes, doneBytes); err != nil {
		return resp, err
	}

	return toHexString(doneBytes), nil
}

func toHexString(input []byte) (output string) {
	return fmt.Sprintf("%x", input)
}
