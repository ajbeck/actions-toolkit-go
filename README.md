# GitHub Actions Toolkit for Go

The Actions Toolkit for Go is a library that simplifies building GitHub Actions using the Go programming language. It provides utilities to interact with GitHub Actions runtime, manage workflow commands, and handle file-based commands for environment variables, outputs, and state. The toolkit is based on the official GitHub Actions toolkit for Javascript/TypeScript, its code is here: https://github.com/actions/toolkit.

## Usage

```bash
go get github.com/ajbeck/actions-toolkit-go
```

```go
import (
    "github.com/ajbeck/actions-toolkit-go/actioncontext"
    "github.com/ajbeck/actions-toolkit-go/core"
)
```

### actioncontext

This provides context values which are loaded from environment variables set by the GitHub Runner application. Details on the environment variables are here: https://docs.github.com/en/actions/reference/workflows-and-actions/variables

### core

This package provides an interface to use GitHub Workflow Commands, detials on the commands are here: https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions
