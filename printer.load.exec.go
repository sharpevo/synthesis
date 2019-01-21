package instruction

import (
	"fmt"
	"posam/dao"
)

func init() {
	dao.InstructionMap.Set("LOADEXEC", InstructionPrinterLoadExec{})
}

type InstructionPrinterLoadExec struct {
	InstructionPrinterLoad
}

func (i *InstructionPrinterLoadExec) Execute(args ...string) (resp interface{}, err error) {
	if len(args) < 3 {
		return resp, fmt.Errorf("not enough arguments")
	}
	bin, err := i.ParseBin(args[0])
	if err != nil {
		return resp, err
	}
	cycleIndex, err := i.ParseIndex(args[1])
	if err != nil {
		return resp, err
	}
	if cycleIndex > bin.CycleCount-1 || cycleIndex < 0 {
		return resp, fmt.Errorf(
			"invalid cycle index %v (%v)", cycleIndex, bin.CycleCount)
	}
	formations := bin.Formations[cycleIndex]
	groupIndex, err := i.ParseIndex(args[2])
	if err != nil {
		return resp, err
	}
	if groupIndex > len(formations)-1 || groupIndex < 0 {
		return resp, fmt.Errorf(
			"invalid cycle index %v (%v)", cycleIndex, bin.CycleCount)
	}

	moveArgs := formations[groupIndex].Move
	fmt.Println(
		"move",
		moveArgs.MotorConf.Path,
		moveArgs.PositionX,
		moveArgs.PositionY,
		moveArgs.MotorConf.Speed,
		moveArgs.MotorConf.Accel,
	)
	moveIns := InstructionTMLMoveAbs{}
	resp, err = moveIns.Execute(
		moveArgs.MotorConf.Path,
		moveArgs.PositionX,
		moveArgs.PositionY,
		moveArgs.MotorConf.Speed,
		moveArgs.MotorConf.Accel,
	)
	if err != nil {
		return resp, err
	}
	lineBuffers := ""
	for _, printArgs := range formations[groupIndex].Prints {
		fmt.Println(
			"print",
			printArgs.PrintheadConf.Path,
			printArgs.PrintheadConf.BitsPerPixel,
			printArgs.PrintheadConf.Width,
			printArgs.PrintheadConf.BufferSize,
			printArgs.LineBuffer,
		)
		lineBuffers += " | " + printArgs.LineBuffer
		printIns := InstructionPrinterHeadPrintData{}
		resp, err = printIns.Execute(
			printArgs.PrintheadConf.Path,
			printArgs.PrintheadConf.BitsPerPixel,
			printArgs.PrintheadConf.Width,
			printArgs.PrintheadConf.BufferSize,
			printArgs.LineBuffer,
		)
		if err != nil {
			return resp, err
		}
	}
	return fmt.Sprintf(
		"cycle %v group %v: move (%v, %v), print %v",
		cycleIndex,
		groupIndex,
		moveArgs.PositionX,
		moveArgs.PositionY,
		lineBuffers,
	), err
}
