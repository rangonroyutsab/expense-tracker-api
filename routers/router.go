package routers

import (
	"expense-tracker-api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/health",
			beego.NSInclude(
				&controllers.HealthController{},
			),
			beego.NSRouter("", &controllers.HealthController{}),
		),

		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.AuthController{},
			),
			beego.NSRouter("/register", &controllers.AuthController{}, "post:Register"),
			beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
		),

		beego.NSNamespace("/expenses",
			beego.NSInclude(
				&controllers.ExpenseController{},
			),
			beego.NSRouter("", &controllers.ExpenseController{}, "post:CreateExpense;get:ListExpenses"),
			beego.NSRouter("/summary", &controllers.ExpenseController{}, "get:Summary"),
			beego.NSRouter("/:id", &controllers.ExpenseController{}, "get:GetExpense;put:UpdateExpense;delete:DeleteExpense"),
		),
	)

	beego.AddNamespace(ns)
}
