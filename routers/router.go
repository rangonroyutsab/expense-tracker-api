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

		beego.NSRouter("/expenses", &controllers.ExpenseController{}, "post:CreateExpense;get:ListExpenses"),
		beego.NSRouter("/expenses/summary", &controllers.ExpenseController{}, "get:Summary"),
		beego.NSRouter("/expenses/:id", &controllers.ExpenseController{}, "get:GetExpense;put:UpdateExpense;delete:DeleteExpense"),
	)

	beego.AddNamespace(ns)
}
