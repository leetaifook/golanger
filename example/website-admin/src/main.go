package main

import (
	"controllers"
	. "golanger/middleware"
	"golanger/utils"
)

func main() {
	if mongoDns, ok := controllers.Page.Config.Database["MongoDB"]; ok && mongoDns != "" {
		mgoServer := utils.NewMongo(mongoDns)
		defer mgoServer.Close()
		Middleware.Add("db", mgoServer)
	}

	startApp()
}
