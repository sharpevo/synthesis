package instree

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"synthesis/dao/alientek"
	"synthesis/dao/aoztech"
	"synthesis/dao/canalystii"
	"synthesis/dao/ricoh_g5"
	"synthesis/gui/tree"
	"time"
)

type Node struct {
	tree.Node
	DevicePath  string
	DeviceType  string
	Instruction string
	Arguments   string
	Children    []Node
}

func (n *Node) Generate() (string, error) {
	dir := filepath.Join(
		os.TempDir(),
		"igenetech",
		time.Now().Format("2006-01-02"),
	)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	file, err := ioutil.TempFile(dir, "")
	if err != nil {
		return "", err
	}
	filePath := file.Name()

	defer file.Close()
	if err != nil {
		return filePath, err
	}
	offset := 0
	for _, child := range n.Children {
		offset += 1
		nodeType, err := child.Type()
		if err != nil {
			return filePath, err
		}
		switch nodeType {
		case TYPE_INS:
			arguments := child.Arguments
			switch child.DeviceType {
			case ricoh_g5.NAME:
				arguments = fmt.Sprintf(
					"%s %s",
					path.Join(child.DevicePath, "CONN", ricoh_g5.IDNAME),
					arguments,
				)
			case alientek.NAME:
				arguments = fmt.Sprintf(
					"%s %s",
					path.Join(child.DevicePath, "CONN", alientek.IDNAME),
					arguments,
				)
			case aoztech.NAME:
				arguments = fmt.Sprintf(
					"%s %s",
					path.Join(child.DevicePath, "CONN", aoztech.IDNAME),
					arguments,
				)
			case canalystii.NAME:
				arguments = fmt.Sprintf(
					"%s %s",
					path.Join(child.DevicePath, "CONN", canalystii.IDNAME),
					arguments,
				)
			}
			fmt.Println(child.DeviceType, arguments)
			file.WriteString(fmt.Sprintf("%s %s\n", child.Instruction, arguments))
			break
		case TYPE_SET_ONCE:
			setPath, err := child.Generate()
			if err != nil {
				return filePath, err
			}
			file.WriteString(fmt.Sprintf(
				"%s %s\n", child.Instruction, setPath))
			break
		case TYPE_SET_LOOP:
			setPath, err := child.Generate()
			if err != nil {
				return filePath, err
			}
			file.WriteString(
				fmt.Sprintf("%s %s\n", child.Instruction, setPath))
			file.WriteString(
				fmt.Sprintf("LOOP %d %s\n", offset, child.ArgumentList()[0]))
			offset += 1
			break
		case TYPE_SET_COND:
			setPath, err := child.Generate()
			if err != nil {
				return filePath, err
			}
			var1 := child.ArgumentList()[0]
			opsb := child.ArgumentList()[1]
			var2 := child.ArgumentList()[2]
			var opst string
			switch opsb {
			case ">":
				opst = "GTGOTO"
				break
			case "<":
				opst = "LTGOTO"
				break
			case "!=":
				opst = "NEGOTO"
				break
			case "==":
				opst = "EQGOTO"
				break
			default:
				return filePath, fmt.Errorf(
					"invalid operator in %q",
					child.Title,
				)
			}

			file.WriteString(
				fmt.Sprintf("CMPVAR %s %s\n", var1, var2))
			file.WriteString(
				fmt.Sprintf("%s %d\n", opst, offset+3))
			file.WriteString(
				fmt.Sprintf("GOTO %d\n", offset+4))
			file.WriteString(
				fmt.Sprintf("%s %s\n", child.Instruction, setPath))
			offset += 3
			break
		}
	}
	file.Sync()
	return filePath, nil
}

func shouldBeInstructionSet(instruction string) bool {
	return instruction == INST_SET_SYNC ||
		instruction == INST_SET_ASYN
}

func (n *Node) ArgumentList() (result []string) {
	if n.Arguments == "" {
		return result
	}
	result = strings.Split(n.Arguments, " ")
	return
}

func (n *Node) Type() (string, error) {

	if shouldBeInstructionSet(n.Instruction) {
		if len(n.Children) == 0 {
			return "", fmt.Errorf(
				"instruction set %q has no instructions",
				n.Title,
			)
		}
		switch len(n.ArgumentList()) {
		case 0:
			return TYPE_SET_ONCE, nil
		case 1:
			return TYPE_SET_LOOP, nil
		case 3:
			return TYPE_SET_COND, nil
		default:
			return "", fmt.Errorf(
				"instruction %q is not valid instruction set",
				n.Title,
			)
		}
	} else {
		if len(n.Children) > 0 {
			return "", fmt.Errorf(
				"instruction %q should be instruction set",
				n.Title,
			)
		}
		return TYPE_INS, nil
	}
}
