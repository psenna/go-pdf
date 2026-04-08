github repository: psenna/go-pdf

# Core Philosophy
Spec-Driven Development: No code is written without a corresponding specification/test.

Small Incremental Changes: Each iteration must be the smallest possible unit of work.

Separation of Concerns: Maintain a clean architecture (e.g., separating business logic from interface/database layers).

# Workflow Rules
Test-First Mandate: You MUST write a failing test before writing any production code.

Sequential Execution: Break the plan into a task list. Complete exactly one task at a time.

Comprehensive Coverage: Create tests for every possible decision path (if/else, switch cases, error handling).

Verification: After every change, the agent must verify the test passes before moving to the next task.

Create one branch (issue/{issue-id}) to work. After finish issue implementation, create the PR to the main branch.

# Golang Best Practices
Idiomatic Go: Use gofmt, follow standard naming conventions (camelCase), and prefer composition over inheritance.

Explicit Error Handling: Never ignore errors. Handle them explicitly and return them up the stack where appropriate.

Concurrency: Use goroutines and channels only when necessary and ensure they are properly synchronized/closed.

Documentation: Provide concise comments for exported functions and complex logic.
