package controllers

import (
	"errors"
	"strconv"

	"expense-tracker-api/models"

	beego "github.com/beego/beego/v2/server/web"
)

// Response represents the common API response format.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// BaseController provides shared response helpers for all controllers.
type BaseController struct {
	beego.Controller
}

// Success sends a successful JSON response.
func (c *BaseController) Success(status int, message string, data interface{}) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.ServeJSON()
}

// Error sends an error JSON response.
func (c *BaseController) Error(status int, message string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = Response{
		Success: false,
		Message: message,
	}
	c.ServeJSON()
}

// GetCurrentUserID reads and validates the X-User-ID header.
func (c *BaseController) GetCurrentUserID() (int, error) {
	userIDHeader := c.Ctx.Input.Header("X-User-ID")
	if userIDHeader == "" {
		c.Error(401, "Unauthorized")
		return -1, errors.New("missing X-User-ID header")
	}

	userID, err := strconv.Atoi(userIDHeader)
	if err != nil || userID <= 0 {
		c.Error(401, "Unauthorized")
		return -1, errors.New("invalid X-User-ID header")
	}

	if _, err := models.GetUserByID(userID); err != nil {
		c.Error(401, "Unauthorized")
		return -1, err
	}

	return userID, nil
}
