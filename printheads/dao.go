package printheads

import (
	"fmt"
)

const (
	// units
	NM = 1
	UM = 1e3
	MM = 1e6
	// printhead
	SPACE_ROW_A        = 0.55 * MM
	SPACE_ROW_B        = 11.81 * MM
	WIDTH_ROW          = 54.1 * MM
	SPACE_NOZZLE_PITCH = 169.3 * UM
	ROW_OFFSET_X       = 84.65 * UM
	ROW_OFFSET_Y_A     = 0.55 * MM
	ROW_OFFSET_Y_B     = 11.81 * MM
	// axis
	RESOLUTION_X = ROW_OFFSET_X
	RESOLUTION_Y = 0.8 * UM
	// the first nozzle position of row for given row d
	POS_FIRST_NOZZLE_ROW_D_X = 0
	POS_FIRST_NOZZLE_ROW_D_Y = 0
	POS_FIRST_NOZZLE_ROW_C_X = POS_FIRST_NOZZLE_ROW_D_X - ROW_OFFSET_X
	POS_FIRST_NOZZLE_ROW_C_Y = POS_FIRST_NOZZLE_ROW_D_Y - ROW_OFFSET_Y_A
	POS_FIRST_NOZZLE_ROW_B_X = POS_FIRST_NOZZLE_ROW_D_X
	POS_FIRST_NOZZLE_ROW_B_Y = POS_FIRST_NOZZLE_ROW_D_Y - ROW_OFFSET_Y_B
	POS_FIRST_NOZZLE_ROW_A_X = POS_FIRST_NOZZLE_ROW_D_X - ROW_OFFSET_X
	POS_FIRST_NOZZLE_ROW_A_Y = POS_FIRST_NOZZLE_ROW_D_Y - ROW_OFFSET_Y_A - ROW_OFFSET_Y_B
)

type Nozzle struct {
	Index     int
	Row       int
	PositionX int
	PositionY int
}

func NewNozzle(index int) (*Nozzle, error) {
	n := &Nozzle{}
	n.Index = index
	row, err := n.ArrangeRow()
	n.Row = row
	return n, err
}

func (n *Nozzle) ArrangeRow() (int, error) {
	mod := n.Index % 4
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

type Row struct {
	Index   int
	Nozzles []*Nozzle
}

type PrintHead struct {
	Rows []*Row
}

func NewPrintHead(row int, nozzle int) (*PrintHead, error) {
	h := &PrintHead{}
	h.Rows = []*Row{}
	for index := 0; index < row; index++ {
		h.Rows = append(h.Rows, &Row{Index: index})
	}
	for index := 0; index < nozzle; index++ {
		n, err := NewNozzle(index)
		if err != nil {
			return h, err
		}
		h.Rows[n.Row].Nozzles = append(h.Rows[n.Row].Nozzles, n)
	}
	return h, nil
}
