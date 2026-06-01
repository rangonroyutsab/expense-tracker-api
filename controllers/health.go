package controllers

// HealthController handles API health checks.
type HealthController struct {
	BaseController
}

// Get returns the server health status.
// @Title Health Check
// @Description Check whether the API server is running.
// @Success 200 {object} controller.BasicResponse
// @router /health [get]
func (c *HealthController) Get() {
	c.Success(200, "Server is running", nil)
}
