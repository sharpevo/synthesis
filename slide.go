package slide

import (
	"posam/util/geometry"
)

const (
	SPACE_X_FROM_CENTROID = 10 * geometry.MM
	SPACE_Y_FROM_CENTROID = 25 * geometry.MM
)

type Slide struct {
	CentroidPosition *geometry.Position
	AvailableArea    *geometry.Area
}

func NewSlide(posx int, posy int) *Slide {
	return &Slide{
		CentroidPosition: &geometry.Position{
			X: posx,
			Y: posy,
		},
		AvailableArea: &geometry.Area{
			Top:    posy + SPACE_Y_FROM_CENTROID,
			Right:  posx + SPACE_X_FROM_CENTROID,
			Bottom: posy - SPACE_Y_FROM_CENTROID,
			Left:   posx - SPACE_X_FROM_CENTROID,
		},
	}
}
