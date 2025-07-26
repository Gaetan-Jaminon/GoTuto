# Go Packages vs Modules and Project Initialization

## Overview
Understanding the difference between packages and modules is fundamental to Go development. This note covers the key concepts and how to properly initialize a Go project.

## Modules vs Packages

### Module
- A **module** is a collection of Go packages stored in a file tree with a `go.mod` file at its root
- A module defines the module path (import path prefix) and dependency requirements
- One module = one repository typically
- Introduced in Go 1.11 to replace GOPATH-based development
- Enables versioned dependency management

### Package
- A **package** is a collection of Go source files in the same directory
- All files in a package must have the same package name (declared with `package` keyword)
- Packages are the unit of compilation and encapsulation in Go
- A module can contain multiple packages in subdirectories

## Project Structure Example
```
myproject/                 <- Module root
├── go.mod                 <- Module definition
├── go.sum                 <- Dependency checksums
├── main.go                <- Package main
├── utils/                 <- Package utils
│   └── helper.go
└── api/                   <- Package api
    ├── handler.go
    └── routes.go
```

## Initializing a Go Project

### Step 1: Create Project Directory
```bash
mkdir myproject
cd myproject
```

### Step 2: Initialize Module
```bash
go mod init <module-name>
```

**Examples:**
- `go mod init github.com/username/myproject` (for GitHub)
- `go mod init example.com/myproject` (generic)
- `go mod init myproject` (local development)

### Step 3: Create main.go
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

### Step 4: Run the Project
```bash
go run main.go
# or
go run .
```

## Key Files Created

### go.mod File
```go
module github.com/username/myproject

go 1.21
```

- Defines the module path
- Specifies Go version
- Lists dependencies (added automatically when you import packages)

### go.sum File (created automatically)
- Contains checksums of dependencies
- Ensures reproducible builds
- Should be committed to version control

## Important Commands

| Command | Purpose |
|---------|---------|
| `go mod init <name>` | Initialize a new module |
| `go mod tidy` | Add missing and remove unused modules |
| `go mod download` | Download modules to local cache |
| `go mod verify` | Verify dependencies |
| `go build` | Compile packages and dependencies |
| `go run .` | Run the main package |

## Best Practices

1. **Module naming**: Use a domain-based path (e.g., `github.com/user/repo`)
2. **Package naming**: Use short, lowercase names
3. **One package per directory**: All `.go` files in a directory must belong to the same package
4. **Main package**: Entry point of executable programs
5. **Commit go.sum**: Always commit both `go.mod` and `go.sum` to version control

## Common Gotchas

- Package name should match directory name (except for `main` package)
- Cannot have multiple packages in the same directory
- Module path should be importable by others if you plan to share your code
- Use `go mod tidy` regularly to keep dependencies clean