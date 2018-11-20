package reagent

import (
	"image/color"
)

const (
	REAGENT_A = "A"
	REAGENT_C = "C"
	REAGENT_G = "G"
	REAGENT_T = "T"
	REAGENT_Z = "Z"
	REAGENT_N = "-"
)

type Reagent struct {
	Name  string
	Color *color.NRGBA
}

var (
	red     = &color.NRGBA{0xff, 0x00, 0x00, 0xff}
	green   = &color.NRGBA{0x00, 0xff, 0x00, 0xff}
	blue    = &color.NRGBA{0x00, 0x00, 0xff, 0xff}
	magenta = &color.NRGBA{0xff, 0x00, 0xff, 0xff}
	yellow  = &color.NRGBA{0xff, 0xff, 0x00, 0xff}
	black   = &color.NRGBA{0x00, 0x00, 0x00, 0x00}
	white   = &color.NRGBA{0xff, 0xff, 0xff, 0xff}

	Activator = &Reagent{
		REAGENT_Z,
		yellow,
	}

	Nil = &Reagent{
		REAGENT_N,
		white,
	}

	reagentMap = map[string]*Reagent{
		REAGENT_A: &Reagent{
			REAGENT_A,
			red,
		},
		REAGENT_C: &Reagent{
			REAGENT_C,
			green,
		},
		REAGENT_G: &Reagent{
			REAGENT_G,
			blue,
		},
		REAGENT_T: &Reagent{
			REAGENT_T,
			magenta,
		},
		REAGENT_Z: Activator,
		REAGENT_N: &Reagent{
			REAGENT_N,
			white,
		},
	}
)

func NewReagent(name string) *Reagent {
	reagent, ok := reagentMap[name]
	if !ok {
		return &Reagent{
			name,
			black,
		}
	}
	return reagent
}
