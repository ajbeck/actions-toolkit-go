Project Overview

The GitHub Actions Toolkit for Go is a library for building GitHub Actions in Go, inspired by the official JavaScript/TypeScript toolkit (https://github.com/actions/toolkit). This is a **private module** at `github.com/smartpension/actions-toolkit-go`.

# Module Architecture

The module consists of 2 packages:

## `actioncontext`

Exposes strongly typed GitHub Actions runtime context from environment variables. The package:

- Auto-loads context on `init()` using environment variables
- Provides package-level variables for all GitHub Actions environment variables (e.g., `actioncontext.Actor`, `actioncontext.Ref`, `actioncontext.EventName`).
- Provides an `ActionContext` type
- In addition to each individual variable loaded on `init()` there is a `CurrentContext` variable which contains an `ActionContext` typed value

## `core`

Provides GitHub Actions workflow commands and file commands. Three main areas:

1. **Workflow Commands** (`workflowcommands.go`): Stdout-based commands that GitHub Actions recognizes
2. **File Commands** (`filecommands.go`): Write to runner-provided files to set values
3. **Step Summary** (`summary.go`): Build HTML job summaries

# Development Workflow and Tools

## Reading Go Documentation

@docs/go-doc-cmd-help.md

## Testing

All commands must be run from the repository root.

```bash
# Run all tests with local cache (avoids permission issues on shared systems)
GOCACHE=$(pwd)/.gocache/build GOMODCACHE=$(pwd)/.gocache/mod go test ./...

# Run a single test
GOCACHE=$(pwd)/.gocache/build GOMODCACHE=$(pwd)/.gocache/mod go test ./actioncontext -run TestSpecificTest
```

### Building and Linting

```bash
# Format code (always run before committing)
go fmt -json ./...

# Static analysis
go vet -json ./...

# Build all packages (verifies imports and types)
GOCACHE=$(pwd)/.gocache/build GOMODCACHE=$(pwd)/.gocache/mod go build ./...

# Clean dependencies after adding/removing imports
go mod tidy
```

## Testing Patterns

- Use **table-driven tests** for variations
- Use **afero** (`github.com/spf13/afero`) for filesystem operations to enable in-memory testing
- **Dependency injection**: Functions accept `afero.Fs` and `LookupEnvFunc` interfaces for testing
    - Example: `loadActionContext(fs afero.Fs, lookupEnv lookupEnvFunc)`
    - Example: `buildFileCommand(fs, lookupEnv, uuidFunc, ...)`
- Prefer **interfaces over concrete types** for function parameters (enables mocking)
- Test both success and error paths (empty env vars, invalid input, UUID generation errors)