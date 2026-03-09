# AGENTS.md - Development Guide for AI Agents

## Project Overview

Go CLI (`synctl`) for managing Syncloud resources using kubernetes, docker, dnsmasq.
- **Module**: `synctl`
- **Go Version**: 1.25.0
- **Dependencies**: `github.com/spf13/cobra`, `github.com/google/uuid`

---

## Build Commands

```bash
go build -o synctl .
go run .
go run . install
go run . get nodes
GOOS=linux GOARCH=amd64 go build -o synctl-linux-amd64 .
```

---

## Test Commands

```bash
go test ./...
go test -v ./...
go test -v -run TestName ./...
go test -v ./internal/domain/...
go test -cover ./...
go test -race ./...
```

---

## Lint and Format

```bash
go fmt ./...
go vet ./...
golangci-lint run ./...
go build ./... 2>&1 | head -50
```

---

## Code Style

### Import Organization

Three groups separated by blank lines:
1. Standard library
2. External packages
3. Internal packages (synctl/...)

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

Use aliases for conflicts:
```go
import (
    nodes "synctl/cmd/get/nodes"
)
```

### Naming

- **Files**: snake_case (`install_service.go`)
- **Packages**: lowercase (`install`, `dns`)
- **Types/Interfaces**: PascalCase (`InstallService`)
- **Functions**: PascalCase (`NewInstallService`)
- **Variables**: camelCase (`nodeInspector`)
- **Interfaces**: Suffix with `er` (`Installer`, `Inspector`)

### Package Structure

```
cmd/                     # Cobra commands
internal/
  ├── application/       # Services & interfaces
  │   ├── interfaces/   # Port interfaces
  │   ├── services/     # Service implementations
  │   └── outputs/      # Output formatters
  ├── domain/           # Entities, DTOs
  ├── infrastructure/   # Adapters (docker, k8s, dns, system)
  └── logger/           # Logging
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

Early returns with wrapped errors:

```go
func (s *InstallService) Install() error {
    for _, runtime := range s.runtimes {
        snap, err := runtime.Install()
        if err != nil {
            return err
        }
    }
    return s.repo.Save(state)
}
```

### Constructor Functions

```go
func NewInstaller(
    runner interfaces.CommandRunner,
    sudo interfaces.PrivilegedRunner,
    logger interfaces.Logger,
) *Installer {
    return &Installer{
        runner: runner,
        sudo:   sudo,
        logger: logger,
    }
}
```

### Struct Tags

```go
type State struct {
    Version   string     `json:"version"`
    Cluster   *Cluster   `json:"cluster,omitempty"`
    Nodes     []Node     `json:"nodes"`
}
```

### Logging

Use `synctl/internal/application/interfaces/logger.go`. Use `logger.Info()`, `logger.Error()`. Don't log and return errors.

### Cobra Commands

```go
var Cmd = &cobra.Command{...}

func init() {
    RootCmd.AddCommand(Cmd)
}
```

Keep logic in services, not cmd packages.

---

## Workflow

Before committing: `go fmt ./...` and `go vet ./...`. Test: `go test ./...`. Build check: `go build -o /dev/null .`

---

## Adding New Commands

1. Create `cmd/newcommand/`
2. Add `var Cmd = &cobra.Command{...}`
3. Register in parent's `init()`
4. Keep logic in `internal/application/services/`

## Adding New Services

1. Define interface in `internal/application/interfaces/`
2. Implement in `internal/infrastructure/`
3. Use `New*()` constructor returning `*Type`
4. Inject dependencies through constructor
