package main

import (
	"controllers"
	. "golanger/middleware"
	"golanger/utils"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()*2 + 1)

	if mongoDns, ok := controllers.Page.Config.Database["MongoDB"]; ok && mongoDns != "" {
		mgoServer := utils.NewMongo(mongoDns)
		defer mgoServer.Close()
		Middleware.Add("db", mgoServer)
	}

	startApp()
}
