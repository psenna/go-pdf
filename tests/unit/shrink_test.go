package main

import (
	"os"
	"testing"

	"github.com/psenna/go-pdf/pkg/pdf"
)

func TestOptimize(t *testing.T) {
	// Create a test PDF file
	testPath := "/tmp/test.pdf"
	testData := []byte("%PDF-1.4\n")
	err := os.WriteFile(testPath, testData, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testPath)

	// Test optimization
	optimizedPath, err := pdf.Optimize(testPath)
	if err != nil {
		t.Fatal(err)
	}

	if optimizedPath == "" {
		t.Fatal("Expected non-empty optimized path")
	}

	// Verify optimized file exists
	if _, err := os.Stat(optimizedPath); os.IsNotExist(err) {
		t.Error("Optimized file should exist")
	}

	// Cleanup
	os.Remove(optimizedPath)
}
