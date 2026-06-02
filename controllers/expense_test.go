package controllers_test

import (
	"net/http"
	"path/filepath"
	"testing"

	appcontrollers "expense-tracker-api/controllers"
	"expense-tracker-api/models"
)

func setupExpenseControllerTest(t *testing.T) {
	t.Helper()

	tempDir := t.TempDir()

	models.UsersCSVPath = filepath.Join(tempDir, "users.csv")
	models.ExpensesCSVPath = filepath.Join(tempDir, "expenses.csv")

	if err := models.EnsureUserFile(); err != nil {
		t.Fatalf("EnsureUserFile() error = %v", err)
	}

	if err := models.EnsureExpenseFile(); err != nil {
		t.Fatalf("EnsureExpenseFile() error = %v", err)
	}

	user := &models.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret123",
	}

	if err := models.CreateUser(user); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
}

func TestCreateExpenseController(t *testing.T) {
	tests := []struct {
		name       string
		body       appcontrollers.CreateExpenseRequest
		headers    map[string]string
		wantStatus int
	}{
		{
			name: "successful create",
			body: appcontrollers.CreateExpenseRequest{
				Title:       "Lunch",
				Amount:      350.50,
				Category:    "Food",
				Note:        "Team lunch",
				ExpenseDate: "2025-06-10",
			},
			headers:    map[string]string{"X-User-ID": "1"},
			wantStatus: http.StatusCreated,
		},
		{
			name: "missing auth",
			body: appcontrollers.CreateExpenseRequest{
				Title:       "Lunch",
				Amount:      350.50,
				Category:    "Food",
				ExpenseDate: "2025-06-10",
			},
			headers:    nil,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "missing title",
			body: appcontrollers.CreateExpenseRequest{
				Amount:      350.50,
				Category:    "Food",
				ExpenseDate: "2025-06-10",
			},
			headers:    map[string]string{"X-User-ID": "1"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "negative amount",
			body: appcontrollers.CreateExpenseRequest{
				Title:       "Lunch",
				Amount:      -10,
				Category:    "Food",
				ExpenseDate: "2025-06-10",
			},
			headers:    map[string]string{"X-User-ID": "1"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid category",
			body: appcontrollers.CreateExpenseRequest{
				Title:       "Lunch",
				Amount:      350.50,
				Category:    "Invalid",
				ExpenseDate: "2025-06-10",
			},
			headers:    map[string]string{"X-User-ID": "1"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid date",
			body: appcontrollers.CreateExpenseRequest{
				Title:       "Lunch",
				Amount:      350.50,
				Category:    "Food",
				ExpenseDate: "bad-date",
			},
			headers:    map[string]string{"X-User-ID": "1"},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			setupExpenseControllerTest(t)

			recorder := performJSONRequest(
				t,
				http.MethodPost,
				"/api/v1/expenses",
				testCase.body,
				testCase.headers,
			)

			if recorder.Code != testCase.wantStatus {
				t.Fatalf(
					"status = %d, want %d, body = %s",
					recorder.Code,
					testCase.wantStatus,
					recorder.Body.String(),
				)
			}
		})
	}
}

func TestExpenseControllerCRUD(t *testing.T) {
	setupExpenseControllerTest(t)

	headers := map[string]string{"X-User-ID": "1"}

	createBody := appcontrollers.CreateExpenseRequest{
		Title:       "Lunch",
		Amount:      350.50,
		Category:    "Food",
		Note:        "Team lunch",
		ExpenseDate: "2025-06-10",
	}

	createResponse := performJSONRequest(t, http.MethodPost, "/api/v1/expenses", createBody, headers)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d, body = %s", createResponse.Code, http.StatusCreated, createResponse.Body.String())
	}

	listResponse := performJSONRequest(t, http.MethodGet, "/api/v1/expenses", nil, headers)
	if listResponse.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d, body = %s", listResponse.Code, http.StatusOK, listResponse.Body.String())
	}

	getResponse := performJSONRequest(t, http.MethodGet, "/api/v1/expenses/1", nil, headers)
	if getResponse.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d, body = %s", getResponse.Code, http.StatusOK, getResponse.Body.String())
	}

	updateBody := appcontrollers.UpdateExpenseRequest{
		Title:       "Dinner",
		Amount:      500,
		Category:    "Food",
		Note:        "Dinner with friends",
		ExpenseDate: "2025-06-11",
	}

	updateResponse := performJSONRequest(t, http.MethodPut, "/api/v1/expenses/1", updateBody, headers)
	if updateResponse.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d, body = %s", updateResponse.Code, http.StatusOK, updateResponse.Body.String())
	}

	deleteResponse := performJSONRequest(t, http.MethodDelete, "/api/v1/expenses/1", nil, headers)
	if deleteResponse.Code != http.StatusOK {
		t.Fatalf("delete status = %d, want %d, body = %s", deleteResponse.Code, http.StatusOK, deleteResponse.Body.String())
	}

	getDeletedResponse := performJSONRequest(t, http.MethodGet, "/api/v1/expenses/1", nil, headers)
	if getDeletedResponse.Code != http.StatusNotFound {
		t.Fatalf("get deleted status = %d, want %d, body = %s", getDeletedResponse.Code, http.StatusNotFound, getDeletedResponse.Body.String())
	}
}

