package instree

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"posam/gui/tree"
	"strings"
	"time"
)

type Node struct {
	tree.Node
	Children []Node
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
			file.WriteString(fmt.Sprintf("%s\n", child.Data))
			break
		case TYPE_SET_ONCE:
			setPath, err := child.Generate()
			if err != nil {
				return filePath, err
			}
			file.WriteString(fmt.Sprintf(
				"%s %s\n", child.Instruction(), setPath))
			break
		case TYPE_SET_LOOP:
			setPath, err := child.Generate()
			if err != nil {
				return filePath, err
			}
			file.WriteString(
				fmt.Sprintf("%s %s\n", child.Instruction(), setPath))
			file.WriteString(
				fmt.Sprintf("LOOP %d %s\n", offset, child.Arguments()[0]))
			offset += 1
			break
		case TYPE_SET_COND:
			setPath, err := child.Generate()
			if err != nil {
				return filePath, err
			}
			var1 := child.Arguments()[0]
			opsb := child.Arguments()[1]
			var2 := child.Arguments()[2]
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
					n.Title,
				)
			}

			file.WriteString(
				fmt.Sprintf("CMPVAR %s %s\n", var1, var2))
			file.WriteString(
				fmt.Sprintf("%s %d\n", opst, offset+3))
			file.WriteString(
				fmt.Sprintf("GOTO %d\n", offset+4))
			file.WriteString(
				fmt.Sprintf("%s %s\n", child.Instruction(), setPath))
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

func (n *Node) Instruction() string {
	return strings.Split(n.Data, " ")[0]
}

func (n *Node) Arguments() []string {
	return strings.Split(n.Data, " ")[1:]
}

func (n *Node) Type() (string, error) {
	dataList := strings.Split(strings.Trim(n.Data, "\" "), " ")
	instruction := dataList[0]
	argumentList := dataList[1:]

	if shouldBeInstructionSet(instruction) {
		if len(n.Children) == 0 {
			return "", fmt.Errorf(
				"instruction set %q has no instructions",
				n.Title,
			)
		}
		switch len(argumentList) {
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
