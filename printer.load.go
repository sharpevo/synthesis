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

func (i *InstructionPrinterLoad) ParseIndex(variableName string) (int, error) {
	indexVariable, err := i.ParseIntVariable(variableName)
	if err != nil {
		return 0, err
	}
	index64, _ := indexVariable.Value.(int64)
	return int(index64), nil
}
