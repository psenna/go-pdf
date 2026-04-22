package api

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/psenna/go-pdf/internal/pdf"
	"github.com/psenna/go-pdf/internal/tempfile"
)

func ExtractHandler(w http.ResponseWriter, r *http.Request) {
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
	text, err := pdf.ExtractText(tmpPath)
	if err != nil {
		if errors.Is(err, pdf.ErrTextExtractionFailed) {
			http.Error(w, "Text extraction failed", http.StatusInternalServerError)
			return
		}
		if errors.Is(err, pdf.ErrNoTextFound) {
			http.Error(w, "No text found in PDF", http.StatusBadRequest)
			return
		}
		if errors.Is(err, pdf.ErrInvalidPDF) {
			http.Error(w, "Invalid PDF format", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(text)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
}
