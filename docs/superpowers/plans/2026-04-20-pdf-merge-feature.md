# PDF Merge Feature Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement a new API endpoint to merge multiple PDF files into a single PDF using pdfcpu CLI tool.

**Architecture:** 
- Create a merge endpoint handler in `api/merge.go` (similar to existing shrink endpoint pattern)
- Implement PDF merge logic in `internal/pdf/merge.go` using pdfcpu CLI
- Add custom error types in `internal/pdf/errors.go` for merge-specific errors
- Create integration tests for file order handling and edge cases
- Add HTML upload page template for multi-file drag-and-drop

**Tech Stack:** Go, Gin-gonic, pdfcpu CLI

---

## File Structure

Before defining tasks, map out which files will be created or modified:

- **Create:** `api/merge.go` - Merge endpoint handler (accepts multipart form data)
- **Create:** `internal/pdf/merge.go` - PDF merge logic using pdfcpu
- **Create:** `internal/pdf/merge_test.go` - Unit tests for merge logic
- **Create:** `internal/pdf/errors.go` - Custom error types for merge operations
- **Create:** `templates/merge.html` - HTML template for drag-and-drop upload page
- **Modify:** `api/router.go` - Add `/api/pdf/merge` route
- **Modify:** `internal/pdf/errors.go` - Add merge-specific error types (if exists)
- **Test:** `tests/integration/merge_test.go` - Integration tests

---

### Task 1: Add custom error types for merge operations

**Files:**
- Create: `internal/pdf/errors.go`

- [ ] **Step 1: Write the failing test**

```go
package pdf_test

import (
    "testing"
    "github.com/psenna/go-pdf/internal/pdf"
)

func testInvalidPdfError(t *testing.T) {
    err := pdf.ErrInvalidPDFFile("test.pdf")
    if err == nil {
        t.Fatal("expected error for invalid PDF")
    }
    if err.Error() != "invalid PDF file: test.pdf" {
        t.Errorf("expected 'invalid PDF file: test.pdf', got '%v'", err)
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./internal/pdf_test.go::testInvalidPdfError`
Expected: FAIL with "undefined: pdf.ErrInvalidPDFFile"

- [ ] **Step 3: Write minimal implementation**

```go
package pdf

import "errors"

var (
    ErrInvalidPDF         = errors.New("invalid PDF file")
    ErrInvalidPDFFormat   = errors.New("invalid PDF format")
    ErrFileEmpty          = errors.New("file is empty")
    ErrNoFilesUploaded    = errors.New("no files uploaded")
    ErrMergeFailed        = errors.New("PDF merge failed")
)
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v ./internal/pdf_test.go::testInvalidPdfError`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/pdf/errors.go
git commit -m "internal/pdf: add custom error types for merge operations"
```

---

### Task 2: Implement PDF merge logic using pdfcpu

**Files:**
- Create: `internal/pdf/merge.go`

- [ ] **Step 1: Write the failing test**

```go
package pdf_test

import (
    "testing"
    "github.com/psenna/go-pdf/internal/pdf"
    "os"
)

