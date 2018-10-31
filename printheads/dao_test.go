package printheads_test

import (
	"fmt"
	"posam/dao/printheads"
	"reflect"
	"testing"
)

func TestOne(t *testing.T) {
	fmt.Println(printheads.SPACE_ROW_B)
	fmt.Println(printheads.SPACE_NOZZLE_PITCH)
}

func TestPosition(t *testing.T) {
	index := []int{
		0, 1, 2, 3,
		4, 5, 6, 7,
		8, 9, 10, 11,
		12, 13, 14, 15,
		16, 17, 18, 19,
		1276, 1277, 1278, 1279,
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
		n, err := printheads.NewNozzle(v)
		if err != nil || n.Row != row[k] {
			t.Errorf(
				"\n%d: EXPECT: %v\nGET: %v\n",
				k,
				row[k],
				n.Row,
			)
		}
	}
}

func TestPrintHeadLineD(t *testing.T) {
	h, _ := printheads.NewPrintHeadLineD(
		4,
		1280,
		169.3*printheads.UM,
		84.65*printheads.UM,
		550.3*printheads.UM,
		11.811*printheads.MM,
		0,
		0,
	)
	nozzleMap := map[int]printheads.Nozzle{
		1: printheads.Nozzle{
			Index:     0,
			Row:       0,
			PositionX: -84650,
			PositionY: -12361300,
		},
		3: printheads.Nozzle{
			Index:     2,
			Row:       1,
			PositionX: 0,
			PositionY: -11811000,
		},
		2: printheads.Nozzle{
			Index:     1,
			Row:       2,
			PositionX: -84650,
			PositionY: -550300,
		},
		4: printheads.Nozzle{
			Index:     3,
			Row:       3,
			PositionX: 0,
			PositionY: 0,
		},
		1277: printheads.Nozzle{
			Index:     1276,
			Row:       0,
			PositionX: 53922050,
			PositionY: -12361300,
		},
		1279: printheads.Nozzle{
			Index:     1278,
			Row:       1,
			PositionX: 54006700,
			PositionY: -11811000,
		},
		1278: printheads.Nozzle{
			Index:     1277,
			Row:       2,
			PositionX: 53922050,
			PositionY: -550300,
		},
		1280: printheads.Nozzle{
			Index:     1279,
			Row:       3,
			PositionX: 54006700,
			PositionY: 0,
		},
	}
	for _, row := range h.Rows {
		for _, nozzle := range row.Nozzles {
			//fmt.Printf("%#v\n", nozzle)
			if nozzle.Row != row.Index {
				t.Errorf(
					"\nEXPECT: %v\nGET: %v\n",
					row.Index,
					nozzle.Row,
				)
			}
			if n, ok := nozzleMap[nozzle.Index+1]; ok {
				if !reflect.DeepEqual(n, *nozzle) {
					t.Errorf(
						"\n%d EXPECT: %#v\nGET: %#v\n",
						nozzle.Index+1,
						n,
						*nozzle,
					)
				}
			}
		}
	}
}

func TestPrintHead(t *testing.T) {
	h, _ := printheads.NewPrintHead(
		4,
		1280,
		169.3*printheads.UM,
		84.65*printheads.UM,
		550.3*printheads.UM,
		11.811*printheads.MM,
	)
	nozzleMap := map[int]printheads.Nozzle{
		1: printheads.Nozzle{
			Index:     0,
			Row:       0,
			PositionX: 27088000,
			PositionY: -6180650,
		},
		3: printheads.Nozzle{
			Index:     2,
			Row:       1,
			PositionX: 27172650,
			PositionY: -5630350,
		},
		2: printheads.Nozzle{
			Index:     1,
			Row:       2,
			PositionX: 27088000,
			PositionY: 5630350,
		},
		4: printheads.Nozzle{
			Index:     3,
			Row:       3,
			PositionX: 27172650,
			PositionY: 6180650,
		},
		1277: printheads.Nozzle{
			Index:     1276,
			Row:       0,
			PositionX: 81094700,
			PositionY: -6180650,
		},
		1279: printheads.Nozzle{
			Index:     1278,
			Row:       1,
			PositionX: 81179350,
			PositionY: -5630350,
		},
		1278: printheads.Nozzle{
			Index:     1277,
			Row:       2,
			PositionX: 81094700,
			PositionY: 5630350,
		},
		1280: printheads.Nozzle{
			Index:     1279,
			Row:       3,
			PositionX: 81179350,
			PositionY: 6180650,
		},
	}
	for _, v := range nozzleMap {
		fmt.Printf("%s\n", v.String())
	}
	for _, row := range h.Rows {
		for _, nozzle := range row.Nozzles {
			//fmt.Printf("%#v\n", nozzle)
			if nozzle.Row != row.Index {
				t.Errorf(
					"\nEXPECT: %v\nGET: %v\n",
					row.Index,
					nozzle.Row,
				)
			}
			if n, ok := nozzleMap[nozzle.Index+1]; ok {
				if !reflect.DeepEqual(n, *nozzle) {
					t.Errorf(
						"\n%d EXPECT: %#v\nGET: %#v\n",
						nozzle.Index+1,
						n,
						*nozzle,
					)
				}
			}
		}
	}
}
