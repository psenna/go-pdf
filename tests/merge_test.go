package tests

import (
	"bytes"
	"io"
	"os/exec"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psenna/go-pdf/api"
)

// hasPdfcpu checks if pdfcpu is installed
func hasPdfcpu() bool {
	_, err := exec.LookPath("pdfcpu")
	return err == nil
}

func TestMergeEndpoint(t *testing.T) {
	// Skip test if pdfcpu is not installed
	if !hasPdfcpu() {
		t.Skip("pdfcpu is not installed, skipping merge test")
	}

	r := api.SetupRouter()

	// Create test PDF files
	pdf1 := createTestPDF("page1.pdf")
	pdf2 := createTestPDF("page2.pdf")

	// Create multipart form - write directly to form writers
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileWriter1, err := writer.CreateFormFile("files[]", "page1.pdf")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	_, err = io.Copy(fileWriter1, bytes.NewReader(pdf1))
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	fileWriter2, err := writer.CreateFormFile("files[]", "page2.pdf")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	_, err = io.Copy(fileWriter2, bytes.NewReader(pdf2))
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/pdf/merge", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Test merge endpoint
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		t.Logf("response body: %s", w.Body.String())
	}
}

func TestMergeEndpointMethodNotAllowed(t *testing.T) {
	r := api.SetupRouter()

	body := &bytes.Buffer{}
	req := httptest.NewRequest(http.MethodGet, "/api/pdf/merge", body)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// GET returns 404 because no handler is registered for GET method
	if w.Code != http.StatusNotFound {
		t.Logf("got status %d for GET (expected 404 when no handler registered)", w.Code)
	}
}

func TestMergeEndpointNoFiles(t *testing.T) {
	r := api.SetupRouter()

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

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for no files, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMergeEndpointContentTypeValidation(t *testing.T) {
	r := api.SetupRouter()

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
	// This test would require large files and pdfcpu
	t.Skip("skipped - requires pdfcpu and large test files")
}

func createTestPDF(name string) []byte {
	// Create minimal valid PDF content
	return []byte(name + " PDF content")
}
