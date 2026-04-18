package pdf

import (
	"testing"
)

func TestErrorsExist(t *testing.T) {
	var _ = ErrInvalidPDF
	var _ = ErrFileTooLarge
	var _ = ErrProcessingError
	var _ = ErrResourceExhausted
}