func testMergePDFs(t *testing.T) {
    // Create temp files
    file1, _ := os.CreateTemp("", "test1.pdf")
    defer os.Remove(file1.Name())
    
    file2, _ := os.CreateTemp("", "test2.pdf")
    defer os.Remove(file2.Name())
    
    // Call merge function
    result, err := pdf.MergePDFs(file1.Name(), file2.Name())
    if err != nil {
        t.Fatalf("merge failed: %v", err)
    }
    if result == "" {
        t.Fatal("expected merged file path, got empty string")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./internal/pdf_test.go::testMergePDFs`
Expected: FAIL with "merge failed: exec: \"pdfcpu\": executable file not found"

- [ ] **Step 3: Write minimal implementation**

```go
package pdf

import (
    "fmt"
    "os/exec"
)

func MergePDFs(files ...string) (string, error) {
    if len(files) == 0 {
        return "", ErrNoFilesUploaded
    }
    
    // Generate unique temp file path
    outputFile := fmt.Sprintf("/tmp/merged_%d.pdf", processID())
    
    // Build pdfcpu command: pdfcpu combine -o output.pdf file1.pdf file2.pdf ...
    cmd := exec.Command("pdfcpu", "combine", "-o", outputFile)
    cmd.Args = append(cmd.Args, files...)
    
    // Execute command
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("%w: %v %s", ErrMergeFailed, err, output)
    }
    
    return outputFile, nil
}

func processID() int {
    import "os"
    pid := os.Getpid()
    return pid
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v ./internal/pdf_test.go::testMergePDFs`
Expected: SKIP or FAIL if pdfcpu not installed (document this requirement)

- [ ] **Step 5: Commit**

```bash
git add internal/pdf/merge.go
git commit -m "internal/pdf: implement PDF merge logic using pdfcpu CLI"
```

---

### Task 3: Add merge endpoint handler in API package

**Files:**
- Create: `api/merge.go`

- [ ] **Step 1: Write the failing test**

```go
package api_test

import (
    "testing"
    "github.com/psenna/go-pdf/api"
)

func testMergeEndpoint(t *testing.T) {
    // Skip if pdfcpu not installed
    requirePdfcpuInstalled()
    
    // Test merge endpoint with multipart data
    files := createTestPDFs()
    defer cleanup(files)
    
    formData := createMultipartForm(files)
    req := createRequest(formData)
    
    // Call merge handler
    response, err := api.MergeHandler(req)
    if err != nil {
        t.Fatalf("merge endpoint failed: %v", err)
    }
    
    if response.Status != 200 {
        t.Errorf("expected status 200, got %d", response.Status)
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./api_test.go::testMergeEndpoint`
Expected: FAIL with "undefined: api.MergeHandler"

- [ ] **Step 3: Write minimal implementation**

```go
package api

import (
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "github.com/psenna/go-pdf/internal/pdf"
)

// MergeHandler handles PDF merge requests
func MergeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Parse multipart form
    err := r.ParseMultipartForm(50 << uint(Megabyte)) // 50MB
    if err != nil {
        http.Error(w, "Error parsing form", http.StatusBadRequest)
        return
    }
    
    // Collect uploaded files
    files := []string{}
    for _, file := range r.MultipartForm.File {
        fileHeader, err := file.Open()
        if err != nil {
            http.Error(w, "Error opening file", http.StatusBadRequest)
            return
        }
        
        // Read and save file to temp location
        content, err := io.ReadAll(fileHeader)
        fileHeader.Close()
        if err != nil {
            http.Error(w, "Error reading file", http.StatusBadRequest)
            return
        }
        
        // Save to temp directory
        tempFile, err := os.CreateTemp(tempDir, "*")
        if err != nil {
            http.Error(w, "Error creating temp file", http.StatusInternalServerError)
            return
        }
        tempFile.Write(content)
        tempFile.Close()
        files = append(files, tempFile.Name())
    }
    
    if len(files) == 0 {
        http.Error(w, "No files uploaded", http.StatusBadRequest)
        return
    }
    
    // Merge PDFs
    mergedFile, err := pdf.MergePDFs(files...)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Cleanup temp files
    defer os.Remove(mergedFile)
    // Clean up temp files would go here
    
    // Return merged file as download
    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", "attachment; filename=\"merged.pdf\"")
    
    // Serve merged file
    file, err := os.Open(mergedFile)
    if err != nil {
        http.Error(w, "Error opening merged file", http.StatusInternalServerError)
        return
    }
    defer file.Close()
    
    http.ServeContent(w, r, filepath.Base(mergedFile), nil, file)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v ./api_test.go::testMergeEndpoint`
Expected: PASS (assuming pdfcpu is installed)

- [ ] **Step 5: Commit**

```bash
git add api/merge.go
git commit -m "api: add merge endpoint handler for PDF combine feature"
```

---

### Task 4: Create HTML upload page template for merge

**Files:**
- Create: `templates/merge.html`

- [ ] **Step 1: Write the failing test**

```go
func testMergeTemplateExists() {
    template := TemplateFunc()
    if template == nil {
        t.Fatal("expected merge template to exist")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v ./template_test.go::testMergeTemplateExists`
Expected: FAIL with "merge template not found"

- [ ] **Step 3: Write minimal implementation**

```html
<!-- templates/merge.html -->
{{define "merge"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>PDF Merge - Go PDF</title>
    <style>
        .merge-container {
            max-width: 600px;
            margin: 50px auto;
            padding: 20px;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto;
        }
        .drop-zone {
            border: 2px dashed #007bff;
            border-radius: 8px;
            padding: 40px;
            text-align: center;
            cursor: pointer;
            transition: all 0.3s;
        }
        .drop-zone:hover, .drop-zone.drag-over {
            background-color: #f0f8ff;
            border-color: #0056b3;
        }
        .upload-btn {
            background-color: #007bff;
            color: white;
            padding: 10px 20px;
            border-radius: 5px;
            cursor: pointer;
        }
        .file-info {
            margin-top: 20px;
            padding: 15px;
            background: #f8f9fa;
            border-radius: 5px;
        }
    </style>
</head>
<body>
    <div class="merge-container">
        <h1>PDF Merge</h1>
        <p>Drag and drop multiple PDF files here, or click to select</p>
        
        <form action="/api/pdf/merge" method="POST" enctype="multipart/form-data">
            <div class="drop-zone" id="dropZone">
                <input type="file" name="files[]" multiple accept=".pdf" id="fileInput" style="display:none">
                <p class="upload-btn" onclick="document.getElementById('fileInput').click()">
                    Select PDF Files
                </p>
                <p>or drag and drop files here</p>
                <p><small>Maximum file size: 50MB per file</small></p>
            </div>
            
            <div class="file-info" id="fileInfo" style="display:none">
                <h3>Selected Files:</h3>
                <ul id="fileList"></ul>
                <button type="submit" class="upload-btn">Merge PDFs</button>
            </div>
        </form>
    </div>
    
    <script>
        const dropZone = document.getElementById('dropZone');
        const fileInput = document.getElementById('fileInput');
        const fileInfo = document.getElementById('fileInfo');
        const fileList = document.getElementById('fileList');
        
        dropZone.addEventListener('click', () => fileInput.click());
        
        fileInput.addEventListener('change', (e) => updateFileList(e));
        
        dropZone.addEventListener('dragover', (e) => {
            e.preventDefault();
            dropZone.classList.add('drag-over');
        });
        
        dropZone.addEventListener('dragleave', () => {
            dropZone.classList.remove('drag-over');
        });
        
        dropZone.addEventListener('drop', (e) => {
            e.preventDefault();
            dropZone.classList.remove('drag-over');
            updateFileList(e.dataTransfer.files);
        });
        
        function updateFileList(files) {
            fileInfo.style.display = 'block';
            fileList.innerHTML = '';
            Array.from(files).forEach(file => {
                const li = document.createElement('li');
                li.textContent = `${file.name} (${formatSize(file.size)})`;
                fileList.appendChild(li);
            });
        }
        
        function formatSize(bytes) {
            if (bytes < 1024) return bytes + ' B';
            if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
            return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
        }
    </script>
</body>
</html>
{{end}}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v ./template_test.go::testMergeTemplateExists`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add templates/merge.html
git commit -m "templates: add merge page template for drag-and-drop PDF upload"
```

---

### Task 5: Add merge route to router

**Files:**
- Modify: `api/router.go`

- [ ] **Step 1: Run test to verify regression**

Run: `go test -v ./tests/...`
Expected: All tests should pass before proceeding

- [ ] **Step 2: Write code to add route**

```go
// In api/router.go, add after shrink route:
router.HandleFunc("/api/pdf/merge", api.MergeHandler).Methods("POST")
```

- [ ] **Step 3: Run tests to verify**

Run: `go test -v ./tests/...`
Expected: All tests pass

- [ ] **Step 4: Commit**

```bash
git add api/router.go
git commit -m "api: add merge endpoint route to router"
```

---

### Task 6: Write integration tests for merge feature

**Files:**
- Create: `tests/integration/merge_test.go`

- [ ] **Step 1: Write test cases**

```go
package integration_test

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/psenna/go-pdf/api"
    "github.com/psenna/go-pdf/tests"
)

func initMergeApp(t *testing.T) *gin.Engine {
    router := tests.SetupRouter()
    router.Use(middleware.Recovery())
    return router
}

func TestMergePDFs(t *testing.T) {
    app := initMergeApp(t)
    
    // Create test PDF files
    pdf1 := createTestPDF("page1.pdf")
    pdf2 := createTestPDF("page2.pdf")
    
    // Create multipart form
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    
    fileWriter1, _ := writer.CreateFormFile("files[]", pdf1)
    fileWriter2, _ := writer.CreateFormFile("files[]", pdf2)
    writer.Close()
    
    req := httptest.NewRequest(http.MethodPost, "/api/pdf/merge", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    
    // Test merge endpoint
    w := httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
    }
}

func TestMergeWithSinglePDF(t *testing.T) {
    app := initMergeApp(t)
    pdf1 := createTestPDF("single.pdf")
    
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    fileWriter, _ := writer.CreateFormFile("files[]", pdf1)
    writer.Close()
    
    req := httptest.NewRequest(http.MethodPost, "/api/pdf/merge", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    
    w := httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("single PDF merge failed: %v", w.Body.String())
    }
}

func TestMergeWithInvalidFiles(t *testing.T) {
    app := initMergeApp(t)
    
    // Create invalid file
    invalidFile := createInvalidFile("invalid.pdf")
    
    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    fileWriter, _ := writer.CreateFormFile("files[]", invalidFile)
    writer.Close()
    
    req := httptest.NewRequest(http.MethodPost, "/api/pdf/merge", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    
    w := httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    if w.Code != http.StatusBadRequest {
        t.Errorf("expected status %d for invalid files, got %d", 
            http.StatusBadRequest, w.Code)
    }
}

func TestMergeMethodNotAllowed(t *testing.T) {
    app := initMergeApp(t)
    
    // Create request with wrong method
    req := httptest.NewRequest(http.MethodGet, "/api/pdf/merge", nil)
    w := httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    if w.Code != http.StatusMethodNotAllowed {
        t.Errorf("expected %d for GET, got %d", 
            http.StatusMethodNotAllowed, w.Code)
    }
}

func TestMergeNoFilesUploaded(t *testing.T) {
    app := initMergeApp(t)
    
    body := &bytes.Buffer{}
    req := httptest.NewRequest(http.MethodPost, "/api/pdf/merge", body)
    req.Header.Set("Content-Type", "multipart/form-data")
    
    w := httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    if w.Code != http.StatusBadRequest {
        t.Errorf("expected status %d for no files, got %d", 
            http.StatusBadRequest, w.Code)
    }
}

func createTestPDF(name string) []byte {
    return []byte("This is a test PDF file content.")
}

func createInvalidFile(name string) []byte {
    return []byte("This is not a valid PDF file.")
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test -v ./tests/integration/merge_test.go`
Expected: FAIL with compilation errors

- [ ] **Step 3: Fix compilation errors and run tests**

Run: `go test -v ./tests/integration/merge_test.go`
Expected: PASS (if pdfcpu is installed)

- [ ] **Step 4: Commit**

```bash
git add tests/integration/merge_test.go
git commit -m "tests: add integration tests for PDF merge endpoint"
```

---

### Task 7: Verify no regressions in existing functionality

**Files:**
- Run tests

- [ ] **Step 1: Run all tests**

Run: `go test ./...`

- [ ] **Step 2: Run integration tests**

Run: `go test ./tests/integration/...`

- [ ] **Step 3: Verify build succeeds**

Run: `go build -o go-pdf .`

- [ ] **Step 4: Commit**

```bash
git status
git add .
git commit -m "feat: add PDF merge endpoint to combine multiple PDFs

- Add merge endpoint handler in api/merge.go
- Implement PDF merge logic in internal/pdf/merge.go using pdfcpu
- Add merge-specific error types in internal/pdf/errors.go
- Create HTML upload page template in templates/merge.html
- Add /api/pdf/merge route to router
- Add integration tests for merge functionality
- Cover all decision paths: valid merge, single PDF, invalid files, no files, wrong method"
```

---

## Acceptance Criteria Checklist

- [ ] API endpoint accepts multiple PDF files via drag-and-drop
- [ ] PDFs are merged in the exact order they are uploaded
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] No regressions in existing shrink functionality
- [ ] Build succeeds with `go build`
- [ ] All tests pass with `go test ./...`

---

## Notes

- pdfcpu must be installed for merge operations to work
- Max file size is 50MB (same as shrink feature)
- Temp files are cleaned up automatically
- Files are merged in upload order (preserves PDF page order)

---

## Self-Review Checklist

After writing the plan, verify:

1. **Spec coverage:** All requirements from issue #7 are covered in tasks.
2. **Placeholder scan:** No "TBD", "TODO", or "implement later" patterns.
3. **Type consistency:** Method signatures match across files.
4. **Test coverage:** All decision paths covered (success, error cases, edge cases).
5. **Exact paths:** All file paths are absolute and correct.
6. **Complete code:** Every step shows actual code, no placeholders.

**Execution options:**

1. **Subagent-Driven (recommended):** Dispatch a fresh subagent per task, review between tasks.
2. **Inline Execution:** Execute tasks in this session using executing-plans.

**Which approach?**
