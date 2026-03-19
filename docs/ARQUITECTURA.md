# Arquitectura

## Visión General

Syncloud **NO administra Docker o Kubernetes directamente**. Administra **Syncloud Resources** que luego se reconcilian hacia distintos runtimes.

```
┌─────────────────────────────────────────────────────────────────┐
│                    syncloud-state.json                          │
│                    (FUENTE DE VERDAD)                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                    ┌─────────┴─────────┐
                    │                   │
                    ▼                   ▼
         ┌──────────────────┐  ┌──────────────────┐
         │  Estado Deseado  │  │  Estado Actual   │
         │  (declared)      │  │  (observed)      │
         └────────┬─────────┘  └────────┬─────────┘
                  │                      │
                  └──────────┬───────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │      DIFF        │
                    │ (desired vs actual)│
                    └────────┬─────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │   Reconciler     │
                    └────────┬─────────┘
                             │
              ┌──────────────┴──────────────┐
              │                             │
              ▼                             ▼
     ┌─────────────────┐          ┌─────────────────┐
     │ DockerReconciler│          │  K8sReconciler  │
     └────────┬────────┘          └────────┬────────┘
              │                             │
              ▼                             ▼
     ┌─────────────────┐          ┌─────────────────┐
     │  docker run /   │          │  kubectl apply  │
     │  docker APIs    │          │    k8s APIs     │
     └─────────────────┘          └─────────────────┘
```

**syncloud-state.json es el centro de todo:**
- Define el **estado deseado** de todos los recursos
- Es leído por el **Reconciler** para comparar con el estado actual
- Se actualiza después de cada operación (`apply`, `deploy`, `delete`)
- Persiste entre sesiones para sobrevida de reinicios

## Estructura del Proyecto

```
synctl/
├── cmd/                    → Comandos CLI (Cobra)
├── internal/
│   ├── app/                → Componentes centralizados (factory)
│   ├── application/        → Interfaces y servicios
│   │   ├── interfaces/     → Puertos (contratos)
│   │   ├── services/       → Lógica de negocio
│   │   └── outputs/        → Formateadores de salida
│   ├── domain/             → Modelos universales
│   ├── infrastructure/     → Adaptadores de runtime
│   │   ├── docker/         → Adaptador Docker
│   │   ├── k8s/            → Adaptador Kubernetes
│   │   └── dns/            → Sistema DNS
│   ├── executor/           → Ejecutor de comandos
│   └── logger/             → Logging
├── docs/                   → Documentación
└── examples/               → Ejemplos de recursos
```

## Capas

### cmd/ - Capa de Presentación

Punto de entrada CLI con comandos Cobra:

| Comando | Descripción |
|---------|-------------|
| `install` | Instala Syncloud y sus dependencias |
| `apply` | Aplica Syncloud Resources desde YAML |
| `inspect` | Inspecciona el estado actual |
| `get` | Consulta recursos y nodos |
| `daemon` | Modo daemon (futuro) |
| `version` | Muestra versión |

### internal/domain/ - Modelo de Recursos (Universal)

Define el modelo abstracto de recursos que es independiente del runtime:

```go
type Resource struct {
    Id        string                 // Identificador único
    Name      string                 // Nombre del recurso
    Kind      string                 // Tipo abstracto (Container, Deployment, etc.)
    Runtime   string                 // Runtime destino (docker/kubernetes)
    NodeId    string                 // Nodo donde se ejecuta
    Spec      map[string]interface{} // Especificación universal
    Status    map[string]interface{} // Estado observado
    Ownership Ownership              // Propiedad del recurso
    CreatedAt string
    UpdatedAt string
}
```

**Principios del modelo:**
- Los recursos son **declarativos** y **abstractos**
- Un mismo recurso puede targetear Docker o Kubernetes
- El `Kind` define la semántica, no la implementación

### internal/application/interfaces/ - Puertos

Contratos que permiten la abstracción de runtimes:

```go
// RuntimeReconciler: Interfaz para adaptadores de runtime
type RuntimeReconciler interface {
    Runtime() string                    // "docker" | "kubernetes"
    Observe(resource *Resource) (*RuntimeResourceState, error)
    Reconcile(ctx ReconcileContext) error
}

// ResourceParser: Parsea YAML a recursos del dominio
type ResourceParser interface {
    Parse(path string) ([]*Resource, error)
}
```

### internal/application/services/ - Servicios

Contiene la lógica de negocio organizada por flujo:

