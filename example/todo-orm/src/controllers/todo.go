package controllers

import (
	"database/sql"
	"fmt"
	. "golanger/middleware"
	"golanger/utils"
	"models"
	"net/http"
	"strconv"
	"time"
)

type PageTodo struct {
	*Application
}

func init() {
	App.RegisterController("todo/", &PageTodo{App})
}

func (p *PageTodo) Index() {
	todos, err := models.GetTodoLists()
	if err != nil {
		fmt.Println(err)
	}

	p.Body = todos
}

func (p *PageTodo) New() {
	if p.Request.Method == "POST" {
		if title, ok := p.POST["title"]; ok {
			todo := models.Todo{
				Title:    title,
				PostDate: utils.NewTime().GetTimeToStr(time.Now().Unix()),
			}
			_, err := models.SaveTodo(todo)
			if err == nil {
				http.Redirect(p.ResponseWriter, p.Request, "/", http.StatusFound)
			} else {
				p.Body = "数据库错误：" + fmt.Sprintf("%v", err)
				p.Template = "todo/error.html"
			}
		}

	}
}

func (p *PageTodo) Edit() {
	if p.Request.Method == "GET" {
		if sid, ok := p.GET["id"]; ok {
			id, _ := strconv.Atoi(sid)
			todo := models.GetTodo(int64(id))
			p.Body = todo
		} else {
			p.Template = "todo/error.html"
		}
	}

	if p.Request.Method == "POST" {
		sid, iok := p.POST["id"]
		title, tok := p.POST["title"]
		if iok && tok {
			id, _ := strconv.Atoi(sid)
			todo := models.Todo{
				Id:       int64(id),
				Title:    title,
				PostDate: utils.NewTime().GetTimeToStr(time.Now().Unix()),
			}
			_, err := models.UpdateTodo(todo)
			if err == nil {
				http.Redirect(p.ResponseWriter, p.Request, "/", http.StatusFound)
			} else {
				p.Body = "数据库错误：" + fmt.Sprintf("%v", err)
				p.Template = "todo/error.html"
			}
		}
	}
}

func (p *PageTodo) Delete() {
	if sid, ok := p.GET["id"]; ok {
		id, _ := strconv.Atoi(sid)
		_, err := models.DeleteTodo(int64(id))
		if err == nil {
			http.Redirect(p.ResponseWriter, p.Request, "/", http.StatusFound)
		} else {
			p.Body = "数据库错误：" + fmt.Sprintf("%v", err)
			p.Template = "todo/error.html"
		}

	} else {
		p.Body = "错误的请求"
		p.Template = "todo/error.html"
	}
}

func (p *PageTodo) Finish() {
	if p.Request.Method == "GET" {
		id, _ := strconv.Atoi(p.GET["id"])
		status := p.GET["status"]
		if id > 0 && (status == "yes" || status == "no") {
			finished := 0
			if status == "yes" {
				finished = 1
			}

			postDate := utils.NewTime().GetTimeToStr(time.Now().Unix())
			var db = Middleware.Get("db").(*sql.DB)
			sql := "UPDATE todo SET finished = " + strconv.Itoa(finished) + ", post_date = \"" + postDate + "\" WHERE id = " + strconv.Itoa(id)
			_, err := db.Exec(sql)
			if err == nil {
				http.Redirect(p.ResponseWriter, p.Request, "/", http.StatusFound)
			} else {
				p.Body = "数据库错误：" + fmt.Sprintf("%v", err)
				p.Template = "todo/error.html"
			}
		} else {
			p.Body = "错误的请求"
			p.Template = "todo/error.html"
		}
	}
}

func (p *PageTodo) Newtodo() {
	models.InitTodo()
	http.Redirect(p.ResponseWriter, p.Request, "/", http.StatusFound)
}
