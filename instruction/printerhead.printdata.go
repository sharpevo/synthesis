package instruction

import (
	"fmt"
	"synthesis/dao/ricoh_g5"
)

func init() {
	ricoh_g5.InstructionMap.Set("PRINTDATA", InstructionPrinterHeadPrintData{})
}

type InstructionPrinterHeadPrintData struct {
	Instruction
}

func (i *InstructionPrinterHeadPrintData) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 6 {
		return resp, fmt.Errorf("not enough arguments")
	}
	cm, found := i.Env.Get(args[0])
	if !found {
		return resp, fmt.Errorf("device %q is not defined", args[0])
	}
	cm.Lock()
	variable, _ := i.GetVarLockless(cm, args[0])
	address := variable.GetValue().(string)
	cm.Unlock()
	bitsPerPixel := args[1]
	width := args[2]
	lineBufferSize := args[3]
	lineBuffer0 := args[4]
	lineBuffer1 := args[5]
	instance := ricoh_g5.Instance(address)
	if instance == nil {
		return resp, fmt.Errorf("device %q is not initialized", args[0])
	}
	resp, err = instance.PrintData(
		bitsPerPixel, width, lineBufferSize, lineBuffer0, lineBuffer1,
	)
	if err != nil {
		return resp, err
	}
	return
}
