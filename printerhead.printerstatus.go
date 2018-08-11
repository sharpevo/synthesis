package instruction

import (
	"fmt"
	"posam/dao/ricoh_g5"
)

type InstructionPrinterHeadPrinterStatus struct {
	Instruction
}

func (i *InstructionPrinterHeadPrinterStatus) Execute(args ...string) (resp interface{}, err error) {
	if len(args) != 1 {
		return resp, fmt.Errorf("not enough arguments")
	}

	variable, _ := i.ParseVariable(args[0])
	resp, err = ricoh_g5.Instance("").QueryPrinterStatus()
	variable.Value = resp
	if err != nil {
		return variable.Value, err
	}
	return
}
