package printhead

import (
	"posam/util/geometry"
	"posam/util/reagent"
)

type Printhead struct {
	Index    int
	Reagents []*reagent.Reagent

	Path    string
	Reverse bool

	OffsetX float64
	OffsetY float64
}

func NewPrinthead(
	index int,
	reagents []*reagent.Reagent,
	path string,
	reverse bool,
	offsetx float64,
	offsety float64,
) *Printhead {
	p := &Printhead{
		Index:    index,
		Reagents: reagents,
		Path:     path,
		Reverse:  reverse,
		OffsetX:  offsetx,
		OffsetY:  offsety,
	}
	return p
}

func (p *Printhead) MakeNozzles(
	posbx int,
	posby int,
) []*Nozzle {
	nozzles := []*Nozzle{}
	for index := 0; index < 1280; index++ {
		nozzle, _ := NewNozzle(index)
		nozzle.Reagent = p.Reagents[nozzle.RowIndex]
		nozzle.Printhead = p
		posx := posbx - 1
		posy := posby
		switch nozzle.RowIndex {
		case 0:
			nozzle.Pos = geometry.NewPosition(
				posx+nozzle.Index,
				posy+292,
			)
		case 1:
			nozzle.Pos = geometry.NewPosition(
				posx+nozzle.Index-2,
				posy+13,
			)
		case 2:
			nozzle.Pos = geometry.NewPosition(
				posx+nozzle.Index,
				posy+279,
			)
		case 3:
			nozzle.Pos = geometry.NewPosition(
				posx+nozzle.Index-2,
				posy,
			)
		}
		nozzles = append(nozzles, nozzle)
	}
	return nozzles
}

func (p *Printhead) MakeNozzlesMH5440(
	posx int,
	posy int,
) []*Nozzle {
	nozzles := []*Nozzle{}
	for index := 0; index < 1280; index++ {
		nozzle, _ := NewNozzle(index)
		nozzle.Reagent = p.Reagents[nozzle.RowIndex]
		nozzle.Printhead = p
		switch nozzle.RowIndex {
		case 0:
			nozzle.Pos = geometry.NewPosition(
				posx-3+nozzle.Index,
				posy+292,
			)
		case 1:
			nozzle.Pos = geometry.NewPosition(
				posx-3+nozzle.Index,
				posy+13,
			)
		case 2:
			nozzle.Pos = geometry.NewPosition(
				posx-3+nozzle.Index,
				posy+279,
			)
		case 3:
			nozzle.Pos = geometry.NewPosition(
				posx-3+nozzle.Index,
				posy,
			)
		}
		nozzles = append(nozzles, nozzle)
	}
	return nozzles
}
