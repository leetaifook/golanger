package main

import (
	"fmt"
	. "golanger/middleware"
	"golanger/utils"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()*2 + 1)

	sqlite, err := utils.NewSqlite("./data/todo.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	defer sqlite.Close()
	Middleware.Add("db", sqlite)

	startApp()
}
