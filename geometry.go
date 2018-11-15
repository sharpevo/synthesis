package geometry

import ()

const (
	NM = 1
	UM = 1e3
	MM = 1e6
)

type Position struct {
	X int
	Y int
}

type Area struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}
