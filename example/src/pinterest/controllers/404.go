package controllers

import (
	"net/http"
)

type Page404 struct {
	*App
}

func (p *Page404) Init() {
	p.Document.GenerateHtml = false
	p.Template = "_notfound/404.html"
	p.ResponseWriter.WriteHeader(http.StatusNotFound)
}
