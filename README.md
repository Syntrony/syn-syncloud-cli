# Syncloud CLI

CLI para gestionar aplicaciones en Docker y Kubernetes usando un modelo de recursos declarativo y universal.

## Quick Start

```bash
# Instalar Syncloud
synctl install

# Aplicar recursos
synctl apply -f examples/resources.yaml

# Ver recursos
synctl get resources

# Inspeccionar estado
synctl inspect
```

## Arquitectura

Syncloud no administra Docker/Kubernetes directamente. Administra **Syncloud Resources** que se reconcilian hacia el runtime apropiado:

```
synctl apply → Syncloud Resource → Runtime Adapter → Docker/Kubernetes
```

Ver [docs/ARQUITECTURA.md](docs/ARQUITECTURA.md) para detalles.

## Documentación

| Documento | Descripción |
|-----------|-------------|
| [docs/README.md](docs/README.md) | Guía de uso y comandos |
| [docs/ARQUITECTURA.md](docs/ARQUITECTURA.md) | Arquitectura del sistema |
| [docs/REFERENCIA.md](docs/REFERENCIA.md) | Referencia de código y tipos |
| [docs/PENDIENTES.md](docs/PENDIENTES.md) | Tareas pendientes y roadmap |
| [AGENTS.md](AGENTS.md) | Guía para desarrollo |

## Recursos Soportados

**Docker:** Container, Network, Image

**Kubernetes:** Deployment, Service, Ingress, Namespace, Secret, ConfigMap, ServiceAccount, Role, RoleBinding, ClusterRole, Pod, Endpoints, K8sResource

## Build

```bash
go build -o synctl .
./synctl --help
```
