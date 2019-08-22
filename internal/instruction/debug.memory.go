package instruction

import (
	"synthesis/internal/dao"
)

func init() {
	dao.InstructionMap.Set("MEMBOUND", InstructionMemBound{})
}

type InstructionMemBound struct {
	Instruction
}

func (i *InstructionMemBound) Execute(args ...string) (resp interface{}, err error) {
	var result [][]int
	for i := 0; i < 15; i++ {
		a := make([]int, 0, 9999999)
		result = append(result, a)
	}
	return resp, err
}
