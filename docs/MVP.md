# Flujo de MVP - Syncloud Platform

## Overview

El flujo MVP de Syncloud permite a los usuarios desplegar una plataforma completa en 3 pasos:

```
1. synctl install   → Instalar runtimes (k8s, docker, dnsmasq)
2. synctl deploy    → Desplegar la plataforma
3. Acceder via DNS → Abrir en el navegador
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
- Detecta el entorno (Container/VPS)
- Instala Docker o K3s/K3d según el entorno
- Instala y configura dnsmasq
- Crea el archivo `syncloud-state.json`
- Configura la red interna

**Estado generado (`syncloud-state.json`):**
```json
{
  "version": "1.0.0",
  "cluster": {
    "id": "uuid",
    "name": "syncloud",
    "mode": "container|vps",
    "dns": "syncloud.local"
  },
  "nodes": [...],
  "resources": [...]
}
```

### Paso 3: Deploy - Desplegar la Plataforma

```bash
synctl deploy
```

**Qué hace:**
- Lee los recursos de `deploy/resources.yaml`
- Valida y parsea los recursos
- Persiste en `syncloud-state.json`
- Reconcilia contra Docker/Kubernetes
- Muestra el DNS configurado

**Stack desplegado:**
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

### Paso 4: Acceder

```
Syncloud platform deployed successfully!
Access your platform at: http://syncloud.local
```

El usuario copia la URL en su navegador.

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
                    │    │ Docker/K3s      │   │
                    │    │ dnsmasq         │   │
                    │    │ syncloud-state  │   │
                    │    └─────────────────┘   │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │   synctl deploy          │
                    │   ┌─────────────────┐   │
                    │   │ deploy/resources │   │
                    │   │  ↓               │   │
                    │   │ Parse & Validate │   │
                    │   │  ↓               │   │
                    │   │ Update State     │   │
                    │   │  ↓               │   │
                    │   │ Reconcile        │   │
                    │   │  ↓               │   │
                    │   │ docker run       │   │
                    │   └─────────────────┘   │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │  http://syncloud.local   │
                    │    (Mostrar en terminal) │
                    └──────────────────────────┘
```

## syncloud-state.json - Centro del Sistema

El archivo `syncloud-state.json` es **la médula espinal** de Syncloud:

```json
{
  "version": "1.0.0",
  "cluster": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "syncloud",
    "mode": "container",
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
      "spec": {...},
      "status": {...},
      "ownership": {...},
      "createdAt": "...",
      "updatedAt": "..."
    }
  ]
}
```

**Funciones del state:**
1. **Fuente de verdad** - Estado deseado de todos los recursos
2. **Persistencia** - Sobrevive reinicios
3. **Sincronización** - Reconciler lo usa para comparar desired vs actual
4. **Auditoría** - Historial de cambios

## Comandos del Flujo MVP

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl install` | Instala runtimes | ✅ Implementado |
| `synctl deploy` | Despliega plataforma | ✅ Implementado |
| `synctl apply -f <file>` | Aplica recursos custom | ✅ Implementado |
| `synctl inspect` | Muestra estado | ✅ Implementado |
| `synctl get resources` | Lista recursos | ✅ Implementado |
| `synctl get nodes` | Lista nodos | ✅ Implementado |

## Recursos del Deploy

Los recursos de deploy están en: `internal/application/services/apply/deploy/resources.yaml`

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

## Siguientes Pasos Post-MVP

Una vez completado el flujo MVP:

1. **Persistencia de estado robusta**
2. **Rollback y historial**
3. **Multi-tenant**
4. **Escalado automático**
5. **Monitoreo y métricas**
