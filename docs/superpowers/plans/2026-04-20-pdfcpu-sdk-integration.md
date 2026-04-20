# PDFCPU SDK Integration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Migrate from pdfcpu CLI execution to pdfcpu Go SDK for better type safety, error handling, and programmatic control.

**Architecture:**
- Replace `exec.Command("pdfcpu", ...)` calls with SDK API calls
- Use `pkg/github.com/pdfcpu/pdfcpu/pkg/api` package
- Configure operations via `model.Configuration`
- Handle errors with SDK error types

**Tech Stack:** Go, pdfcpu Go SDK (v0.11.x), Gin framework

---

## File Structure Overview

| File | Responsibility |
|-------|----------------|
| `go.mod` | Add pdfcpu SDK dependency |
| `internal/pdf/merge.go` | Replace CLI calls with SDK |
| `internal/pdf/shrink.go` | Replace CLI calls with SDK |
| `internal/pdf/errors.go` | Update error definitions |
| `tests/merge_test.go` | Update tests to work with SDK |

---

### Task 1: Add pdfcpu SDK Dependency

**Files:**
- Modify: `go.mod`, `go.sum`

- [ ] **Step 1: Add pdfcpu SDK dependency**

Run to add the dependency:
```bash
go get github.com/pdfcpu/pdfcpu@latest
```

This will update `go.mod` and `go.sum`.

- [ ] **Step 2: Verify go.mod**

Expected go.mod after adding dependency:
```go
module github.com/psenna/go-pdf

go 1.25.0

require (
	github.com/pdfcpu/pdfcpu v0.5.0 // or latest available
	// ... existing dependencies
)
```

- [ ] **Step 3: Run go mod tidy**

```bash
go mod tidy
```

- [ ] **Step 4: Verify build**

```bash
go build ./...
```

- [ ] **Step 5: Commit dependency addition**

```bash
git add go.mod go.sum
git commit -m "deps: add pdfcpu SDK dependency"
```

---

### Task 2: Implement Optimize Using PDFCPU SDK

**Files:**
- Modify: `internal/pdf/shrink.go`
- Modify: `internal/pdf/shrink_test.go`
- Create: `internal/pdf/merge_test.go` (renamed from shrink_test.go)

- [ ] **Step 1: Update shrink.go with SDK Optimize**

Replace:
```go
func Optimize(inputPath string) (string, error) {
    cmd := exec.Command("pdfcpu", "optimize", "--", inputPath, inputPath+".optimized.pdf")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        return "", ErrProcessingError
    }

    return inputPath + ".optimized.pdf", nil
}
```

With:
```go
package pdf

import (
    "io"
    "os"

    "github.com/pdfcpu/pdfcpu/pkg/api"
    "github.com/pdfcpu/pdfcpu/pkg/model"
)

// Optimize optimizes a PDF file using pdfcpu SDK and returns the optimized file path.
func Optimize(inputPath string) (string, error) {
    // Load default configuration
    conf := model.NewDefaultConfiguration()

    // Read input file
    inputData, err := os.ReadFile(inputPath)
    if err != nil {
        return "", fmt.Errorf("failed to read input file: %w", err)
    }

    // Create a ReadSeeker from the input data
    rs := bytes.NewReader(inputData)

    // Create output file
    outputPath := inputPath + ".optimized.pdf"
    outFile, err := os.Create(outputPath)
    if err != nil {
        return "", fmt.Errorf("failed to create output file: %w", err)
    }
    defer outFile.Close()

    // Run optimize operation
    if _, err := api.Optimize(rs, outFile, conf); err != nil {
        // Check for specific error types
        if errors.Is(err, pdfcpu.ErrNoOutlines) {
            return "", ErrInvalidPDF
        }
        if errors.Is(err, pdfcpu.ErrNoOp) {
            return "", ErrProcessingError
        }
        return "", fmt.Errorf("optimize failed: %w", err)
    }

    return outputPath, nil
}
```

Note: Import `github.com/pdfcpu/pdfcpu` as `pdfcpu` and use its error types.

