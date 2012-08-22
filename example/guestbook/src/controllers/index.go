package controllers

import (
	"fmt"
	. "golanger/middleware"
	"golanger/utils"
	. "models"
	"net/http"
)

type PageIndex struct {
	*App
}

func init() {
	Page.RegisterController("index/", &PageIndex{Page})
}

func (p *PageIndex) Index() {
	mgo := Middleware.Get("db").(*utils.Mongo)
	coll := mgo.C(ColGuestBook)

	query := coll.Find(nil).Sort("-timestamp")

	var entries []ModelGuestBook
	if err := query.All(&entries); err != nil {
		http.Error(p.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	p.Body = entries
}

func (p *PageIndex) Sign() {
	if p.Request.Method != "POST" {
		p.Body = "不支持这种请求方式: " + fmt.Sprintf("%v", p.Request.Method)
		p.Template = "index/error.html"
		return
	}

	entry := NewGuestBook()
	entry.Name = p.POST["name"]
	entry.Message = p.POST["message"]

	if entry.Name == "" {
		entry.Name = "Some dummy who forgot a name"
	}
	if entry.Message == "" {
		entry.Message = "Some dummy who forgot a message."
	}

	mgo := Middleware.Get("db").(*utils.Mongo)
	coll := mgo.C(ColGuestBook)

	if err := coll.Insert(entry); err != nil {
		p.Body = "数据库错误：" + fmt.Sprintf("%v", err)
		p.Template = "index/error.html"
		return
	} else {
		http.Redirect(p.ResponseWriter, p.Request, "/", http.StatusFound)
	}

}
