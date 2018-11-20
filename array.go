package slide

import ()

type Array struct {
	Slides []*Slide
}

func NewArray(slides ...*Slide) *Array {
	a := &Array{}
	for _, s := range slides {
		a.AddSlide(s)
	}
	return a
}

func NewDefaultArray(spacex int, spacey int, count int) *Array {
	a := &Array{}
	for i := 0; i < count; i++ {
		slide := NewSlide(
			i*(WIDTH+SPACE),
			0,
			spacex,
			spacey,
		)
		a.AddSlide(slide)
	}
	return a
}

func (a *Array) AddSlide(slide *Slide) {
	a.Slides = append(a.Slides, slide)
}

func (a *Array) Top() int {
	top := 0
	for _, s := range a.Slides {
		if s.Top() > top {
			top = s.Top()
		}
	}
	return top
}

func (a *Array) Bottom() int {
	bottom := 0
	for _, s := range a.Slides {
		if bottom > s.Bottom() {
			bottom = s.Bottom()
		}
	}
	return bottom
}

func (a *Array) Left() int {
	left := 0
	for _, s := range a.Slides {
		if left > s.Left() {
			left = s.Left()
		}
	}
	return left
}

func (a *Array) AddSpot(spot *Spot) bool {
	for k, _ := range a.Slides {
		if a.Slides[k].IsFull() {
			continue
		}
		return a.Slides[k].AddSpot(spot)
	}
	return false
}

func (a *Array) NextSpotInVert(cycleIndex int) *Spot {
	var target *Spot
	sum := 0
	for _, slide := range a.Slides {
		for _, row := range slide.Spots {
			for _, spot := range row {
				if spot == nil {
					continue
				}
				sum += 1
				reagent := spot.NextReagent(cycleIndex)
				if reagent != nil {
					if target == nil {
						target = spot
					} else {
						if spot.Pos.X < target.Pos.X {
							target = spot
						}
					}
				}
			}
		}
	}
	return target
}

func (a *Array) SpotCount() int {
	count := 0
	for _, slide := range a.Slides {
		c := slide.SpotCount()
		if c > count {
			count = c
		}
	}
	return count
}

func (a *Array) ReagentCount() int {
	count := 0
	for _, slide := range a.Slides {
		count += slide.ReagentCount()
	}
	return count
}

func (a *Array) AvailableSpots() []*Spot {
	spots := []*Spot{}
	for _, slide := range a.Slides {
		spots = append(spots, slide.AvailableSpots()...)
	}
	return spots
}

func (a *Array) SpotsIn(top int, right int, bottom int, left int) []*Spot {
	spots := []*Spot{}
	for _, slide := range a.Slides {
		spots = append(spots, slide.SpotsIn(top, right, bottom, left)...)
	}
	return spots
}
