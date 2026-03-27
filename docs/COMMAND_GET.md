# Comando: get

## Descripción

El comando `get` permite listar y consultar recursos y nodos del sistema Syncloud. Es la herramienta principal para obtener visibilidad sobre el estado actual de la plataforma.

## Subcomandos

### get resources
Lista los recursos desplegados en la plataforma, con soporte para filtros.

### get nodes
Lista los nodos que conforman el cluster de Syncloud.

## Uso

```bash
# Listar todos los recursos
synctl get resources

# Listar recursos con filtros
synctl get resources --name <nombre>
synctl get resources --kind <tipo>
synctl get resources --runtime <docker|k8s>
synctl get resources --id <id>

# Listar nodos
synctl get nodes
```

### Parámetros para get resources

| Parámetro | Alias | Descripción |
|-----------|-------|-------------|
| `--name` | `-n` | Filtrar por nombre del recurso |
| `--kind` | `-k` | Filtrar por tipo de recurso (Deployment, Service, etc.) |
| `--runtime` | `-r` | Filtrar por runtime (docker, k8s) |
| `--id` | `-i` | Filtrar por ID específico del recurso |
| `--node-id` | - | Filtrar por ID del nodo |

## Flujo de Ejecución

```
┌─────────────────────────────────────────────────────────────────────┐
│                     synctl get <subcomando>                         │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  1. Inicialización de Componentes                                  │
│     - Logger: Inicializar sistema de logs                          │
│     - Repository: Cargar estado actual del sistema                 │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
              ┌─────────────────────────────────┐
              │         Selector de             │
              │         Subcomando              │
              └─────────────────────────────────┘
                          │
          ┌───────────────┴───────────────┐
          ▼                               ▼
┌─────────────────────┐       ┌─────────────────────┐
│   get resources    │       │    get nodes        │
└─────────────────────┘       └─────────────────────┘
          │                               │
          ▼                               ▼
```

### Flujo: get resources

```
┌─────────────────────────────────────────────────────────────────────┐
│  2. Construcción del Filtro                                       │
│     - Crear objeto GetResourceFilter                               │
│     - Aplicar filtros especificados por usuario                    │
│     - Valores por defecto si no se especifican                      │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  3. Ejecución del Servicio                                         │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  GetResourceService                                      │   │
│     │  - Consultar repositorio de estado                       │   │
│     │  - Aplicar filtros                                      │   │
│     │  - Retornar lista de recursos                           │   │
│     └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  4. Formateo de Salida                                             │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  Output Formatter (outputs.PrintResources)              │   │
│     │  - Tabla con columnas: Name, Kind, Runtime, Status      │   │
│     └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  5. Presentación                                                   │
│     - Si hay resultados: Mostrar tabla de recursos                │
│     - Si no hay resultados: Mostrar mensaje informativo           │
└─────────────────────────────────────────────────────────────────────┘
```

### Flujo: get nodes

```
┌─────────────────────────────────────────────────────────────────────┐
│  2. Ejecución del Servicio                                         │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  GetNodeService                                         │   │
│     │  - Consultar repositorio de estado                       │   │
│     │  - Retornar lista de nodos                              │   │
│     └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  3. Formateo de Salida                                             │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  Output Formatter (outputs.PrintNodes)                   │   │
│     │  - Tabla con columnas: Name, Status, Roles, Version      │   │
│     └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  4. Presentación                                                   │
│     - Si hay resultados: Mostrar tabla de nodos                    │
│     - Si no hay resultados: Mostrar mensaje informativo           │
└─────────────────────────────────────────────────────────────────────┘
```

## Estructura de Datos

### Recurso
- **Name**: Nombre identificador del recurso
- **Kind**: Tipo de recurso (Deployment, Service, Container, etc.)
- **Runtime**: Origen del recurso (docker, k8s)
- **Status**: Estado actual del recurso
- **NodeId**: Nodo donde está desplegado
- **Id**: Identificador único

### Nodo
- **Name**: Nombre del nodo
- **Status**: Estado del nodo (Ready, NotReady)
- **Roles**: Rol del nodo (master, worker)
- **Version**: Versión de Kubernetes/Docker

## Casos de Uso

### Listar todos los recursos
```bash
synctl get resources
```

### Buscar recursos por nombre
```bash
synctl get resources --name my-app
```

### Listar solo recursos de Kubernetes
```bash
synctl get resources --runtime k8s
```

### Listar solo deployments
```bash
synctl get resources --kind Deployment
```

### Ver nodos del cluster
```bash
synctl get nodes
```

## Consideraciones

- **Lectura del Estado**: El comando lee del archivo de estado local, no consulta directamente los runtimes.
- **Filtros**: Los filtros son opcionales y pueden combinarse.
- **Salida Formateada**: La salida siempre está formateada en tabla para facilitar la lectura.
