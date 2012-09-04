package main

import (
	"code.google.com/p/go.net/websocket"
	"controllers"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()*2 + 1)

	http.Handle("/chat", websocket.Handler(controllers.BuildConnection))

	go controllers.InitChatRoom()
	startApp()
}
