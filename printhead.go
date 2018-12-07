package printhead

import (
	"posam/util/geometry"
	"posam/util/reagent"
)

type Printhead struct {
	Index    int
	Reagents []*reagent.Reagent
}

func NewPrinthead(
	index int,
	reagents []*reagent.Reagent,
) *Printhead {
	p := &Printhead{
		Index:    index,
		Reagents: reagents,
	}
	return p
}

func (p *Printhead) MakeNozzles(
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
				posx+nozzle.Index,
				posy,
			)
		case 1:
			nozzle.Pos = geometry.NewPosition(
				posx+nozzle.Index,
				posy+279,
			)
		case 2:
			nozzle.Pos = geometry.NewPosition(
				posx+nozzle.Index,
				posy+13,
			)
		case 3:
			nozzle.Pos = geometry.NewPosition(
				posx+nozzle.Index,
				//posy+305,
				posy+292,
			)
		}
		nozzles = append(nozzles, nozzle)
	}
	return nozzles
}
