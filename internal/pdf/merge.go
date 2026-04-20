package pdf

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Merge combines multiple PDF files into a single PDF using pdfcpu.
// The PDFs are merged in the exact order they are passed.
func Merge(files ...string) (string, error) {
	if len(files) == 0 {
		return "", ErrNoFilesUploaded
	}

	// Generate unique temp file path
	outputPath := fmt.Sprintf("/tmp/merged_%d.pdf", getProcessID())

	// Build pdfcpu combine command
	cmd := exec.Command("pdfcpu", "combine")
	args := append(cmd.Args, "-o", outputPath)
	args = append(args, files...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %v", ErrMergeFailed, err)
	}

	return outputPath, nil
}

// getProcessID returns a unique ID based on current time
func getProcessID() int {
	return int(time.Now().UnixNano())
}
