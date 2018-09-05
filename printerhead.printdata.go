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
	address := args[0]
	bitsPerPixel := args[1]
	width := args[2]
	lineBufferSize := args[3]
	lineBuffer := args[4]
	resp, err = ricoh_g5.Instance(address).PrintData(
		bitsPerPixel, width, lineBufferSize, lineBuffer,
	)
	if err != nil {
		return resp, err
	}
	return
}
