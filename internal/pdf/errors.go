package pdf

import "errors"

var (
	// Shrink errors
	ErrInvalidPDF       = errors.New("invalid PDF format")
	ErrFileTooLarge     = errors.New("file too large, exceeds maximum size")
	ErrProcessingError  = errors.New("PDF processing failed")
	ErrResourceExhausted = errors.New("resource exhausted, service unavailable")

	// Merge errors
	ErrNoFilesUploaded      = errors.New("no files uploaded")
	ErrMergeFailed          = errors.New("PDF merge failed")
	ErrInvalidPDFFormat     = errors.New("invalid PDF format for merge")

	// Text extraction errors
	ErrTextExtractionFailed = errors.New("text extraction failed")
	ErrNoTextFound          = errors.New("no text found in PDF")
)
