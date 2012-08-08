package controllers

import (
	. "../models"
	. "golanger/middleware"
	"golanger/utils"
	"golanger/web"
	"net/http"
	"net/url"
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

	a.checkLogin()
}

func (a *App) checkLogin() {
	if a.checkUser() {
		if a.Request.URL.Path == "/login.html" {
			http.Redirect(a.ResponseWriter, a.Request, "/index.html", http.StatusFound)
		}
	} else {
		if a.Request.URL.Path != "/login.html" {
			http.Redirect(a.ResponseWriter, a.Request, "/login.html?back_url="+url.QueryEscape(a.Request.URL.String()), http.StatusFound)
		}
	}
}

func (a *App) checkUser() (res bool) {
	username, uok := a.SESSION[a.M["SESSION_UNAME"].(string)]
	ukey, ukok := a.SESSION[a.M["SESSION_UKEY"].(string)]

	if uok && ukok {
		mgoServer := Middleware.Get("db").(*utils.Mongo)

		colQuerier := utils.M{"name": username, "password": ukey, "status": 1}
		colSelecter := utils.M{"name": 1}
		col := ModelUser{}
		err := mgoServer.C(ColUser).Find(colQuerier).Select(colSelecter).One(&col)

		if err == nil && col.Name != "" {
			res = true
		}
	}

	return res
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
