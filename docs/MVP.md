# Flujo de MVP - Syncloud Platform v1.0

## Overview

El objetivo es una **versión 1.0 completa y liberable** que permita gestionar todo el ciclo de vida de la plataforma Syncloud desde la terminal.

### Flujo de Usuario

```
1. synctl install     → Configurar cluster (Docker/K3s, dnsmasq)
2. synctl deploy      → Desplegar plataforma (DB + Backend + WebApp)
3. synctl get dns     → Ver registros DNS configurados
4. synctl get nodes   → Ver nodos del cluster
5. synctl get resources → Ver recursos desplegados
6. Acceder via URL    → http://syncloud.local
```

### Gestión Completa

```
┌─────────────────────────────────────────────────────────────────────┐
│                         syncloud-state.json                         │
│                     (FUENTE DE VERDAD CENTRAL)                      │
│                                                                      │
│  ┌──────────┐  ┌──────────┐  ┌───────────┐  ┌──────────────────┐   │
│  │ Cluster  │  │  Nodes    │  │ Resources │  │    DNS Records   │   │
│  │  Config  │  │  Manage   │  │   CRUD    │  │      CRUD        │   │
│  └──────────┘  └──────────┘  └───────────┘  └──────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                              │
          ┌───────────────────┼───────────────────┐
          │                   │                   │
          ▼                   ▼                   ▼
   ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
   │   Docker    │     │   K8s      │     │   dnsmasq   │
   │  Runtime    │     │  Runtime   │     │    DNS      │
   └─────────────┘     └─────────────┘     └─────────────┘
```

## Flujo Detallado

### Paso 1: Install

```bash
synctl install
```

**Qué hace:**
1. Detecta el entorno de ejecución (Container/VPS)
2. Instala Docker o K3s/K3d según el entorno
3. Instala y configura dnsmasq para resolución DNS
4. Genera `syncloud-state.json` con la configuración inicial
5. Registra el DNS provisional en `cluster.dns`

**Output esperado:**
```
Detecting runtime... container
Installing Docker... done
Installing dnsmasq... done
Generating cluster configuration... done
Creating syncloud-state.json... done

Cluster configured:
  - Mode: container
  - DNS: syncloud.local
  - Node: syncloud-node

Syncloud installed successfully!
```

### Paso 2: Deploy

```bash
synctl deploy
```

**Qué hace:**
1. Lee los recursos de `deploy/resources.yaml`
2. Valida y parsea al modelo Syncloud
3. Persiste en `syncloud-state.json`
4. Ejecuta el Reconciler para crear contenedores
5. Muestra el DNS configurado

**Output esperado:**
```
Deploying Syncloud platform...
Creating network: syncloud-network... done
Creating container: postgres-db... done
Creating container: backend... done
Creating container: webapp... done

Syncloud platform deployed successfully!
Access your platform at: http://syncloud.local
```

### Paso 3: Gestionar

```bash
# Ver estado del cluster
synctl inspect

# Ver nodos
synctl get nodes

# Ver recursos
synctl get resources

# Ver DNS records
synctl get dns

# Ver detalle de recurso
synctl resource describe postgres-db

# Ver logs
synctl resource logs backend

# Agregar DNS record
synctl dns add api --server 192.168.1.101

# Eliminar recurso
synctl resource delete webapp
```

## syncloud-state.json - Médula del Sistema

El archivo `syncloud-state.json` es **el contrato central** que representa el estado deseado de la plataforma. Toda operación de synctl lee y escribe en este archivo.

### Estructura Completa

```json
{
  "version": "1.0.0",
  "cluster": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "syncloud",
    "mode": "container|vps",
    "dns": "syncloud.local",
    "createdAt": "2024-01-15T10:30:00Z"
  },
  "nodes": [
    {
      "id": "node-1",
      "hostname": "syncloud-node",
      "role": "control-plane",
      "ip": "192.168.1.100"
    }
  ],
  "resources": [
    {
      "id": "res-uuid",
      "name": "postgres-db",
      "kind": "Container",
      "runtime": "docker",
      "nodeId": "node-1",
      "spec": {
        "image": "postgres:15-alpine",
        "ports": [...],
        "env": [...],
        "volumes": [...]
      },
      "status": {
        "state": "running",
        "containerId": "abc123..."
      },
      "ownership": {
        "ownerId": "syncloud",
        "ownerKind": "platform",
        "ownerName": "syncloud"
      },
      "createdAt": "...",
      "updatedAt": "..."
    }
  ],
  "dns": {
    "records": [
      {"name": "app", "server": "192.168.1.100"},
      {"name": "api", "server": "192.168.1.100"}
    ]
  }
}
```

