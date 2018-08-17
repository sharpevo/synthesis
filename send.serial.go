package instruction

import (
	"encoding/hex"
	"fmt"
	"posam/dao/alientek"
)

type InstructionSendSerial struct {
	Instruction
}

func (c *InstructionSendSerial) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 2 {
		return resp, fmt.Errorf("not enough arguments")
	}

	instruction := args[0]
	doneResp := args[1]
	sentResp := ""
	if len(args) == 3 {
		sentResp = args[2]
	}

	output, err := send(instruction, sentResp, doneResp)

	resp = output
	return
}

func send(
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

	if _, err = alientek.Instance(string(0x01)).SerialClient.Send(data, sentBytes, doneBytes); err != nil {
		return resp, err
	}

	return toHexString(doneBytes), nil
}

func toHexString(input []byte) (output string) {
	return fmt.Sprintf("%x", input)
}
