package controllers

import (
	. "../models"
	. "golanger/middleware"
	"golanger/utils"
	"golanger/web"
	"labix.org/v2/mgo/bson"
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
	OffLogin bool
	OffRight bool
}

func (a *App) Init() {
	a.Page.Init()

	if a.OffLogin || a.checkLogin() {
		checkRight, _ := strconv.ParseBool(a.Environment["CheckRight"])

		if checkRight {
			a.getRole()
			a.getModule(false)

			if !a.OffRight {
				a.checkRight()
			}
		} else {
			a.getModule(true)
		}
	}

	a.OffLogin = false
	a.OffRight = false
}

func (a *App) getRole() {
	if username, nok := a.SESSION[a.Page.Config.M["SESSION_UNAME"].(string)]; nok {
		if _, ok := a.SESSION["role"]; !ok {
			mgoServer := Middleware.Get("db").(*utils.Mongo)
			cols := []ModelRole{}
			colSelector := utils.M{}
			colQuerier := utils.M{"delete": 0, "users": username}
			colSorter := []string{"-right.scope", "-status"}

			query := mgoServer.C(ColRole).Find(colQuerier).Select(colSelector).Sort(colSorter...)
			iter := query.Iter()
			for {
				col := ModelRole{}
				b := iter.Next(&col)
				if b != true {
					break
				}

				cols = append(cols, col)
			}

			a.SESSION["role"] = cols
		}
	}

}

func (a *App) getModule(showAll bool) {
	if _, ok := a.SESSION["modules"]; !ok {
		mgoServer := Middleware.Get("db").(*utils.Mongo)
		cols := []ModelModule{}
		colSelector := utils.M{"name": 1, "path": 1}
		colQuerier := utils.M{"status": 1}
		colSorter := []string{"-order", "-create_time"}
		hasModule := map[string]bool{}
		roles := a.SESSION["role"]
		query := mgoServer.C(ColModule).Find(colQuerier).Select(colSelector).Sort(colSorter...)
		iter := query.Iter()
		for {
			col := ModelModule{}
			b := iter.Next(&col)
			if b != true {
				break
			}

			if showAll {
				cols = append(cols, col)
			} else {
				for _, role := range roles.([]ModelRole) {
					if _, ok := hasModule[col.Path]; !ok {
						switch role.Right["scope"].(string) {
						case "3": //site
							cols = append(cols, col)
							hasModule[col.Path] = true
						case "2": //app
							cols = append(cols, col)
							hasModule[col.Path] = true
						case "1": //module
							if role.Right["modules"] != nil {
								for _, module := range role.Right["modules"].([]interface{}) {
									mod := module.(bson.M)
									if mod["module"].(string) == col.Path {
										cols = append(cols, col)
										hasModule[col.Path] = true
									}
								}
							}
						case "0": //action
							if role.Right["modules"] != nil {
								for _, module := range role.Right["modules"].([]interface{}) {
									mod := module.(bson.M)
									if len(mod["actions"].([]interface{})) > 0 {
										if mod["module"].(string) == col.Path {
											cols = append(cols, col)
											hasModule[col.Path] = true
										}
									}
								}
							}
						}
					}
				}
			}
		}

		a.SESSION["modules"] = cols
	}
}

func (a *App) checkRight() {
	var hasRight bool
	reqModule := a.CurrentController
	reqAction := a.CurrentAction
	for _, module := range a.SESSION["modules"].([]ModelModule) {
		if module.Path == reqModule {
			for _, role := range a.SESSION["role"].([]ModelRole) {
				if role.Right["scope"].(string) != "0" {
					hasRight = true
					goto Check_Right
				} else {
					if role.Right["modules"] != nil {
						for _, mod := range role.Right["modules"].([]interface{}) {
							m := mod.(bson.M)
							for _, action := range m["actions"].([]interface{}) {
								if reqAction == action.(string) {
									hasRight = true
									goto Check_Right
								}
							}
						}
					}
				}
			}
		}
	}

Check_Right:
	if !hasRight {
		a.ResponseWriter.WriteHeader(http.StatusForbidden)
		a.ResponseWriter.Write([]byte("无权限"))
		a.Close = true
	}
}

func (a *App) checkLogin() bool {
	var b bool
	if a.checkUser() {
		b = true
		if a.Request.URL.Path == "/login.html" {
			http.Redirect(a.ResponseWriter, a.Request, "/index.html", http.StatusFound)
			a.Close = true
		}
	} else {
		if a.Request.URL.Path != "/login.html" {
			http.Redirect(a.ResponseWriter, a.Request, "/login.html?back_url="+url.QueryEscape(a.Request.URL.String()), http.StatusFound)
			a.Close = true
		}
	}

	return b
}

func (a *App) checkUser() (res bool) {
	username, uok := a.SESSION[a.Page.Config.M["SESSION_UNAME"].(string)]
	ukey, ukok := a.SESSION[a.Page.Config.M["SESSION_UKEY"].(string)]

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
			Func:  template.FuncMap{},
		},
	},
}
