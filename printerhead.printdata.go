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

	variable, _ := i.ParseVariable(args[0])
	bitsPerPixel := args[1]
	width := args[2]
	lineBufferSize := args[3]
	lineBuffer := args[4]
	resp, err = ricoh_g5.Instance("").PrintData(
		bitsPerPixel, width, lineBufferSize, lineBuffer,
	)
	variable.Value = resp
	if err != nil {
		return variable.Value, err
	}
	return
}
