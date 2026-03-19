# Syncloud CLI

CLI para desplegar y gestionar aplicaciones en Docker y Kubernetes con un modelo de recursos declarativo y universal.

## Quick Start

```bash
# 1. Instalar Syncloud (Docker, K3s, dnsmasq)
synctl install

# 2. Desplegar la plataforma completa
synctl deploy

# 3. Acceder (URL mostrada en terminal)
# http://syncloud.local
```

## Flujo MVP v1.0

```
┌─────────────────────────────────────────────────────────────────────┐
│  synctl install → Configurar runtimes (Docker/K3s + DNS)           │
│                         ↓                                           │
│  synctl deploy  → Desplegar DB + Backend + WebApp                  │
│                         ↓                                           │
│  synctl get dns / nodes / resources → Gestionar plataforma         │
│                         ↓                                           │
│  http://syncloud.local → Abrir en navegador                        │
└─────────────────────────────────────────────────────────────────────┘
```

El archivo `syncloud-state.json` es el centro del sistema:
- **Estado deseado**: Recursos declarados
- **Estado actual**: Recursos observados
- **Diff**: Diferencia para reconciliar
- **DNS Records**: Registros para acceso
- **Nodes**: Nodos del cluster

## Gestión Completa de Objetos

### Resources (CRUD)
```bash
synctl apply -f resources.yaml      # Crear/actualizar
synctl get resources                # Listar
synctl resource describe <name>     # Ver detalle
synctl resource delete <name>      # Eliminar
synctl resource logs <name>        # Ver logs
synctl resource exec <name> -- <cmd> # Ejecutar comando
synctl resource scale <name> --replicas=3 # Escalar
```

### DNS Records (CRUD)
```bash
synctl get dns                     # Listar registros
synctl dns add api --server 192.168.1.101  # Agregar
synctl dns delete api              # Eliminar
synctl dns update api --server 192.168.1.102 # Actualizar
```

### Nodes (CRUD)
```bash
synctl get nodes                   # Listar nodos
synctl node describe <name>        # Ver detalle
synctl node add -f node.yaml       # Agregar nodo
synctl node remove <name>          # Remover nodo
synctl node cordon <name>          # Deshabilitar scheduling
synctl node uncordon <name>       # Habilitar scheduling
```

## Comandos

| Comando | Descripción |
|---------|-------------|
| `synctl install` | Instala runtimes y configura el sistema |
| `synctl deploy` | Despliega la plataforma completa |
| `synctl apply -f <file>` | Aplica recursos desde YAML |
| `synctl inspect` | Inspecciona el estado actual |
| `synctl get resources` | Lista recursos con filtros |
| `synctl get nodes` | Lista nodos del sistema |
| `synctl get dns` | Lista registros DNS |
| `synctl version` | Muestra versión |

## Documentación

| Documento | Descripción |
|-----------|-------------|
| [docs/MVP.md](docs/MVP.md) | Flujo MVP y arquitectura del state |
| [docs/ARQUITECTURA.md](docs/ARQUITECTURA.md) | Arquitectura completa del sistema |
| [docs/PENDIENTES.md](docs/PENDIENTES.md) | Roadmap de desarrollo v1.0 |
| [docs/README.md](docs/README.md) | Guía de uso detallada |
| [AGENTS.md](AGENTS.md) | Guía para desarrollo |

## Recursos Soportados

**Docker:** Container, Network, Image

**Kubernetes:** Deployment, Service, Ingress, Namespace, Secret, ConfigMap, ServiceAccount, Role, RoleBinding, ClusterRole, Pod, Endpoints, K8sResource

## Build

```bash
go build -o synctl .
./synctl --help
```
