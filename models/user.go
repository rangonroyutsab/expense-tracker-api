package models

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// User represents a registered user in the system.
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

// UsersCSVPath stores the configurable users CSV file path.
var UsersCSVPath = "data/users.csv"

var userCSVHeader = []string{"id", "name", "email", "password", "created_at"}

// EnsureUserFile creates the users CSV file with a header if needed.
func EnsureUserFile() error {
	if err := os.MkdirAll(filepath.Dir(UsersCSVPath), 0755); err != nil {
		return err
	}

	fileInfo, err := os.Stat(UsersCSVPath)
	if err == nil && fileInfo.Size() > 0 {
		return nil
	}

	file, err := os.Create(UsersCSVPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if err := writer.Write(userCSVHeader); err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}

// GetAllUsers returns all users from the users CSV file.
func GetAllUsers() ([]User, error) {
	file, err := os.Open(UsersCSVPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	for index, record := range records {
		if index == 0 {
			continue
		}

		if len(record) != 5 {
			return nil, errors.New("invalid users CSV row")
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		users = append(users, User{
			ID:        id,
			Name:      record[1],
			Email:     record[2],
			Password:  record[3],
			CreatedAt: record[4],
		})
	}

	return users, nil
}

// GetUserByEmail returns a user matching the provided email.
func GetUserByEmail(email string) (*User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if strings.EqualFold(user.Email, email) {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

// GetUserByID returns a user matching the provided ID.
func GetUserByID(id int) (*User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, errors.New("user not found")
}

// CreateUser creates a user and appends it to the users CSV file.
func CreateUser(user *User) error {
	id, err := getNextUserID()
	if err != nil {
		return err
	}

	user.ID = id
	user.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	file, err := os.OpenFile(UsersCSVPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if err := writer.Write([]string{
		strconv.Itoa(user.ID),
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
	}); err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}

func getNextUserID() (int, error) {
	users, err := GetAllUsers()
	if err != nil {
		return 0, err
	}

	maxID := 0

	for _, user := range users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}

	return maxID + 1, nil
}

func writeAllUsers(users []User) error {
	if err := os.MkdirAll(filepath.Dir(UsersCSVPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(UsersCSVPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if err := writer.Write(userCSVHeader); err != nil {
		return err
	}

	for _, user := range users {
		if err := writer.Write([]string{
			strconv.Itoa(user.ID),
			user.Name,
			user.Email,
			user.Password,
			user.CreatedAt,
		}); err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}
