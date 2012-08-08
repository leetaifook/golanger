package main

import (
	"./controllers"
	"golanger/utils"
)

var (
	mgoServer *utils.Mongo
)

func main() {
	if mongoDns, ok := controllers.Page.Config.Database["MongoDB"]; ok && mongoDns != "" {
		mgoServer = utils.NewMongo(mongoDns)
		defer mgoServer.Close()
	}

	startApp()
}
