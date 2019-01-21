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
	cm, found := i.Env.Get(args[0])
	if !found {
		return resp, fmt.Errorf("device %q is not defined", args[0])
	}
	cm.Lock()
	variable, _ := i.GetVarLockless(cm, args[0])
	address := variable.GetValue().(string)
	cm.Unlock()
	instance := ricoh_g5.Instance(address)
	if instance == nil {
		return resp, fmt.Errorf("device %q is not initialized", args[0])
	}
	resp, err = instance.QueryErrorCode()
	if err != nil {
		return resp, err
	}
	return
}
