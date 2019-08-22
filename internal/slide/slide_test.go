package slide_test

import (
	//"fmt"
	"synthesis/internal/geometry"
	"synthesis/internal/slide"
	"testing"
)

func TestNewSlide(t *testing.T) {
	s := slide.NewSlide(10*geometry.MM, 20*geometry.MM)
	if s.AvailableArea.Top != 45*geometry.MM {
		t.Errorf("invalid top: %d\n", s.AvailableArea.Top)
	}
	if s.AvailableArea.Right != 20*geometry.MM {
		t.Errorf("invalid right: %d\n", s.AvailableArea.Right)
	}
	if s.AvailableArea.Bottom != -5*geometry.MM {
		t.Errorf("invalid bottom: %d\n", s.AvailableArea.Bottom)
	}
	if s.AvailableArea.Left != 0*geometry.MM {
		t.Errorf("invalid left: %d\n", s.AvailableArea.Left)
	}
}
