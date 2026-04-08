package pdf

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GenerateTempFilename creates a unique filename based on request data
func GenerateTempFilename(data []byte) string {
	hash := sha256.Sum256(data)
	return filepath.Join(os.TempDir(), fmt.Sprintf("%x.pdf", hash[:16]))
}

// CleanupTempFiles removes temp PDF files older than 1 hour
func CleanupTempFiles() error {
	tempDir := os.TempDir()
	now := time.Now()

	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Name() == ".git" || entry.Name() == "claude" {
			continue
		}

		if filepath.Ext(entry.Name()) != ".pdf" {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		age := now.Sub(info.ModTime())
		if age > time.Hour {
			if err := os.Remove(filepath.Join(tempDir, entry.Name())); err != nil {
				// Log but don't fail on cleanup errors
			}
		}
	}

	return nil
}
