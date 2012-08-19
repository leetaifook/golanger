package controllers

import (
	. "../models"
	. "golanger/middleware"
	"golanger/utils"
	"io/ioutil"
	"net/http"
	"os"
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
	p.OffRight = true
	p.App.Init()
}

func (p *PageInstall) Index() {
	fileInstallLock := "./data/install.lock"

	if _, err := os.Stat(fileInstallLock); err == nil {
		p.ResponseWriter.Write([]byte("程序已经安装过，如需要重新安装，请删除data目录下的install.lock文件后重试"))
	} else {
		mgoServer := Middleware.Get("db").(*utils.Mongo)
		email := "root@admin.com"
		username := "root"
		password := utils.Strings("123456").Md5()
		tnow := time.Now()
		mgoServer.C(ColUser).Insert(&ModelUser{
			Email:       email,
			Name:        username,
			Password:    password,
			Status:      1,
			Create_time: tnow.Unix(),
			Update_time: tnow.Unix(),
		})

		mgoServer.C(ColModule).Insert(&ModelModule{
			Name:        "模块管理",
			Path:        "module/",
			Order:       0,
			Status:      1,
			Create_time: tnow.Unix(),
			Update_time: tnow.Unix(),
		})

		mgoServer.C(ColModule).Insert(&ModelModule{
			Name:        "用户管理",
			Path:        "user/",
			Order:       1,
			Status:      1,
			Create_time: tnow.Unix(),
			Update_time: tnow.Unix(),
		})

		mgoServer.C(ColModule).Insert(&ModelModule{
			Name:        "角色管理",
			Path:        "role/",
			Order:       2,
			Status:      1,
			Create_time: tnow.Unix(),
			Update_time: tnow.Unix(),
		})

		mgoServer.C(ColRole).Insert(&ModelRole{
			Name:   "超级管理员",
			Users:  []string{"root"},
			Status: 1,
			Right: utils.M{
				"scope":   "3",
				"modules": []utils.M{},
			},
			Create_time: tnow.Unix(),
			Update_time: tnow.Unix(),
		})

		ioutil.WriteFile(fileInstallLock, []byte("installed"), 0777)

		sessionSign := p.COOKIE[p.SessionName]
		if sessionSign != "" {
			p.ClearSession(p.COOKIE[p.SessionName])
		}

		p.ResponseWriter.Write([]byte("安装成功...<br/>用户名:root,密码:123456"))
	}

	http.Redirect(p.ResponseWriter, p.Request, "/login.html", http.StatusFound)

	p.Close = true
}
