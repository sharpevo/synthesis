package instruction

import (
	"fmt"
	"math/big"
	"posam/interpreter/vrb"
)

type InstructionTMLPosition struct {
	InstructionTML
}

func (i *InstructionTMLPosition) ParseVariable(
	variableString string,
) (variable *vrb.Variable, err error) {
	variable, found := i.Env.Get(variableString)
	if !found {
		newVariable, err := vrb.NewVariable(variableString, "0.0")
		if err != nil {
			return variable, err
		}
		variable = i.Env.Set(newVariable)
		fmt.Printf("%#v\n", variable)
	}
	return variable, nil
}

func (i *InstructionTMLPosition) PositionToBigFloat(
	positionFloat float64,
) (positionBigFloat *big.Float, err error) {
	positionBigFloat = big.NewFloat(positionFloat)
	if err != nil {
		return big.NewFloat(0), err
	}
	return positionBigFloat, nil
}
