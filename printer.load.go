package instruction

import (
	"fmt"
	"posam/util/formation"
)

type InstructionPrinterLoad struct {
	Instruction
}

func (i *InstructionPrinterLoad) ParseFormations(filePath string) (*formation.Bin, error) {
	variable, found := i.Env.Get(filePath)
	if !found {
		return formation.ParseBin(filePath)
	}
	filePathString, ok := variable.Value.(string)
	if !ok {
		return nil, fmt.Errorf("invalid string variable", filePath)
	}
	return formation.ParseBin(filePathString)
}
