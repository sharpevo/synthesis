package printheads

import (
	"fmt"
	"math"
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

type Position struct {
	X int
	Y int
}

type Nozzle struct {
	Position
	Index int
	Row   int
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

func (n *Nozzle) String() string {
	return fmt.Sprintf("%d at line %d: %v, %v",
		n.Index+1,
		n.Row,
		float64(n.X)/UM,
		float64(n.Y)/UM,
	)
}

func (n *Nozzle) IsAvailable(posx int, posy int, tolerance int) bool {
	return math.Abs(float64(n.X-posx)) < float64(tolerance) &&
		math.Abs(float64(n.Y-posy)) < float64(tolerance)
}

type Row struct {
	Position
	Index       int
	Reagent     string
	NozzleSpace int
	Nozzles     []*Nozzle
}

func NewRow(
	index int,
	posx int,
	posy int,
	nozzleSpace int,
	reagents []string,
) *Row {
	return &Row{
		Index: index,
		Position: Position{
			X: posx,
			Y: posy,
		},
		Reagent:     reagents[index],
		NozzleSpace: nozzleSpace,
	}
}

func (r *Row) CalcNozzlePosition(index int) (int, int) {
	factor := index
	return r.X + factor*r.NozzleSpace, r.Y
}

func (r *Row) AddNozzle(nozzle *Nozzle) {
	nozzle.X, nozzle.Y = r.CalcNozzlePosition(len(r.Nozzles))
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
) (*PrintHead, error) {
	dposy := (rowSpaceB + rowSpaceA) / 2
	dposx := nozzleSpace*160 + rowOffset
	fmt.Println(dposx, dposy)
	return NewPrintHeadLineD(
		rowCount,
		nozzleCount,
		nozzleSpace,
		rowOffset,
		rowSpaceA,
		rowSpaceB,
		dposx,
		dposy,
	)
}

func NewPrintHeadLineD(
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
	reagents := []string{"A", "C", "G", "T"}
	for index := 0; index < rowCount; index++ {
		posx, posy := h.CalcRowPosition(index, dposx, dposy)
		row := NewRow(index, posx, posy, nozzleSpace, reagents)
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

func (h *PrintHead) UpdatePosition(dposx int, dposy int) {
	for index, row := range h.Rows {
		row.X, row.Y = h.CalcRowPosition(index, dposx, dposy)
		for index, nozzle := range row.Nozzles {
			nozzle.X, nozzle.Y = row.CalcNozzlePosition(index)
		}
	}
}

func (h *PrintHead) UpdatePositionLeftBottom(aposx int, aposy int) {
	dposx := aposx + h.RowOffset
	dposy := aposy + h.RowSpaceA + h.RowSpaceB
	h.UpdatePosition(dposx, dposy)
}

func (h *PrintHead) UpdatePositionStar(x int, y int) {
	dposx := x
	dposy := y + h.RowSpaceA + h.RowSpaceB
	h.UpdatePosition(dposx, dposy)
}
