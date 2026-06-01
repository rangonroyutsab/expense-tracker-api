package models

import beego "github.com/beego/beego/v2/server/web"

// InitPaths loads CSV file paths from Beego config.
// Dynamic path offers test isolation.
func InitPaths() {
	if usersPath, err := beego.AppConfig.String("users_csv"); err == nil && usersPath != "" {
		UsersCSVPath = usersPath
	}

	if expensesPath, err := beego.AppConfig.String("expenses_csv"); err == nil && expensesPath != "" {
		ExpensesCSVPath = expensesPath
	}
}
