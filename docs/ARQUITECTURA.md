# Arquitectura

## Vision General

synctl sigue los principios de Clean Architecture con separacion en capas:

```
cmd/               → Comandos CLI (Cobra)
internal/
├── application/   → Interfaces y servicios
├── domain/        → Modelos y DTOs
├── infrastructure/ → Implementaciones
├── executor/      → Ejecucion de comandos
└── logger/       → Logging
```

## Capas

### cmd/ - Capa de Presentacion

Contiene los comandos Cobra que son el punto de entrada de la CLI.

- **install**: Instala Syncloud en el cluster
- **inspect**: Inspecciona el estado actual
- **get**: Obtiene recursos (nodes, resources)
- **version**: Muestra version

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
}

type Logger interface {
    Info(message string)
    Error(message string)
    Debug(message string)
}

type PrivilegedRunner interface {
    Run(name string, args ...string) (*dto.CommandResult, error)
}
```

### internal/application/services/ - Servicios

Contiene la logica de negocio:

- **InstallService**: Orquesta la instalacion
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

### internal/infrastructure/ - Adaptadores

Implementaciones concretas de las interfaces:

- **k8s/**: Instalador K3s/K3d, detector, inspector
- **docker/**: Instalador Docker, detector, inspector
- **system/**: Detector de entorno, inspector

### internal/executor/ - Ejecutor

- **ExecRunner**: Ejecuta comandos del sistema
- **SudoRunner**: Ejecuta comandos con sudo

## Flujo de instalacion

1. **Detectar entorno**: Container vs VPS
2. **Detectar si requiere sudo**
3. **Instalar Docker** (si no existe, VPS only)
4. **Instalar K3s/K3d** segun tipo de entorno
5. **Configurar kubectl**
6. **Guardar estado** en `helpers/syncloud-state.json`

## Dependencias

```
github.com/spf13/cobra  → CLI
github.com/google/uuid  → UUIDs
```

## Extension

Para agregar un nuevo comando:

1. Crear directorio en `cmd/<comando>/`
2. Definir `Cmd` como `*cobra.Command`
3. Registrar en `cmd/root.go` `init()`

Para agregar nueva funcionalidad:

1. Definir interfaz en `internal/application/interfaces/`
2. Implementar en `internal/infrastructure/`
3. Consumir en servicios
