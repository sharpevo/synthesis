package substrate

func (s *Substrate) SpotsPerSlide() int {
	return s.spotsPerSlide()
}

func (s *Substrate) IsOverloaded() error {
	return s.isOverloaded()
}

var ParseLines = parseLines
