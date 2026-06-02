package main

import (
	"expense-tracker-api/models"
	_ "expense-tracker-api/routers"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {

	models.InitPaths()

	if err := models.EnsureUserFile(); err != nil {
		logs.Error("failed to initialize users CSV: %v", err)
		return
	}

	if err := models.EnsureExpenseFile(); err != nil {
		logs.Error("failed to initialize expenses CSV: %v", err)
		return
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	logs.Info("starting expense-tracker-api")
	beego.Run()
}
