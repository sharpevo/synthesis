package instruction

import (
	"fmt"
	"log"
	"synthesis/internal/interpreter/vrb"
	"synthesis/pkg/concurrentmap"
	"strconv"
)

type Instructioner interface {
	Execute(args ...string) (interface{}, error)
	// TODO: rb
}

type Instruction struct {
	Env   *Stack
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

func (i *Instruction) ParseVariable(name string) (*concurrentmap.ConcurrentMap, error) {
	cm, found := i.Env.Get(name)
	if !found {
		newVariable, err := vrb.NewVariable(name, name)
		if err != nil {
			return cm, err
		}
		i.Env.Set(newVariable)
		cm, found = i.Env.Get(name)
	}
	return cm, nil
}

func (i *Instruction) ParseIntVariable(name string) (cm *concurrentmap.ConcurrentMap, err error) {
	cm, found := i.Env.Get(name)
	if !found {
		newVariable, err := vrb.NewVariable(name, "0")
		if err != nil {
			return cm, err
		}
		i.Env.Set(newVariable)
		cm, found = i.Env.Get(name)
	}
	cm.Lock()
	variable, _ := i.GetVarLockless(cm, name)
	if variable.Type != vrb.INT {
		err = fmt.Errorf("invalid type of int variable %s", vrb.INT)
	}
	cm.Unlock()
	return cm, err
}

func (i *Instruction) ParseFloat64Variable(name string) (cm *concurrentmap.ConcurrentMap, err error) {
	cm, found := i.Env.Get(name)
	if !found {
		newVariable, err := vrb.NewVariable(name, "0.0")
		if err != nil {
			return cm, err
		}
		i.Env.Set(newVariable)
		cm, found = i.Env.Get(name)
	}
	cm.Lock()
	variable, _ := i.GetVarLockless(cm, name)
	if variable.Type != vrb.FLOAT {
		err = fmt.Errorf("invalid type of float variable %s", vrb.FLOAT)
	}
	cm.Unlock()
	return cm, err
}

func (i *Instruction) ParseInt(input string) (output int, err error) {
	cm, found := i.Env.Get(input)
	if !found {
		output, err = strconv.Atoi(input)
		if err != nil {
			return output, err
		}
	} else {
		cm.Lock()
		defer cm.Unlock()
		outputVar, _ := i.GetVarLockless(cm, input)
		output, err = strconv.Atoi(fmt.Sprintf("%v", outputVar.GetValue()))
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
	cm, found := i.Env.Get(input)
	if !found {
		output, err = strconv.ParseFloat(input, 64)
		if err != nil {
			return output, err
		}
	} else {
		cm.Lock()
		defer cm.Unlock()
		outputVar, _ := i.GetVarLockless(cm, input)
		output, err = strconv.ParseFloat(fmt.Sprintf("%v", outputVar.GetValue()), 64)
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
	cm, found := i.Env.Get("SYS_ERR")
	if !found {
		log.Printf("invalid variable ERR")
	}
	cm.Lock()
	defer cm.Unlock()
	varErr, _ := i.GetVarLockless(cm, "SYS_ERR")
	if message != "" {
		varErr.SetValue(message)
	} else {
		varErr.SetValue("")
	}
}

func (i *Instruction) GetVarLockless(
	cm *concurrentmap.ConcurrentMap,
	key string,
) (
	variable *vrb.Variable,
	err error,
) {
	variablei, _ := cm.GetLockless(key)
	variable, _ = variablei.(*vrb.Variable)
	return variable, nil
}
