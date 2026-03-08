# Referencia de Codigo

## Domain Types

### RuntimeType

```go
type RuntimeType string

const (
    Container RuntimeType = "container"  // k3d
    VPS       RuntimeType = "vps"        // k3s
)
```

### State

```go
type State struct {
    Version   string     `json:"version"`
    Cluster   *Cluster   `json:"cluster,omitempty"`
    Nodes     []Node     `json:"nodes"`
    Resources []Resource `json:"resources"`
}
```

### Cluster

```go
type Cluster struct {
    Id        string `json:"id"`
    Name      string `json:"name"`
    Mode      string `json:"mode"`
    CreatedAt string `json:"createdAt"`
}
```

### Node

```go
type Node struct {
    Id       string `json:"id"`
   Hostname string `json:"hostname"`
    Role     string `json:"role"`
    Ip       string `json:"ip"`
}
```

### Resource

```go
type Resource struct {
    Id        string             `json:"id"`
    Name      string             `json:"name"`
    Kind      string             `json:"kind"`
    Runtime   string             `json:"runtime"`
    NodeId    string             `json:"nodeId"`
    Spec      resource.Spec      `json:"spec"`
    Status    resource.Status    `json:"status"`
    Ownership resource.Ownership `json:"ownership"`
    CreatedAt string             `json:"createdAt"`
    UpdatedAt string             `json:"updatedAt"`
}
```

## DTOs

### CommandResult

```go
type CommandResult struct {
    Stdout string `json:"stdout,omitempty"`
    Stderr string `json:"stderr,omitempty"`
    Error  string `json:"error,omitempty"`
}
```

### InspectDto

```go
type InspectDto struct {
    Found     bool
    Version   string
    Mode      string
    NodeCount int
    ResCount  int
}
```

## Interfaces

### CommandRunner

```go
type CommandRunner interface {
    Run(name string, args ...string) (*dto.CommandResult, error)
}
```

### PrivilegedRunner

```go
type PrivilegedRunner interface {
    Run(name string, args ...string) (*dto.CommandResult, error)
}
```

### Repository

```go
type Repository interface {
    Save(state *domain.State) error
    Load() (*domain.State, error)
    Exists() (bool, error)
}
```

### Logger

```go
type Logger interface {
    Info(message string)
    Error(message string)
    Debug(message string)
}
```

## Services

### InstallService

```go
type InstallService struct {
    repo          interfaces.Repository
    runtimes      []run_install.RuntimeInstaller
    nodeInspector system.Inspector
    builder       *StateBuilder
    logger        interfaces.Logger
}

func NewInstallService(
    repo interfaces.Repository,
    runtimes []run_install.RuntimeInstaller,
    nodeInspector system.Inspector,
    builder *StateBuilder,
    logger interfaces.Logger,
) *InstallService

func (s *InstallService) Install() error
```

### RuntimeInstaller

```go
type RuntimeInstaller interface {
    Install() (*domain.Snapshot, error)
    Inspect() ([]*domain.Node, error)
    Detect() (bool, error)
}
```

## Infraestructura

### K8s Installer

```go
type Installer struct {
    runner  interfaces.CommandRunner
    sudo    interfaces.PrivilegedRunner
    runtime domain.RuntimeType
    logger  interfaces.Logger
}

func NewInstaller(runner, sudo, runtime, logger) *Installer

func (i *Installer) InstallKubectl() error
func (i *Installer) InstallCluster() error   // K3s o K3d segun runtime
func (i *Installer) ConfigureCluster() error // Configura kubectl
```

### Docker Installer

```go
type Installer struct {
    runner  interfaces.CommandRunner
    sudo    interfaces.PrivilegedRunner
    runtime domain.RuntimeType
    logger  interfaces.Logger
}

func (i *Installer) InstallDocker() error
func (i *Installer) StartDocker() error
func (i *Installer) ConfigureDocker(user string) error
```

### ExecRunner

```go
type ExecRunner struct {
    Logger interfaces.Logger
}

func NewExecRunner(logger interfaces.Logger) *ExecRunner

func (r *ExecRunner) Run(name string, args ...string) (*dto.CommandResult, error)
```
