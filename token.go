package commandparser

import ()

type Command struct {
	Name     string
	Function FunctionType
}

type FunctionType func(...string) (string, error)

var CommandMap = map[string]FunctionType{
	"TEST": CmdTest,
}

func CmdTest(args ...string) (string, error) {
	return args[0] + "_test", nil
}
