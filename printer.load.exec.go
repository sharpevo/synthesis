package instruction

import (
	"fmt"
	"posam/dao"
	"strconv"
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
	bin, err := i.ParseFormations(args[0])
	if err != nil {
		return resp, err
	}
	cycleIndex, err := strconv.Atoi(args[1])
	if err != nil {
		return resp, err
	}
	if cycleIndex > bin.CycleCount-1 || cycleIndex < 0 {
		return resp, fmt.Errorf(
			"invalid cycle index %v (%v)", cycleIndex, bin.CycleCount)
	}
	formations := bin.Formations[cycleIndex]
	groupIndex, err := strconv.Atoi(args[2])
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
	for _, printArgs := range formations[groupIndex].Prints {
		fmt.Println(
			"print",
			printArgs.PrintheadConf.Path,
			printArgs.PrintheadConf.BitsPerPixel,
			printArgs.PrintheadConf.Width,
			printArgs.PrintheadConf.BufferSize,
			printArgs.LineBuffer,
		)
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
		"group %v of cycle %v completed",
		groupIndex,
		cycleIndex,
	), err
}
