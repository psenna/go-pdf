package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psenna/go-pdf/pkg/middleware"
)

func TestFileSizeMiddleware(t *testing.T) {
	// Create a handler that returns 200
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create middleware
	maxSize := int64(1024) // 1KB
	mw := middleware.FileSizeMiddleware(maxSize)
	wrapped := mw(handler)

	// Test with small file (should pass)
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Length", "512")
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Small file should pass, got %d", w.Code)
	}

	// Test with large file (should fail)
	req = httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Length", "2048")
	w = httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("Large file should fail, got %d", w.Code)
	}
}

func TestFileSizeMiddleware_Multipart(t *testing.T) {
	// Create a handler that returns 200
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse multipart form
		err := r.ParseMultipartForm(10<<20)
		if err != nil {
			t.Error(err)
		}
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	maxSize := int64(1024) // 1KB
	mw := middleware.FileSizeMiddleware(maxSize)
	wrapped := mw(handler)

	// Test with multipart form (no Content-Length check)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("file", "test file content")
	writer.Close()

	req := httptest.NewRequest("POST", "/test", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	// Multipart should pass through to handler
	if w.Code != http.StatusOK {
		t.Errorf("Multipart should pass, got %d", w.Code)
	}
}
