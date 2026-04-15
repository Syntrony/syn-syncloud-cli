# Comando: describe

## Descripción

El comando `describe` muestra información detallada de un recurso específico en la plataforma Syncloud. Combina datos del estado guardado con información en tiempo real del runtime (Docker o Kubernetes), proporcionando una vista completa del recurso.

## Uso

```bash
# Describir un recurso por nombre
synctl describe -n <nombre>

# Describir con kind específico
synctl describe -n <nombre> -k <tipo>

# Describir indicando runtime
synctl describe -n <nombre> -r <docker|k8s>

# Describir por ID
synctl describe -I <id>

# Salida en formato JSON
synctl describe -n <nombre> -o json
```

### Parámetros

| Parámetro | Alias | Descripción |
|-----------|-------|-------------|
| `--name` | `-n` | Nombre del recurso a describir |
| `--kind` | `-k` | Tipo de recurso (e.g., docker.container, k8s.deployment) |
| `--runtime` | `-r` | Runtime del recurso (docker, kubernetes) |
| `--id` | `-I` | ID único del recurso |
| `--output` | `-o` | Formato de salida (table, json) |

## Flujo de Ejecución

```
┌─────────────────────────────────────────────────────────────────────┐
│                     synctl describe                                 │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  1. Inicialización de Componentes                                  │
│     - Logger: Inicializar sistema de logs                            │
│     - Repository: Cargar estado actual del sistema                  │
│     - Runner: Inicializar ejecutor de comandos                      │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  2. Construcción de Filtros                                         │
│     - Crear GetResourceFilter                                       │
│     - Aplicar criterios: name, kind, id, runtime                   │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  3. Búsqueda en Estado (cached)                                    │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  GetResourceService.GetResourceBy()                        │   │
│     │  - Consultar syncloud-state.json                         │   │
│     │  - Aplicar filtros                                      │   │
│     └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │ RECURSO ENCONTRADO?        │
                    └─────────────┬─────────────┘
                        Sí       │       No
                        ▼        ▼
┌──────────────────┐    ┌────────────────────────────────────────┐
│ 4. Observar      │    │ 4b. Buscar en Runtime (fallback)        │
│    Runtime       │    │    - docker inspect / kubectl describe │
│    (live data)   │    │    - Retornar datos si existe          │
└──────────────────┘    └────────────────────────────────────────┘
        │                              │
        └──────────────┬──────────────┘
                       ▼
┌─────────────────────────────────────────────────────────────────────┐
│  5. Construcción del DTO                                            │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  ResourceDetailDto                                       │   │
│     │  - Id, Name, Kind, Runtime                               │   │
│     │  - CreatedAt, UpdatedAt                                  │   │
│     │  - Status (cached)                                       │   │
│     │  - Spec (cached)                                         │   │
│     │  - LiveSpec (runtime actual)                            │   │
│     └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  6. Formateo de Salida                                             │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  PrintDescribe()                                        │   │
│     │  - Table: formato legible para humanos                 │   │
│     │  - JSON: estructura completa para scripts              │   │
│     └──────────────────────────────���───────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

## Estrategia de Datos Mixtos

El comando describe utiliza una estrategia de datos mixtos que prioriza la información disponible:

### 1. Estado Guardado (Cached)
- **Fuente**: `syncloud-state.json`
- **Contenido**: Spec, Status, metadatos guardados
- **Ventaja**: Rápido, disponible offline

### 2. Datos del Runtime (Live)
- **Fuente**: `docker inspect` o `kubectl describe`
- **Contenido**: Estado actual del recurso en el runtime
- **Ventaja**: Información actualizada en tiempo real

### 3. Combinación
- Si existe en estado: obtiene datos cached + live del runtime
- Si NO existe en estado: busca directamente en runtime
- Siempre intenta obtener datos vivos para complementar

## Tipos de Recursos Soportados

### Docker
- `docker.container`
- `docker.network`
- `docker.volume`
- `docker.image`

### Kubernetes
- `k8s.namespace`
- `k8s.deployment`
- `k8s.service`
- `k8s.pod`
- `k8s.ingress`
- `k8s.configmap`
- `k8s.secret`

## Formato de Salida

### Table (default)
```
Name:        my-app
Kind:        k8s.deployment
Runtime:     kubernetes
ID:          deploy/my-app
Created:     2024-01-15T10:30:00Z
Updated:    2024-01-15T10:30:00Z

Status:
  replicas: 3
  ready: 3

Spec:
  image: myregistry/my-app:v1.0
  replicas: 3

Live Data:
  availableReplicas: 3
  readyReplicas: 3
  updatedReplicas: 3
```

### JSON
```json
{
  "Id": "deploy/my-app",
  "Name": "my-app",
  "Kind": "k8s.deployment",
  "Runtime": "kubernetes",
  "CreatedAt": "2024-01-15T10:30:00Z",
  "UpdatedAt": "2024-01-15T10:30:00Z",
  "Status": {
    "replicas": 3,
    "ready": 3
  },
  "Spec": {
    "image": "myregistry/my-app:v1.0",
    "replicas": 3
  },
  "LiveSpec": {
    "availableReplicas": 3,
    "readyReplicas": 3,
    "updatedReplicas": 3
  },
  "Events": []
}
```

## Casos de Uso

### Ver detalles de un deployment
```bash
synctl describe -n my-app -k k8s.deployment
```

### Ver estado actual de un container
```bash
synctl describe -n web-container -k docker.container -r docker
```

### Obtener salida JSON para scripts
```bash
synctl describe -n my-app -o json > my-app-status.json
```

### Ver información viva del runtime
```bash
synctl describe -n my-app -k k8s.deployment -o json
```

## Consideraciones

- **Live Data**: Por defecto intenta obtener datos del runtime para complementar la información cached.
- **Fallback**: Si el recurso no está en el estado guardado, busca directamente en el runtime.
- **FormatoJSON**: Útil para integración con otros sistemas o automatización.
- **Errores**: Si el recurso no existe en ningún lado, retorna error informativo.