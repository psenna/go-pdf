package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShrinkHandler_FileTooLarge(t *testing.T) {
	handler := ShrinkHandler

	largeFile := make([]byte, 60*1024*1024) // 60MB
	req := httptest.NewRequest("POST", "/api/pdf/shrink", bytes.NewReader(largeFile))
	req.Header.Set("Content-Type", "application/pdf")

	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestShrinkHandler_NoFileField(t *testing.T) {
	handler := ShrinkHandler

	// Create a request without file field
	req := httptest.NewRequest("POST", "/api/pdf/shrink", nil)
	req.Header.Set("Content-Type", "multipart/form-data")

	w := httptest.NewRecorder()
	handler(w, req)

	// Should return bad request
	if w.Code != http.StatusBadRequest {
		t.Logf("expected %d, got %d: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}
