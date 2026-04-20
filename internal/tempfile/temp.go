package tempfile

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
)

const (
	tempDir      = "/tmp"
	maxFilenameLen = 32
)

// GenerateHashFilename generates a unique filename using SHA256 hash.
func GenerateHashFilename() (string, error) {
	data := fmt.Sprintf("%d%d", os.Getpid(), rand.Int())
	hash := sha256.Sum256([]byte(data))
	hashStr := hex.EncodeToString(hash[:])
	return hashStr[:maxFilenameLen], nil
}

// CreateTempFile creates a temporary file with the hash-based filename.
func CreateTempFile(content []byte) (string, error) {
	filename, err := GenerateHashFilename()
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("%s/%s", tempDir, filename)
	if err := os.WriteFile(path, content, 0600); err != nil {
		return "", err
	}
	return path, nil
}

// CleanupTempFile deletes a temporary file.
func CleanupTempFile(path string) error {
	return os.Remove(path)
}
