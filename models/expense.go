package models

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

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

// ExpensesCSVPath stores the configurable expenses CSV file path.
var ExpensesCSVPath = "data/expenses.csv"

// AllowedCategories contains valid expense categories.
var AllowedCategories = []string{
	"Food",
	"Transport",
	"Housing",
	"Entertainment",
	"Shopping",
	"Healthcare",
	"Education",
	"Utilities",
	"Other",
}

var expenseCSVHeader = []string{
	"id",
	"user_id",
	"title",
	"amount",
	"category",
	"note",
	"expense_date",
	"created_at",
}

// EnsureExpenseFile creates the expenses CSV file with a header if needed.
func EnsureExpenseFile() error {
	if err := os.MkdirAll(filepath.Dir(ExpensesCSVPath), 0755); err != nil {
		return err
	}

	fileInfo, err := os.Stat(ExpensesCSVPath)
	if err == nil && fileInfo.Size() > 0 {
		return nil
	}

	file, err := os.Create(ExpensesCSVPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if err := writer.Write(expenseCSVHeader); err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}

// GetExpensesByUserID returns all expenses belonging to a user.
func GetExpensesByUserID(userID int) ([]Expense, error) {
	expenses, err := getAllExpenses()
	if err != nil {
		return nil, err
	}

	userExpenses := make([]Expense, 0)

	for _, expense := range expenses {
		if expense.UserID == userID {
			userExpenses = append(userExpenses, expense)
		}
	}

	return userExpenses, nil
}

// GetExpenseByID returns one expense by expense ID and user ID.
func GetExpenseByID(id int, userID int) (*Expense, error) {
	expenses, err := GetExpensesByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, expense := range expenses {
		if expense.ID == id {
			return &expense, nil
		}
	}

	return nil, errors.New("expense not found")
}

// CreateExpense creates a new expense in the expenses CSV file.
func CreateExpense(expense *Expense) error {
	id, err := getNextExpenseID()
	if err != nil {
		return err
	}

	expense.ID = id
	expense.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	file, err := os.OpenFile(ExpensesCSVPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if err := writer.Write([]string{
		strconv.Itoa(expense.ID),
		strconv.Itoa(expense.UserID),
		expense.Title,
		strconv.FormatFloat(expense.Amount, 'f', 2, 64),
		expense.Category,
		expense.Note,
		expense.ExpenseDate,
		expense.CreatedAt,
	}); err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}

// UpdateExpense updates an existing expense owned by the user.
func UpdateExpense(updatedExpense *Expense) error {
	expenses, err := getAllExpenses()
	if err != nil {
		return err
	}

	found := false

	for index, expense := range expenses {
		if expense.ID == updatedExpense.ID && expense.UserID == updatedExpense.UserID {
			updatedExpense.CreatedAt = expense.CreatedAt
			expenses[index] = *updatedExpense
			found = true
			break
		}
	}

	if !found {
		return errors.New("expense not found")
	}

	return writeAllExpenses(expenses)
}

// DeleteExpense deletes an expense owned by the user.
func DeleteExpense(id int, userID int) error {
	expenses, err := getAllExpenses()
	if err != nil {
		return err
	}

	updatedExpenses := make([]Expense, 0, len(expenses))
	found := false

	for _, expense := range expenses {
		if expense.ID == id && expense.UserID == userID {
			found = true
			continue
		}

		updatedExpenses = append(updatedExpenses, expense)
	}

	if !found {
		return errors.New("expense not found")
	}

	return writeAllExpenses(updatedExpenses)
}

// IsAllowedCategory checks whether a category is valid.
func IsAllowedCategory(category string) bool {
	for _, allowedCategory := range AllowedCategories {
		if category == allowedCategory {
			return true
		}
	}

	return false
}

func getAllExpenses() ([]Expense, error) {
	file, err := os.Open(ExpensesCSVPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	expenses := make([]Expense, 0)

	for index, record := range records {
		if index == 0 {
			continue
		}

		if len(record) != 8 {
			return nil, errors.New("invalid expenses CSV row")
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		userID, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}

		amount, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, Expense{
			ID:          id,
			UserID:      userID,
			Title:       record[2],
			Amount:      amount,
			Category:    record[4],
			Note:        record[5],
			ExpenseDate: record[6],
			CreatedAt:   record[7],
		})
	}

	return expenses, nil
}

func getNextExpenseID() (int, error) {
	expenses, err := getAllExpenses()
	if err != nil {
		return 0, err
	}

	maxID := 0

	for _, expense := range expenses {
		if expense.ID > maxID {
			maxID = expense.ID
		}
	}

	return maxID + 1, nil
}

func writeAllExpenses(expenses []Expense) error {
	if err := os.MkdirAll(filepath.Dir(ExpensesCSVPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(ExpensesCSVPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if err := writer.Write(expenseCSVHeader); err != nil {
		return err
	}

	for _, expense := range expenses {
		if err := writer.Write([]string{
			strconv.Itoa(expense.ID),
			strconv.Itoa(expense.UserID),
			expense.Title,
			strconv.FormatFloat(expense.Amount, 'f', 2, 64),
			expense.Category,
			expense.Note,
			expense.ExpenseDate,
			expense.CreatedAt,
		}); err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}

// CategorySummary represents total spending for a category.
type CategorySummary struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
	Count    int     `json:"count"`
}

// ExpenseSummary represents expense summary response data.
type ExpenseSummary struct {
	DateFrom    string            `json:"date_from,omitempty"`
	DateTo      string            `json:"date_to,omitempty"`
	TotalAmount float64           `json:"total_amount"`
	TotalCount  int               `json:"total_count"`
	ByCategory  []CategorySummary `json:"by_category"`
}

// FilterExpensesByDate filters expenses by optional date range.
func FilterExpensesByDate(expenses []Expense, dateFrom, dateTo string) ([]Expense, error) {
	if dateFrom != "" {
		if _, err := time.Parse("2006-01-02", dateFrom); err != nil {
			return nil, errors.New("date_from must be valid YYYY-MM-DD")
		}
	}

	if dateTo != "" {
		if _, err := time.Parse("2006-01-02", dateTo); err != nil {
			return nil, errors.New("date_to must be valid YYYY-MM-DD")
		}
	}

	filteredExpenses := make([]Expense, 0)

	for _, expense := range expenses {
		if _, err := time.Parse("2006-01-02", expense.ExpenseDate); err != nil {
			return nil, errors.New("expense_date must be valid YYYY-MM-DD")
		}

		if dateFrom != "" && expense.ExpenseDate < dateFrom {
			continue
		}

		if dateTo != "" && expense.ExpenseDate > dateTo {
			continue
		}

		filteredExpenses = append(filteredExpenses, expense)
	}

	return filteredExpenses, nil
}

// SortExpenses sorts expenses by amount or expense_date.
func SortExpenses(expenses []Expense, sortBy, sortOrder string) error {
	if sortOrder != "" && sortBy == "" {
		return errors.New("sort_by is required when sort_order is provided")
	}

	if sortBy == "" {
		return nil
	}

	if sortBy != "amount" && sortBy != "expense_date" {
		return errors.New("sort_by must be amount or expense_date")
	}

	if sortOrder == "" {
		sortOrder = "desc"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		return errors.New("sort_order must be asc or desc")
	}

	sort.Slice(expenses, func(i, j int) bool {
		if sortBy == "amount" {
			if sortOrder == "asc" {
				return expenses[i].Amount < expenses[j].Amount
			}

			return expenses[i].Amount > expenses[j].Amount
		}

		if sortOrder == "asc" {
			return expenses[i].ExpenseDate < expenses[j].ExpenseDate
		}

		return expenses[i].ExpenseDate > expenses[j].ExpenseDate
	})

	return nil
}

// BuildExpenseSummary builds a spending summary from expenses.
func BuildExpenseSummary(expenses []Expense, dateFrom, dateTo string) ExpenseSummary {
	summary := ExpenseSummary{
		DateFrom:   dateFrom,
		DateTo:     dateTo,
		ByCategory: make([]CategorySummary, 0),
	}

	categoryMap := make(map[string]*CategorySummary)

	for _, expense := range expenses {
		summary.TotalAmount += expense.Amount
		summary.TotalCount++

		if _, exists := categoryMap[expense.Category]; !exists {
			categoryMap[expense.Category] = &CategorySummary{
				Category: expense.Category,
			}
		}

		categoryMap[expense.Category].Total += expense.Amount
		categoryMap[expense.Category].Count++
	}

	for _, categorySummary := range categoryMap {
		summary.ByCategory = append(summary.ByCategory, *categorySummary)
	}

	sort.Slice(summary.ByCategory, func(i, j int) bool {
		return summary.ByCategory[i].Category < summary.ByCategory[j].Category
	})

	return summary
}