func TestExpenseListQueryValidation(t *testing.T) {
	setupExpenseControllerTest(t)

	headers := map[string]string{"X-User-ID": "1"}

	expenses := []appcontrollers.CreateExpenseRequest{
		{
			Title:       "Lunch",
			Amount:      350.50,
			Category:    "Food",
			ExpenseDate: "2025-06-10",
		},
		{
			Title:       "Bus",
			Amount:      50,
			Category:    "Transport",
			ExpenseDate: "2025-06-11",
		},
	}

	for _, expense := range expenses {
		response := performJSONRequest(t, http.MethodPost, "/api/v1/expenses", expense, headers)
		if response.Code != http.StatusCreated {
			t.Fatalf("create status = %d, want %d, body = %s", response.Code, http.StatusCreated, response.Body.String())
		}
	}

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "date filter",
			path:       "/api/v1/expenses?date_from=2025-06-01&date_to=2025-06-30",
			wantStatus: http.StatusOK,
		},
		{
			name:       "sort amount desc",
			path:       "/api/v1/expenses?sort_by=amount&sort_order=desc",
			wantStatus: http.StatusOK,
		},
		{
			name:       "sort date asc",
			path:       "/api/v1/expenses?sort_by=expense_date&sort_order=asc",
			wantStatus: http.StatusOK,
		},
		{
			name:       "sort order without sort by",
			path:       "/api/v1/expenses?sort_order=asc",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid sort by",
			path:       "/api/v1/expenses?sort_by=category",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid date",
			path:       "/api/v1/expenses?date_from=invalid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := performJSONRequest(t, http.MethodGet, testCase.path, nil, headers)

			if recorder.Code != testCase.wantStatus {
				t.Fatalf(
					"status = %d, want %d, body = %s",
					recorder.Code,
					testCase.wantStatus,
					recorder.Body.String(),
				)
			}
		})
	}
}

func TestExpenseSummaryController(t *testing.T) {
	setupExpenseControllerTest(t)

	headers := map[string]string{"X-User-ID": "1"}

	expense := appcontrollers.CreateExpenseRequest{
		Title:       "Lunch",
		Amount:      350.50,
		Category:    "Food",
		ExpenseDate: "2025-06-10",
	}

	createResponse := performJSONRequest(t, http.MethodPost, "/api/v1/expenses", expense, headers)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d, body = %s", createResponse.Code, http.StatusCreated, createResponse.Body.String())
	}

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "full summary",
			path:       "/api/v1/expenses/summary",
			wantStatus: http.StatusOK,
		},
		{
			name:       "date range summary",
			path:       "/api/v1/expenses/summary?date_from=2025-06-01&date_to=2025-06-30",
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing date_to",
			path:       "/api/v1/expenses/summary?date_from=2025-06-01",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid date",
			path:       "/api/v1/expenses/summary?date_from=invalid&date_to=2025-06-30",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := performJSONRequest(t, http.MethodGet, testCase.path, nil, headers)

			if recorder.Code != testCase.wantStatus {
				t.Fatalf(
					"status = %d, want %d, body = %s",
					recorder.Code,
					testCase.wantStatus,
					recorder.Body.String(),
				)
			}
		})
	}
}

func TestExpenseControllerInvalidUserID(t *testing.T) {
	setupExpenseControllerTest(t)

	body := appcontrollers.CreateExpenseRequest{
		Title:       "Lunch",
		Amount:      350.50,
		Category:    "Food",
		ExpenseDate: "2025-06-10",
	}

	tests := []struct {
		name       string
		headers    map[string]string
		wantStatus int
	}{
		{
			name:       "non existing user id",
			headers:    map[string]string{"X-User-ID": "999"},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid user id",
			headers:    map[string]string{"X-User-ID": "abc"},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "zero user id",
			headers:    map[string]string{"X-User-ID": "0"},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := performJSONRequest(
				t,
				http.MethodPost,
				"/api/v1/expenses",
				body,
				testCase.headers,
			)

			if recorder.Code != testCase.wantStatus {
				t.Fatalf(
					"status = %d, want %d, body = %s",
					recorder.Code,
					testCase.wantStatus,
					recorder.Body.String(),
				)
			}
		})
	}
}

func TestExpenseControllerInvalidExpenseID(t *testing.T) {
	setupExpenseControllerTest(t)

	headers := map[string]string{"X-User-ID": "1"}

	updateBody := appcontrollers.UpdateExpenseRequest{
		Title:       "Dinner",
		Amount:      500,
		Category:    "Food",
		ExpenseDate: "2025-06-11",
	}

	tests := []struct {
		name       string
		method     string
		path       string
		body       interface{}
		wantStatus int
	}{
		{
			name:       "get invalid expense id",
			method:     http.MethodGet,
			path:       "/api/v1/expenses/abc",
			body:       nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "put invalid expense id",
			method:     http.MethodPut,
			path:       "/api/v1/expenses/abc",
			body:       updateBody,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "delete invalid expense id",
			method:     http.MethodDelete,
			path:       "/api/v1/expenses/abc",
			body:       nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "get zero expense id",
			method:     http.MethodGet,
			path:       "/api/v1/expenses/0",
			body:       nil,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := performJSONRequest(
				t,
				testCase.method,
				testCase.path,
				testCase.body,
				headers,
			)

			if recorder.Code != testCase.wantStatus {
				t.Fatalf(
					"status = %d, want %d, body = %s",
					recorder.Code,
					testCase.wantStatus,
					recorder.Body.String(),
				)
			}
		})
	}
}
