# Actions Toolkit for Go

The Actions Toolkit for Go is a library that simplifies building GitHub Actions using the Go programming language. It provides utilities to interact with GitHub Actions runtime, manage workflow commands, and handle file-based commands for environment variables, outputs, and state. The toolkit is based on the official GitHub Actions toolkit for Javascript/TypeScript, its code is here: https://github.com/actions/toolkit. Github provides 10 distinct packages in the toolkit, to understand the full scope of functionality, see their documentation. Each package is in its own folder within this path https://github.com/actions/toolkit/tree/main/packages.

## Project Structure & Module Organization

- `actioncontext/`: helpers to read GitHub Actions runtime env vars and expose strongly typed context data.
- `core/`: workflow commands (notice, warning, error, grouping), file commands (env/output/path/state), and shared utilities.
- `go.mod` / `go.sum`: module definition for `github.com/ajbeck/actions-toolkit-go`.
- Tests live alongside source files using the Go `_test.go` convention; add new tests in the same package folders.

## Coding Style & Naming Conventions

- Follow standard Go style: tabs for indentation, CamelCase for exported names, and short, descriptive receiver names.
- Keep packages cohesive and prefer small, focused files.
- Avoid panics in library code; return errors and use `fmt.Errorf` with context.
- Dependency injection should be used where possible to facilitate testing and flexibility.
- Where possible prefer Interfaces over Types for function parameters to allow for easier mocking and testing.

## Code Comments

Code added to the Actions Toolkit for Go should follow these guidelines to ensure consistency, maintainability, and quality. Code should be documented following the Go Lang guidelines here: https://go.dev/doc/comment. Documentation should be validated by using the `go doc` command to generate and review the output. To understand how to use the `go doc` read the usage documentation output by the command `go help doc`.

## Testing Guidelines

- Add `_test.go` files next to the code under test; prefer table-driven tests for variations.
- Whenever possible mock external dependencies (file system, env vars) using interfaces.
- When interacting with the file system use the `afero` library to allow for in-memory testing.
- Aim for coverage on new branches/edge cases (empty env vars, invalid input, UUID errors).
- Run `go test ./...` locally; include failing scenarios and assertions for printed workflow command output where relevant.

## Developer Workflow Commands

Commands for development workflow must always be executed from the root of the repository. Explore the go cli documentation with the help command to better understand each command and its flags: https://go.dev/doc/cmd.

- `go fmt ./...`: format all Go files; run before committing.
- `go vet ./...`: static checks for common mistakes.
- `GOCACHE=$(pwd)/.gocache/build GOMODCACHE=$(pwd)/.gocache/mod go test ./...`: run the full unit test suite with a repo-local build cache to avoid permission issues on shared systems.
- `GOCACHE=$(pwd)/.gocache/build GOMODCACHE=$(pwd)/.gocache/mod go build ./...`: compile all packages to ensure imports and types stay sound using the local cache path.
- `go mod tidy`: clean dependencies when you add/remove imports.

## Tools

- When working with github repositories use the `gh` cli tool.
