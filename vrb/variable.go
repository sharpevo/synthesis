package vrb

import (
	"fmt"
	"math/big"
	"strconv"
)

type VariableType int

const (
	STRING VariableType = iota
	INT
	FLOAT
	BYTEARRAY
)

func (v VariableType) String() string {
	switch v {
	case STRING:
		return "string"
	case INT:
		return "int"
	case FLOAT:
		return "float"
	case BYTEARRAY:
		return "bytes"
	default:
		return "unknown"
	}
}

type ComparisonType int

const (
	UNKNOWN ComparisonType = iota + 10
	GREATER                // 11
	LESS                   // 12
	EQUAL                  // 13
	UNEQUAL                // 14
)

type Variable struct {
	Name  string
	Value interface{}
	Type  VariableType
}

// TODO: value out of range error processing
func NewVariable(name string, input string) (*Variable, error) {
	variable := &Variable{
		Name: name,
	}
	variable.Value, variable.Type, _ = ParseValue(input)
	return variable, nil
}

func ParseValue(input string) (interface{}, VariableType, error) {
	if output, err := strconv.ParseInt(input, 0, 64); err == nil {
		return output, INT, nil
	}
	if output, _, err := big.ParseFloat(
		input,
		10,
		53,
		big.ToNearestEven,
	); err == nil {
		return output, FLOAT, nil
	}
	return input, STRING, nil
}

func Compare(var1 *Variable, var2 *Variable) (ComparisonType, error) {
	if var1.Type != var2.Type {
		return UNKNOWN, fmt.Errorf(
			"mismatched type comparison is not allowed: %s(%T) and %s(%T)",
			var1.Name,
			var1.Type,
			var2.Name,
			var2.Type,
		)
	}
	switch var1.Type {
	case INT:
		value1 := var1.Value.(int64)
		value2 := var2.Value.(int64)
		if value1 == value2 {
			return EQUAL, nil
		}
		if value1 > value2 {
			return GREATER, nil
		}
		if value1 < value2 {
			return LESS, nil
		}
		return UNKNOWN, fmt.Errorf(
			"failed to compare int variables: %q and %q",
			var1.Name,
			var2.Name,
		)
	case FLOAT:
		value1 := var1.Value.(*big.Float)
		value2 := var2.Value.(*big.Float)
		result := value1.Cmp(value2)
		if result == 0 {
			return EQUAL, nil
		}
		if result > 0 {
			return GREATER, nil
		}
		if result < 0 {
			return LESS, nil
		}
		return UNKNOWN, fmt.Errorf(
			"failed to compare float variables: %q and %q",
			var1.Name,
			var2.Name,
		)
	case STRING:
		if var1.Value == var2.Value {
			return EQUAL, nil
		} else {
			return UNEQUAL, nil
		}
		return UNKNOWN, fmt.Errorf(
			"failed to compare string variables: %q and %q",
			var1.Name,
			var2.Name,
		)
	}
	return UNKNOWN, fmt.Errorf(
		"failed to compare variables: %q and %q",
		var1.Name,
		var2.Name,
	)
}
