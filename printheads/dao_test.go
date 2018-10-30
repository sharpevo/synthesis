package printheads_test

import (
	"fmt"
	"posam/dao/printheads"
	"testing"
)

func TestOne(t *testing.T) {
	fmt.Println(printheads.SPACE_ROW_B)
	fmt.Println(printheads.SPACE_NOZZLE_PITCH)
}

func TestPosition(t *testing.T) {
	index := []int{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
		17, 18, 19, 20,
		1277, 1278, 1279, 1280,
	}
	row := []int{
		0, 2, 1, 3,
		0, 2, 1, 3,
		0, 2, 1, 3,
		0, 2, 1, 3,
		0, 2, 1, 3,
		0, 2, 1, 3,
	}

	for k, v := range index {
		position, err := printheads.Position(v)
		if err != nil || position != row[k] {
			t.Errorf(
				"\n%d: EXPECT: %v\nGET: %v\n",
				k,
				row[k],
				position,
			)
		}
	}
}
