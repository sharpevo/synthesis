package geometry

import (
//"fmt"
)

const (
	NM = 1
	UM = 1e3
	MM = 1e6

	MPI = 25.4
	DPI = 600

	UNIT = MPI * MM / DPI
)

type Position struct {
	X int
	Y int
}

func NewPosition(posx int, posy int) *Position {
	return &Position{
		X: posx,
		Y: posy,
	}
}

func (p *Position) Equal(q *Position) bool {
	return p.X == q.X && p.Y == q.Y
}

func (p *Position) AtLeft(q *Position) bool {
	return p.X < q.X
}

func (p *Position) Sub(q *Position) *Position {
	return &Position{
		X: p.X - q.X,
		Y: p.Y - q.Y,
	}
}

func Millimeter(input float64) int {
	return int(input/MPI*DPI + 0.5)
}

func RoundedDot(input float64, dpi int) int {
	dot := Millimeter(input)
	return RoundDot(dot, dpi)
}

func RoundDot(input int, dpi int) int {
	div := DPI / dpi
	if rem := input % div; rem != 0 {
		input -= rem
	}
	return input
}

func Dot2Millimeter(input int) float64 {
	output := float64(input) / DPI * MPI
	return output
}

func Raw(value int, offset float64) float64 {
	return offset - Dot2Millimeter(value)
}
