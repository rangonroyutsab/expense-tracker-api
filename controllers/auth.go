package controllers

import (
	"encoding/json"
	"net/mail"
	"strings"

	"expense-tracker-api/models"

	"github.com/beego/beego/v2/core/logs"
)

// AuthController handles user registration and login.
type AuthController struct {
	BaseController
}

// Register creates a new user account.
// @Title Register User
// @Description Register a new user.
// @Param body body controllers.RegisterRequest true "Register request body"
// @Success 201 {object} controllers.BasicResponse
// @Failure 400 {object} controllers.BasicResponse
// @Failure 409 {object} controllers.BasicResponse
// @Failure 500 {object} controllers.BasicResponse
// @router /api/v1/auth/register [post]
func (c *AuthController) Register() {
	var request RegisterRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		c.Error(400, "Invalid request body")
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)

	if request.Name == "" {
		c.Error(400, "Name is required")
		return
	}

	if request.Email == "" {
		c.Error(400, "Email is required")
		return
	}

	if _, err := mail.ParseAddress(request.Email); err != nil {
		c.Error(400, "Email must be valid")
		return
	}

	if request.Password == "" {
		c.Error(400, "Password is required")
		return
	}

	if len(request.Password) < 6 {
		c.Error(400, "Password must be at least 6 characters")
		return
	}

	if _, err := models.GetUserByEmail(request.Email); err == nil {
		c.Error(409, "Email already exists")
		return
	}

	user := &models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	if err := models.CreateUser(user); err != nil {
		logs.Error("failed to create user: %v", err)
		c.Error(500, "Internal server error")
		return
	}

	c.Success(201, "User registered successfully", nil)
}

// Login authenticates a user using email and password.
// @Title Login User
// @Description Login using email and password.
// @Param body body controllers.LoginRequest true "Login request body"
// @Success 200 {object} controllers.Response
// @Failure 400 {object} controllers.BasicResponse
// @Failure 401 {object} controllers.BasicResponse
// @Failure 500 {object} controllers.BasicResponse
// @router /api/v1/auth/login [post]
func (c *AuthController) Login() {
	var request LoginRequest

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		c.Error(400, "Invalid request body")
		return
	}

	request.Email = strings.TrimSpace(request.Email)
	request.Password = strings.TrimSpace(request.Password)

	if request.Email == "" {
		c.Error(400, "Email is required")
		return
	}

	if request.Password == "" {
		c.Error(400, "Password is required")
		return
	}

	user, err := models.GetUserByEmail(request.Email)
	if err != nil || user.Password != request.Password {
		c.Error(401, "Invalid email or password")
		return
	}

	c.Success(200, "Login successful", LoginData{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	})
}
