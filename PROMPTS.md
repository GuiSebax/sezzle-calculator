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

## Step 2 - Frontend

Now build the frontend in /frontend. The backend in /backend is done and exposes POST /api/v1/calculate/{operation} (add, subtract, multiply, divide, power, sqrt, percentage) with body {"a": number, "b": number} (sqrt only uses "a"). It returns {"operation": string, "result": number} on success and {"error": string} on failure (400 for invalid input, 422 for math errors). Read the backend code first to confirm the contract.

Requirements:

1. Vite + React + TypeScript (strict mode). No UI library; plain CSS or CSS modules. Keep dependencies minimal.
2. A small typed API client module (src/api) that calls the backend, maps error responses to typed errors, and reads the base URL from VITE_API_URL (default http://localhost:8080). Components must not call fetch directly.
3. UI: two number inputs, an operation selector, a Calculate button, and a result area. When the operation is sqrt, hide or disable the second input. Show a loading state while the request is in flight.
4. Input validation on the client: required fields, valid numbers, and a clear inline message. Disable the button while invalid. Still rely on the backend as the source of truth and display its error messages clearly.
5. Handle network errors (backend down) with a friendly message.
6. Responsive layout that works well on mobile (basic support), accessible labels, and keyboard submit with Enter.
7. Unit tests with Vitest and React Testing Library: the API client (mocking fetch) and the main component (valid calculation, validation message, backend error, sqrt hiding the second input). Add a coverage script.
8. npm scripts: dev, build, test, coverage.

Keep components small and the code idiomatic. After finishing, explain the main design decisions briefly so I can defend them in an interview. Don't create Docker files or the README yet.

## Step 3 - Docker and documentation

Final step. The Go backend (/backend, entry point in cmd/server) and the React + TypeScript frontend (/frontend) are done and tested. Read both projects first.

1. Docker:
   - backend/Dockerfile: multi-stage build (golang image to build, small final image such as distroless or alpine), running as a non-root user.
   - frontend/Dockerfile: multi-stage build (node to build, nginx to serve the static files).
   - docker-compose.yml at the root running both. The frontend must reach the backend correctly (configure VITE_API_URL at build time or proxy /api through nginx, and explain the choice). Backend CORS must allow the frontend origin.
   - Add .dockerignore files.
2. Coverage report: make sure both layers can generate a coverage report with a single command, and add the generated report folders to .gitignore. Add a short "Test coverage" section to the README with the current numbers (backend: run go test ./... -cover to get them; frontend: 45 tests, ~98% statements).
3. README.md at the root, in clear English, with:
   - Overview of the project
   - Tech stack
   - Project structure
   - Setup and how to run (with Docker Compose, and locally without Docker for both backend and frontend)
   - How to run tests and coverage
   - API documentation: endpoints, request/response examples with curl, status codes (200, 400, 422), and error format
   - Design decisions and trade-offs (separation of the pure calculator package from the HTTP layer, one endpoint per operation, 400 vs 422, standard library only, typed API client in the frontend, client-side validation with the backend as source of truth, etc.)
   - Assumptions (for example, how percentage and sqrt are interpreted)
   - Possible improvements with more time
   - A note that AI tooling was used, pointing to PROMPTS.md
4. Make sure .gitignore covers node_modules, build outputs, and coverage files.

Keep the README concise and well organized. Don't change application behavior.
