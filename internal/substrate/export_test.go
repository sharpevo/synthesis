package substrate

func (s *Substrate) SpotsPerSlide() int {
	return s.spotsPerSlide()
}

func (s *Substrate) IsOverloaded() error {
	return s.isOverloaded()
}

func (s *Substrate) SetCurColumn(curColumn int) {
	s.curColumn = curColumn
}

func (s *Substrate) CurColumn() int {
	return s.curColumn
}

func (s *Substrate) SetCurSlide(curSlide int) {
	s.curSlide = curSlide
}

func (s *Substrate) CurSlide() int {
	return s.curSlide
}

func (s *Substrate) MarginRightu() int {
	return s.marginRightu()
}

func (s *Substrate) MarginBottomu() int {
	return s.marginBottomu()
}

func (s *Substrate) Allocate(x, y, right, bottom int) (int, int, error) {
	return s.allocate(x, y, right, bottom)
}

var SplitByLine = splitByLine
var SplitByChar = splitByChar
