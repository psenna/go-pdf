package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Merge combines multiple PDF files into a single PDF using pdfcpu SDK.
// The PDFs are merged in the exact order they are passed.
func Merge(files ...string) (string, error) {
	if len(files) == 0 {
		return "", ErrNoFilesUploaded
	}

	// Load default configuration
	conf := model.NewDefaultConfiguration()

	// Generate unique temp directory
	tempDir := fmt.Sprintf("/tmp/merged_%d", getProcessID())
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy files to temp directory with standardized names
	copiedFiles := make([]string, len(files))
	for i, f := range files {
		baseName := filepath.Base(f)
		// Ensure .pdf extension
		if filepath.Ext(baseName) == "" {
			baseName = fmt.Sprintf("page%d", i+1)
		}

		// Create a name with unique suffix
		name := fmt.Sprintf("merged_%d_%s.pdf", i, filepath.Base(baseName))
		dstPath := filepath.Join(tempDir, name)

		// Read and copy file
		content, err := os.ReadFile(f)
		if err != nil {
			return "", fmt.Errorf("failed to copy file %s: %w", f, err)
		}

		if err := os.WriteFile(dstPath, content, 0644); err != nil {
			return "", fmt.Errorf("failed to write temp file %s: %w", dstPath, err)
		}

		copiedFiles[i] = dstPath
	}

	// Use pdfcpu MergeCreate to merge PDFs
	// Merge creates a new PDF file with merged content
	outPath := filepath.Join(tempDir, "merged.pdf")
	if err := api.MergeCreateFile(copiedFiles, outPath, conf); err != nil {
		return "", fmt.Errorf("%w: %v", ErrMergeFailed, err)
	}

	return outPath, nil
}

// getProcessID returns a unique ID based on current time
func getProcessID() int {
	return int(time.Now().UnixNano())
}
