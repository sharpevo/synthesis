package substrate_test

import (
	"fmt"
	"synthesis/internal/substrate"
	"testing"
)

func TestNewSubstrate(t *testing.T) { // {{{
	spots10x10 := []*substrate.Spot{}
	for _ = range [900]int{} {
		spots10x10 = append(spots10x10, substrate.NewSpot())
	}
	spots11x14 := []*substrate.Spot{}
	for _ = range [1638]int{} {
		spots11x14 = append(spots11x14, substrate.NewSpot())
	}
	spots4x21 := []*substrate.Spot{}
	for _ = range [27000]int{} {
		spots4x21 = append(spots4x21, substrate.NewSpot())
	}
	// cases{{{
	type input struct {
		slidecounth int
		slidecountv int

		slidewidth  float64
		slideheight float64

		slidespaceh float64
		slidespacev float64

		spots []*substrate.Spot

		spotspace int
		leftmost  int
	}

	type output struct {
		SpotCount  int
		SpotSpaceu int

		MaxSpotsh int
		MaxSpotsv int
		Width     int
		Height    int
		LeftMostu int

		SlideHeight  float64
		SlideWidth   float64
		SlideHeightu int
		SlideWidthu  int
		SlideNumh    int
		SlideNumv    int
		SlideSpacehu int
		SlideSpacevu int

		Top    int
		Bottom int
	}
	cases := []struct {
		title  string
		input  input
		output output
	}{
		{
			"10x10",
			input{
				slidecounth: 3,
				slidecountv: 3,
				slidewidth:  1.524,
				slideheight: 1.524,
				slidespaceh: 1,
				slidespacev: 1,
				spots:       spots10x10,
				spotspace:   4,
				leftmost:    0,
			},
			output{
				Top:          156,
				Bottom:       0,
				SpotCount:    900,
				SpotSpaceu:   4,
				MaxSpotsh:    156,
				MaxSpotsv:    156,
				Width:        157,
				Height:       157,
				LeftMostu:    0,
				SlideHeight:  1.524,
				SlideWidth:   1.524,
				SlideHeightu: 36,
				SlideWidthu:  36,
				SlideNumh:    3,
				SlideNumv:    3,
				SlideSpacehu: 24,
				SlideSpacevu: 24,
			},
		},
		{
			"11x14",
			input{
				slidecounth: 3,
				slidecountv: 3,
				slidewidth:  2.201333,
				slideheight: 2.032,
				slidespaceh: 1,
				slidespacev: 1,
				spots:       spots11x14,
				spotspace:   4,
				leftmost:    0,
			},
			output{
				Top:          192,
				Bottom:       0,
				SpotCount:    1638,
				SpotSpaceu:   4,
				MaxSpotsh:    204,
				MaxSpotsv:    192,
				Width:        205,
				Height:       193,
				LeftMostu:    0,
				SlideWidth:   2.201333,
				SlideHeight:  2.032,
				SlideWidthu:  52,
				SlideHeightu: 48,
				SlideNumh:    3,
				SlideNumv:    3,
				SlideSpacehu: 24,
				SlideSpacevu: 24,
			},
		},
		{
			"4x21",
			input{
				slidecounth: 3,
				slidecountv: 3,
				slidewidth:  4,
				slideheight: 21,
				slidespaceh: 2,
				slidespacev: 4,
				spots:       spots4x21,
				spotspace:   4,
				leftmost:    1652,
			},
			output{
				Top:          1678,
				Bottom:       2,
				SpotCount:    27000,
				SpotSpaceu:   4,
				MaxSpotsh:    379,
				MaxSpotsv:    1678,
				Width:        380,
				Height:       1679,
				LeftMostu:    1652,
				SlideWidth:   4,
				SlideHeight:  21,
				SlideWidthu:  92,
				SlideHeightu: 496,
				SlideNumh:    3,
				SlideNumv:    3,
				SlideSpacehu: 44,
				SlideSpacevu: 92,
			},
		},
	} // }}}
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			s, err := substrate.NewSubstrate(
				c.input.slidecounth,
				c.input.slidecountv,
				c.input.slidewidth,
				c.input.slideheight,
				c.input.slidespaceh,
				c.input.slidespacev,
				c.input.spots,
				c.input.spotspace,
				c.input.leftmost,
			)
			if err != nil {
				t.Fatal(err)
			}
			if s.Top() != c.output.Top ||
				s.Bottom() != c.output.Bottom ||
				s.SpotCount != c.output.SpotCount ||
				s.SpotSpaceu != c.output.SpotSpaceu ||
				s.MaxSpotsh != c.output.MaxSpotsh ||
				s.MaxSpotsv != c.output.MaxSpotsv ||
				s.Width != c.output.Width ||
				s.Height != c.output.Height ||
				s.LeftMostu != c.output.LeftMostu ||
				s.SlideWidth != c.output.SlideWidth ||
				s.SlideHeight != c.output.SlideHeight ||
				s.SlideWidthu != c.output.SlideWidthu ||
				s.SlideHeightu != c.output.SlideHeightu ||
				s.SlideNumh != c.output.SlideNumh ||
				s.SlideNumv != c.output.SlideNumv ||
				s.SlideSpacehu != c.output.SlideSpacehu ||
				s.SlideSpacevu != c.output.SlideSpacevu {
				t.Error(
					err,
					"\nEXPECT:\n",
					c.output.Top,
					c.output.Bottom,
					c.output.SpotCount,
					c.output.SpotSpaceu,
					c.output.MaxSpotsh,
					c.output.MaxSpotsv,
					c.output.Width,
					c.output.Height,
					c.output.LeftMostu,
					c.output.SlideWidth,
					c.output.SlideHeight,
					c.output.SlideWidthu,
					c.output.SlideHeightu,
					c.output.SlideNumh,
					c.output.SlideNumv,
					c.output.SlideSpacehu,
					c.output.SlideSpacevu,
					"\nACTUAL:\n",
					s.Top(),
					s.Bottom(),
					s.SpotCount,
					s.SpotSpaceu,
					s.MaxSpotsh,
					s.MaxSpotsv,
					s.Width,
					s.Height,
					s.LeftMostu,
					s.SlideWidth,
					s.SlideHeight,
					s.SlideWidthu,
					s.SlideHeightu,
					s.SlideNumh,
					s.SlideNumv,
					s.SlideSpacehu,
					s.SlideSpacevu,
				)
			}

		})
	}
} // }}}

