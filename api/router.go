package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes for the application.
func SetupRouter() *gin.Engine {
	// Set gin test mode
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	r.GET("/health", HealthHandler)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Go PDF")
	})
	r.POST("/api/pdf/shrink", gin.WrapF(ShrinkHandler))
	r.GET("/shrink", func(c *gin.Context) {
		data, _ := os.ReadFile("templates/shrink.html")
		c.HTML(http.StatusOK, "text/html", data)
	})

	r.POST("/api/pdf/merge", gin.WrapF(MergeHandler))

	return r
}

// SetupRouterForTests creates a router for testing.
func SetupRouterForTests() *gin.Engine {
	// Set gin test mode
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	r.GET("/health", HealthHandler)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Go PDF")
	})
	r.POST("/api/pdf/shrink", gin.WrapF(ShrinkHandler))
	r.GET("/shrink", func(c *gin.Context) {
		data, _ := os.ReadFile("templates/shrink.html")
		c.HTML(http.StatusOK, "text/html", data)
	})

	r.POST("/api/pdf/merge", gin.WrapF(MergeHandler))

	return r
}