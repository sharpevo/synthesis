package devtree

import (
	//"fmt"
	//"io/ioutil"
	//"os"
	//"path/filepath"
	"posam/gui/tree"
	//"strings"
	//"time"
	"path"
	"posam/interpreter"
	"posam/interpreter/vrb"
)

const (
	DEVICE_CONF_FILE = "devices.bin"
)

type Node struct {
	tree.Node
	Children []Node
}

func InitStack(stack *interpreter.Stack) error {
	node := new(Node)
	err := tree.ImportNode(node, DEVICE_CONF_FILE)
	if err != nil {
		return err
	}
	cpath := path.Join("/", node.Title)
	for _, v := range node.Children {
		err = setVar(cpath, v, stack)
	}
	return nil
}

func setVar(ppath string, node Node, stack *interpreter.Stack) error {
	cpath := path.Join(ppath, node.Title)
	variable, err := vrb.NewVariable(cpath, node.Data)
	if err != nil {
		return err
	}
	stack.Set(variable)
	for _, v := range node.Children {
		err = setVar(cpath, v, stack)
	}
	return nil
}
