package controllers

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"expense-tracker-api/models"

	"github.com/beego/beego/v2/core/logs"
)

// ExpenseController handles expense CRUD operations.
type ExpenseController struct {
	BaseController
}

// CreateExpense creates a new expense.
// @Title Create Expense
// @Description Create a new expense for the authenticated user.
// @Param X-User-ID header int true "Authenticated user ID"
// @Param body body controllers.CreateExpenseRequest true "Create expense request body"
// @Success 201 {object} controllers.Response
// @Failure 400 {object} controllers.BasicResponse
// @Failure 401 {object} controllers.BasicResponse
// @Failure 500 {object} controllers.BasicResponse
// @router /expenses [post]
func (c *ExpenseController) CreateExpense() {
	userID, err := c.GetCurrentUserID()
	if err != nil {
		return
	}

	var request CreateExpenseRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		c.Error(400, "Invalid request body")
		return
	}

	if err := validateExpenseInput(request.Title, request.Amount, request.Category, request.ExpenseDate); err != nil {
		c.Error(400, err.Error())
		return
	}

	expense := &models.Expense{
		UserID:      userID,
		Title:       strings.TrimSpace(request.Title),
		Amount:      request.Amount,
		Category:    strings.TrimSpace(request.Category),
		Note:        strings.TrimSpace(request.Note),
		ExpenseDate: strings.TrimSpace(request.ExpenseDate),
	}

	if err := models.CreateExpense(expense); err != nil {
		logs.Error("failed to create expense: %v", err)

		c.Error(500, "Internal server error")
		return
	}

	c.Success(201, "Expense created successfully", expense)
}

// ListExpenses returns all expenses for the authenticated user.
// @Title List Expenses
// @Description List all expenses for the authenticated user with optional date filtering and sorting.
// @Param X-User-ID header int true "Authenticated user ID"
// @Param date_from query string false "Filter expenses on or after this date. Format: YYYY-MM-DD"
// @Param date_to query string false "Filter expenses on or before this date. Format: YYYY-MM-DD"
// @Param sort_by query string false "Sort field: amount or expense_date"
// @Param sort_order query string false "Sort order: asc or desc"
// @Success 200 {object} controllers.Response
// @Failure 400 {object} controllers.BasicResponse
// @Failure 401 {object} controllers.BasicResponse
// @Failure 500 {object} controllers.BasicResponse
// @router /expenses [get]
func (c *ExpenseController) ListExpenses() {
	userID, err := c.GetCurrentUserID()
	if err != nil {
		return
	}

	dateFrom := strings.TrimSpace(c.GetString("date_from"))
	dateTo := strings.TrimSpace(c.GetString("date_to"))
	sortBy := strings.TrimSpace(c.GetString("sort_by"))
	sortOrder := strings.TrimSpace(c.GetString("sort_order"))

	expenses, err := models.GetExpensesByUserID(userID)
	if err != nil {
		logs.Error("failed to list expenses: %v", err)

		c.Error(500, "Internal server error")
		return
	}

	expenses, err = models.FilterExpensesByDate(expenses, dateFrom, dateTo)
	if err != nil {
		c.Error(400, err.Error())
		return
	}

	if err := models.SortExpenses(expenses, sortBy, sortOrder); err != nil {
		c.Error(400, err.Error())
		return
	}

	c.Success(200, "Expenses retrieved", expenses)
}

// GetExpense returns one expense.
// @Title Get Expense
// @Description Get one expense owned by the authenticated user.
// @Param X-User-ID header int true "Authenticated user ID"
// @Param id path int true "Expense ID"
// @Success 200 {object} controllers.Response
// @Failure 400 {object} controllers.BasicResponse
// @Failure 401 {object} controllers.BasicResponse
// @Failure 404 {object} controllers.BasicResponse
// @router /expenses/:id [get]
func (c *ExpenseController) GetExpense() {
	userID, err := c.GetCurrentUserID()
	if err != nil {
		return
	}

	expenseID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil || expenseID <= 0 {
		c.Error(400, "Invalid expense ID")
		return
	}

	expense, err := models.GetExpenseByID(expenseID, userID)
	if err != nil {
		c.Error(404, "Expense not found")
		return
	}

	c.Success(200, "Expense retrieved", expense)
}

