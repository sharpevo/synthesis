package instruction

import (
	"fmt"
	"posam/dao/ricoh_g5"
)

func init() {
	ricoh_g5.InstructionMap.Set("ERRORCODE", InstructionPrinterHeadErrorCode{})
}

type InstructionPrinterHeadErrorCode struct {
	Instruction
}

func (i *InstructionPrinterHeadErrorCode) Execute(args ...string) (resp interface{}, err error) {
	if len(args) != 1 {
		return resp, fmt.Errorf("not enough arguments")
	}

	variable, _ := i.ParseVariable(args[0])
	resp, err = ricoh_g5.Instance("").QueryErrorCode()
	variable.Value = resp
	if err != nil {
		return resp, err
	}
	return
}
