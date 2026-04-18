package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler responds with a JSON health status.
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}