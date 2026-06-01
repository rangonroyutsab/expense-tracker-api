package models

import "errors"

// User represents a registered user in the system.
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

// UsersCSVPath stores the path to the users CSV file.
var UsersCSVPath = "data/users.csv"

// GetUserByID returns a user by ID.
func GetUserByID(id int) (*User, error) {
	return nil, errors.New("user not found")
}
