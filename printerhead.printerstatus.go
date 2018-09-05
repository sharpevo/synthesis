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
	variable, found := i.Env.Get(args[0])
	if !found {
		return resp, fmt.Errorf("device %q is not defined", args[0])
	}
	address := variable.Value.(string)
	fmt.Printf("%q\n", address)
	instance := ricoh_g5.Instance(address)
	if instance == nil {
		return resp, fmt.Errorf("device %q is not initialized", args[0])
	}
	resp, err = instance.QueryPrinterStatus()
	if err != nil {
		return resp, err
	}
	return
}
