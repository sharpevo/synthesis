package instruction

import (
	"fmt"
	"synthesis/dao/ricoh_g5"
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
	cm, found := i.Env.Get(args[0])
	if !found {
		return resp, fmt.Errorf("device %q is not defined", args[0])
	}
	cm.Lock()
	variable, _ := i.GetVarLockless(cm, args[0])
	address := variable.GetValue().(string)
	cm.Unlock()

	headBoardIndex := args[1]
	rowIndexOfHeadBoard := args[2]
	voltagePercentage := args[3]
	segmentCount := args[4]
	segmentArgumentList := args[5:]
	instance := ricoh_g5.Instance(address)
	if instance == nil {
		return resp, fmt.Errorf("device %q is not initialized", args[0])
	}
	resp, err = instance.SendWaveform(
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
