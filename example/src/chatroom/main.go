package main

import (
	"./add-on/code.google.com/p/go.net/websocket"
	"./controllers"
	"golanger/utils"
	"net/http"
)

var (
	mgoServer *utils.Mongo
)

func main() {
	if mongoDns, ok := controllers.Page.Config.Database["MongoDB"]; ok && mongoDns != "" {
		mgoServer = utils.NewMongo(mongoDns)
		defer mgoServer.Close()
	}

	http.Handle("/chat", websocket.Handler(controllers.BuildConnection))

	go controllers.InitChatRoom()
	startApp()
}
