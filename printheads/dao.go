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
	PositionX int
	PositionY int
}

type Row struct {
	Nozzles []Nozzle
}

type PrintHead struct {
	RowD Row
	RowC Row
	RowB Row
	RowA Row
}

func Position(index int) (int, error) {
	mod := index % 4
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
		return -1, fmt.Errorf("invalid index %v", index)
	}
}
