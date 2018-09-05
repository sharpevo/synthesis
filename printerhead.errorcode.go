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
	if len(args) < 1 {
		return resp, fmt.Errorf("not enough arguments")
	}
	address := args[0]
	resp, err = ricoh_g5.Instance(address).QueryErrorCode()
	if err != nil {
		return resp, err
	}
	return
}
