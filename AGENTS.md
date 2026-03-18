# AGENTS.md - Development Guide for AI Agents

## Project Overview

Go CLI (`synctl`) for managing Syncloud resources using kubernetes, docker, dnsmasq.
- **Module**: `synctl`
- **Go Version**: 1.25.0
- **Dependencies**: `github.com/spf13/cobra`, `github.com/google/uuid`, `go.yaml.in/yaml/v3`

---

## Build Commands

```bash
go build -o synctl .              # Build binary
go run .                          # Run without building
go run . install                  # Run install command
go run . get nodes                # Run get nodes command
GOOS=linux GOARCH=amd64 go build -o synctl-linux-amd64 .  # Cross-compile
go build -o /dev/null .           # Build check only
```

---

## Test Commands

```bash
go test ./...                     # Run all tests
go test -v ./...                  # Verbose output
go test -v -run TestName ./...   # Run single test by name
go test -v ./internal/domain/... # Run tests in specific package
go test -cover ./...              # With coverage
go test -race ./...               # Race detector
```

---

## Lint and Format

```bash
go fmt ./...          # Format code
go vet ./...          # Vet checks
golangci-lint run ./...  # Full linting (if installed)
```

---

## Code Style

### Import Organization

Three groups separated by blank lines: stdlib, external, internal.

```go
import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/google/uuid"

    "synctl/cmd/get"
    "synctl/internal/domain"
)
```

Use aliases for conflicts: `nodes "synctl/cmd/get/nodes"`

### Naming Conventions

- **Files**: snake_case (`install_service.go`)
- **Packages**: lowercase (`install`, `dns`)
- **Types/Interfaces**: PascalCase (`InstallService`)
- **Functions**: PascalCase (`NewInstallService`)
- **Variables**: camelCase (`nodeInspector`)
- **Interfaces**: Suffix with `er` (`Installer`, `Inspector`)

### Package Structure

```
cmd/               # Cobra commands
internal/
  ├── application/ # Services & interfaces
  │   ├── interfaces/
  │   ├── services/
  │   └── outputs/
  ├── domain/     # Entities, DTOs
  ├── infrastructure/  # Adapters (docker, k8s, dns, system)
  └── logger/    # Logging
```

### Interface Definition

Define in `internal/application/interfaces/`, implement in `internal/infrastructure/`:

```go
type DnsInstaller interface {
    InstallDns() error
    ConfigureDns() error
}
```

### Error Handling

Early returns with wrapped errors. Never log and return errors.

```go
func (s *InstallService) Install() error {
    for _, runtime := range s.runtimes {
        if err := runtime.Install(); err != nil {
            return fmt.Errorf("installing runtime: %w", err)
        }
    }
    return s.repo.Save(state)
}
```

### Constructor Functions

```go
func NewInstaller(runner interfaces.CommandRunner, logger interfaces.Logger) *Installer {
    return &Installer{runner: runner, logger: logger}
}
```

### Logging

Use `synctl/internal/application/interfaces/logger.go`. Prefer `logger.Info()`, `logger.Error()`.

### Cobra Commands

```go
var Cmd = &cobra.Command{...}
func init() { RootCmd.AddCommand(Cmd) }
```

Keep business logic in services, not cmd packages.

---

## Workflow

Before committing:
```bash
go fmt ./... && go vet ./... && go test ./... && go build -o /dev/null .
```

---

## Adding New Commands

1. Create `cmd/newcommand/command.go`
2. Add `var Cmd = &cobra.Command{...}`
3. Register in parent's `init()`: `ParentCmd.AddCommand(Cmd)`
4. Keep logic in `internal/application/services/`

## Adding New Services

1. Define interface in `internal/application/interfaces/`
2. Implement in `internal/infrastructure/`
3. Use `New*()` constructor returning `*Type`
4. Inject dependencies through constructor
