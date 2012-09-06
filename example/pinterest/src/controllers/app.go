package controllers

import (
	"golanger/web"
)

type Application struct {
	*web.Page
}

func (a *Application) Init() {
	a.Page.Init()
}

var App = &Application{
	Page: web.NewPage(web.PageParam{}),
}
