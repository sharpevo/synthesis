package instruction

import (
	"fmt"
	"posam/dao"
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
	variable, err := i.ParseIntVariable(args[0])
	if err != nil {
		return resp, err
	}
	bin, err := i.ParseFormations(args[1])
	if err != nil {
		return resp, err
	}
	variable.Value = int64(bin.CycleCount)
	return bin.CycleCount, err
}