#### Flujo de Apply
```
ApplyService
  ├─► ResourceParser     → Parsea YAML a recursos universales
  ├─► KindResolver       → Mapea Kind a runtime específico
  ├─► Validator          → Valida recursos antes de aplicar
  ├─► Repository         → Persiste estado
  └─► RuntimeReconciler → Reconcilia contra runtime
        ├─► DockerReconciler
        └─► K8sReconciler
```

#### Flujo de Install
```
InstallService
  ├─► RuntimeInstaller   → Instala dependencias
  │     ├─► K8sInstallService  (K3s/K3d)
  │     ├─► DockerInstallService
  │     └─► DnsInstallService
  └─► StateBuilder       → Construye estado inicial
```

### internal/infrastructure/ - Adaptadores de Runtime

Implementaciones concretas de los puertos:

| Paquete | Responsabilidad |
|---------|-----------------|
| `docker/` | docker run, docker API para containers, networks, images |
| `k8s/` | kubectl apply, k8s API para deployments, services, etc. |
| `dns/` | dnsmasq para resolución de dominios |
| `system/` | Detección de entorno, inspector del sistema |

## Flujo de Apply

```
┌─────────────────────────────────────────────────────────────────┐
│                      synctl apply -f resources.yaml              │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│  1. Parse YAML → ResourceYAML[]                                   │
│     └─► apiVersion: syncloud/v1, kind: Container, spec: {...}   │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│  2. KindResolver.Resolve(kind) → {runtime, kind, mode}          │
│     └─► "Container" → {runtime: "docker", kind: "docker.container"}│
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│  3. ResourceBuilder.Build() → *domain.Resource                    │
│     └─► Modelo universal con Spec abstracto                      │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│  4. Validator.Validate() → Validación semántica                │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│  5. Repository.Upsert() → Persistencia de estado                 │
└─────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│  6. RuntimeReconciler.Reconcile()                                │
│                                                                  │
│     ┌─────────────┐         ┌─────────────┐                      │
│     │   Docker    │         │    K8s      │                      │
│     │ Reconciler  │         │ Reconciler  │                      │
│     └──────┬──────┘         └──────┬──────┘                      │
│            │                       │                             │
│            ▼                       ▼                             │
│     ┌─────────────┐         ┌─────────────┐                      │
│     │ docker run  │         │kubectl apply│                      │
│     │ docker api │         │  k8s api   │                      │
│     └─────────────┘         └─────────────┘                      │
└─────────────────────────────────────────────────────────────────┘
```

## Recursos Soportados

### Modelo Universal (syncloud/v1)

| Kind | Runtime | Descripción |
|------|---------|-------------|
| `Container` | Docker | Contenedor abstracto |
| `Network` | Docker | Red de contenedores |
| `Image` | Docker | Imagen Docker |
| `Deployment` | Kubernetes | Deployment K8s |
| `Service` | Kubernetes | Servicio K8s |
| `Ingress` | Kubernetes | Ingress K8s |
| `Namespace` | Kubernetes | Namespace K8s |
| `Secret` | Kubernetes | Secret K8s |
| `ConfigMap` | Kubernetes | ConfigMap K8s |
| `ServiceAccount` | Kubernetes | ServiceAccount K8s |
| `Role` | Kubernetes | Role RBAC |
| `RoleBinding` | Kubernetes | RoleBinding RBAC |
| `ClusterRole` | Kubernetes | ClusterRole RBAC |
| `Endpoints` | Kubernetes | Endpoints |
| `Pod` | Kubernetes | Pod individual |
| `K8sResource` | Kubernetes | Recurso raw (passthrough) |

## Extension

### Agregar nuevo comando

1. Crear `cmd/<comando>/<comando>.go`
2. Definir `var Cmd = &cobra.Command{...}`
3. Usar `app.NewComponents()` para inicializar dependencias
4. Registrar en `cmd/root.go`

### Agregar nuevo Kind

1. Definir mapeo en `KindResolver` (`parser/kind_resolver.go`)
2. Implementar reconciler si es necesario
3. Agregar validación en `Validator`

### Agregar nuevo Runtime

1. Implementar `RuntimeReconciler` en `internal/application/services/apply/reconciler/<runtime>/`
2. Registrar en `ApplyService`
3. Agregar mapeo en `KindResolver`

## Dependencias

```
github.com/spf13/cobra    → CLI framework
github.com/google/uuid    → Generación de UUIDs
go.yaml.in/yaml/v3        → Parsing YAML
modelcontextprotocol/go-sdk → Protocolo MCP (futuro)
```
