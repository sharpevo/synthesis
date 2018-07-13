package instruction

import (
	"fmt"
	"posam/dao/alientek"
)

type InstructionLed struct {
	Instruction
}

var Led InstructionLed

func (i *InstructionLed) Execute(args ...string) (resp interface{}, err error) {
	if len(args) != 1 {
		return resp, fmt.Errorf("not enough arguments")
	}

	alientekDao := *alientek.Instance(string(0x01))
	switch args[0] {
	case "on":
		output, err := alientekDao.TurnOnLed()
		resp = output
		return resp, err
	case "off":
		output, err := alientekDao.TurnOffLed()
		resp = output
		return resp, err
	}
	return
}