// UpdateExpense updates one expense.
// @Title Update Expense
// @Description Update one expense owned by the authenticated user.
// @Param X-User-ID header int true "Authenticated user ID"
// @Param id path int true "Expense ID"
// @Param body body controllers.UpdateExpenseRequest true "Update expense request body"
// @Success 200 {object} controllers.Response
// @Failure 400 {object} controllers.BasicResponse
// @Failure 401 {object} controllers.BasicResponse
// @Failure 404 {object} controllers.BasicResponse
// @Failure 500 {object} controllers.BasicResponse
// @router /expenses/:id [put]
func (c *ExpenseController) UpdateExpense() {
	userID, err := c.GetCurrentUserID()
	if err != nil {
		return
	}

	expenseID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil || expenseID <= 0 {
		c.Error(400, "Invalid expense ID")
		return
	}

	var request UpdateExpenseRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		c.Error(400, "Invalid request body")
		return
	}

	if err := validateExpenseInput(request.Title, request.Amount, request.Category, request.ExpenseDate); err != nil {
		c.Error(400, err.Error())
		return
	}

	expense := &models.Expense{
		ID:          expenseID,
		UserID:      userID,
		Title:       strings.TrimSpace(request.Title),
		Amount:      request.Amount,
		Category:    strings.TrimSpace(request.Category),
		Note:        strings.TrimSpace(request.Note),
		ExpenseDate: strings.TrimSpace(request.ExpenseDate),
	}

	if err := models.UpdateExpense(expense); err != nil {
		if err.Error() == "expense not found" {
			c.Error(404, "Expense not found")
			return
		}

		logs.Error("failed to update expense: %v", err)

		c.Error(500, "Internal server error")
		return
	}

	c.Success(200, "Expense updated successfully", expense)
}

// DeleteExpense deletes one expense.
// @Title Delete Expense
// @Description Delete one expense owned by the authenticated user.
// @Param X-User-ID header int true "Authenticated user ID"
// @Param id path int true "Expense ID"
// @Success 200 {object} controllers.BasicResponse
// @Failure 400 {object} controllers.BasicResponse
// @Failure 401 {object} controllers.BasicResponse
// @Failure 404 {object} controllers.BasicResponse
// @Failure 500 {object} controllers.BasicResponse
// @router /expenses/:id [delete]
func (c *ExpenseController) DeleteExpense() {
	userID, err := c.GetCurrentUserID()
	if err != nil {
		return
	}

	expenseID, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil || expenseID <= 0 {
		c.Error(400, "Invalid expense ID")
		return
	}

	if err := models.DeleteExpense(expenseID, userID); err != nil {
		if err.Error() == "expense not found" {
			c.Error(404, "Expense not found")
			return
		}

		logs.Error("failed to delete expense: %v", err)

		c.Error(500, "Internal server error")
		return
	}

	c.Success(200, "Expense deleted successfully", nil)
}

// Summary returns spending summary for the authenticated user.
// @Title Expense Summary
// @Description Generate a full or date-range spending summary.
// @Param X-User-ID header int true "Authenticated user ID"
// @Param date_from query string false "Summary start date. Format: YYYY-MM-DD"
// @Param date_to query string false "Summary end date. Format: YYYY-MM-DD"
// @Success 200 {object} controllers.Response
// @Failure 400 {object} controllers.BasicResponse
// @Failure 401 {object} controllers.BasicResponse
// @Failure 500 {object} controllers.BasicResponse
// @router /expenses/summary [get]
func (c *ExpenseController) Summary() {
	userID, err := c.GetCurrentUserID()
	if err != nil {
		return
	}

	dateFrom := strings.TrimSpace(c.GetString("date_from"))
	dateTo := strings.TrimSpace(c.GetString("date_to"))

	if (dateFrom == "" && dateTo != "") || (dateFrom != "" && dateTo == "") {
		c.Error(400, "Both date_from and date_to are required for date range summary")
		return
	}

	expenses, err := models.GetExpensesByUserID(userID)
	if err != nil {
		logs.Error("failed to load expenses for summary: %v", err)

		c.Error(500, "Internal server error")
		return
	}

	if dateFrom != "" && dateTo != "" {
		expenses, err = models.FilterExpensesByDate(expenses, dateFrom, dateTo)
		if err != nil {
			c.Error(400, err.Error())
			return
		}
	}

	summary := models.BuildExpenseSummary(expenses, dateFrom, dateTo)

	c.Success(200, "Summary generated", summary)
}

func validateExpenseInput(title string, amount float64, category string, expenseDate string) error {
	title = strings.TrimSpace(title)
	category = strings.TrimSpace(category)
	expenseDate = strings.TrimSpace(expenseDate)

	if title == "" {
		return errMessage("Title is required")
	}

	if amount <= 0 {
		return errMessage("Amount must be positive")
	}

	if category == "" {
		return errMessage("Category is required")
	}

	if !models.IsAllowedCategory(category) {
		return errMessage("Category is invalid")
	}

	if expenseDate == "" {
		return errMessage("Expense date is required")
	}

	if _, err := time.Parse("2006-01-02", expenseDate); err != nil {
		return errMessage("Expense date must be valid YYYY-MM-DD")
	}

	return nil
}

type errMessage string

func (e errMessage) Error() string {
	return string(e)
}
