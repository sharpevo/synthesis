package instruction

import (
	"fmt"
	"posam/util/formation"
)

type InstructionPrinterLoad struct {
	Instruction
}

func (i *InstructionPrinterLoad) ParseFormations(filePath string) (*formation.Bin, error) {
	cm, found := i.Env.Get(filePath)
	if !found {
		return formation.ParseBin(filePath)
	}
	cm.Lock()
	defer cm.Unlock()
	variable, _ := i.GetVarLockless(cm, filePath)
	filePathString, ok := variable.GetValue().(string)
	if !ok {
		return nil, fmt.Errorf("invalid string variable", filePath)
	}
	return formation.ParseBin(filePathString)
}

func (i *InstructionPrinterLoad) ParseBin(name string) (*formation.Bin, error) {
	cm, found := i.Env.Get(name)
	if !found {
		return nil, fmt.Errorf("failed to parse bin %q", name)
	}
	cm.Lock()
	defer cm.Unlock()
	variable, _ := i.GetVarLockless(cm, name)
	bin, ok := variable.GetValue().(*formation.Bin)
	if !ok {
		return nil, fmt.Errorf("invalid bin variable", name)
	}
	fmt.Printf("Load bin file from variable %q", name)
	return bin, nil
}

func (i *InstructionPrinterLoad) ParseIndex(variableName string) (int, error) {
	cm, err := i.ParseIntVariable(variableName)
	if err != nil {
		return 0, err
	}
	cm.Lock()
	indexVariable, _ := i.GetVarLockless(cm, variableName)
	index64, _ := indexVariable.GetValue().(int64)
	cm.Unlock()
	return int(index64), nil
}
