package models

// Expense represents a user expense stored in CSV.
type Expense struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	Title       string  `json:"title"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Note        string  `json:"note"`
	ExpenseDate string  `json:"expense_date"`
	CreatedAt   string  `json:"created_at"`
}

// ExpensesCSVPath stores the path to the expenses CSV file.
var ExpensesCSVPath = "data/expenses.csv"
