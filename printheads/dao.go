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
	row, err := n.GetRowByIndex()
	n.Row = row
	return n, err
}

func (n *Nozzle) GetRowByIndex() (int, error) {
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

type Row struct {
	Index       int
	PositionX   int
	PositionY   int
	NozzleSpace int
	Nozzles     []*Nozzle
}

func NewRow(
	index int,
	posx int,
	posy int,
	nozzleSpace int,
) *Row {
	return &Row{
		Index:       index,
		PositionX:   posx,
		PositionY:   posy,
		NozzleSpace: nozzleSpace,
	}
}

func (r *Row) CalcNozzlePosition(index int) (int, int) {
	factor := index
	return r.PositionX + factor*r.NozzleSpace, r.PositionY
}

func (r *Row) AddNozzle(nozzle *Nozzle) {
	nozzle.PositionX, nozzle.PositionY = r.CalcNozzlePosition(len(r.Nozzles))
	r.Nozzles = append(r.Nozzles, nozzle)
}

type PrintHead struct {
	Index     int
	Rows      []*Row
	RowOffset int
	RowSpaceA int
	RowSpaceB int
}

func NewPrintHead(
	rowCount int,
	nozzleCount int,
	nozzleSpace int,
	rowOffset int,
	rowSpaceA int,
	rowSpaceB int,
	dposx int,
	dposy int,
) (*PrintHead, error) {
	h := &PrintHead{
		RowOffset: rowOffset,
		RowSpaceA: rowSpaceA,
		RowSpaceB: rowSpaceB,
	}
	h.Rows = []*Row{}
	for index := 0; index < rowCount; index++ {
		posx, posy := h.CalcRowPosition(index, dposx, dposy)
		row := NewRow(index, posx, posy, nozzleSpace)
		h.Rows = append(h.Rows, row)
	}
	for index := 0; index < nozzleCount; index++ {
		n, err := NewNozzle(index)
		if err != nil {
			return h, err
		}
		h.Rows[n.Row].AddNozzle(n)
	}
	return h, nil
}

func (h *PrintHead) CalcRowPosition(index int, dposx int, dposy int) (int, int) {
	switch index {
	case 0:
		return dposx - h.RowOffset, dposy - h.RowSpaceA - h.RowSpaceB
	case 1:
		return dposx, dposy - h.RowSpaceB
	case 2:
		return dposx - h.RowOffset, dposy - h.RowSpaceA
	case 3:
		return dposx, dposy
	}
	return 0, 0
}
