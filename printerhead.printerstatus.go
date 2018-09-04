package instruction

import (
	"posam/dao/ricoh_g5"
)

func init() {
	ricoh_g5.InstructionMap.Set("PRINTERSTATUS", InstructionPrinterHeadPrinterStatus{})
}

type InstructionPrinterHeadPrinterStatus struct {
	Instruction
}

func (i *InstructionPrinterHeadPrinterStatus) Execute(args ...string) (resp interface{}, err error) {
	if len(args) == 1 {
		variable, _ := i.ParseVariable(args[0])
		resp, err = ricoh_g5.Instance("").QueryPrinterStatus()
		variable.Value = resp
		if err != nil {
			return variable.Value, err
		}
	} else {
		resp, err = ricoh_g5.Instance("").QueryPrinterStatus()
		return
	}
	return
}
