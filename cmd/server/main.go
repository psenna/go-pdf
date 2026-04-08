package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/psenna/go-pdf/internal/handlers"
)

func main() {
	r := gin.Default()

	r.GET("/health", handlers.HealthHandler)
	r.GET("/", handlers.IndexHandler)

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
