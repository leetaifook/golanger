package models

import (
	"golanger/utils"
)

var (
	ColUser = utils.M{
		"name":  "user",
		"index": []string{"email", "name", "password", "status", "createtime"},
	}
)

/*
用户表
user
{
    "email": <email>,
    "name" : <name>,
    "password" : <password>,
    "status" : <status>,
    "createtime" : <createtime>,
    "updatetime" : <updatetime>
}
*/
type ModelUser struct {
	Email      string
	Name       string
	Password   string
	Status     byte
	Createtime int64
	Updatetime int64
}
