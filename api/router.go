package api

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed templates
var templateFS embed.FS

// SetupRouter configures all routes for the application.
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	// Parse embedded templates and register them with the router.
	tmpl, err := gin.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		panic(err)
	}
	r.SetTemplate(tmpl)

	r.GET("/health", HealthHandler)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	return r
}