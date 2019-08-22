package printhead

import (
	"fmt"
	"synthesis/internal/geometry"
	"synthesis/internal/reagent"
)

type Nozzle struct {
	Index     int
	Printhead *Printhead
	Pos       *geometry.Position
	Reagent   *reagent.Reagent

	RowIndex int
}

func NewNozzle(
	index int,
) (*Nozzle, error) {
	nozzle := &Nozzle{}
	nozzle.Index = index
	rowIndex, err := nozzle.CalRowByIndex()
	if err != nil {
		return nil, err
	}
	nozzle.RowIndex = rowIndex
	return nozzle, err
}

func (n *Nozzle) CalRowByIndex() (int, error) {
	mod := (n.Index + 1) % 4
	switch mod {
	case 0:
		return 3, nil
	case 1:
		return 0, nil
	case 2:
		return 1, nil
	case 3:
		return 2, nil
	default:
		return -1, fmt.Errorf("invalid index %v", n.Index)
	}
}

func (n *Nozzle) CalRowByIndex2() (int, error) {
	mod := (n.Index + 1) % 4
	switch mod {
	case 0:
		return 3, nil
	case 1:
		return 0, nil
	case 2:
		return 2, nil
	case 3:
		return 1, nil
	default:
		return -1, fmt.Errorf("invalid index %v", n.Index)
	}
}
