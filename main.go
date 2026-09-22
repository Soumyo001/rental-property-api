package main

import (
	"os"
	_ "rental-property-api/routers"
	"rental-property-api/services"

	logs "github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	datapath := beego.AppConfig.DefaultString("datapath", "data/rental_properties.json")
	if err := services.Init(datapath); err != nil {
		logs.Critical("could not load data from %s: %v", datapath, err)
		os.Exit(1)
	}

	logs.Info("Loaded %d properties from file %s", services.GetStore().Count(), datapath)

	logs.Info("Server running at http://localhost:8080...")
	beego.Run()
}