### Objetos Gestionados

| Objeto | Descripción | CRUD |
|--------|-------------|------|
| `cluster` | Configuración global del cluster | R/U |
| `nodes` | Nodos físicos/virtuales | C/R/U/D |
| `resources` | Recursos abstractos (Container, Deployment, etc.) | C/R/U/D |
| `dns.records` | Registros DNS | C/R/U/D |

### Ciclo de Vida del State

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   State     │ ──►  │  Reconciler │ ──►  │   Runtime   │
│  (desired)  │      │   Compare   │      │  (docker/   │
│             │      │   Diff      │      │   k8s)      │
└─────────────┘      └──────┬──────┘      └─────────────┘
       ▲                    │
       │                    ▼
       │             ┌─────────────┐
       └─────────────│   Observe   │
        (actual)     │  (runtime)  │
                     └─────────────┘
```

## Stack Desplegado

```
┌─────────────────────────────────────────────────────┐
│                  syncloud-network                    │
│                                                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │
│  │  postgres   │  │   backend   │  │   webapp    │ │
│  │    (DB)     │◄─┤  (API)      │◄─┤   (nginx)   │ │
│  │   :5432     │  │   :8080     │  │    :80      │ │
│  └─────────────┘  └─────────────┘  └─────────────┘ │
└─────────────────────────────────────────────────────┘
         │               │                │
         └───────────────┴────────────────┘
                            │
                            ▼
                   ┌─────────────────┐
                   │  syncloud.local  │
                   │   (DNS dnsmasq)  │
                   └─────────────────┘
```

## Comandos v1.0 - Gestión Completa

### Instalación y Despliegue

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl install` | Instalar y configurar cluster | ✅ |
| `synctl deploy` | Desplegar plataforma completa | ✅ |
| `synctl inspect` | Inspeccionar estado general | ✅ |

### Gestión de Recursos (CRUD)

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl apply -f <file>` | Crear/actualizar recursos | ✅ |
| `synctl get resources` | Leer recursos | ✅ |
| `synctl resource describe <name>` | Ver detalle de recurso | 🔴 |
| `synctl resource delete <name>` | Eliminar recurso | 🔴 |
| `synctl resource logs <name>` | Ver logs | 🔴 |
| `synctl resource exec <name> -- <cmd>` | Ejecutar comando | 🔴 |
| `synctl resource scale <name> --replicas=N` | Escalar recurso | 🔴 |
| `synctl resource restart <name>` | Reiniciar recurso | 🔴 |

### Gestión de DNS (CRUD)

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl get dns` | Listar registros DNS | 🔴 |
| `synctl dns add <name> --server <ip>` | Agregar registro | 🔴 |
| `synctl dns delete <name>` | Eliminar registro | 🔴 |
| `synctl dns update <name> --server <ip>` | Actualizar registro | 🔴 |

### Gestión de Nodos (CRUD)

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl get nodes` | Listar nodos | ✅ |
| `synctl node describe <name>` | Ver detalle de nodo | 🔴 |
| `synctl node add -f <file>` | Agregar nodo | 🔴 |
| `synctl node remove <name>` | Remover nodo | 🔴 |
| `synctl node cordon <name>` | Deshabilitar scheduling | 🔴 |
| `synctl node uncordon <name>` | Habilitar scheduling | 🔴 |

### Validación y Diagnóstico

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl validate -f <file>` | Validar recursos | 🔴 |
| `synctl diff -f <file>` | Comparar con estado | 🔴 |
| `synctl apply --dry-run -f <file>` | Validar sin aplicar | 🔴 |
| `synctl apply --preview -f <file>` | Preview de cambios | 🔴 |

### Output

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl get resources --output json` | Formato JSON | 🔴 |
| `synctl get resources --output yaml` | Formato YAML | 🔴 |

## Recursos del Deploy

Ubicación: `internal/application/services/apply/deploy/resources.yaml`

| Nombre | Imagen | Tipo |
|--------|--------|------|
| `postgres-db` | postgres:15-alpine | Container |
| `backend` | syncloud/backend:latest | Container |
| `webapp` | syncloud/webapp:latest | Container |
| `syncloud-network` | - | Network |

## Post-v1.0

Una vez liberada la v1.0:

1. **Plugins** - Cargar runtimes adicionales
2. **MCP Protocol** - Integración con Model Context Protocol
3. **API Server** - Modo daemon con REST API
4. **Multi-cluster** - Gestión de múltiples clusters
5. **CI/CD Integration** - Hooks para pipelines

## Notas

- ✅ = Implementado
- 🔴 = Por implementar para v1.0
