package controllers

import (
	"encoding/json"
	"fmt"
	. "golanger/middleware"
	"golanger/utils"
	"helper"
	. "models"
	"net/http"
)

type PageIndex struct {
	Application
}

func init() {
	App.RegisterController("index/", PageIndex{})
}

func (p *PageIndex) Init(w http.ResponseWriter, r *http.Request) {
	p.OffRight = true
	p.Application.Init(w, r)
}

func (p *PageIndex) Index(w http.ResponseWriter, r *http.Request) {
}

func (p *PageIndex) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if _, ok := p.POST["ajax"]; ok {
			p.Hide = true
			mgoServer := Middleware.Get("db").(*helper.Mongo)

			m := utils.M{
				"status":  1,
				"message": "",
			}

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Cache-Control", "no-store")

			username := p.POST["username"]
			password := p.POST["password"]
			passwordMd5 := utils.Strings(password).Md5()

			colQuerier := utils.M{"name": username, "status": 1, "delete": 0}
			colSelecter := utils.M{"password": 1}
			col := ModelUser{}
			var jres []byte
			err := mgoServer.C(ColUser).Find(colQuerier).Select(colSelecter).One(&col)
			if err != nil || col.Password == "" {
				m["status"] = -1
				m["message"] = "无此用户"
			} else {
				if passwordMd5 != col.Password {
					m["status"] = 0
					m["message"] = "密码错误"
				} else {
					m["back_url"] = ""
					if _, ok := p.GET["back_url"]; ok {
						m["back_url"] = p.GET["back_url"]
					}

					p.SESSION[p.M["SESSION_UNAME"].(string)] = username
					p.SESSION[p.M["SESSION_UKEY"].(string)] = passwordMd5
				}
			}

			jres, _ = json.Marshal(m)
			w.Write(jres)
			return
		}
	}
}

func (p *PageIndex) Logout(w http.ResponseWriter, r *http.Request) {
	sessionSign := p.COOKIE[p.Session.CookieName]
	if sessionSign != "" {
		p.Session.Clear(sessionSign)
	}

	http.Redirect(w, r, "/login.html", http.StatusFound)
}
