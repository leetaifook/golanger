package main

import (
	"code.google.com/p/go.net/websocket"
	"controllers"
	"net/http"
)

func main() {
	http.Handle("/chat", websocket.Handler(controllers.BuildConnection))

	go controllers.InitChatRoom()
	startApp()
}
