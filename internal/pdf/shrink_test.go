package pdf

import (
	"os"
	"os/exec"
	"testing"

	"github.com/psenna/go-pdf/internal/tempfile"
)

func TestOptimize(t *testing.T) {
	// Skip if pdfcpu not installed
	if _, err := exec.LookPath("pdfcpu"); err != nil {
		t.Skip("pdfcpu not installed")
	}

	tmpFile, err := tempfile.CreateTempFile([]byte("%PDF-1.4..."))
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile)

	result, err := Optimize(tmpFile)
	if err != nil {
		t.Error(err)
	}
	if result != tmpFile+".optimized" {
		t.Error("expected optimized file path")
	}
}
