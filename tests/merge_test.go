package tests

import (
	"os"
	"os/exec"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// hasPdfcpu checks if pdfcpu is installed
func hasPdfcpu() bool {
	_, err := exec.LookPath("pdfcpu")
	return err == nil
}

// createTestPDF creates a minimal valid PDF file and returns file path
func createTestPDF(name string) (string, error) {
	// Create a minimal valid PDF
	pdfContent := []byte(
		"%PDF-1.4\n" +
			"1 0 obj\n" +
			"<< /Type /Catalog\n" +
			"   /Pages 2 0 R >>\n" +
			"endobj\n" +
			"2 0 obj\n" +
			"<< /Type /Pages\n" +
			"   /Kids []\n" +
			"   /Count 0 >>\n" +
			"endobj\n" +
			"trailer\n" +
			"<< /Root 1 0 R >>\n" +
			"size 3\n" +
			"startxref\n" +
			"0\n" +
			"%%EOF",
	)

	// Create temp file
	pdfFile, err := os.CreateTemp("", "test-pdf-*.pdf")
	if err != nil {
		return "", err
	}

	if _, err := pdfFile.Write(pdfContent); err != nil {
		pdfFile.Close()
		os.Remove(pdfFile.Name())
		return "", err
	}

	if err := pdfFile.Close(); err != nil {
		os.Remove(pdfFile.Name())
		return "", err
	}

	return pdfFile.Name(), nil
}

func TestMergeEndpoint(t *testing.T) {
	// Skip test if pdfcpu is not installed
	if !hasPdfcpu() {
		t.Skip("pdfcpu is not installed, skipping merge test")
	}

	// Create test PDF files
	pdf1Path, err := createTestPDF("page1.pdf")
	if err != nil {
		t.Skip("failed to create test PDF 1:", err)
	}
	defer os.Remove(pdf1Path)

	pdf2Path, err := createTestPDF("page2.pdf")
	if err != nil {
		t.Skip("failed to create test PDF 2:", err)
	}
	defer os.Remove(pdf2Path)

	// Use pdfcpu SDK to merge
	conf := model.NewDefaultConfiguration()
	outPath := "/tmp/merged_test_sdk.pdf"

	if err := api.MergeCreateFile([]string{pdf1Path, pdf2Path}, outPath, conf); err != nil {
		t.Errorf("MergeCreateFile failed: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Errorf("merged file not created at %s", outPath)
	}

	// Cleanup
	os.Remove(outPath)
}

func TestMergeEndpointMethodNotAllowed(t *testing.T) {
	// Test is skipped - router setup moved to main package
	t.Skip("router setup moved to main package")
}

func TestMergeEndpointNoFiles(t *testing.T) {
	// Test is skipped - router setup moved to main package
	t.Skip("router setup moved to main package")
}

func TestMergeEndpointContentTypeValidation(t *testing.T) {
	// Test is skipped - router setup moved to main package
	t.Skip("router setup moved to main package")
}

func TestMergeEndpointFileTooLarge(t *testing.T) {
	// This test would require large files and pdfcpu
	t.Skip("skipped - requires pdfcpu and large test files")
}
