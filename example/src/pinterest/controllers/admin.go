package controllers

import (
	. "../models"
	"golanger/utils"
	"net/http"
	"path/filepath"
	"strconv"
)

type PageAdmin struct {
	*App
}

func init() {
	Page.RegisterController("admin/", &PageAdmin{Page})
}

func (p *PageAdmin) Init() {
	p.App.Init()
	_, fileName := filepath.Split(p.Request.URL.Path)
	if fileName != "login.html" {
		if _, ok := p.SESSION["user"]; !ok {
			http.Redirect(p.ResponseWriter, p.Request, "/admin/login.html", http.StatusFound)
		}
	}
}

func (p *PageAdmin) Login() {
	if p.Request.Method == "POST" {
		username := p.POST["username"]
		password := p.POST["password"]
		if username == p.Config.Environment["Username"] && password == p.Config.Environment["Password"] {
			p.SESSION["user"] = username
			http.Redirect(p.ResponseWriter, p.Request, "/admin/index.html", http.StatusFound)
		}
	}
}

func (p *PageAdmin) Logout() {
	delete(p.SESSION, "user")
	http.Redirect(p.ResponseWriter, p.Request, "/admin/index.html", http.StatusFound)
}

func (p *PageAdmin) Index() {
	body := utils.M{}
	body["invalidImages"], _ = GetInvalidImages()
	body["classes"], _ = GetClasses()

	p.Body = body
}

func (p *PageAdmin) Delete() {
	if p.Request.Method == "POST" {
		if _, ok := p.POST["ajax"]; ok {
			var id = p.POST["id"]
			idValue, _ := strconv.ParseInt(id, 0, 64)
			go DeleteImageWithId(idValue)
		}
	}
}

func (p *PageAdmin) Recover() {
	if p.Request.Method == "POST" {
		if _, ok := p.POST["ajax"]; ok {
			var id = p.POST["id"]
			idValue, _ := strconv.Atoi(id)
			image, err := GetImage(idValue)
			if err == nil {
				image.Status = 1
				go SaveImages(*image)
			}

		}
	}
}

func (p *PageAdmin) Adclass() {
	if p.Request.Method == "POST" {
		if _, ok := p.POST["ajax"]; ok {
			var name = p.POST["className"]
			go AddClass(Class{
				Name: name,
			})
		}
	}
}

func (p *PageAdmin) Declass() {
	if p.Request.Method == "POST" {
		if _, ok := p.POST["ajax"]; ok {
			var id = p.POST["classId"]
			idValue, _ := strconv.ParseInt(id, 0, 64)
			go DeleteClassWithId(idValue)
		}
	}
}

func (p *PageAdmin) Edclass() {
	if p.Request.Method == "POST" {
		if _, ok := p.POST["ajax"]; ok {
			var id = p.POST["classId"]
			var name = p.POST["className"]
			idValue, _ := strconv.ParseInt(id, 0, 64)
			go EditClassWithId(idValue, name)
		}
	}
}
