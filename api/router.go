package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes for the application.
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	r.GET("/health", HealthHandler)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Go PDF")
	})

	return r
}