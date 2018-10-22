package instruction

import (
	"fmt"
	"posam/dao"
	"strconv"
	"time"
)

func init() {
	dao.InstructionMap.Set("SLEEP", InstructionSleep{})
}

type InstructionSleep struct {
	Instruction
}

func (i *InstructionSleep) Execute(args ...string) (interface{}, error) {
	seconds, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return nil, err
	}
	duration := time.Duration(seconds*1000) * time.Millisecond
	time.Sleep(duration)
	return fmt.Sprintf("sleep in %v", duration), nil
}
