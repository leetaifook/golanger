package web

import (
	"fmt"
	"strings"
	"sync"
	"text/template"
)

type Page struct {
	PageRLock sync.RWMutex
	PageLock  sync.Mutex
	*Site
	DefaultController   interface{}
	NotFoundtController interface{}
	Controller          map[string]interface{}
	CurrentController   string
	CurrentAction       string
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
		Func:          template.FuncMap{},
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
	var relUrlPath string
	if strings.HasPrefix(urlPath, p.Site.Root) {
		relUrlPath = urlPath[len(p.Site.Root):]
	} else {
		relUrlPath = urlPath
	}

	i, ok := p.Controller[relUrlPath]
	if !ok {
		i = p.NotFoundtController
	}

	return
}

func (p *Page) AddTemplateFunc(name string, i interface{}) {
	_, ok := p.TemplateFunc[name]
	if !ok {
		p.TemplateFunc[name] = i
	} else {
		fmt.Println("func:" + name + " be added,do not reepeat to add")
	}
}

func (p *Page) DelTemplateFunc(name string) {
	if _, ok := p.TemplateFunc[name]; ok {
		delete(p.TemplateFunc, name)
	}
}
