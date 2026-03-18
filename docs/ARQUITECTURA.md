# Arquitectura

## Vision General

synctl sigue los principios de Clean Architecture con separacion en capas:

```
cmd/               → Comandos CLI (Cobra)
internal/
 ├── app/          → Componentes centralizados
 ├── application/  → Interfaces y servicios
 ├── domain/       → Modelos y DTOs
 ├── infrastructure/→ Implementaciones
 ├── executor/     → Ejecucion de comandos
 └── logger/       → Logging
```

## Capas

### cmd/ - Capa de Presentacion

Contiene los comandos Cobra que son el punto de entrada de la CLI.

- **install**: Instala Syncloud en el cluster
- **apply**: Aplica recursos desde YAML
- **inspect**: Inspecciona el estado actual
- **get**: Obtiene recursos (nodes, resources)
- **version**: Muestra version

### internal/app/ - Componentes Centralizados

Factory para inicializar componentes comunes:

```go
components := app.NewComponents()
components.Init()
components.WithDefaultRepo()
components.DetectSudo()
```

**Beneficios:**
- DRY: Elimina código duplicado de inicialización
- Single Source: Un solo lugar para cambiar rutas/configuraciones
- Testability: Fácil mock de dependencias

### internal/application/interfaces/ - Puertos

Define los contratos (interfaces) que la aplicacion necesita:

```go
type CommandRunner interface {
    Run(name string, args ...string) (*dto.CommandResult, error)
}

type Repository interface {
    Save(state *domain.State) error
    Load() (*domain.State, error)
    Exists() (bool, error)
    Upsert(resources []*domain.Resource) (*domain.Snapshot, error)
}

type Logger interface {
    Info(message string)
    Error(message string)
    Debug(message string)
}

type RuntimeReconciler interface {
    Runtime() string
    Observe(resource *domain.Resource) (*domain.RuntimeResourceState, error)
    Reconcile(context domain.ReconcileContext) error
}
```

### internal/application/services/ - Servicios

Contiene la logica de negocio:

**Install Flow:**
- **InstallService**: Orquesta la instalacion de runtimes
- **k8s.InstallService**: Instala K3s/K3d
- **docker.InstallService**: Instala Docker
- **dns.InstallService**: Instala dnsmasq

**Apply Flow:**
- **ApplyService**: Orquesta la aplicacion de recursos
- **ResourceParser**: Parsea YAML a Domain Resources
- **KindResolver**: Mapea tipos a runtimes
- **Validator**: Valida recursos
- **DockerReconciler**: Reconcilia recursos Docker
- **K8sReconciler**: Reconcilia recursos Kubernetes

**Query Flow:**
- **InspectService**: Inspecciona estado
- **GetNodeService**: Obtiene nodos
- **GetResourceService**: Obtiene recursos con filtros

### internal/domain/ - Entidades

Modelos del dominio:

- **State**: Estado de Syncloud (version, cluster, nodes, resources)
- **Cluster**: Informacion del cluster (id, name, mode)
- **Node**: Nodo del sistema (id, hostname, role, ip)
- **Resource**: Recurso (id, name, kind, runtime, nodeId, spec, status, ownership)
- **RuntimeType**: Tipo de entorno (Container, VPS)
- **ReconcileContext**: Contexto para reconciliacion

### internal/infrastructure/ - Adaptadores

Implementaciones concretas de las interfaces:

- **k8s/**: Instalador K3s/K3d, detector, inspector
- **docker/**: Instalador Docker, detector, inspector
- **system/**: Detector de entorno, inspector
- **dns/**: Instalador dnsmasq, detector, inspector

### internal/executor/ - Ejecutor

- **ExecRunner**: Ejecuta comandos del sistema
- **SudoRunner**: Ejecuta comandos con sudo

## Flujo de Apply (synctl apply -f resources.yaml)

```
1. Parse YAML → ResourceYAML[]
       ↓
2. KindResolver.Resolve() → {runtime, kind, mode}
       ↓
3. ResourceBuilder.Build() → *domain.Resource
       ↓
4. Validator.Validate() → spec validation
       ↓
5. Repository.Upsert() → Snapshot
       ↓
6. StateBuilder.Build() → *State
       ↓
7. Repository.Save(state)
       ↓
8. Para cada RuntimeReconciler:
   ├── Observe(resource) → RuntimeResourceState
   ├── Diff.Compare() → Action
   └── Reconcile(ctx) → Apply changes
```

### Recursos Soportados

**Docker:**
- Container
- Network
- Image

**Kubernetes:**
- Deployment, Namespace, Service, Ingress
- ServiceAccount, Role, RoleBinding, ClusterRole
- Secret, ConfigMap, Endpoints, Pod
- K8sResource (raw mode)

## Dependencias

```
github.com/spf13/cobra  → CLI
github.com/google/uuid  → UUIDs
go.yaml.in/yaml/v3      → YAML parsing
```

## Extension

Para agregar un nuevo comando:

1. Crear directorio en `cmd/<comando>/`
2. Usar `app.NewComponents()` para inicializar
3. Definir `Cmd` como `*cobra.Command`
4. Registrar en `cmd/root.go` `init()`

Para agregar un nuevo runtime:

1. Implementar `RuntimeReconciler` en `internal/application/services/apply/reconciler/<runtime>/`
2. Registrar en ApplyService
3. Agregar mapeo en `KindResolver`
