# Flujo de MVP - Syncloud Platform

## Overview

El flujo MVP permite a los usuarios desplegar una plataforma completa en 3 pasos mínimos:

```
1. synctl install   → Configurar runtimes (k8s, docker, dnsmasq)
2. synctl deploy    → Desplegar la plataforma (DB + Backend + WebApp)
3. Acceder via DNS → Abrir http://syncloud.local
```

## Flujo Detallado

### Paso 1: Descargar e Instalar

```bash
# Descargar synctl (futuro: curl script o package manager)
curl -fsSL https://get.syncloud.io | sh

# O compilar desde código
go build -o synctl .
```

### Paso 2: Install - Configurar Runtimes

```bash
synctl install
```

**Qué hace:**
1. Detecta el entorno de ejecución (Container/VPS)
2. Instala Docker o K3s/K3d según el entorno detectado
3. Instala y configura dnsmasq para resolución DNS
4. Genera el archivo `syncloud-state.json` con la configuración inicial del cluster
5. Configura la red interna y registra el DNS provisional

**Output esperado:**
```
Installing Docker... done
Installing dnsmasq... done
Configuring DNS: syncloud.local
Syncloud installed successfully!
```

### Paso 3: Deploy - Desplegar la Plataforma

```bash
synctl deploy
```

**Qué hace:**
1. Lee los recursos de `deploy/resources.yaml` (postgres-db, backend, webapp, network)
2. Valida y parsea los recursos al modelo Syncloud
3. Persiste el estado deseado en `syncloud-state.json`
4. Ejecuta el Reconciler para crear los contenedores en Docker
5. Recarga el state y muestra el DNS configurado

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

### Paso 4: Acceder

El usuario copia la URL `http://syncloud.local` en su navegador.

## Diagrama de Flujo

```
┌─────────────────────────────────────────────────────────────────────┐
│                          USUARIO                                     │
└─────────────────────────────────────────────────────────────────────┘
                                   │
                                   ▼
                    ┌──────────────────────────┐
                    │    synctl install        │
                    │    ┌─────────────────┐   │
                    │    │ Detect Runtime   │   │
                    │    │ Install Docker   │   │
                    │    │ Install K3s     │   │
                    │    │ Setup dnsmasq   │   │
                    │    │ Generate State  │   │
                    │    └─────────────────┘   │
                    │           │               │
                    │           ▼               │
                    │    syncloud-state.json    │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │   synctl deploy          │
                    │   ┌─────────────────┐   │
                    │   │ Parse Resources │   │
                    │   │ Validate       │   │
                    │   │ Update State   │◄──┼── syncloud-state.json
                    │   │ Reconcile      │   │     (desired state)
                    │   └───────┬───────┘   │
                    │           │           │
                    │           ▼           │
                    │   ┌─────────────────┐ │
                    │   │ Docker Runtime  │ │
                    │   │ docker run ...  │ │
                    │   └─────────────────┘ │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │  http://syncloud.local   │
                    │    (Mostrar DNS terminal)│
                    └──────────────────────────┘
```

## syncloud-state.json - Médula del Sistema

El archivo `syncloud-state.json` es **el contrato central** que representa el estado deseado de la plataforma. Toda operación de synctl lee y escribe en este archivo.

### Estructura Técnica

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
      "kind": "docker.container",
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
        "managed": true,
        "owner": "syncloud"
      },
      "createdAt": "...",
      "updatedAt": "..."
    }
  ]
}
```

### Ciclo de Vida del State

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   State     │ ──►  │  Reconciler │ ──►  │   Runtime   │
│  (desired)  │      │   Compare   │      │  (docker/  │
│             │      │   Diff     │      │   k8s)     │
└─────────────┘      └──────┬──────┘      └─────────────┘
       ▲                   │
       │                   ▼
       │            ┌─────────────┐
       └────────────│   Observe   │
        (actual)    │  (runtime)  │
                    └─────────────┘
```

### Funciones del State

| Función | Descripción |
|---------|-------------|
| **Fuente de verdad** | Define qué recursos deben existir |
| **Persistencia** | Sobrevive reinicios del sistema |
| **Sincronización** | Reconciler lo usa para comparar desired vs actual |
| **Auditoría** | Historial de cambios y timestamps |
| **Contrato DNS** | El campo `cluster.dns` define la URL de acceso |

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
```

### Contenedores

| Nombre | Imagen | Puertos | Descripción |
|--------|--------|---------|-------------|
| `postgres-db` | postgres:15-alpine | 5432 | Base de datos PostgreSQL |
| `backend` | syncloud/backend:latest | 8080 | API backend |
| `webapp` | syncloud/webapp:latest | 80 | Frontend Nginx |

### Redes

| Nombre | Driver | Descripción |
|--------|--------|-------------|
| `syncloud-network` | bridge | Red interna para la plataforma |

## Comandos del Flujo MVP

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl install` | Instala runtimes y genera state inicial | ✅ Implementado |
| `synctl deploy` | Despliega plataforma y muestra DNS | ✅ Implementado |
| `synctl apply -f <file>` | Aplica recursos custom | ✅ Implementado |
| `synctl inspect` | Muestra estado actual | ✅ Implementado |
| `synctl get resources` | Lista recursos | ✅ Implementado |
| `synctl get nodes` | Lista nodos | ✅ Implementado |

## Recursos del Deploy

Los recursos de deploy están en: `internal/application/services/apply/deploy/resources.yaml`

## Siguientes Pasos Post-MVP

Una vez completado el flujo MVP:

1. **Validación avanzada** - dry-run, preview de cambios
2. **Operaciones CRUD** - delete, update de recursos
3. **Observabilidad** - logs, describe, exec
4. **Orquestación** - dependencias, rollback
5. **Multi-tenant** - múltiples plataformas
