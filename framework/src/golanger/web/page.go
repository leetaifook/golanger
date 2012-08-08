package web

import (
	"text/template"
)

type Page struct {
	*Site
	DefaultController   interface{}
	NotFoundtController interface{}
	Controller          map[string]interface{}
	Template            string
	TemplateFunc        template.FuncMap
	Config
	*Document
}

func (p *Page) Init() {
	p.Site.Init()
	p.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (p *Page) Reset() {
	document := p.Document
	globalCss, okCss := document.Css["global"]
	globalJs, okJs := document.Js["global"]
	globalImg, okImg := document.Img["global"]
	p.Document = &Document{
		Static:        document.Static,
		Theme:         document.Theme,
		GlobalCssFile: document.GlobalCssFile,
		GlobalJsFile:  document.GlobalJsFile,
		Css:           map[string]string{},
		Js:            map[string]string{},
		Img:           map[string]string{},
	}

	if okCss {
		p.Document.Css["global"] = globalCss
	}

	if okJs {
		p.Document.Js["global"] = globalJs
	}

	if okImg {
		p.Document.Img["global"] = globalImg
	}
}

func (p *Page) SetDefaultController(c interface{}) *Page {
	p.DefaultController = c

	return p
}

func (p *Page) SetNotFoundController(c interface{}) *Page {
	p.NotFoundtController = c

	return p
}

func (p *Page) SetTemplate(template string) *Page {
	p.Template = template

	return p
}

func (p *Page) RegisterController(relUrlPath string, i interface{}) *Page {
	if _, ok := p.Controller[relUrlPath]; !ok {
		p.Controller[relUrlPath] = i
	}

	return p
}

func (p *Page) GetController(urlPath string) (i interface{}) {
	relUrlPath := urlPath[len(p.Site.Root):]
	i, ok := p.Controller[relUrlPath]
	if !ok {
		i = p.NotFoundtController
	}

	return
}
