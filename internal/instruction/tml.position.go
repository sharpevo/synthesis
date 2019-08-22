package instruction

import (
	"math/big"
	"synthesis/internal/interpreter/vrb"
	"synthesis/pkg/concurrentmap"
)

type InstructionTMLPosition struct {
	InstructionTML
}

func (i *InstructionTMLPosition) ParseVariable(
	variableString string,
) (cm *concurrentmap.ConcurrentMap, err error) {
	cm, found := i.Env.Get(variableString)
	if !found {
		newVariable, err := vrb.NewVariable(variableString, "0.0")
		if err != nil {
			return cm, err
		}
		i.Env.Set(newVariable)
		cm, found = i.Env.Get(variableString)
	}
	return cm, nil
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
