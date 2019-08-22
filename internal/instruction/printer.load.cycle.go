package instruction

import (
	"fmt"
	"synthesis/internal/dao"
)

func init() {
	dao.InstructionMap.Set("LOADCYCLE", InstructionPrinterLoadCycle{})
}

type InstructionPrinterLoadCycle struct {
	InstructionPrinterLoad
}

func (i *InstructionPrinterLoadCycle) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 2 {
		return resp, fmt.Errorf("not enough arguments")
	}
	cm, err := i.ParseIntVariable(args[0])
	if err != nil {
		return resp, err
	}
	bin, err := i.ParseFormations(args[1])
	if err != nil {
		return resp, err
	}

	// save bin to the variable with the type of string
	cm2, err := i.ParseVariable(args[1])
	if err != nil {
		return resp, err
	}
	isSameCM := cm == cm2

	cm.Lock()
	variable, _ := i.GetVarLockless(cm, args[0])
	variable.SetValue(int64(bin.CycleCount))

	if !isSameCM {
		cm2.Lock()
	}
	binVar, _ := i.GetVarLockless(cm, args[1])
	binVar.SetValue(bin)
	if !isSameCM {
		cm2.Unlock()
	}
	cm.Unlock()

	i.Env.Set(variable)

	return bin.CycleCount, err
}
