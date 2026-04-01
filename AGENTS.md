# AGENTS.md - Lightning Framework

## Commands

```bash
# Build
go build ./...

# Run all tests
go test ./...

# Run a single test
go test -v -run TestName ./

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Run tests for a specific file
go test -v -run TestName context_test.go context.go request.go response.go consts.go json.go context_data.go cookie.go

# Makefile
make test          # runs tests with coverage, generates coverage.html
```

**Note:** Go 1.20+ required. The `make test` target generates `coverage.out` and `coverage.html` (both gitignored).

## Code Style

### Imports
- Two groups separated by blank line: stdlib first, then third-party.
- No vendoring. Use fully qualified import paths.
```go
import (
	"encoding/json"
	"os"
	"strings"

	"github.com/valyala/fasthttp"
)
```

### Formatting
- Use `gofmt` (tabs for indentation, standard Go style).
- No line length limit enforced.

### Types & Structs
- Exported fields: `PascalCase`; unexported fields: `camelCase`.
- Each field on its own line.
- Function types over interfaces where possible (`type Middleware = HandlerFunc`).
- Custom map types for simple wrappers: `type cookiesMap map[string]string`.

### Naming Conventions
- **Constructors:** `NewXxx()` (exported), `newXxx()` (unexported).
- **Methods:** Verb-first: `Get()`, `Post()`, `SetData()`, `SetHeader()`, `AddRoute()`.
- **Constants:** `PascalCase` with category prefix: `StatusOK`, `MethodGet`, `HeaderContentType`, `MIMEApplicationJSON`.
- **Package:** `lightning` (lowercase, single word).

### Error Handling
- Return `error` for recoverable failures: `ParamInt() (int, error)`, `File() error`.
- Silent failure on marshal errors in `JSON()`/`XML()` (no error propagation).
- Panic only in `resolveAddress()` (too many params) and `LoadHTMLGlob()` (via `template.Must`).
- `Recovery()` middleware catches panics and returns 500.
- No custom error types; use standard `error` interface.

### Comments
- Every exported function/type/method must have a doc comment starting with its name.
- Unexported functions should also have descriptive comments.
- Use section comments in `consts.go` to group constants.

## Architecture

### Key Patterns
- **Handlers:** `func(ctx *Context)` — never expose `*fasthttp.RequestCtx` directly.
- **Middleware:** Same signature as handlers (`type Middleware = HandlerFunc`). Chain via `ctx.Next()`.
- **Context pooling:** `sync.Pool` with `acquireContext()`/`releaseContext()`. Always call `reset()` before returning.
- **Router:** Trie-based with per-method roots. Supports `:param` and `*wildcard` patterns.
- **Request serving:** `app.serveRequest(ctx)` for testing; `app.RequestHandler()` for fasthttp server.

### Response Helpers
- `ctx.JSON(code, obj)` — sets Content-Type to `application/json`.
- `ctx.Text(code, text)` — sets Content-Type to `text/plain`.
- `ctx.HTML(code, name, data)` — renders named template.
- `ctx.XML(code, obj)` — sets Content-Type to `application/xml`.
- `ctx.Success(data)` — returns `{"code":0,"message":"ok","data":...}`.
- `ctx.Fail(code, msg)` — returns `{"code":N,"message":"..."}` with 200 status.

## Testing

- All tests in `package lightning` (same package).
- Use `createTestContext(method, path, body)` for full Context with request/response.
- Use `newTestCtx(method, path)` for bare `*fasthttp.RequestCtx`.
- Use `app.serveRequest(ctx)` to simulate requests without a real server.
- Table-driven tests with `t.Run()` for parameterized cases.
- No third-party test libraries (no testify).
- Target coverage: **≥90%** (currently ~90.1%).

## Project Structure

All source files at root level (single package):

| File | Purpose |
|------|---------|
| `lightning.go` | Application struct, config, server startup, context pooling |
| `context.go` | Context API — params, queries, headers, response methods |
| `router.go` | Trie-based router with `:param` and `*wildcard` support |
| `group.go` | Route grouping with prefix inheritance and middleware |
| `request.go` | Internal request wrapper (fasthttp delegation) |
| `response.go` | Internal response wrapper (status, body, redirect, file) |
| `consts.go` | HTTP constants: status codes, methods, MIME types, headers |
| `logger.go` / `recovery.go` | Built-in middleware |
| `cookie.go` / `context_data.go` | Simple map wrappers for cookies and per-request data |

## CI

GitHub Actions (`.github/workflows/go.yml`): runs `go build` and `go test` on push/PR to `main`.
No linting configured (no golangci-lint, go vet, or staticcheck).
