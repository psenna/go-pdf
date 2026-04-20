package pdf

import (
	"os"
	"os/exec"
)

// Optimize processes a PDF file using pdfcpu and returns the optimized file path.
func Optimize(inputPath string) (string, error) {
	// Generate unique output path
	outputPath := inputPath + ".optimized"

	// Run pdfcpu optimize command
	cmd := exec.Command("pdfcpu", "optimize", inputPath, outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", ErrProcessingError
	}

	return outputPath, nil
}
