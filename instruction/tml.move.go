package instruction

import (
	"fmt"
	"synthesis/dao/aoztech"
)

type InstructionTMLMove struct {
	InstructionTML
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
	instance, err = i.ParseDevice(args[0])
	if err != nil {
		return instance, pos, speed, accel, err
	}
	pos, err = i.ParseFloat(args[1])
	if err != nil {
		return instance, pos, speed, accel, err
	}
	speed, err = i.ParseFloat(args[2])
	if err != nil {
		return instance, pos, speed, accel, err
	}
	accel, err = i.ParseFloat(args[3])
	if err != nil {
		return instance, pos, speed, accel, err
	}

	pos, err = i.ParseFloat(args[1])
	if err != nil {
		return instance, pos, speed, accel, err
	}
	speed, err = i.ParseFloat(args[2])
	if err != nil {
		return instance, pos, speed, accel, err
	}
	accel, err = i.ParseFloat(args[3])
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
	instance, err = i.ParseDevice(args[0])
	if err != nil {
		return instance, posx, posy, speed, accel, err
	}
	posx, err = i.ParseFloat(args[1])
	if err != nil {
		return instance, posx, posy, speed, accel, err
	}
	posy, err = i.ParseFloat(args[2])
	if err != nil {
		return instance, posx, posy, speed, accel, err
	}
	speed, err = i.ParseFloat(args[3])
	if err != nil {
		return instance, posx, posy, speed, accel, err
	}
	accel, err = i.ParseFloat(args[4])
	if err != nil {
		return instance, posx, posy, speed, accel, err
	}
	return
}
