# AGENTS.md - Development Guide for AI Agents

## Project Overview

Go CLI (`synctl`) for managing Syncloud resources using Cobra and Clean Architecture.

- **Module**: `synctl`
- **Go Version**: 1.25.0
- **Dependencies**: `github.com/spf13/cobra`, `github.com/google/uuid`

---

## Build, Lint, and Test Commands

```bash
# Build
go build -o synctl .
go run .
go install .

# Test (all packages)
go test ./...
go test -v ./...
go test -cover ./...

# Test specific function
go test -v -run TestFunctionName ./path/to/package

# Test specific file
go test -v ./path/to/package -run TestFunctionName

# Lint
go vet ./...
go fmt ./...
golangci-lint run ./...
```

---

## Project Structure

```
cmd/                         # CLI commands (Cobra)
├── install/                 # Install command
├── inspect/                 # Inspect command
├── get/                     # Get commands (nodes, resources)
└── version/                 # Version command

internal/
├── application/
│   ├── interfaces/          # Port interfaces
│   │   ├── install/         # Install-related interfaces
│   │   ├── system/         # System interfaces
│   │   ├── command_runner.go
│   │   ├── logger.go
│   │   └── repository.go
│   └── services/            # Application services
│       └── install/         # Install service implementations
├── domain/                  # Models, DTOs, types
│   ├── dto/                 # Data transfer objects
│   └── resource/            # Resource models
├── infrastructure/          # Adapters (implementations)
│   ├── docker/
│   ├── k8s/
│   └── system/
├── executor/                # Command execution
└── logger/                  # Logging implementation
```

---

## Code Style Guidelines

### Naming Conventions

- **Packages**: lowercase, short names (e.g., `cmd`, `k8s`, `dto`)
- **Types/Interfaces**: PascalCase (`InstallService`, `RuntimeInstaller`)
- **Variables/Fields**: camelCase (`runtimeType`, `nodeInspector`)
- **Interfaces**: Use `er` suffix (`Reader`, `Installer`, `Detector`)
- **Constructor Functions**: `New<TypeName>` pattern
- **Files**: lowercase with underscores for multi-word packages (`install_service.go`)

### Imports

Group imports with blank lines, ordered by category:

```go
import (
    "fmt"
    "os"

    // domain
    "synctl/internal/domain"
    dto "synctl/internal/domain/dto"

    // application interfaces
    "synctl/internal/application/interfaces"
    run_install "synctl/internal/application/interfaces/install"

    // infrastructure
    "synctl/internal/infrastructure/k8s"

    // external
    "github.com/spf13/cobra"
)
```

### Types and Constants

Use type-specific constants and string types for enums:

```go
type RuntimeType string

const (
    K3s RuntimeType = "k3s"
    K3d RuntimeType = "k3d"
)

type InstallMode string

const (
    System InstallMode = "system"
    User   InstallMode = "user"
)
```

### Struct Tags

Use JSON tags for DTOs and command results:

```go
type CommandResult struct {
    Stdout string `json:"stdout,omitempty"`
    Stderr string `json:"stderr,omitempty"`
    Error  string `json:"error,omitempty"`
}
```

### Interfaces and Dependency Injection

Define interfaces in `internal/application/interfaces`, implement in infrastructure:

```go
// internal/application/interfaces/command_runner.go
type CommandRunner interface {
    Run(name string, args ...string) (*dto.CommandResult, error)
}

// internal/infrastructure/k8s/installer.go
type K8sInstaller struct {
    runner interfaces.CommandRunner
}
```

### Constructor Pattern

Always use constructor functions with explicit dependencies:

```go
func NewInstallService(
    repo interfaces.Repository,
    runtimes []run_install.RuntimeInstaller,
    nodeInspector system.Inspector,
    builder *StateBuilder,
    logger interfaces.Logger,
) *InstallService {
    return &InstallService{
        repo:          repo,
        runtimes:      runtimes,
        nodeInspector: nodeInspector,
        builder:       builder,
        logger:        logger,
    }
}
```

### Error Handling

Return errors directly, wrap with context using `fmt.Errorf`:

```go
func (s *InstallService) Install() error {
    snapshot, err := runtime.Install()
    if err != nil {
        return fmt.Errorf("runtime install failed: %w", err)
    }
    return nil
}
```

### Logging

Use the Logger interface from `internal/application/interfaces/logger.go`:

```go
type Logger interface {
    Info(message string)
    Error(message string)
    Debug(message string)
}

// Usage
s.logger.Info("Syncloud Platform installed successfully!!!")
```

### Cobra Commands

Define commands as package-level variables with `RunE`:

```go
var Cmd = &cobra.Command{
    Use:   "install",
    Short: "Install Syncloud",
    RunE: func(cmd *cobra.Command, args []string) error {
        service := installservice.NewInstallService(repo, ...)
        return service.Install()
    },
}
```

Register commands in `cmd/root.go` `init()` function:

```go
func init() {
    RootCmd.AddCommand(version.Cmd)
    RootCmd.AddCommand(install.Cmd)
}
```

---

## Development

- Dev container: `.devcontainer/` (Go 1.25+, Docker for k3d)
- Run locally: `go mod download && go build -o synctl . && ./synctl --help`

---

## Key Files

- `main.go` - Entry point
- `cmd/root.go` - Root command, command registration
- `go.mod` / `go.sum` - Dependencies

---

## Common Patterns

### Repository Pattern

```go
type Repository interface {
    Save(state *domain.State) error
    Load() (*domain.State, error)
}
```

### Service Layer

Services orchestrate domain logic and use interfaces for external dependencies:

```go
type ServiceName struct {
    dependency1 Interface1
    dependency2 Interface2
}
```