func TestEdge(t *testing.T) { // {{{
	s := &substrate.Substrate{
		MaxSpotsv: 5,
	}
	if s.Top() != 5 {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			2,
			s.Top(),
		)
	}
	if s.Bottom() != 1 {
		t.Errorf(
			"\nEXPECT: %v\n GET: %v\n\n",
			2,
			s.Top(),
		)
	}
} // }}}

func TestSpotsPerSlide(t *testing.T) { // {{{
	cases := []struct {
		width  float64
		height float64
		expect int
	}{
		{10, 10, 3600},
		{11, 14, 5478},
		{4, 21, 3000},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%vx%v", c.width, c.height), func(t *testing.T) {
			s := &substrate.Substrate{
				SlideWidth:  c.width,
				SlideHeight: c.height,
			}
			actual := s.SpotsPerSlide()
			if actual != c.expect {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.expect,
					actual,
				)
			}
		})
	}
} // }}}

func TestIsOverloaded(t *testing.T) { // {{{
	cases := []struct {
		slidenumh   int
		slidenumv   int
		spotcount   int
		slidewidth  float64
		slideheight float64
		expectError bool
	}{
		{
			3,
			3,
			900,
			10,
			10,
			false,
		},
		{
			3,
			3,
			52302, // 52302/5478 + 0.5 = 10
			11,
			14,
			true,
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%vx%v", c.slidewidth, c.slideheight), func(t *testing.T) {
			s := &substrate.Substrate{
				SlideNumh:   c.slidenumh,
				SlideNumv:   c.slidenumv,
				SpotCount:   c.spotcount,
				SlideWidth:  c.slidewidth,
				SlideHeight: c.slideheight,
			}
			err := s.IsOverloaded()
			fmt.Println(err)
			if err != nil && !c.expectError {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					"no error",
					err,
				)
			}
		})
	}
} // }}}

func TestStrip(t *testing.T) { // {{{
	cases := []struct {
		maxspotsh int
		leftmostu int
		expect    int
	}{
		{
			1279,
			0,
			1, // should print even less than 1280
		},
		{
			1280,
			0,
			1,
		},
		{
			1281,
			0,
			2, // quo + 1
		},
		{
			1279,
			2,
			2, // quo + 1
		},
		{
			1280,
			1280,
			2, // quo + 1
		},
		{
			1281,
			1280,
			3, // quo + 1
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%v %v", c.maxspotsh, c.leftmostu), func(t *testing.T) {
			s := &substrate.Substrate{}
			s.MaxSpotsh = c.maxspotsh
			s.LeftMostu = c.leftmostu
			actual := s.Strip()
			if actual != c.expect {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.expect,
					actual,
				)
			}
		})
	}
} // }}}

func TestMarginRight(t *testing.T) { // {{{
	cases := []struct {
		slidewidthu  int
		slidespacehu int
		curcolumn    int
		expect       int
	}{
		{
			36,
			24,
			1,
			36,
		},
		{
			36,
			24,
			2,
			96,
		},
		{
			36,
			24,
			3,
			156,
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("column-%v", c.curcolumn), func(t *testing.T) {
			s := &substrate.Substrate{}
			s.SlideWidthu = c.slidewidthu
			s.SlideSpacehu = c.slidespacehu
			s.SetCurColumn(c.curcolumn)
			actual := s.MarginRightu()
			if actual != c.expect {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.expect,
					actual,
				)
			}
		})
	}
} // }}}

