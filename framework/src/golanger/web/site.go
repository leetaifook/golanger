package web

type Site struct {
	*Base
	Root    string
	Version string
}

func (s *Site) Init() *Site {
	s.Base.Init()

	return s
}
