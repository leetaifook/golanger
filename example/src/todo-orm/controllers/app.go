package controllers

import (
	"golanger/web"
	"strconv"
	"text/template"
	"time"
)

const (
	MAX_FORM_SIZE = 2 << 20 // 2MB => 2的20次方 乘以 2 =》 2 * 1024 * 1024
)

type App struct {
	*web.Page
}

func (a *App) Init() {
	a.Page.Init()
}

var Config = web.Config{
	TemplateDirectory:       "./view/",
	TemporaryDirectory:      "./tmp/",
	StaticDirectory:         "./static/",
	ThemeDirectory:          "theme/",
	Theme:                   "default",
	StaticCssDirectory:      "css/",
	StaticJsDirectory:       "js/",
	StaticImgDirectory:      "img/",
	HtmlDirectory:           "html/",
	UploadDirectory:         "upload/",
	TemplateGlobalDirectory: "_global/",
	TemplateGlobalFile:      "*",
	IndexDirectory:          "index/",
	IndexPage:               "index.html",
	SiteRoot:                "/",
	Environment:             map[string]string{},
	Database:                map[string]string{},
}

var Page = &App{
	Page: &web.Page{
		Site: &web.Site{
			Base: &web.Base{
				MAX_FORM_SIZE: MAX_FORM_SIZE,
			},
			Version: strconv.Itoa(time.Now().Year()),
		},
		Controller:   map[string]interface{}{},
		TemplateFunc: template.FuncMap{},
		Config:       Config,
		Document: &web.Document{
			Theme: Config.Theme,
			Css:   map[string]string{},
			Js:    map[string]string{},
			Img:   map[string]string{},
		},
	},
}
