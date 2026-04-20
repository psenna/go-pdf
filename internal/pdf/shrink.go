package pdf

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Optimize optimizes a PDF file using pdfcpu SDK and returns the optimized file path.
func Optimize(inputPath string) (string, error) {
	// Load default configuration
	conf := model.NewDefaultConfiguration()

	// Read input file
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to read input file: %w", err)
	}

	// Create a ReadSeeker from the input data
	rs := bytes.NewReader(inputData)

	// Create output file
	outputPath := inputPath + ".optimized.pdf"
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Run optimize operation
	if err := api.Optimize(rs, outFile, conf); err != nil {
		// Check for specific error types based on error message
		errMsg := err.Error()
		if strings.Contains(strings.ToLower(errMsg), "invalid") {
			return "", ErrInvalidPDF
		}
		if strings.Contains(strings.ToLower(errMsg), "too large") || strings.Contains(errMsg, "requestentitytoolarge") {
			return "", ErrFileTooLarge
		}
		if strings.Contains(strings.ToLower(errMsg), "resource") || strings.Contains(errMsg, "out of memory") {
			return "", ErrResourceExhausted
		}
		return "", fmt.Errorf("optimize failed: %w", err)
	}

	return outputPath, nil
}
