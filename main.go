package main

import (
	"expense-tracker-api/models"
	_ "expense-tracker-api/routers"
	"expense-tracker-api/utils"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func main() {
	if err := utils.InitSentry(); err != nil {
		logs.Error("failed to initialize Sentry: %v", err)
	}
	defer utils.FlushSentry()

	beego.BConfig.RecoverFunc = func(ctx *context.Context, cfg *beego.Config) {
		if r := recover(); r != nil {
			utils.CapturePanicValue(r)

			ctx.Output.SetStatus(500)
			ctx.Output.JSON(map[string]interface{}{
				"success": false,
				"message": "Internal server error",
			}, false, false)
		}
	}

	models.InitPaths()

	logs.Info("starting expense-tracker-api")
	beego.Run()
}
