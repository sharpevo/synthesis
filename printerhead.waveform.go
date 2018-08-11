package instruction

import (
	"fmt"
	"posam/dao/ricoh_g5"
)

type InstructionPrinterHeadWaveform struct {
	Instruction
}

func (i *InstructionPrinterHeadWaveform) Execute(args ...string) (resp interface{}, err error) {
	if len(args) != 6 {
		return resp, fmt.Errorf("not enough arguments")
	}

	variable, _ := i.ParseVariable(args[0])
	headBoardIndex := args[1]
	rowIndexOfHeadBoard := args[2]
	voltagePercentage := args[3]
	segmentCount := args[4]
	segment := args[5]

	resp, err = ricoh_g5.Instance("").SendWaveform(
		headBoardIndex,
		rowIndexOfHeadBoard,
		voltagePercentage,
		segmentCount,
		segment,
	)
	variable.Value = resp
	if err != nil {
		return variable.Value, err
	}
	return
}
