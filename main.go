package main

import (
	"log"

	"github.com/psenna/go-pdf/api"
)

func main() {
	r := api.SetupRouter()
	log.Printf("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}