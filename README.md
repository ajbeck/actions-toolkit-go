# GitHub Actions Toolkit for Go

A Go library for building GitHub Actions, inspired by the official [JavaScript/TypeScript toolkit](https://github.com/actions/toolkit).

[![GoDoc](https://pkg.go.dev/badge/mod/github.com/ajbeck/actions-toolkit-go)](https://pkg.go.dev/mod/github.com/ajbeck/actions-toolkit-go)

## Why?

GitHub Actions have become a popular tool, with most actions being either JavaScript/TypeScript or Docker-based. I've found a more natural fit (and more performant during a workflow) to be publishing Go binaries in the action repo and triggering them via a JavaScript action and script to call the binary. Given I'm doing that often, I wanted a more idiomatic way to interact with GitHub Actions from Go code. There are some existing libraries, such as [https://github.com/sethvargo/go-githubactions](https://github.com/sethvargo/go-githubactions) which are a great fit for some. However, I wanted something which has roots more in line with the official toolkit, making it easier to switch between languages if needed.

### 🔑 Key Features

- **Simplicity**: Easy-to-use API, which is patterned after the official GitHub Actions toolkit.
- **Comprehensive**: Covers a wide range of GitHub Actions functionalities, including inputs, outputs, environment variables, logging, and more.
- **Go Native**: Leverages Go's strengths, making it a natural choice for Go developers building GitHub Actions.

## Installation

```bash
go get github.com/ajbeck/actions-toolkit-go
```

## Documentation

For complete API documentation, see [pkg.go.dev/github.com/ajbeck/actions-toolkit-go](https://pkg.go.dev/github.com/ajbeck/actions-toolkit-go).
For contributor workflow and testing conventions, see `AGENTS.md`.

## License

MIT
