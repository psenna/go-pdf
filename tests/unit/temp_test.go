package main

import (
	"os"
	"testing"

	"github.com/psenna/go-pdf/pkg/pdf"
)

func TestGenerateTempFilename(t *testing.T) {
	filename := pdf.GenerateTempFilename([]byte("test data"))
	if filename == "" {
		t.Fatal("Expected non-empty filename")
	}

	// Verify it's a valid hash-based name
	if len(filename) < 16 {
		t.Errorf("Filename too short: %d", len(filename))
	}
}

func TestCleanupTempFiles(t *testing.T) {
	// Create a temp file
	f, err := os.CreateTemp("", "test-*.pdf")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer func() {
		if err := os.Remove(f.Name()); err != nil {
			t.Logf("Could not remove temp file: %v", err)
		}
	}()

	// Verify file exists
	if _, err := os.Stat(f.Name()); os.IsNotExist(err) {
		t.Error("Temp file should exist")
	}

	// Cleanup - should not error even if no old files exist
	if err := pdf.CleanupTempFiles(); err != nil {
		t.Error(err)
	}

	// File should still exist (it's not older than 1 hour)
	if _, err := os.Stat(f.Name()); os.IsNotExist(err) {
		t.Error("Temp file should still exist (not old enough)")
	}
}
