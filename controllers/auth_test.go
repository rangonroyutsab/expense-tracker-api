package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	appcontrollers "expense-tracker-api/controllers"
	"expense-tracker-api/models"
	_ "expense-tracker-api/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.RunMode = "test"
}

func setupAuthControllerTest(t *testing.T) {
	t.Helper()

	tempDir := t.TempDir()
	models.UsersCSVPath = filepath.Join(tempDir, "users.csv")

	if err := models.EnsureUserFile(); err != nil {
		t.Fatalf("EnsureUserFile() error = %v", err)
	}
}

func performJSONRequest(
	t *testing.T,
	method string,
	path string,
	body interface{},
	headers map[string]string,
) *httptest.ResponseRecorder {
	t.Helper()

	var requestBody bytes.Buffer

	if body != nil {
		if err := json.NewEncoder(&requestBody).Encode(body); err != nil {
			t.Fatalf("failed to encode request body: %v", err)
		}
	}

	request := httptest.NewRequest(method, path, &requestBody)
	request.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	recorder := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(recorder, request)

	return recorder
}

func TestAuthRegister(t *testing.T) {
	tests := []struct {
		name       string
		body       appcontrollers.RegisterRequest
		wantStatus int
	}{
		{
			name: "successful registration",
			body: appcontrollers.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "secret123",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "missing name",
			body: appcontrollers.RegisterRequest{
				Email:    "john@example.com",
				Password: "secret123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing email",
			body: appcontrollers.RegisterRequest{
				Name:     "John Doe",
				Password: "secret123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			body: appcontrollers.RegisterRequest{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "secret123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			body: appcontrollers.RegisterRequest{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "short password",
			body: appcontrollers.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			setupAuthControllerTest(t)

			recorder := performJSONRequest(
				t,
				http.MethodPost,
				"/api/v1/auth/register",
				testCase.body,
				nil,
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

func TestAuthDuplicateRegister(t *testing.T) {
	setupAuthControllerTest(t)

	body := appcontrollers.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret123",
	}

	firstResponse := performJSONRequest(
		t,
		http.MethodPost,
		"/api/v1/auth/register",
		body,
		nil,
	)

	if firstResponse.Code != http.StatusCreated {
		t.Fatalf(
			"first register status = %d, want %d, body = %s",
			firstResponse.Code,
			http.StatusCreated,
			firstResponse.Body.String(),
		)
	}

	secondResponse := performJSONRequest(
		t,
		http.MethodPost,
		"/api/v1/auth/register",
		body,
		nil,
	)

	if secondResponse.Code != http.StatusConflict {
		t.Fatalf(
			"duplicate register status = %d, want %d, body = %s",
			secondResponse.Code,
			http.StatusConflict,
			secondResponse.Body.String(),
		)
	}
}

func TestAuthLogin(t *testing.T) {
	setupAuthControllerTest(t)

	registerBody := appcontrollers.RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret123",
	}

	registerResponse := performJSONRequest(
		t,
		http.MethodPost,
		"/api/v1/auth/register",
		registerBody,
		nil,
	)

	if registerResponse.Code != http.StatusCreated {
		t.Fatalf(
			"register status = %d, want %d, body = %s",
			registerResponse.Code,
			http.StatusCreated,
			registerResponse.Body.String(),
		)
	}

	tests := []struct {
		name       string
		body       appcontrollers.LoginRequest
		wantStatus int
	}{
		{
			name: "successful login",
			body: appcontrollers.LoginRequest{
				Email:    "john@example.com",
				Password: "secret123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "wrong password",
			body: appcontrollers.LoginRequest{
				Email:    "john@example.com",
				Password: "wrongpass",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "missing email",
			body: appcontrollers.LoginRequest{
				Password: "secret123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			body: appcontrollers.LoginRequest{
				Email: "john@example.com",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "unknown user",
			body: appcontrollers.LoginRequest{
				Email:    "missing@example.com",
				Password: "secret123",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := performJSONRequest(
				t,
				http.MethodPost,
				"/api/v1/auth/login",
				testCase.body,
				nil,
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
