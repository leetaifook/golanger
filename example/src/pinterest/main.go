package main

import (
	"./controllers"
	"fmt"
	. "golanger/database/activerecord"
	. "golanger/middleware"
	"golanger/utils"
	"os"
)

func main() {
	if sqliteDns, ok := controllers.Page.Config.Database["Sqlite"]; ok && sqliteDns != "" {
		sqlite, err := utils.NewSqlite(sqliteDns)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		defer sqlite.Close()
		orm := NewActiveRecord(sqlite)
		Middleware.Add("orm", orm)
		Middleware.Add("db", sqlite)
	}

	startApp()
}
