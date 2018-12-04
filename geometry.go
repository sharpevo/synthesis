package geometry

import (
	"fmt"
)

const (
	NM = 1
	UM = 1e3
	MM = 1e6

	UNIT = 25.4 * MM / 600
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

type Area struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

func Unit(input float64) int {
	return int(input*600/25.4 + 0.5)
}

func Mm2(input int) string {
	output := fmt.Sprintf("%.6f", float64(input)*25.4/600)
	fmt.Println("convert", input, output)
	return output
}

func Mm(input int) float64 {
	output := float64(input) * 25.4 / 600
	fmt.Println("convert", input, output)
	return output
}
