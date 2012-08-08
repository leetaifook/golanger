package models

import (
	. "golanger/database/activerecord"
	. "golanger/middleware"
)

type Images struct {
	Id         int64  `index:"PK" field:"id"`
	Name       string `field:"name"`
	Ext        string `field:"ext"`
	Path       string `field:"path"`
	CreateTime int64  `field:"create_time"`
}

func GetImagesLists(page ...int) (*[]Images, error) {
	var pg int
	if len(page) > 0 {
		pg = page[0] - 1
	}

	if pg < 1 {
		pg = 0
	}

	num := 20
	start := pg * num
	//OnDebug = true
	var orm = Middleware.Get("orm").(ActiveRecord)
	images := []Images{}
	err := orm.Limit(num, start).OrderBy("id DESC").FindAll(&images)

	return &images, err
}

func SaveImages(images Images) (int64, error) {
	var orm = Middleware.Get("orm").(ActiveRecord)
	resNum, err := orm.Save(&images)

	return resNum, err
}
