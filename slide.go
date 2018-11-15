package printing

import ()

const (
	SPACE_X_FROM_CENTROID = 10 * MM
	SPACE_Y_FROM_CENTROID = 25 * MM
)

type Slide struct {
	CentroidPosition *Position
	AvailableArea    *Area
}

func NewSlide(posx int, posy int) *Slide {
	return &Slide{
		CentroidPosition: &Position{
			X: posx,
			Y: posy,
		},
		AvailableArea: &Area{
			Top:    posy + SPACE_Y_FROM_CENTROID,
			Right:  posx + SPACE_X_FROM_CENTROID,
			Bottom: posy - SPACE_Y_FROM_CENTROID,
			Left:   posx - SPACE_X_FROM_CENTROID,
		},
	}
}
