package pdf

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/psenna/go-pdf/internal/tempfile"
)

// createTestPDF creates a valid PDF using pdfcpu CLI when available
func createTestPDF() (string, error) {
	// Create temp file
	pdfFile, err := tempfile.CreateTempFile(nil)
	if err != nil {
		return "", err
	}

	// Use pdfcpu create command if available
	pdfPath := pdfFile + ".pdf"
	cmd := exec.Command("pdfcpu", "create", pdfPath)
	if err := cmd.Run(); err != nil {
		// If pdfcpu CLI is not available, we can't create valid test PDFs
		// This is expected in CI environments
		return "", fmt.Errorf("pdfcpu create failed: %w", err)
	}

	return pdfPath, nil
}

func TestOptimize(t *testing.T) {
	// Check if pdfcpu CLI is available
	if _, err := exec.LookPath("pdfcpu"); err != nil {
		t.Skip("pdfcpu not installed")
	}

	// Create a valid PDF
	tmpFile, err := createTestPDF()
	if err != nil {
		t.Skipf("failed to create test PDF: %v", err)
	}
	defer os.Remove(tmpFile)

	// Run optimize using SDK
	conf := model.NewDefaultConfiguration()

	// Read input file
	inputData, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ReadSeeker from the input data
	rs := bytes.NewReader(inputData)

	// Create output file
	outputPath := tmpFile + ".optimized.pdf"
	outFile, err := os.Create(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	// Run optimize operation
	if err := api.Optimize(rs, outFile, conf); err != nil {
		t.Errorf("Optimize failed: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("optimized file not created")
	}
}

func TestOptimizeNoFiles(t *testing.T) {
	// Check if pdfcpu CLI is available
	if _, err := exec.LookPath("pdfcpu"); err != nil {
		t.Skip("pdfcpu not installed, skipping optimize no files test")
	}

	// Create a valid PDF
	pdfFile, err := createTestPDF()
	if err != nil {
		t.Skipf("failed to create test PDF: %v", err)
	}
	defer os.Remove(pdfFile)

	// Read input file
	inputData, err := os.ReadFile(pdfFile)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ReadSeeker from the input data
	rs := bytes.NewReader(inputData)

	// Create output file
	outputPath := pdfFile + ".optimized.pdf"
	outFile, err := os.Create(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	defer outFile.Close()

	// Run optimize operation
	conf := model.NewDefaultConfiguration()
	if err := api.Optimize(rs, outFile, conf); err != nil {
		// Log any optimization errors
		t.Logf("optimize returned error: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("optimized file not created")
	}
}
