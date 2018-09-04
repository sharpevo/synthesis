package devtree

import (
	"fmt"
	"path"
	"posam/gui/tree"
)

const (
	DEVICE_CONF_FILE = "devices.bin"
)

var ConfMap map[string]string
var ConnMap map[string][]string

type Node struct {
	tree.Node
	Enabled  bool
	Type     string
	Children []Node
}

func ParseDeviceConf() {
	ConfMap = make(map[string]string)
	ConnMap = make(map[string][]string)
	node := new(Node)
	err := tree.ImportNode(node, DEVICE_CONF_FILE)
	if err != nil {
		fmt.Println(err)
	}
	cpath := path.Join("/", node.Title)
	for _, v := range node.Children {
		err = parseDeviceConf(cpath, v)
	}
	//fmt.Println(ConfMap)
	//fmt.Println(ConnMap)
}

func parseDeviceConf(ppath string, node Node) error {
	cpath := path.Join(ppath, node.Title)
	ConfMap[cpath] = node.Data
	if isAvailableNode(node) {
		if node.Type == "" {
			node.Type = DEV_TYPE_UNK
		}
		ConnMap[node.Type] = append(ConnMap[node.Type], cpath)
	}
	for _, v := range node.Children {
		err := parseDeviceConf(cpath, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func isAvailableNode(node Node) bool {
	return node.Type == DEV_TYPE_UNK || node.Enabled
}

func ComposeVarName(args ...string) string {
	return path.Join(args...)
}

func ParseConnList() map[string]string {
	ParseDeviceConf()
	return GetConnMap()
}

func GetConnMap() map[string]string {
	m := make(map[string]string)
	m[DEV_TYPE_UNK] = DEV_TYPE_UNK
	for _, v := range []string{
		DEV_TYPE_ALT,
		DEV_TYPE_RCG,
		DEV_TYPE_CAN,
	} {
		for _, s := range ConnMap[v] {
			m[s] = v
		}
	}
	return m
}
