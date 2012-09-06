package main

import (
	. "controllers"
	"flag"
	"fmt"
	. "golanger/middleware"
	"golanger/utils"
	"os"
	"path/filepath"
	"runtime"
	_ "templateFunc"
)

var (
	addr       = flag.String("addr", ":80", "Server port")
	configPath = flag.String("config", "./config/site.yaml", "site filepath of config")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()*2 + 1)

	flag.Parse()
	os.Chdir(filepath.Dir(os.Args[0]))
	fmt.Println("Listen server address: " + *addr)
	fmt.Println("Read configuration file success, fithpath: " + filepath.Join(filepath.Dir(os.Args[0]), *configPath))

	App.Load(*configPath)

	if mongoDns, ok := App.Database["MongoDB"]; ok && mongoDns != "" {
		mgoServer := utils.NewMongo(mongoDns)
		defer mgoServer.Close()
		Middleware.Add("db", mgoServer)
	}

	App.HandleFavicon()
	App.HandleStatic()
	App.ListenAndServe(*addr)
}