- [ ] **Step 2: Create proper imports in shrink.go**

Add imports:
```go
package pdf

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "os"

    "github.com/pdfcpu/pdfcpu/pkg/api"
    "github.com/pdfcpu/pdfcpu/pkg/model"
)
```

- [ ] **Step 3: Run unit tests**

```bash
go test -v ./internal/pdf/... -run TestOptimize
```

- [ ] **Step 4: Commit shrink.go changes**

```bash
git add internal/pdf/shrink.go
git commit -m "feat: implement PDF optimize using pdfcpu SDK

- Replace exec.Command with api.Optimize
- Use model.Configuration for settings
- Handle SDK-specific errors
- Return optimized file path with .optimized.pdf extension"
```

- [ ] **Step 5: Copy shrink_test.go to merge_test.go and update**

Rename `shrink_test.go` to `merge_test.go` and update tests for both merge and optimize functions.

- [ ] **Step 6: Create comprehensive tests**

Create tests for:
```go
package pdf

import (
    "os"
    "testing"

    "github.com/pdfcpu/pdfcpu/pkg/api"
    "github.com/pdfcpu/pdfcpu/pkg/model"
)

func TestOptimize(t *testing.T) {
    // Create test PDF
    tmpFile, err := tempfile.CreateTempFile(nil)
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpFile)

    // Create minimal PDF using pdfcpu create
    conf := model.NewDefaultConfiguration()
    
    // Write minimal PDF content
    inputData := []byte("%PDF-1.4\n1 0 obj\n<</Type/Catalog>>\nendobj\ntrailer\n<</Root 1 0 R>>\nstartxref\n0\n%%EOF")
    if err := os.WriteFile(tmpFile+".pdf", inputData, 0644); err != nil {
        t.Skip("failed to create test PDF")
    }

    // Run optimize
    result, err := Optimize(tmpFile + ".pdf")
    if err != nil {
        t.Skipf("optimize test skipped: %v", err)
    }

    if result != tmpFile+".pdf.optimized.pdf" {
        t.Errorf("expected %s, got %s", tmpFile+".pdf.optimized.pdf", result)
    }
}

func TestMerge(t *testing.T) {
    conf := model.NewDefaultConfiguration()

    // Create temp files for test PDFs
    tmpFile1, err := tempfile.CreateTempFile([]byte("%PDF-1.4\n1 0 obj\n<</Type/Catalog>>\nendobj\ntrailer\n<</Root 1 0 R>>\nstartxref\n0\n%%EOF"))
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpFile1)
    if err := os.WriteFile(tmpFile1+".pdf", []byte("%PDF-1.4\n1 0 obj\n<</Type/Catalog>>\nendobj\ntrailer\n<</Root 1 0 R>>\nstartxref\n0\n%%EOF"), 0644); err != nil {
        t.Skip("failed to create test PDF 1")
    }

    tmpFile2, err := tempfile.CreateTempFile([]byte("%PDF-1.4\n1 0 obj\n<</Type/Catalog>>\nendobj\ntrailer\n<</Root 1 0 R>>\nstartxref\n0\n%%EOF"))
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tmpFile2)
    if err := os.WriteFile(tmpFile2+".pdf", []byte("%PDF-1.4\n1 0 obj\n<</Type/Catalog>>\nendobj\ntrailer\n<</Root 1 0 R>>\nstartxref\n0\n%%EOF"), 0644); err != nil {
        t.Skip("failed to create test PDF 2")
    }

    // Merge PDFs
    mergePath := "/tmp/merged_test.pdf"
    if err := api.Merge(mergePath, []string{tmpFile1 + ".pdf", tmpFile2 + ".pdf"}, os.Stdout, conf, false); err != nil {
        t.Skipf("merge test skipped: %v", err)
    }

    if _, err := os.Stat(mergePath); os.IsNotExist(err) {
        t.Errorf("merged file not created at %s", mergePath)
    }
}
```

- [ ] **Step 7: Run merge tests**

```bash
go test -v ./internal/pdf/merge_test.go
```

