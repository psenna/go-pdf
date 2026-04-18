# Base Project Structure Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create the base Go API project structure using gin-gonic with a health endpoint, HTML template page, multi-stage Docker build, and GitHub Actions CI.

**Architecture:** Standard Go project layout with `main.go` as the entry point, gin-gonic for HTTP routing, Go embed for HTML templates. Multi-stage Dockerfile builds the binary in a builder stage and runs it from a scratch image. GitHub Actions runs tests on every PR and push to main.

**Tech Stack:** Go 1.26, gin-gonic, embed (stdlib), docker, GitHub Actions

---

### File Structure

```
go-pdf/
├── main.go                    # Application entry point
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── Dockerfile                 # Multi-stage Docker build
├── .gitignore                 # Git ignore rules
├── .github/workflows/ci.yml   # GitHub Actions CI
├── api/
│   ├── router.go              # Gin router setup with embedded templates
│   └── health.go              # Health check HTTP handler
├── templates/
│   └── index.html             # Home page template (embedded into binary)
└── tests/
    └── api_test.go            # HTTP handler integration tests
```

---

### Task 1: Initialize Go module

**Files:**
- Create: `go.mod`
- Create: `go.sum`

- [ ] **Step 1: Initialize the Go module**

Run:
```bash
cd /home/devops/src/github.com/psenna/go-pdf
go mod init github.com/psenna/go-pdf
```

- [ ] **Step 2: Commit**

```bash
git add go.mod
git commit -m "feat: initialize go module"
```

---

### Task 2: Create project directory structure

**Files:**
- Create dirs: `api/`, `templates/`, `tests/`, `.github/workflows/`

- [ ] **Step 1: Create directories**

Run:
```bash
cd /home/devops/src/github.com/psenna/go-pdf
mkdir -p api templates tests .github/workflows
```

- [ ] **Step 2: Commit**

```bash
git add api/ templates/ tests/ .github/
git commit -m "chore: create project directory structure"
```

---

### Task 3: Create health handler and tests (TDD)

**Files:**
- Create: `tests/api_test.go`
- Create: `api/health.go`

- [ ] **Step 1: Write the failing test**

Create `tests/api_test.go` with the following content:

```go
package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psenna/go-pdf/api"
)

func TestHealthEndpoint(t *testing.T) {
	r := api.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Errorf("expected status 'ok', got %v", resp["status"])
	}
}

func TestHomePage(t *testing.T) {
	r := api.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./tests/ -v`
Expected: FAIL — package api not found

- [ ] **Step 3: Create the health handler**

Create `api/health.go` with the following content:

```go
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler responds with a JSON health status.
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
```

- [ ] **Step 4: Create the router (without templates)**

Create `api/router.go` with the following content:

```go
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes for the application.
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	r.GET("/health", HealthHandler)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Go PDF")
	})

	return r
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./tests/ -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add tests/api_test.go api/health.go api/router.go
git commit -m "feat: add health endpoint and home page handler"
```

---

### Task 4: Add HTML template and embed into binary

**Files:**
- Modify: `api/router.go`
- Create: `templates/index.html`

- [ ] **Step 1: Create the HTML template**

Create `templates/index.html` with the following content:

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go PDF</title>
</head>
<body>
    <h1>Go PDF</h1>
</body>
</html>
```

- [ ] **Step 2: Update the router to embed and serve the template**

Replace the entire content of `api/router.go` with:

```go
package api

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed ../templates/*.html
var templateFS embed.FS

// SetupRouter configures all routes for the application.
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	// Parse embedded templates and register them with the router.
	tmplFiles, _ := fs.Glob(templateFS, "../templates/*.html")
	tmpl, err := gin.ParseFS(tmplFS, "../templates/*.html")
	if err != nil {
		panic(err)
	}
	r.SetTemplate(tmpl)

	r.GET("/health", HealthHandler)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	return r
}
```

- [ ] **Step 3: Run tests to verify everything works**

Run: `go test ./tests/ -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add templates/index.html api/router.go
git commit -m "feat: embed HTML template and serve via gin"
```

---

### Task 5: Create main.go entry point

**Files:**
- Create: `main.go`

- [ ] **Step 1: Create the main entry point**

Create `main.go` with the following content:

```go
package main

import (
	"log"
	"net/http"

	"github.com/psenna/go-pdf/api"
)

func main() {
	r := api.SetupRouter()
	log.Printf("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

- [ ] **Step 2: Verify it builds**

Run: `go build -o go-pdf .`

- [ ] **Step 3: Commit**

```bash
git add main.go
git commit -m "feat: add main entry point"
```

---

### Task 6: Create Dockerfile with multi-stage build

**Files:**
- Create: `Dockerfile`

- [ ] **Step 1: Create the multi-stage Dockerfile**

Create `Dockerfile` with the following content:

```dockerfile
# Build stage
FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o go-pdf .

# Runtime stage
FROM scratch
WORKDIR /app
COPY --from=builder /app/go-pdf .
COPY --from=builder /app/templates ./templates
EXPOSE 8080
CMD ["./go-pdf"]
```

- [ ] **Step 2: Commit**

```bash
git add Dockerfile
git commit -m "feat: add multi-stage Dockerfile with scratch base"
```

---

### Task 7: Create GitHub Actions CI workflow

**Files:**
- Create: `.github/workflows/ci.yml`

- [ ] **Step 1: Create the CI workflow**

Create `.github/workflows/ci.yml` with the following content:

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.26'

      - name: Download dependencies
        run: go mod download

      - name: Run tests
        run: go test ./... -v

      - name: Build
        run: go build -o go-pdf .
```

- [ ] **Step 2: Commit**

```bash
git add .github/workflows/ci.yml
git commit -m "feat: add GitHub Actions CI workflow"
```

---

### Task 8: Add .gitignore

**Files:**
- Create: `.gitignore`

- [ ] **Step 1: Create the .gitignore file**

Create `.gitignore` with the following content:

```
# Binaries
go-pdf

# Runtime data
*.log
```

- [ ] **Step 2: Final commit**

```bash
git add .gitignore
git commit -m "chore: add .gitignore"
```

---

## Self-Review Checklist

- [ ] Spec coverage: Health endpoint - Task 3
- [ ] Spec coverage: HTML page with "Go PDF" - Task 3 + Task 4
- [ ] Spec coverage: Docker multi-stage with scratch - Task 6
- [ ] Spec coverage: GitHub Actions on PRs and main - Task 7
- [ ] Spec coverage: Gin-gonic - Task 3
- [ ] Placeholder scan: No TBD/TODO/placeholder patterns found
- [ ] Type consistency: All file paths, imports, and function signatures are complete and match

---

Plan complete and saved. Two execution options:

1. **Subagent-Driven (recommended)** - Dispatch a fresh subagent per task, review between tasks, fast iteration

2. **Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints

Which approach?