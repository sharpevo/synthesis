package instruction

import (
	"fmt"
	"posam/dao/ricoh_g5"
)

func init() {
	ricoh_g5.InstructionMap.Set("PRINTDATA", InstructionPrinterHeadPrintData{})
}

type InstructionPrinterHeadPrintData struct {
	Instruction
}

func (i *InstructionPrinterHeadPrintData) Execute(args ...string) (resp interface{}, err error) {
	if len(args) != 5 {
		return resp, fmt.Errorf("not enough arguments")
	}
	variable, found := i.Env.Get(args[0])
	if !found {
		return resp, fmt.Errorf("device %q is not defined", args[0])
	}
	address := variable.Value.(string)
	bitsPerPixel := args[1]
	width := args[2]
	lineBufferSize := args[3]
	lineBuffer := args[4]
	instance := ricoh_g5.Instance(address)
	if instance == nil {
		return resp, fmt.Errorf("device %q is not initialized", args[0])
	}
	resp, err = instance.PrintData(
		bitsPerPixel, width, lineBufferSize, lineBuffer,
	)
	if err != nil {
		return resp, err
	}
	return
}
