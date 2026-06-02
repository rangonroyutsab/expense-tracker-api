package models

import (
	"os"
	"path/filepath"
	"testing" //built-in testing package of Go
)

// creates a temporary user.csv
func setupUserTestFile(t *testing.T) {
	t.Helper() // test helper

	tempDir := t.TempDir() // creates a temporary folder for testing
	UsersCSVPath = filepath.Join(tempDir, "users.csv")

	if err := EnsureUserFile(); err != nil {
		t.Fatalf("EnsureUserFile() error = %v", err)
	}
}

func TestEnsureUserFile(t *testing.T) {
	tempDir := t.TempDir()
	UsersCSVPath = filepath.Join(tempDir, "users.csv")

	if err := EnsureUserFile(); err != nil {
		t.Fatalf("EnsureUserFile() error = %v", err)
	}

	if _, err := os.Stat(UsersCSVPath); err != nil {
		t.Fatalf("expected users CSV file to exist, got error = %v", err)
	}
}

// creates a user, then checks whether the user can be found by email and by ID
func TestCreateUserAndLookup(t *testing.T) {
	setupUserTestFile(t)

	user := &User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret123",
	}

	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	// table-driven test
	tests := []struct {
		name      string
		lookup    func() (*User, error)
		wantID    int
		wantEmail string
	}{
		{
			name: "lookup by email",
			lookup: func() (*User, error) {
				return GetUserByEmail("john@example.com")
			},
			wantID:    1,
			wantEmail: "john@example.com",
		},
		{
			name: "lookup by id",
			lookup: func() (*User, error) {
				return GetUserByID(1)
			},
			wantID:    1,
			wantEmail: "john@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.lookup()
			if err != nil {
				t.Fatalf("lookup error = %v", err)
			}

			if got.ID != tt.wantID {
				t.Fatalf("ID = %d, want %d", got.ID, tt.wantID)
			}

			if got.Email != tt.wantEmail {
				t.Fatalf("Email = %s, want %s", got.Email, tt.wantEmail)
			}
		})
	}
}

func TestUserLookupFailures(t *testing.T) {
	setupUserTestFile(t)

	tests := []struct {
		name   string
		lookup func() (*User, error)
	}{
		{
			name: "missing email",
			lookup: func() (*User, error) {
				return GetUserByEmail("missing@example.com")
			},
		},
		{
			name: "missing id",
			lookup: func() (*User, error) {
				return GetUserByID(99)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.lookup(); err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestGetAllUsersMalformedCSV(t *testing.T) {
	tempDir := t.TempDir()
	UsersCSVPath = filepath.Join(tempDir, "users.csv")

	content := []byte("id,name,email,password,created_at\n1,John Doe,john@example.com\n")

	if err := os.WriteFile(UsersCSVPath, content, 0644); err != nil {
		t.Fatalf("failed to write malformed users CSV: %v", err)
	}

	if _, err := GetAllUsers(); err == nil {
		t.Fatal("expected malformed users CSV error, got nil")
	}
}
