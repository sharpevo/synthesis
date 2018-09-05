package instruction

import (
	"fmt"
	"posam/dao/ricoh_g5"
)

func init() {
	ricoh_g5.InstructionMap.Set("PRINTERSTATUS", InstructionPrinterHeadPrinterStatus{})
}

type InstructionPrinterHeadPrinterStatus struct {
	Instruction
}

func (i *InstructionPrinterHeadPrinterStatus) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 1 {
		return resp, fmt.Errorf("not enough arguments")
	}
	address := args[0]
	resp, err = ricoh_g5.Instance(address).QueryPrinterStatus()
	if err != nil {
		return resp, err
	}
	return
}
