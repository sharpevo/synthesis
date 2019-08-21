package instruction

import (
	"fmt"
	"synthesis/dao"
	"time"
)

func init() {
	dao.InstructionMap.Set("SLEEP", InstructionSleep{})
}

type InstructionSleep struct {
	Instruction
}

func (i *InstructionSleep) Execute(args ...string) (interface{}, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("not enough arguments")
	}
	seconds, err := i.ParseFloat(args[0])
	if err != nil {
		return seconds, err
	}
	duration := time.Duration(seconds*1000) * time.Millisecond
	<-time.After(duration)
	return fmt.Sprintf("sleep in %v", duration), nil
}
