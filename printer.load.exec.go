package instruction

import (
	"fmt"
	"posam/dao"
	//"posam/interpreter"
)

//var DEBUG = true

var DEBUG = false

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
	formationCount := len(formations)
	if formationCount == 0 {
		// happens when activator is checked but no activator reagent assigned
		return fmt.Sprintf("no group in the cycle %v", cycleIndex), nil
	}
	groupIndex, err := i.ParseIndex(args[2])
	if err != nil {
		return resp, err
	}
	if groupIndex > formationCount-1 || groupIndex < 0 {
		return resp, fmt.Errorf(
			"invalid group index %v (%v)", groupIndex, formationCount)
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
	if !DEBUG {
		moveIns := InstructionTMLMoveAbs{}
		//moveIns.Env = interpreter.NewStack(i.Env)
		moveIns.Env = NewStack(i.Env)
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
	}

	fmt.Println(
		"print",
		formations[groupIndex].Print.PrintheadConf.Path,
		formations[groupIndex].Print.PrintheadConf.BitsPerPixel,
		formations[groupIndex].Print.PrintheadConf.Width,
		formations[groupIndex].Print.PrintheadConf.BufferSize,
		formations[groupIndex].Print.LineBuffers,
	)
	if !DEBUG {
		printIns := InstructionPrinterHeadPrintData{}
		//printIns.Env = interpreter.NewStack(i.Env)
		printIns.Env = NewStack(i.Env)
		resp, err = printIns.Execute(
			formations[groupIndex].Print.PrintheadConf.Path,
			formations[groupIndex].Print.PrintheadConf.BitsPerPixel,
			formations[groupIndex].Print.PrintheadConf.Width,
			formations[groupIndex].Print.PrintheadConf.BufferSize,
			formations[groupIndex].Print.LineBuffers[0],
			formations[groupIndex].Print.LineBuffers[1],
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
		formations[groupIndex].Print.LineBuffers,
	), err
}
