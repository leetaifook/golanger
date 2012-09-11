package web

import (
	"io/ioutil"
	"net/http"
	"os"
)

type templateCache struct {
	ModTime int64
	Content string
}

type Site struct {
	*Base
	TemplateCache map[string]templateCache
	Root          string
	Version       string
}

func (s *Site) Init(w http.ResponseWriter, r *http.Request) *Site {
	s.Base.Init(w, r)

	return s
}

func (s *Site) SetTemplateCache(tmplKey, tmplPath string) {
	if tmplFi, err := os.Stat(tmplPath); err == nil {
		if b, err := ioutil.ReadFile(tmplPath); err == nil {
			s.TemplateCache[tmplKey] = templateCache{
				ModTime: tmplFi.ModTime().Unix(),
				Content: string(b),
			}
		}
	}

}

func (s *Site) GetTemplateCache(tmplKey string) templateCache {
	if tmpl, ok := s.TemplateCache[tmplKey]; ok {
		return tmpl
	}

	return templateCache{}
}

func (s *Site) DelTemplateCache(tmplKey string) {
	delete(s.TemplateCache, tmplKey)
}
