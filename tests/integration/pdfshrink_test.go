package integration

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/psenna/go-pdf/handlers"
	"github.com/psenna/go-pdf/internal/config"
)

func TestShrinkEndpoint(t *testing.T) {
	// Create config
	cfg := &config.Config{
		MaxFileSize:  50 << 20, // 50MB
		Timeout:       30 * time.Second,
		ConcurrencyLimit: 5,
	}

	// Create handler
	handler := handlers.ShrinkHandler(cfg)

	// Create test PDF file
	testPDF := []byte("%PDF-1.4\n")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// Write file field
	writerFile, err := writer.CreateFormFile("file", "test.pdf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = writerFile.Write(testPDF)
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()

	req := httptest.NewRequest("POST", "/api/pdf/shrink", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should return 200 with optimized PDF
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Verify response contains PDF content
	if len(w.Body.Bytes()) == 0 {
		t.Error("Expected PDF content in response")
	}
}

func TestShrinkEndpoint_MethodNotAllowed(t *testing.T) {
	cfg := &config.Config{
		MaxFileSize:  50 << 20,
		Timeout:       30 * time.Second,
		ConcurrencyLimit: 5,
	}

	handler := handlers.ShrinkHandler(cfg)

	req := httptest.NewRequest("GET", "/api/pdf/shrink", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", w.Code)
	}
}

func TestShrinkEndpoint_WrongContentType(t *testing.T) {
	cfg := &config.Config{
		MaxFileSize:  50 << 20,
		Timeout:       30 * time.Second,
		ConcurrencyLimit: 5,
	}

	handler := handlers.ShrinkHandler(cfg)

	req := httptest.NewRequest("POST", "/api/pdf/shrink", bytes.NewReader([]byte("test")))
	req.Header.Set("Content-Type", "application/pdf")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestShrinkEndpoint_FileTooLarge(t *testing.T) {
	cfg := &config.Config{
		MaxFileSize:  1024, // 1KB limit
		Timeout:       30 * time.Second,
		ConcurrencyLimit: 5,
	}

	handler := handlers.ShrinkHandler(cfg)

	// Create large file
	largeFile := make([]byte, 2048)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// Write file field
	writerFile, err := writer.CreateFormFile("file", "test.pdf")
	if err != nil {
		t.Fatal(err)
	}
	_, err = writerFile.Write(largeFile)
	if err != nil {
		t.Fatal(err)
	}
	writer.Close()

	req := httptest.NewRequest("POST", "/api/pdf/shrink", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Expected 413, got %d", w.Code)
	}
}
