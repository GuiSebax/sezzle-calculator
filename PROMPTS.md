# AI Prompts Used

This file lists the prompts used with Claude Code during development, as requested in the assignment.

## Step 1 - Backend

I'm working on a take-home assignment for Sezzle: a full-stack calculator with a React + TypeScript frontend and a Go backend microservice. They evaluate correctness, clarity, idiomatic code, testability, and documentation. The time budget is 2–4 hours, so prioritize simplicity and quality over extra features.

Analyze the current folder before creating anything. Create a monorepo with:

- /backend (Go 1.22+, standard library only, no frameworks)
- /frontend (will be built in a later step, don't create it yet)

Backend requirements:

1. A `calculator` package with pure functions: Add, Subtract, Multiply, Divide, Power, Sqrt, Percentage. Each returns (float64, error). Use sentinel errors for: division by zero, square root of a negative number, and non-finite results (NaN/Inf).
2. An HTTP layer exposing one endpoint per operation using Go 1.22 ServeMux patterns: POST /api/v1/calculate/{operation}, where operation is add, subtract, multiply, divide, power, sqrt, or percentage. Request body: {"a": 5, "b": 3}. For sqrt, only "a" is required. For percentage, compute a% of b.
3. Validation and errors: malformed JSON, missing required fields, or unknown operation → 400. Math errors (division by zero, sqrt of negative, non-finite result) → 422. Success → 200 with {"operation": "add", "result": 8}. Errors → {"error": "clear message"}. Limit the request body size.
4. GET /health returning 200.
5. CORS middleware allowing the frontend origin, configurable through an environment variable.
6. main.go reading PORT from the environment (default 8080), with graceful shutdown.
7. Table-driven unit tests for the calculator package (normal and edge cases) and handler tests using net/http/httptest.
8. A Makefile with: run, test, coverage (generating coverage.out and an HTML report).

Keep the code idiomatic and small. After you finish, explain each design decision in a few sentences, because I'll need to defend them in an interview. Don't create the README yet.
