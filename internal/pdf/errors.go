package pdf

import "errors"

var (
	ErrInvalidPDF       = errors.New("invalid PDF format")
	ErrFileTooLarge     = errors.New("file too large, exceeds maximum size")
	ErrProcessingError  = errors.New("PDF processing failed")
	ErrResourceExhausted = errors.New("resource exhausted, service unavailable")
)
