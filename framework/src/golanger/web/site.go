package web

import (
	"net/http"
)

type Site struct {
	*Base
	Root    string
	Version string
}

func (s *Site) Init(w http.ResponseWriter, r *http.Request) *Site {
	s.Base.Init(w, r)

	return s
}
