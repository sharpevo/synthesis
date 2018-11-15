package printing_test

import (
	//"fmt"
	"posam/printing"
	"testing"
)

func TestNewSlide(t *testing.T) {
	s := printing.NewSlide(10*printing.MM, 20*printing.MM)
	if s.AvailableArea.Top != 45*printing.MM {
		t.Errorf("invalid top: %d\n", s.AvailableArea.Top)
	}
	if s.AvailableArea.Right != 20*printing.MM {
		t.Errorf("invalid right: %d\n", s.AvailableArea.Right)
	}
	if s.AvailableArea.Bottom != -5*printing.MM {
		t.Errorf("invalid bottom: %d\n", s.AvailableArea.Bottom)
	}
	if s.AvailableArea.Left != 0*printing.MM {
		t.Errorf("invalid left: %d\n", s.AvailableArea.Left)
	}
}
