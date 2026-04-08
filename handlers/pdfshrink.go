package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/psenna/go-pdf/pkg/pdf"
	"github.com/psenna/go-pdf/internal/config"
)

type ShrinkResponse struct {
	OriginalSize     int64   `json:"original_size"`
	OptimizedSize    int64   `json:"optimized_size"`
	ReductionPercent float64 `json:"reduction_percent"`
	ProcessingTimeMs int64   `json:"processing_time_ms"`
}

func ShrinkHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate method
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Validate content type
		contentType := r.Header.Get("Content-Type")
		if contentType == "" || !strings.HasPrefix(contentType, "multipart/form-data") {
			http.Error(w, "Content-Type must be multipart/form-data", http.StatusBadRequest)
			return
		}

		// Parse multipart form
		err := r.ParseMultipartForm(50<<20) // 50MB max memory
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get file from form
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "No file provided", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Validate file size
		if header.Size > cfg.MaxFileSize {
			http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
			return
		}

		// Read file data
		tempData := make([]byte, header.Size)
		_, err = io.ReadFull(file, tempData)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}

		// Generate temp filename
		tempPath := pdf.GenerateTempFilename(tempData)

		// Write to temp file
		out, err := os.Create(tempPath)
		if err != nil {
			http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = out.Write(tempData)
		if err != nil {
			http.Error(w, "Failed to write temp file", http.StatusInternalServerError)
			return
		}

		// Process PDF
		optimizedPath, err := pdf.Optimize(tempPath)
		if err != nil {
			http.Error(w, "Failed to optimize PDF", http.StatusInternalServerError)
			return
		}
		defer os.Remove(optimizedPath)

		// Read optimized PDF
		optimizedData, err := os.ReadFile(optimizedPath)
		if err != nil {
			http.Error(w, "Failed to read optimized PDF", http.StatusInternalServerError)
			return
		}

		// Set headers for streaming
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=\"optimized.pdf\"")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(optimizedData)))

		// Write optimized PDF
		w.Write(optimizedData)

		// Cleanup temp files
		pdf.CleanupTempFiles()
	}
}
