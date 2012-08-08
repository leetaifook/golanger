package models

import (
	"golanger/utils"
)

var (
	ColModule = utils.M{
		"name":  "module",
		"index": []string{"name", "status", "createtime"},
	}
)

/*
模块表
module
{
    "name" : <name>,
    "status" : <status>,
    "createtime" : <createtime>,
    "updatetime" : <updatetime>
}
*/
type ModelModule struct {
	Name       string
	Status     byte
	Createtime int64
	Updatetime int64
}
