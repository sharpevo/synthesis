package printhead

import ()

type Array struct {
	Nozzles        []*Nozzle
	SightBottom    *Nozzle
	SightTop       *Nozzle
	PrintheadCount int
}

func NewArray(
	nozzles []*Nozzle,
	printheadCount int,
) *Array {
	a := &Array{
		Nozzles:        nozzles,
		PrintheadCount: printheadCount,
	}
	a.SightTop, a.SightBottom = a.sights()
	return a
}

func (a *Array) sights() (top *Nozzle, bottom *Nozzle) {
	// assume vertical layout and aligned
	// and the distance is integer multiple of 4, i.e. 25.4/150
	for _, nozzle := range a.Nozzles {
		if nozzle.Index == 0 {
			if bottom == nil {
				bottom = nozzle
			} else {
				// bottom left nozzle of the bottom right printhead
				if nozzle.Pos.Y < bottom.Pos.Y {
					bottom = nozzle
				}
			}
		}
		if nozzle.Index == 3 {
			if top == nil {
				top = nozzle
			} else {
				// top left nozzle of the top right printhead
				if nozzle.Pos.Y > top.Pos.Y {
					top = nozzle
				}
			}
		}
	}
	return top, bottom
}

func (a *Array) MoveBottomRow(rowIndex, posx int, posy int) {
	var deltax, deltay int
	switch rowIndex {
	case 0:
		deltax = posx - a.SightBottom.Pos.X
		deltay = posy - a.SightBottom.Pos.Y
	case 1:
		deltax = posx - (a.SightBottom.Pos.X + 1)
		deltay = posy - (a.SightBottom.Pos.Y + 279)
	case 2:
		deltax = posx - (a.SightBottom.Pos.X + 2)
		deltay = posy - (a.SightBottom.Pos.Y + 13)
	case 3:
		deltax = posx - (a.SightBottom.Pos.X + 3)
		deltay = posy - (a.SightBottom.Pos.Y + 292)
	}
	for _, n := range a.Nozzles {
		n.Pos.X = n.Pos.X + deltax
		n.Pos.Y = n.Pos.Y + deltay
	}
}

func (a *Array) MoveTopRow(rowIndex, posx int, posy int) {
	var deltax, deltay int
	switch rowIndex {
	case 0:
		deltax = posx - (a.SightTop.Pos.X - 3)
		deltay = posy - (a.SightTop.Pos.Y - 292)
	case 1:
		deltax = posx - (a.SightTop.Pos.X - 2)
		deltay = posy - (a.SightTop.Pos.Y - 13)
	case 2:
		deltax = posx - (a.SightTop.Pos.X - 1)
		deltay = posy - (a.SightTop.Pos.Y - 279)
	case 3:
		deltax = posx - a.SightTop.Pos.X
		deltay = posy - a.SightTop.Pos.Y
	}
	for _, n := range a.Nozzles {
		n.Pos.X = n.Pos.X + deltax
		n.Pos.Y = n.Pos.Y + deltay
	}
}

func (a *Array) Top() int {
	return a.SightTop.Pos.Y
}

func (a *Array) Bottom() int {
	return a.SightBottom.Pos.Y
}
