package controllers

import (
	. "../models"
	"encoding/json"
	"fmt"
	. "golanger/middleware"
	"golanger/utils"
	"time"
)

type PageModule struct {
	*App
}

func init() {
	Page.RegisterController("module/", &PageModule{Page})
}

func (p *PageModule) Index() {
	mgoServer := Middleware.Get("db").(*utils.Mongo)
	colQuerier := utils.M{}
	modules := []ModelModule{}
	err := mgoServer.C(ColModule).Find(colQuerier).All(&modules)
	if err != nil {
		fmt.Println("获取模块信息失败！")
		return
	}
	results := []utils.M{}
	for _, module := range modules {
		result := utils.M{}
		result["name"] = module.Name
		result["createtime"] = time.Unix(module.Createtime, 0).String()[0:19]
		result["updatetime"] = time.Unix(module.Updatetime, 0).String()[0:19]
		if module.Status == 1 {
			result["status"] = "停用"
		} else if module.Status == 0 {
			result["status"] = "启用"
		}
		results = append(results, result)
	}
	p.Body = results
}

// Create Module
func (p *PageModule) CreateModule() {
	if p.Request.Method == "POST" {
		if _, ok := p.POST["ajax"]; ok {
			m := utils.M{
				"status":  "1",
				"message": "",
			}
			p.Hide = true
			mgoServer := Middleware.Get("db").(*utils.Mongo)

			colQuerier := utils.M{}
			modules := []ModelModule{}
			err := mgoServer.C(ColModule).Find(colQuerier).All(&modules)
			if err != nil {
				m["status"] = 0
				m["message"] = "获取用户信息失败！"
				return
			}
			modulename := p.POST["modulename"]

			for _, module := range modules {
				if modulename == module.Name {
					m["status"] = 0
					m["message"] = "该模块已存在"
					ret, _ := json.Marshal(m)
					p.ResponseWriter.Write(ret)
					return
				}
			}
			tnow := time.Now()
			mgoServer.C(ColModule).Insert(&ModelModule{
				Name:       modulename,
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

// Delete Module
func (p *PageModule) DeleteModule() {
	if p.Request.Method == "POST" {
		if _, ok := p.POST["ajax"]; ok {
			p.Hide = true
			mgoServer := Middleware.Get("db").(*utils.Mongo)
			modulename := p.POST["modulename"]
			colQuerier, m := utils.M{"name": modulename}, utils.M{"message": ""}
			err := mgoServer.C(ColModule).Remove(colQuerier)
			if err != nil {
				m["message"] = "删除模块失败！"
			} else {
				m["message"] = "成功删除模块！"
			}
			ret, _ := json.Marshal(m)
			p.ResponseWriter.Write(ret)
		}
	}
}

// Stop Module
func (p *PageModule) StopModule() {
	if p.Request.Method == "POST" {
		if _, ok := p.POST["ajax"]; ok {
			p.Hide = true
			m := utils.M{"message": ""}
			if modulename, ok := p.POST["modulename"]; ok {
				mgoServer := Middleware.Get("db").(*utils.Mongo)
				colQuerier := utils.M{"name": modulename}

				result := utils.M{}
				err := mgoServer.C(ColModule).Find(colQuerier).One(&result)
				if err != nil {
					m["message"] = "模块不存在"
				}

				stop := utils.M{}
				if result["status"] == 1 {
					stop["$set"] = utils.M{"status": 0}
				} else {
					stop["$set"] = utils.M{"status": 1}
				}

				err = mgoServer.C(ColModule).Update(colQuerier, stop)
				if err != nil {
					m["message"] = "停用模块失败！"
				} else {
					m["message"] = "成功停用模块！"
				}
			} else {
				m["message"] = "请选择要停用的模块！"
			}
			ret, _ := json.Marshal(m)
			p.ResponseWriter.Write(ret)
		}
	}
}
