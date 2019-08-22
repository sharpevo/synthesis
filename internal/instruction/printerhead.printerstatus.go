package instruction

import (
	"fmt"
	"synthesis/internal/dao/ricoh_g5"
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
	cm, found := i.Env.Get(args[0])
	if !found {
		return resp, fmt.Errorf("device %q is not defined", args[0])
	}
	cm.Lock()
	variable, _ := i.GetVarLockless(cm, args[0])
	address := variable.GetValue().(string)
	cm.Unlock()
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
