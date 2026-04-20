package tempfile

import (
	"os"
	"testing"
)

func TestGenerateHashFilename(t *testing.T) {
	filename, err := GenerateHashFilename()
	if err != nil {
		t.Fatal(err)
	}
	if len(filename) == 0 {
		t.Error("expected non-empty filename")
	}
}

func TestCleanupTempFile(t *testing.T) {
	tmpFile, err := CreateTempFile([]byte("test"))
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile)

	if err := CleanupTempFile(tmpFile); err != nil {
		t.Fatal(err)
	}
}

func TestCreateTempFile(t *testing.T) {
	tmpFile, err := CreateTempFile([]byte("test content"))
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile)

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "test content" {
		t.Error("expected 'test content', got", string(content))
	}
}

func TestGenerateHashFilename_Uniqueness(t *testing.T) {
	filenames := make(map[string]bool)
	for i := 0; i < 100; i++ {
		filename, err := GenerateHashFilename()
		if err != nil {
			t.Fatal(err)
		}
		if filenames[filename] {
			t.Error("expected unique filenames, got duplicate:", filename)
		}
		filenames[filename] = true
	}
}