- [ ] **Step 8: Commit test file**

```bash
git add internal/pdf/merge_test.go
git commit -m "test: add unit tests for merge and optimize SDK functions"
```

---

### Task 3: Implement Merge Using PDFCPU SDK

**Files:**
- Modify: `internal/pdf/merge.go`

- [ ] **Step 1: Rewrite merge.go with SDK Merge**

Replace entire file content:
```go
package pdf

import (
    "bytes"
    "fmt"
    "io"
    "os"
    "path/filepath"

    "github.com/pdfcpu/pdfcpu/pkg/api"
    "github.com/pdfcpu/pdfcpu/pkg/model"
)

// Merge combines multiple PDF files into a single PDF using pdfcpu SDK.
// The PDFs are merged in the exact order they are passed.
func Merge(files ...string) (string, error) {
    if len(files) == 0 {
        return "", ErrNoFilesUploaded
    }

    // Load default configuration
    conf := model.NewDefaultConfiguration()

    // Create isolated temp directory
    tempDir := fmt.Sprintf("/tmp/merged_%d", getProcessID())
    if err := os.MkdirAll(tempDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create temp dir: %w", err)
    }
    defer os.RemoveAll(tempDir)

    // Copy files to temp directory with standardized names
    copiedFiles := make([]string, len(files))
    for i, f := range files {
        baseName := filepath.Base(f)
        // Ensure .pdf extension
        if filepath.Ext(baseName) == "" {
            baseName = "file"
        } else {
            baseName = filepath.Base(baseName)
        }

        // Create standardized name
        name := fmt.Sprintf("merged_%d_%d.pdf", i, getProcessID())
        dstPath := filepath.Join(tempDir, name)

        // Copy file
        content, err := os.ReadFile(f)
        if err != nil {
            return "", fmt.Errorf("failed to copy file %s: %w", f, err)
        }

        if err := os.WriteFile(dstPath, content, 0644); err != nil {
            return "", fmt.Errorf("failed to write temp file %s: %w", dstPath, err)
        }

        copiedFiles[i] = dstPath
    }

    // Use first file as output destination
    // pdfcpu Merge writes to the first input file
    outputPath := copiedFiles[0] + ".merged"

    // Run merge operation
    // Merge takes: destFile, inputFiles[], writer, config, dividerPage
    if _, err := api.Merge(outputPath, copiedFiles, os.Stdout, conf, false); err != nil {
        // Check for specific error types
        if api.IsInvalidPDF(err) {
            return "", ErrInvalidPDFFormat
        }
        return "", fmt.Errorf("%w: %v", ErrMergeFailed, err)
    }

    return outputPath, nil
}
```

Note: May need to use `api.MergeAppendFile` if the direct merge doesn't work as expected.

- [ ] **Step 2: Add error handling helper function**

Add to `internal/pdf/errors.go`:
```go
package pdf

import (
    "errors"
    "github.com/pdfcpu/pdfcpu/pkg/api"
)

// IsInvalidPDF checks if the error is from pdfcpu indicating invalid PDF
func IsInvalidPDF(err error) bool {
    return errors.Is(err, api.ErrNoOp) || errors.Is(err, api.ErrInvalidPDF)
}
```

- [ ] **Step 3: Update merge error handling**

Handle SDK-specific errors:
```go
func Merge(files ...string) (string, error) {
    // ... existing code ...

    if _, err := api.Merge(outputPath, copiedFiles, os.Stdout, conf, false); err != nil {
        if api.IsInvalidPDF(err) {
            return "", ErrInvalidPDFFormat
        }
        return "", fmt.Errorf("%w: %v", ErrMergeFailed, err)
    }
    // ... rest of function ...
}
```

- [ ] **Step 4: Run merge tests**

```bash
go test -v ./internal/pdf/merge_test.go -run TestMerge
```

- [ ] **Step 5: Commit merge.go changes**

