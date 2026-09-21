package main

import (
	_ "rental-property-api/routers"

	logs "github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	logs.Info("Server running at http://localhost:8080...")
	beego.Run()
}
