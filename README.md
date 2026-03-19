# Syncloud CLI

CLI para desplegar y gestionar aplicaciones en Docker y Kubernetes con un modelo de recursos declarativo.

## Quick Start

```bash
# 1. Instalar Syncloud (Docker, K3s, dnsmasq)
synctl install

# 2. Desplegar la plataforma completa
synctl deploy

# 3. Acceder (URL mostrada en terminal)
# http://syncloud.local
```

## Flujo MVP

```
┌─────────────────────────────────────────────────────────────┐
│  synctl install → Configura runtimes (Docker/K3s + DNS)    │
│                         ↓                                   │
│  synctl deploy  → Despliega DB + Backend + WebApp          │
│                         ↓                                   │
│  http://syncloud.local → Abre en navegador                 │
└─────────────────────────────────────────────────────────────┘
```

El archivo `syncloud-state.json` es el centro del sistema:
- **Estado deseado**: Recursos declarados
- **Estado actual**: Recursos observados
- **Diff**: Diferencia para reconciliar

## Comandos

| Comando | Descripción |
|---------|-------------|
| `synctl install` | Instala runtimes y configura el sistema |
| `synctl deploy` | Despliega la plataforma completa |
| `synctl apply -f <file>` | Aplica recursos desde YAML |
| `synctl inspect` | Inspecciona el estado actual |
| `synctl get resources` | Lista recursos con filtros |
| `synctl get nodes` | Lista nodos del sistema |
| `synctl version` | Muestra versión |

## Documentación

| Documento | Descripción |
|-----------|-------------|
| [docs/MVP.md](docs/MVP.md) | Flujo MVP y arquitectura del state |
| [docs/ARQUITECTURA.md](docs/ARQUITECTURA.md) | Arquitectura completa del sistema |
| [docs/README.md](docs/README.md) | Guía de uso detallada |
| [docs/PENDIENTES.md](docs/PENDIENTES.md) | Roadmap de desarrollo |
| [AGENTS.md](AGENTS.md) | Guía para desarrollo |

## Recursos Soportados

**Docker:** Container, Network, Image

**Kubernetes:** Deployment, Service, Ingress, Namespace, Secret, ConfigMap, ServiceAccount, Role, RoleBinding, ClusterRole, Pod, Endpoints, K8sResource

## Build

```bash
go build -o synctl .
./synctl --help
```
