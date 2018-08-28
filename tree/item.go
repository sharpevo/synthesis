package tree

import (
	"encoding/gob"
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"os"
)

type Noder interface {
	Write(string) error
	Read(string) error
}

func DataRole() int {
	return int(core.Qt__UserRole) + 1
}

func GetTreeItemData(item *widgets.QTreeWidgetItem) string {
	return item.Data(0, DataRole()).ToString()
}

func SetTreeItemData(item *widgets.QTreeWidgetItem, data string) {
	item.SetData(
		0,
		DataRole(),
		core.NewQVariant17(data),
	)
}

type Node struct {
	Title    string
	Data     string
	Children []Node
}

func (n *Node) Write(filePath string) error {
	file, err := os.Create(filePath)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	encoder := gob.NewEncoder(file)
	encoder.Encode(n)
	return nil
}

func (n *Node) Read(filePath string) error {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return err
	}
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(n)
	if err != nil {
		return err
	}
	return nil
}
