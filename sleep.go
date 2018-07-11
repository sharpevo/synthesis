package instruction

import (
	"fmt"
	"strconv"
	"time"
)

type InstructionSleep struct {
	Instruction
}

var Sleep InstructionSleep

func (i *InstructionSleep) Execute(args ...string) (interface{}, error) {
	seconds, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Duration(seconds) * time.Second)
	return fmt.Sprintf("sleep %d seconds", seconds), nil
}
