package api

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/psenna/go-pdf/internal/pdf"
	"github.com/psenna/go-pdf/internal/tempfile"
)

// MergeHandler handles PDF merge requests
func MergeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form with 50MB limit
	err := r.ParseMultipartForm(50 << uint(20))
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Collect uploaded files
	uploadedFiles := []string{}
	tempFiles := []string{}

	// Use FormFile to get individual files from the uploaded file array
	fileHeader := r.MultipartForm.File["files[]"]
	if fileHeader == nil {
		http.Error(w, pdf.ErrNoFilesUploaded.Error(), http.StatusBadRequest)
		return
	}

	// Iterate over all uploaded files
	for _, fh := range fileHeader {
		file, err := fh.Open()
		if err != nil {
			http.Error(w, "Error opening file", http.StatusBadRequest)
			return
		}

		// Read file content
		content, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			http.Error(w, "Error reading file", http.StatusBadRequest)
			return
		}

		// Create temp file for pdfcpu
		tmpPath, err := tempfile.CreateTempFile(content)
		if err != nil {
			http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
			return
		}
		tempFiles = append(tempFiles, tmpPath)

		uploadedFiles = append(uploadedFiles, tmpPath)
	}

	if len(uploadedFiles) == 0 {
		http.Error(w, pdf.ErrNoFilesUploaded.Error(), http.StatusBadRequest)
		return
	}

	// Merge PDFs
	mergedPath, err := pdf.Merge(uploadedFiles...)
	if err != nil {
		// Cleanup temp files
		for _, tf := range tempFiles {
			os.Remove(tf)
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Cleanup temp files and merged result
	defer func() {
		for _, tf := range tempFiles {
			os.Remove(tf)
		}
		os.Remove(mergedPath)
	}()

	// Read merged PDF
	mergedContent, err := os.ReadFile(mergedPath)
	if err != nil {
		http.Error(w, "Failed to read merged PDF", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"merged.pdf\"")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(mergedContent)))
	w.WriteHeader(http.StatusOK)
	w.Write(mergedContent)
}
