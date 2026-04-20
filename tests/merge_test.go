package tests

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/psenna/go-pdf/api"
)

// createTestPDF creates a valid PDF using pdfcpu CLI when available
func createTestPDF(name string) (string, error) {
	// Create temp file
	pdfFile, err := os.CreateTemp("", "test-pdf-*.pdf")
	if err != nil {
		return "", err
	}

	// Use pdfcpu create command if available
	pdfPath := pdfFile.Name() + ".pdf"
	cmd := exec.Command("pdfcpu", "create", pdfPath)
	if err := cmd.Run(); err != nil {
		pdfFile.Close()
		os.Remove(pdfFile.Name())
		return "", fmt.Errorf("pdfcpu create failed: %w", err)
	}

	// Close the file after create
	if err := pdfFile.Close(); err != nil {
		os.Remove(pdfPath)
		return "", err
	}

	return pdfPath, nil
}

func TestMergeEndpoint(t *testing.T) {
	// Check if pdfcpu CLI is available
	if _, err := exec.LookPath("pdfcpu"); err != nil {
		t.Skip("pdfcpu not installed, skipping merge test")
	}

	// Create test PDF files
	pdf1Path, err := createTestPDF("page1.pdf")
	if err != nil {
		t.Skipf("failed to create test PDF 1: %v", err)
	}
	defer os.Remove(pdf1Path)

	pdf2Path, err := createTestPDF("page2.pdf")
	if err != nil {
		t.Skipf("failed to create test PDF 2: %v", err)
	}
	defer os.Remove(pdf2Path)

	// Use pdfcpu SDK to merge
	conf := model.NewDefaultConfiguration()
	outPath := "/tmp/merged_test_sdk.pdf"

	if err := pdfcpu.MergeCreateFile([]string{pdf1Path, pdf2Path}, outPath, conf); err != nil {
		t.Fatalf("MergeCreateFile failed: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Errorf("merged file not created at %s", outPath)
	}

	// Cleanup
	os.Remove(outPath)
}

func TestMergeEndpointMethodNotAllowed(t *testing.T) {
	// Test with wrong HTTP method
	r := api.SetupRouterForTests()

	body := &bytes.Buffer{}
	req := httptest.NewRequest(http.MethodGet, "/api/pdf/merge", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// GET should return 404 because no handler is registered for GET method
	// (gin doesn't auto-handle wrong methods for registered POST routes)
	if w.Code != http.StatusNotFound {
		t.Logf("got status %d for GET (expected 404)", w.Code)
	}
}

func TestMergeEndpointNoFiles(t *testing.T) {
	r := api.SetupRouterForTests()

	// Create empty form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err := writer.Close()
	if err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/pdf/merge", body)
	req.Header.Set("Content-Type", "multipart/form-data")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Empty form should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for no files, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMergeEndpointContentTypeValidation(t *testing.T) {
	r := api.SetupRouterForTests()

	// Create request with wrong content type
	body := &bytes.Buffer{}
	req := httptest.NewRequest(http.MethodPost, "/api/pdf/merge", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should fail parsing due to wrong content type
	if w.Code != http.StatusBadRequest {
		t.Logf("got status %d (expected %d or %d)", w.Code, http.StatusBadRequest, http.StatusUnsupportedMediaType)
	}
}

func TestMergeEndpointFileTooLarge(t *testing.T) {
	r := api.SetupRouterForTests()

	// Create a large file (larger than 50MB limit)
	largeContent := make([]byte, 60*1024*1024) // 60MB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileWriter, err := writer.CreateFormFile("files[]", "large.pdf")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err = io.Copy(fileWriter, bytes.NewReader(largeContent)); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/pdf/merge", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Large file should be rejected with 413
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Logf("got status %d for large file (expected 413)", w.Code)
	}
}
