package instruction

import (
	"fmt"
	"posam/dao/aoztech"
	"strconv"
)

type InstructionTMLMove struct {
	Instruction
}

func (i *InstructionTMLMove) Initialize(args ...string) (
	instance *aoztech.Dao,
	pos float64,
	speed float64,
	accel float64,
	err error,
) {
	if len(args) < 4 {
		return instance, pos, speed, accel,
			fmt.Errorf("not enough arguments")
	}
	variable, found := i.Env.Get(args[0])
	if !found {
		return instance, pos, speed, accel,
			fmt.Errorf("device %q is not defined", args[0])
	}
	deviceName := variable.Value.(string)
	instance = aoztech.Instance(deviceName)
	if instance == nil {
		return instance, pos, speed, accel,
			fmt.Errorf("device %q is not initialized", args[0])
	}
	pos, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return instance, pos, speed, accel, err
	}
	speed, err = strconv.ParseFloat(args[2], 64)
	if err != nil {
		return instance, pos, speed, accel, err
	}
	accel, err = strconv.ParseFloat(args[3], 64)
	if err != nil {
		return instance, pos, speed, accel, err
	}
	return
}

func (i *InstructionTMLMove) InitializeDual(args ...string) (
	instance *aoztech.Dao,
	posx float64,
	posy float64,
	speed float64,
	accel float64,
	err error,
) {
	if len(args) < 5 {
		return instance, posx, posy, speed, accel,
			fmt.Errorf("not enough arguments")
	}
	variable, found := i.Env.Get(args[0])
	if !found {
		return instance, posx, posy, speed, accel,
			fmt.Errorf("device %q is not defined", args[0])
	}
	deviceName := variable.Value.(string)
	instance = aoztech.Instance(deviceName)
	if instance == nil {
		return instance, posx, posy, speed, accel,
			fmt.Errorf("device %q is not initialized", args[0])
	}
	posx, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return instance, posx, posy, speed, accel, err
	}
	posy, err = strconv.ParseFloat(args[2], 64)
	if err != nil {
		return instance, posx, posy, speed, accel, err
	}
	speed, err = strconv.ParseFloat(args[3], 64)
	if err != nil {
		return instance, posx, posy, speed, accel, err
	}
	accel, err = strconv.ParseFloat(args[4], 64)
	if err != nil {
		return instance, posx, posy, speed, accel, err
	}
	return
}
