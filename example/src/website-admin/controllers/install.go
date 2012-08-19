package controllers

import (
	. "../models"
	. "golanger/middleware"
	"golanger/utils"
	"net/http"
	"time"
)

type PageInstall struct {
	*App
}

func init() {
	Page.RegisterController("install/", &PageInstall{Page})
}

func (p *PageInstall) Init() {
	p.OffLogin = true
	p.App.Init()
}

func (p *PageInstall) Index() {
	mgoServer := Middleware.Get("db").(*utils.Mongo)
	email := "root@admin.com"
	username := "root"
	password := utils.Strings("123456").Md5()
	tnow := time.Now()
	err := mgoServer.C(ColUser).Insert(&ModelUser{
		Email:       email,
		Name:        username,
		Password:    password,
		Status:      1,
		Create_time: tnow.Unix(),
		Update_time: tnow.Unix(),
	})

	if err != nil {
		p.ResponseWriter.Write([]byte("安装失败"))
	} else {
		p.ResponseWriter.Write([]byte("安装成功...<br/>用户名:root,密码:123456<br/>请修改目录config下的site.yaml文件,将权限控制配置项开启，如:\"CheckRight : true\""))

		http.Redirect(p.ResponseWriter, p.Request, "/login.html", http.StatusFound)
	}

	p.Close = true
}
