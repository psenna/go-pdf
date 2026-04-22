package pdf

import (
	"fmt"

	"github.com/razvandimescu/gopdf/pdf"
)

// ExtractText extracts plain text content from a PDF file.
// Returns all text from all pages joined by newlines.
func ExtractText(inputPath string) (string, error) {
	// Open PDF file directly
	doc, err := pdf.OpenFile(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF document: %w", err)
	}

	// Extract all text from the document
	text, err := doc.Text()
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	if text == "" {
		return "", ErrNoTextFound
	}

	return text, nil
}
