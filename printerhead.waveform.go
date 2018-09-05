package instruction

import (
	"fmt"
	"posam/dao/ricoh_g5"
)

func init() {
	ricoh_g5.InstructionMap.Set("WAVEFORM", InstructionPrinterHeadWaveform{})
}

type InstructionPrinterHeadWaveform struct {
	Instruction
}

func (i *InstructionPrinterHeadWaveform) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 18 {
		return resp, fmt.Errorf("not enough arguments")
	}
	address := args[0]
	headBoardIndex := args[1]
	rowIndexOfHeadBoard := args[2]
	voltagePercentage := args[3]
	segmentCount := args[4]
	segmentArgumentList := args[5:]
	resp, err = ricoh_g5.Instance(address).SendWaveform(
		headBoardIndex,
		rowIndexOfHeadBoard,
		voltagePercentage,
		segmentCount,
		segmentArgumentList,
	)
	if err != nil {
		return resp, err
	}
	return
}
