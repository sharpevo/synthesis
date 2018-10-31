package printheads_test

import (
	"fmt"
	"posam/dao/printheads"
	"reflect"
	"testing"
)

const (
	ROW_COUNT    = 4
	NOZZLE_COUNT = 1280
	NOZZLE_SPACE = 169.3 * printheads.UM
	ROW_OFFSET   = 84.65 * printheads.UM
	ROW_SPACE_A  = 550.3 * printheads.UM
	ROW_SPACE_B  = 11.811 * printheads.MM
)

var NozzleMap0 = map[int]printheads.Nozzle{
	1: printheads.Nozzle{
		Index: 0,
		Row:   0,
		Position: printheads.Position{
			X: -84650,
			Y: -12361300,
		},
	},
	3: printheads.Nozzle{
		Index: 2,
		Row:   1,
		Position: printheads.Position{
			X: 0,
			Y: -11811000,
		},
	},
	2: printheads.Nozzle{
		Index: 1,
		Row:   2,
		Position: printheads.Position{
			X: -84650,
			Y: -550300,
		},
	},
	4: printheads.Nozzle{
		Index: 3,
		Row:   3,
		Position: printheads.Position{
			X: 0,
			Y: 0,
		},
	},
	1277: printheads.Nozzle{
		Index: 1276,
		Row:   0,
		Position: printheads.Position{
			X: 53922050,
			Y: -12361300,
		},
	},
	1279: printheads.Nozzle{
		Index: 1278,
		Row:   1,
		Position: printheads.Position{
			X: 54006700,
			Y: -11811000,
		},
	},
	1278: printheads.Nozzle{
		Index: 1277,
		Row:   2,
		Position: printheads.Position{
			X: 53922050,
			Y: -550300,
		},
	},
	1280: printheads.Nozzle{
		Index: 1279,
		Row:   3,
		Position: printheads.Position{
			X: 54006700,
			Y: 0,
		},
	},
}

var NozzleMap1 = map[int]printheads.Nozzle{
	1: printheads.Nozzle{
		Index: 0,
		Row:   0,
		Position: printheads.Position{
			X: 27088000,
			Y: -6180650,
		},
	},
	3: printheads.Nozzle{
		Index: 2,
		Row:   1,
		Position: printheads.Position{
			X: 27172650,
			Y: -5630350,
		},
	},
	2: printheads.Nozzle{
		Index: 1,
		Row:   2,
		Position: printheads.Position{
			X: 27088000,
			Y: 5630350,
		},
	},
	4: printheads.Nozzle{
		Index: 3,
		Row:   3,
		Position: printheads.Position{
			X: 27172650,
			Y: 6180650,
		},
	},
	1277: printheads.Nozzle{
		Index: 1276,
		Row:   0,
		Position: printheads.Position{
			X: 81094700,
			Y: -6180650,
		},
	},
	1279: printheads.Nozzle{
		Index: 1278,
		Row:   1,
		Position: printheads.Position{
			X: 81179350,
			Y: -5630350,
		},
	},
	1278: printheads.Nozzle{
		Index: 1277,
		Row:   2,
		Position: printheads.Position{
			X: 81094700,
			Y: 5630350,
		},
	},
	1280: printheads.Nozzle{
		Index: 1279,
		Row:   3,
		Position: printheads.Position{
			X: 81179350,
			Y: 6180650,
		},
	},
}

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
		ROW_COUNT,
		NOZZLE_COUNT,
		NOZZLE_SPACE,
		ROW_OFFSET,
		ROW_SPACE_A,
		ROW_SPACE_B,
		0,
		0,
	)
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
			if n, ok := NozzleMap0[nozzle.Index+1]; ok {
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
		ROW_COUNT,
		NOZZLE_COUNT,
		NOZZLE_SPACE,
		ROW_OFFSET,
		ROW_SPACE_A,
		ROW_SPACE_B,
	)
	for _, v := range NozzleMap1 {
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
			if n, ok := NozzleMap1[nozzle.Index+1]; ok {
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

func TestUpdatePosition(t *testing.T) {
	zero := printheads.Position{
		X: 0,
		Y: 0,
	}
	test := printheads.Position{
		X: 27172650,
		Y: 6180650,
	}
	h, _ := printheads.NewPrintHeadLineD(
		ROW_COUNT,
		NOZZLE_COUNT,
		NOZZLE_SPACE,
		ROW_OFFSET,
		ROW_SPACE_A,
		ROW_SPACE_B,
		zero.X,
		zero.Y,
	)
	expected := NozzleMap0[1]
	actual := *h.Rows[0].Nozzles[0]
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"\nEXPECT: %#v\nGET: %#v\n",
			expected,
			actual,
		)
	}
	h.UpdatePosition(test.X, test.Y)
	expected = NozzleMap1[1]
	actual = *h.Rows[0].Nozzles[0]
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"\nEXPECT: %#v\nGET: %#v\n",
			expected,
			actual,
		)
	}
	h.UpdatePosition(zero.X, zero.Y)
	expected = NozzleMap0[1]
	actual = *h.Rows[0].Nozzles[0]
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(
			"\nEXPECT: %#v\nGET: %#v\n",
			expected,
			actual,
		)
	}
}
