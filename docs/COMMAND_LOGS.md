# Comando: logs

## Descripción

El comando `logs` muestra los logs de un recurso específico de la plataforma Syncloud. Soporta recursos Docker (containers) y Kubernetes (pods), con opciones para seguir el output en tiempo real y limitar las líneas mostradas.

Si no se especifica el runtime con `-r`, el comando lo resuelve automáticamente consultando el estado almacenado (`syncloud-state.json`).

## Uso

```bash
# Logs de un recurso (runtime auto-detectado desde estado)
synctl logs -n <nombre>

# Logs especificando runtime explícitamente
synctl logs -n <nombre> -r docker
synctl logs -n <nombre> -r kubernetes

# Seguir logs en tiempo real
synctl logs -n <nombre> -f

# Mostrar solo las últimas N líneas
synctl logs -n <nombre> --tail 100

# Kubernetes con namespace específico
synctl logs -n <nombre> -r kubernetes --namespace production

# Combinado
synctl logs -n my-app -r kubernetes --namespace staging -f --tail 50
```

### Parámetros

| Parámetro | Alias | Descripción |
|-----------|-------|-------------|
| `--name` | `-n` | Nombre del recurso (requerido) |
| `--runtime` | `-r` | Runtime del recurso (`docker` \| `kubernetes`). Auto-detectado si se omite |
| `--namespace` | - | Namespace de Kubernetes (default: del estado, o `default`) |
| `--follow` | `-f` | Seguir el output del log en tiempo real |
| `--tail` | - | Número de líneas finales a mostrar (0 = todas) |

## Flujo de Ejecución

```
┌─────────────────────────────────────────────────────────────────────┐
│                     synctl logs -n <nombre>                         │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  1. Inicialización de Componentes                                   │
│     - Repository: Cargar estado actual del sistema                  │
│     - Runner: Inicializar ejecutor de comandos                      │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  2. Resolución de Runtime                                           │
│     ┌───────────────────────────────────────────────────────────┐   │
│     │  ¿Se proporcionó --runtime?                               │   │
│     │  - Sí: usar directamente                                  │   │
│     │  - No: buscar recurso en syncloud-state.json por nombre   │   │
│     │        y extraer runtime + namespace                       │   │
│     └───────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
              ┌─────────────────────────────────┐
              │      Selector de Runtime        │
              └─────────────────────────────────┘
                          │
          ┌───────────────┴───────────────┐
          ▼                               ▼
┌─────────────────────┐       ┌──────────────────────────┐
│     Docker          │       │      Kubernetes           │
│  docker logs        │       │  kubectl logs             │
│  [-f] [--tail N]    │       │  [-f] [--tail N] [-n ns]  │
│  <nombre>           │       │  <nombre>                 │
└─────────────────────┘       └──────────────────────────┘
          │                               │
          └───────────────┬───────────────┘
                          ▼
┌─────────────────────────────────────────────────────────────────────┐
│  3. Presentación                                                    │
│     - Imprimir output del comando al stdout                         │
└─────────────────────────────────────────────────────────────────────┘
```

## Comandos Subyacentes

### Docker
```bash
docker logs [-f] [--tail <n>] <nombre>
```

### Kubernetes
```bash
kubectl logs [-f] [--tail <n>] [-n <namespace>] <nombre>
```

## Casos de Uso

### Ver logs de un container Docker
```bash
synctl logs -n nginx
```

### Ver últimas 200 líneas de un pod
```bash
synctl logs -n api-server --tail 200
```

### Seguir logs de un deployment en staging
```bash
synctl logs -n my-app -r kubernetes --namespace staging -f
```

### Depurar un container sin conocer su runtime
```bash
# El runtime se resuelve desde el estado guardado
synctl logs -n postgres -f
```

## Consideraciones

- **Auto-detección de runtime**: si el recurso existe en el estado, no es necesario pasar `-r`.
- **Namespace**: para recursos Kubernetes, el namespace se extrae del estado si no se especifica explícitamente.
- **Follow**: con `-f` el proceso queda en foreground hasta que el usuario interrumpa con `Ctrl+C`.
- **Recursos no encontrados**: si el recurso no existe en el estado ni se especifica runtime, el comando retorna error.
