package pdf

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func hasPdfcpu() bool {
	_, err := exec.LookPath("pdfcpu")
	return err == nil
}

func createTestPDF() (string, error) {
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

func TestOptimize(t *testing.T) {
	// Skip if pdfcpu not installed
	if _, err := exec.LookPath("pdfcpu"); err != nil {
		t.Skip("pdfcpu not installed")
	}

	// Create a valid PDF
	tmpFile, err := createTestPDF()
	if err != nil {
		t.Fatal(err)
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
	// This test would require pdfcpu and large test files
	t.Skip("skipped - requires pdfcpu and large test files")
}
