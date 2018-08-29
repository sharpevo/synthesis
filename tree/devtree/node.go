package devtree

import (
	"fmt"
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

var ConnSRLVarNameList []string
var ConnTCPVarNameList []string
var ConnCANVarNameList []string

type Node struct {
	tree.Node
	Enabled  bool
	Type     string
	Children []Node
}

func InitStack(stack *interpreter.Stack) error {
	node := new(Node)
	ConnSRLVarNameList = []string{}
	ConnTCPVarNameList = []string{}
	ConnCANVarNameList = []string{}
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
	if isConnNode(node) {
		switch node.Type {
		case DEV_TYPE_TCP:
			ConnTCPVarNameList = append(ConnTCPVarNameList, cpath)
			break
		case DEV_TYPE_SRL:
			ConnSRLVarNameList = append(ConnSRLVarNameList, cpath)
			break
		case DEV_TYPE_CAN:
			ConnCANVarNameList = append(ConnCANVarNameList, cpath)
			break
		default:
			return fmt.Errorf("invalid node type %q", node.Type)
		}
	}
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

func isConnNode(node Node) bool {
	return node.Enabled && node.Type != DEV_TYPE_UNK
}

func ComposeVarName(args ...string) string {
	return path.Join(args...)
}
