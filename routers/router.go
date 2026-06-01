// @APIVersion 1.0.0
// @Title Expense Tracker API
// @Description A REST API for user authentication and expense tracking using Beego and CSV storage.
// @Contact Rangon Roy Utsab
// @ContactEmail rangonroy@outlook.com
// @License MIT
// @LicenseUrl https://opensource.org/licenses/MIT

package routers

import (
	"expense-tracker-api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/health", &controllers.HealthController{}),

		beego.NSNamespace("/auth",
			beego.NSRouter("/register", &controllers.AuthController{}, "post:Register"),
			beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
		),
	)

	beego.AddNamespace(ns)
}
