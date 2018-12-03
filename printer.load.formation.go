package instruction

import (
	"fmt"
	"posam/dao"
	"strconv"
)

func init() {
	dao.InstructionMap.Set("LOADGROUP", InstructionPrinterLoadFormation{})
}

type InstructionPrinterLoadFormation struct {
	InstructionPrinterLoad
}

func (i *InstructionPrinterLoadFormation) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 3 {
		return resp, fmt.Errorf("not enough arguments")
	}
	variable, err := i.ParseIntVariable(args[0])
	if err != nil {
		return resp, err
	}
	bin, err := i.ParseFormations(args[1])
	if err != nil {
		return 0, err
	}
	cycleIndex, err := strconv.Atoi(args[2])
	if err != nil {
		return 0, err
	}
	if cycleIndex > bin.CycleCount-1 || cycleIndex < 0 {
		return 0, fmt.Errorf(
			"invalid cycle index %v (%v)", cycleIndex, bin.CycleCount)
	}
	variable.Value = int64(len(bin.Formations[cycleIndex]))
	//fmt.Printf("%#v\n", bin.Formations[cycleIndex])
	return len(bin.Formations[cycleIndex]), err
}
