package instruction

import (
	"fmt"
	"log"
	"posam/interpreter"
	"posam/interpreter/vrb"
	"strconv"
)

type Instructioner interface {
	Execute(args ...string) (interface{}, error)
	// TODO: rb
}

type Instruction struct {
	Env   *interpreter.Stack
	title string
	path  string
}

func (i Instruction) Title() string {
	return i.title
}

func (i *Instruction) SetTitle(title string) {
	i.title = title
}

func (i *Instruction) Execute(args ...string) (interface{}, error) {
	return "", nil
}

func (i *Instruction) ParseVariable(name string) (*vrb.Variable, error) {
	variable, found := i.Env.Get(name)
	if !found {
		newVariable, err := vrb.NewVariable(name, name)
		if err != nil {
			return variable, err
		}
		variable = i.Env.Set(newVariable)
	}
	return variable, nil
}

func (i *Instruction) ParseIntVariable(name string) (*vrb.Variable, error) {
	variable, found := i.Env.Get(name)
	if !found {
		newVariable, err := vrb.NewVariable(name, "0")
		if err != nil {
			return variable, err
		}
		variable = i.Env.Set(newVariable)
	}
	if variable.Type != vrb.INT {
		return variable,
			fmt.Errorf("invalid type of int variable %s", vrb.INT)
	}
	return variable, nil
}

func (i *Instruction) ParseFloat64Variable(name string) (*vrb.Variable, error) {
	variable, found := i.Env.Get(name)
	if !found {
		newVariable, err := vrb.NewVariable(name, "0.0")
		if err != nil {
			return variable, err
		}
		variable = i.Env.Set(newVariable)
	}
	if variable.Type != vrb.FLOAT {
		return variable,
			fmt.Errorf("invalid type of float variable %s", vrb.FLOAT)
	}
	return variable, nil
}

func (i *Instruction) ParseInt(input string) (output int, err error) {
	outputVar, found := i.Env.Get(input)
	if !found {
		output, err = strconv.Atoi(input)
		if err != nil {
			return output, err
		}
	} else {
		output, err = strconv.Atoi(fmt.Sprintf("%v", outputVar.Value))
		if err != nil {
			return output,
				fmt.Errorf(
					"failed to parse variable %q to int: %s",
					input,
					err.Error(),
				)
		}
	}
	return output, nil
}

func (i *Instruction) ParseFloat(input string) (output float64, err error) {
	outputVar, found := i.Env.Get(input)
	if !found {
		output, err = strconv.ParseFloat(input, 64)
		if err != nil {
			return output, err
		}
	} else {
		output, err = strconv.ParseFloat(fmt.Sprintf("%v", outputVar.Value), 64)
		if err != nil {
			return output,
				fmt.Errorf(
					"failed to parse variable %q to float: %s",
					input,
					err.Error(),
				)
		}
	}
	return output, nil
}

func (i *Instruction) IssueError(message string) {
	varErr, found := i.Env.Get("SYS_ERR")
	if !found {
		log.Printf("invalid variable ERR")
	}
	if message != "" {
		varErr.Value = message
	} else {
		varErr.Value = ""
	}
}