func TestMarginBottom(t *testing.T) { // {{{
	cases := []struct {
		maxspotsv    int
		slideheightu int
		slidespacevu int
		curslide     int
		expect       int
	}{
		{
			156,
			36,
			24,
			1,
			120,
		},
		{
			157,
			36,
			24,
			2,
			61,
		},
		{
			157,
			36,
			24,
			3,
			1,
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("slide-%v", c.curslide), func(t *testing.T) {
			s := &substrate.Substrate{}
			s.MaxSpotsv = c.maxspotsv
			s.SlideHeightu = c.slideheightu
			s.SlideSpacevu = c.slidespacevu
			s.SetCurSlide(c.curslide)
			actual := s.MarginBottomu()
			if actual != c.expect {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					c.expect,
					actual,
				)
			}
		})
	}
} // }}}

func TestAllocate(t *testing.T) { // {{{
	cases := []struct {
		title        string
		x            int
		y            int
		marginright  int
		marginbottom int

		spotspaceu   int
		slidewidthu  int
		slidespacehu int
		slidespacevu int
		slidenumh    int
		slidenumv    int
		maxspotsv    int
		curslide     int
		curcolumn    int

		expectX         int
		expectY         int
		expectCurSlide  int
		expectCurColumn int
		expectError     bool
	}{
		{
			title:        "current line",
			x:            5, // x < marginright
			y:            0,
			marginright:  10,
			marginbottom: 20,

			expectX:         5, // x
			expectY:         0, // y
			expectCurSlide:  0,
			expectCurColumn: 0,
			expectError:     false,
		},
		{
			title:        "new line",
			x:            11, // x > marginright
			y:            27, // y >= marginbottom + spotspaceu
			marginright:  10,
			marginbottom: 20,

			spotspaceu:  6,
			slidewidthu: 2,

			expectX:         8,  // marginright - slidewidthu
			expectY:         21, // y - spotspaceu
			expectCurSlide:  0,
			expectCurColumn: 0,
			expectError:     false,
		},
		{
			title:        "new slide",
			x:            11, // x > marginright
			y:            25, // y < marginbottom + spotspaceu
			marginright:  10,
			marginbottom: 20,

			spotspaceu: 6,
			slidenumv:  4,
			curslide:   2, // curslide < slidenumv - 1

			slidewidthu:  2,
			slidespacevu: 3,

			expectX:         8,  // marginright - slidewidthu
			expectY:         22, // y - slidespacevu
			expectCurSlide:  3,  // curslide + 1
			expectCurColumn: 0,
			expectError:     false,
		},
		{
			title:        "new column",
			x:            11, // x > marginright
			y:            25, // y < marginbottom + spotspaceu
			marginright:  10,
			marginbottom: 20,

			spotspaceu: 6,
			slidenumv:  4,
			curslide:   4, // curslide > slidenumv - 1

			slidewidthu:  2,
			slidespacehu: 7,
			curcolumn:    2,
			maxspotsv:    50,

			slidenumh: 3,

			expectX:         18, // (slidewidthu + slidespacehu) * curcolumn
			expectY:         50, // maxspotsv
			expectCurSlide:  1,  // 1
			expectCurColumn: 3,  // curcolumn + 1, <= slidenumh
			expectError:     false,
		},
		{
			title:        "new column but overloaded",
			x:            11, // x > marginright
			y:            25, // y < marginbottom + spotspaceu
			marginright:  10,
			marginbottom: 20,

			spotspaceu: 6,
			slidenumv:  4,
			curslide:   4, // curslide > slidenumv - 1

			slidewidthu:  2,
			slidespacehu: 7,
			curcolumn:    2,
			maxspotsv:    50,

			slidenumh: 2,

			expectX:         18, // (slidewidthu + slidespacehu) * curcolumn
			expectY:         50, // maxspotsv
			expectCurSlide:  1,  // 1
			expectCurColumn: 3,  // curcolumn + 1, > slidenumh
			expectError:     true,
		},
	}
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			s := &substrate.Substrate{}
			s.SpotSpaceu = c.spotspaceu
			s.SlideWidthu = c.slidewidthu
			s.SlideSpacehu = c.slidespacehu
			s.SlideSpacevu = c.slidespacevu
			s.SlideNumh = c.slidenumh
			s.SlideNumv = c.slidenumv
			s.MaxSpotsv = c.maxspotsv
			s.SetCurSlide(c.curslide)
			s.SetCurColumn(c.curcolumn)
			actualX, actualY, err := s.Allocate(
				c.x, c.y, c.marginright, c.marginbottom)
			if (!c.expectError && err != nil) ||
				(c.expectError && err == nil) {
				t.Errorf(
					"\nEXPECT: %v\n GET: %v\n\n",
					"no error",
					err,
				)
			}
			if actualX != c.expectX ||
				actualY != c.expectY ||
				s.CurSlide() != c.expectCurSlide ||
				s.CurColumn() != c.expectCurColumn {
				t.Error(
					"\nEXPECT\n",
					c.expectX,
					c.expectY,
					c.expectCurSlide,
					c.expectCurColumn,
					"\nGET\n",
					actualX,
					actualY,
					s.CurSlide(),
					s.CurColumn(),
				)
			}
		})
	}
} // }}}