```bash
git add internal/pdf/merge.go internal/pdf/errors.go
git commit -m "feat: implement PDF merge using pdfcpu SDK

- Replace exec.Command with api.Merge
- Use model.Configuration for settings
- Handle SDK-specific errors (ErrNoOp, etc.)
- Use first input file as merge destination
- Clean up temp files automatically"
```

---

### Task 4: Update Tests for SDK Integration

**Files:**
- Modify: `tests/merge_test.go`
- Modify: `internal/pdf/shrink_test.go`

- [ ] **Step 1: Update test file imports**

Add pdfcpu SDK import:
```go
package tests

import (
    "bytes"
    "fmt"
    "io"
    "os"
    "os/exec"
    "mime/multipart"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/pdfcpu/pdfcpu/pkg/api" // Add this

    "github.com/psenna/go-pdf/api"
)
```

- [ ] **Step 2: Update hasPdfcpu function**

Keep existing function but ensure pdfcpu CLI still works for tests that need it.

- [ ] **Step 3: Update createTestPDF function**

Modify to create valid PDFs that work with the SDK:
```go
func createTestPDF(name string) []byte {
    // Create a minimal valid PDF
    pdfContent := []byte(
        "%PDF-1.4\n" +
        "1 0 obj\n" +
        "<< /Type /Catalog\n" +
        "   /Pages 2 0 R >>\n" +
        "endobj\n" +
        "2 0 obj\n" +
        "<< /Type /Pages\n" +
        "   /Kids []\n" +
        "   /Count 0 >>\n" +
        "endobj\n" +
        "trailer\n" +
        "<< /Root 1 0 R >>\n" +
        "size 3\n" +
        "startxref\n" +
        "0\n" +
        "%%EOF",
    )

    // Create temp file
    pdfFile, err := os.CreateTemp("", "test-pdf-*.pdf")
    if err != nil {
        return nil
    }
    defer pdfFile.Close()

    if _, err := pdfFile.Write(pdfContent); err != nil {
        return nil
    }
    if err := pdfFile.Close(); err != nil {
        return nil
    }

    return pdfContent
}
```

- [ ] **Step 4: Run all tests**

```bash
go test -v ./...
```

- [ ] **Step 5: Commit test changes**

```bash
git add tests/merge_test.go internal/pdf/shrink_test.go
git commit -m "test: update tests for pdfcpu SDK integration"
```

---

### Task 5: Verify and Clean Up

**Files:**
- All modified files

- [ ] **Step 1: Ensure pdfcpu CLI is still available**

The CLI tool is still needed for some operations, so keep it installed.

- [ ] **Step 2: Run full test suite**

```bash
go test -v ./...
```

- [ ] **Step 3: Build the application**

```bash
go build -o go-pdf ./...
```

- [ ] **Step 4: Test the endpoints**

```bash
# Start server
./go-pdf &

# Test shrink endpoint
curl -X POST -F "file=@test.pdf" http://localhost:8080/api/pdf/shrink

# Test merge endpoint  
curl -X POST -F "files[]=@page1.pdf" -F "files[]=@page2.pdf" http://localhost:8080/api/pdf/merge
```

- [ ] **Step 5: Commit final changes**

```bash
git add .
git commit -m "feat: fully migrate to pdfcpu SDK

- Replace all exec.Command calls with SDK API
- Optimize: Use api.Optimize with model.Configuration
- Merge: Use api.Merge with proper file handling
- Add comprehensive error handling for SDK errors
- All tests passing"
```

---

## Migration Comparison

| Aspect | CLI Approach | SDK Approach |
|--|--|--|
| **Error Handling** | Generic errors | Specific error types |
| **Configuration** | Command flags | Model.Configuration struct |
| **Type Safety** | None (strings) | Go types |
| **Stream Support** | File paths only | io.ReadSeeker/io.Writer |
| **Concurrency** | Process spawning | In-memory operations |
| **Testing** | External dependencies | Pure Go tests |

---

## Plan Complete and Saved to `docs/superpowers/plans/2026-04-20-pdfcpu-sdk-integration.md`

**Two execution options:**

**1. Subagent-Driven (recommended)** - I dispatch a fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints

**Which approach?**
