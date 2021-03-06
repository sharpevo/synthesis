package vrb

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
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

func (c ComparisonType) String() string {
	switch c {
	case GREATER:
		return ">"
	case LESS:
		return "<"
	case EQUAL:
		return "=="
	case UNEQUAL:
		return "!="
	default:
		return "unknown"
	}
}

type Variable struct {
	Name  string
	value interface{}
	Type  VariableType

	sync.Mutex
}

func (v *Variable) GetValue() interface{} {
	v.Lock()
	defer v.Unlock()
	return v.value
}

func (v *Variable) SetValue(value interface{}) {
	v.Lock()
	defer v.Unlock()
	v.value = value
}

var PreservedNames = map[string]bool{
	"SYS_CMP": true,
	"SYS_ERR": true,
	"SYS_CUR": true,
}

func NewPreservedVariables() []*Variable {
	variableList := []*Variable{}
	for k, _ := range PreservedNames {
		var variable *Variable
		switch k {
		case "SYS_CMP":
			variable, _ = newVariable(k, fmt.Sprintf("%d", UNKNOWN))
		case "SYS_ERR":
			variable, _ = newVariable(k, "")
		case "SYS_CUR":
			variable, _ = newVariable(k, "0")
		default:
			variable, _ = newVariable(k, "")
		}
		variableList = append(variableList, variable)
	}
	return variableList
}

// TODO: value out of range error processing
func NewVariable(name string, input string) (*Variable, error) {
	if PreservedNames[name] {
		return nil, fmt.Errorf("%q is reserved variable", name)
	}
	return newVariable(name, input)
}

func newVariable(name string, input string) (*Variable, error) {
	variable := &Variable{
		Name: name,
	}
	v, t, _ := ParseValue(input)
	variable.Type = t
	variable.SetValue(v)
	return variable, nil
}

const QUOTE = "\""

func ParseValue(input string) (interface{}, VariableType, error) {
	if strings.HasPrefix(input, QUOTE) && strings.HasSuffix(input, QUOTE) {
		trimed := strings.Trim(input, QUOTE)
		return trimed, STRING, nil
	}
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
			"mismatched type comparison is not allowed: %s(%s) and %s(%s)",
			var1.Name,
			var1.Type,
			var2.Name,
			var2.Type,
		)
	}
	switch var1.Type {
	case INT:
		value1 := var1.GetValue().(int64)
		value2 := var2.GetValue().(int64)
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
		value1 := var1.GetValue().(*big.Float)
		value2 := var2.GetValue().(*big.Float)
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
		if var1.GetValue() == var2.GetValue() {
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
