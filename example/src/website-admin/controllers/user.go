package controllers

import (
	. "../models"
	"encoding/json"
	"fmt"
	. "golanger/middleware"
	"golanger/utils"
	"time"
)

type PageUser struct {
	*App
}

func init() {
	Page.RegisterController("user/", &PageUser{Page})
}

// Show user list
func (p *PageUser) Index() {
	mgoServer := Middleware.Get("db").(*utils.Mongo)
	colQuerier := utils.M{}
	users := []ModelUser{}
	err := mgoServer.C(ColUser).Find(colQuerier).All(&users)
	if err != nil {
		fmt.Println("获取用户信息失败！")
		return
	}

	results := []utils.M{}
	for _, user := range users {
		result := utils.M{}
		result["email"] = user.Email
		result["name"] = user.Name
		result["createtime"] = time.Unix(user.Createtime, 0).String()[0:19]
		result["updatetime"] = time.Unix(user.Updatetime, 0).String()[0:19]
		if user.Status == 1 {
			result["status"] = "停用"
		} else if user.Status == 0 {
			result["status"] = "启用"
		}

		results = append(results, result)
	}

	p.Body = results
}

// Create User
func (p *PageUser) CreateUser() {
	if p.Request.Method == "POST" {
		if _, ok := p.POST["ajax"]; ok {
			mgoServer := Middleware.Get("db").(*utils.Mongo)

			m := utils.M{
				"status":  1,
				"message": "",
			}
			p.Hide = true
			colQuerier := utils.M{}
			users := []ModelUser{}
			err := mgoServer.C(ColUser).Find(colQuerier).All(&users)
			if err != nil {
				fmt.Println("获取用户信息失败！")
				return
			}

			email := p.POST["email"]
			username := p.POST["username"]
			for _, user := range users {
				if username == user.Name {
					m["status"] = 0
					m["message"] = "该用户名已存在"
					ret, _ := json.Marshal(m)
					p.ResponseWriter.Write(ret)
					return
				}
				if email == user.Email {
					m["status"] = -1
					m["message"] = "该邮箱已注册"
					ret, _ := json.Marshal(m)
					p.ResponseWriter.Write(ret)
					return
				}
			}

			password := utils.Strings(p.POST["password"]).Md5()
			tnow := time.Now()
			mgoServer.C(ColUser).Insert(&ModelUser{
				Email:      email,
				Name:       username,
				Password:   password,
				Status:     1,
				Createtime: tnow.Unix(),
				Updatetime: tnow.Unix(),
			})
			ret, _ := json.Marshal(m)
			p.ResponseWriter.Write(ret)

			return
		}
	}
}

// Create Root
func (p *PageUser) CreateRoot() {
	mgoServer := Middleware.Get("db").(*utils.Mongo)
	colQuerier := utils.M{}
	_, err := mgoServer.C(ColUser).RemoveAll(colQuerier)
	if err != nil {
		fmt.Println("删除用户失败！")
		return
	}

	email := "root@admin.com"
	username := "root"
	password := utils.Strings("123456").Md5()
	tnow := time.Now()
	err = mgoServer.C(ColUser).Insert(&ModelUser{
		Email:      email,
		Name:       username,
		Password:   password,
		Status:     1,
		Createtime: tnow.Unix(),
		Updatetime: tnow.Unix(),
	})

	if err != nil {
		fmt.Println("Root用户创建失败！")
	}
}

// Update User
func (p *PageUser) UpdateUser() {
	if p.Request.Method == "POST" {
		if _, ok := p.POST["ajax"]; ok {
			p.Hide = true
			mgoServer := Middleware.Get("db").(*utils.Mongo)

			username, ok := p.POST["username"]
			m := utils.M{
				"status":  "1",
				"message": "",
			}
			if !ok {
				m["message"] = "修改用户时必须指定用户名！"
				return
			}

			colQuerier := utils.M{}

			users := []ModelUser{}
			err := mgoServer.C(ColUser).Find(colQuerier).All(&users)
			if err != nil {
				fmt.Println("获取用户信息失败！")
				return
			}

			email, emailOk := p.POST["email"]
			password, passwordOk := p.POST["password"]

			for _, user := range users {
				if user.Name != username {
					if email == user.Email {
						m["status"] = -1
						m["message"] = "该邮箱已注册"
						ret, _ := json.Marshal(m)
						p.ResponseWriter.Write(ret)
						fmt.Println("该邮箱已注册")
						return
					}
				}
			}

			colQuerier = utils.M{"name": username}

			updateField, change := utils.M{}, utils.M{}

			if emailOk {
				updateField["email"] = email
			}
			if passwordOk {
				password = utils.Strings(password).Md5()
				updateField["password"] = password
			}
			if emailOk || passwordOk {
				updateField["updatetime"] = time.Now().Unix()
				change["$set"] = updateField
				err := mgoServer.C(ColUser).Update(colQuerier, change)
				if err != nil {
					m["message"] = "用户资料更新失败！"
				} else {
					m["message"] = "用户资料更新成功！"
				}
			} else {
				m["message"] = "请输入需要更新的资料！"
			}
			ret, _ := json.Marshal(m)
			p.ResponseWriter.Write(ret)

			return
		}
	}
}

// Delete User
func (p *PageUser) DeleteUser() {
	if p.Request.Method == "POST" {
		if _, ok := p.POST["ajax"]; ok {
			p.Hide = true
			mgoServer := Middleware.Get("db").(*utils.Mongo)
			username := p.POST["username"]
			colQuerier, m := utils.M{"name": username}, utils.M{"message": ""}
			err := mgoServer.C(ColUser).Remove(colQuerier)
			if err != nil {
				m["message"] = "删除用户失败！"
			} else {
				m["message"] = "成功删除用户！"
			}

			ret, _ := json.Marshal(m)
			p.ResponseWriter.Write(ret)
		}
	}
}

// Stop User
func (p *PageUser) StopUser() {
	if p.Request.Method == "POST" {
		if _, ok := p.POST["ajax"]; ok {
			p.Hide = true
			m := utils.M{"message": ""}
			if username, ok := p.POST["username"]; ok {
				mgoServer := Middleware.Get("db").(*utils.Mongo)
				colQuerier := utils.M{"name": username}

				result := utils.M{}
				err := mgoServer.C(ColUser).Find(colQuerier).One(&result)
				if err != nil {
					m["message"] = "用户不存在"
				}

				stop := utils.M{}
				if result["status"] == 1 {
					stop["$set"] = utils.M{"status": 0}
				} else {
					stop["$set"] = utils.M{"status": 1}
				}

				err = mgoServer.C(ColUser).Update(colQuerier, stop)
				if err != nil {
					m["message"] = "停用用户失败！"
				} else {
					m["message"] = "成功停用用户！"
				}
			} else {
				m["message"] = "请选择要停用的用户！"
			}

			ret, _ := json.Marshal(m)
			p.ResponseWriter.Write(ret)
		}
	}
}
