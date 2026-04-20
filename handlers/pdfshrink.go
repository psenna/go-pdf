package handlers

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/psenna/go-pdf/internal/pdf"
	"github.com/psenna/go-pdf/internal/tempfile"
)

func init() {
	Register("POST /api/pdf/shrink", ShrinkHandler)
}

func ShrinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		http.Error(w, "Content-Type is required", http.StatusBadRequest)
		return
	}

	contentType, _, _ = strings.Cut(contentType, ";")
	if contentType != "multipart/form-data" {
		http.Error(w, "Content-Type must be multipart/form-data", http.StatusBadRequest)
		return
	}

	// Handle file size limit (50MB)
	maxSize := int64(50 * 1024 * 1024)
	r.Body = http.MaxBytesReader(w, r.Body, maxSize)

	file, header, err := r.FormFile("file")
	if err != nil {
		_ = header
		if errors.Is(err, io.ErrUnexpectedEOF) {
			http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Create temp file for pdfcpu
	tmpPath, err := tempfile.CreateTempFile(content)
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer tempfile.CleanupTempFile(tmpPath)

	// Process PDF with pdfcpu
	optimizedPath, err := pdf.Optimize(tmpPath)
	if err != nil {
		if errors.Is(err, pdf.ErrInvalidPDF) {
			http.Error(w, "Invalid PDF format", http.StatusBadRequest)
			return
		}
		if errors.Is(err, pdf.ErrFileTooLarge) {
			http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
			return
		}
		if errors.Is(err, pdf.ErrProcessingError) {
			http.Error(w, "Processing failed", http.StatusInternalServerError)
			return
		}
		if errors.Is(err, pdf.ErrResourceExhausted) {
			http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(optimizedPath)

	// Read optimized PDF
	optimizedContent, err := os.ReadFile(optimizedPath)
	if err != nil {
		http.Error(w, "Failed to read optimized PDF", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"optimized.pdf\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(optimizedContent)))
	w.WriteHeader(http.StatusOK)
	w.Write(optimizedContent)
}
