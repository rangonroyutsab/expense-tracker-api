package routers_test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"expense-tracker-api/models"
	_ "expense-tracker-api/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func setupRouterTest(t *testing.T) {
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
		Name:     "Router Test User",
		Email:    "router@example.com",
		Password: "secret123",
	}

	if err := models.CreateUser(user); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
}

func performRouterRequest(method string, path string, userID string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, nil)

	if userID != "" {
		request.Header.Set("X-User-ID", userID)
	}

	recorder := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(recorder, request)

	return recorder
}

func TestRouterHealthRoute(t *testing.T) {
	setupRouterTest(t)

	recorder := performRouterRequest(http.MethodGet, "/api/v1/health", "")

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
}

func TestRouterRegisteredRoutes(t *testing.T) {
	setupRouterTest(t)

	tests := []struct {
		name       string
		method     string
		path       string
		userID     string
		wantStatus int
	}{
		{
			name:       "list expenses route exists",
			method:     http.MethodGet,
			path:       "/api/v1/expenses",
			userID:     "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "summary route exists",
			method:     http.MethodGet,
			path:       "/api/v1/expenses/summary",
			userID:     "1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "get expense by id route exists",
			method:     http.MethodGet,
			path:       "/api/v1/expenses/1",
			userID:     "1",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "unknown route returns not found",
			method:     http.MethodGet,
			path:       "/api/v1/unknown",
			userID:     "",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := performRouterRequest(testCase.method, testCase.path, testCase.userID)

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
