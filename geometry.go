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
