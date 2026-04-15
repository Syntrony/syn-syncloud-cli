# Comando: status

## Descripción

El comando `status` muestra el estado de salud general de la plataforma Syncloud. Verifica la disponibilidad de los runtimes (Docker y Kubernetes) y muestra un resumen del estado instalado: cluster, nodos y recursos gestionados.

## Uso

```bash
synctl status
```

> No requiere parámetros.

## Flujo de Ejecución

```
┌─────────────────────────────────────────────────────────────────────┐
│                     synctl status                                   │
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
│  2. Verificación de Runtimes                                        │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  Docker                                                   │   │
│     │  - docker info --format {{.ServerVersion}}               │   │
│     │  - Estado: running (vX.X.X) | unavailable               │   │
│     └──────────────────────────────────────────────────────────┘   │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  Kubernetes                                               │   │
│     │  - kubectl version --client                              │   │
│     │  - Estado: available | unavailable                       │   │
│     └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  3. Lectura del Estado de Plataforma                                │
│     - Verificar existencia de syncloud-state.json                   │
│     - Si existe: leer versión, cluster, nodos, recursos             │
│     - Si no existe: indicar plataforma no instalada                 │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  4. Presentación                                                    │
│     - Tabla de salud con runtimes y estado de plataforma            │
└─────────────────────────────────────────────────────────────────────┘
```

## Formato de Salida

```
╔══════════════════════════════════════╗
║      Syncloud Platform - Status      ║
╚══════════════════════════════════════╝

  Docker:      running (v27.3.1)
  Kubernetes:  available

  Platform:    installed (v0.1)
  Cluster:     orion [single]
  Nodes:       1
  Resources:   8
```

### Cuando la plataforma no está instalada

```
╔══════════════════════════════════════╗
║      Syncloud Platform - Status      ║
╚══════════════════════════════════════╝

  Docker:      running (v27.3.1)
  Kubernetes:  available

  Platform:    not installed (run: synctl install)
```

## Campos de Salida

| Campo | Descripción |
|-------|-------------|
| `Docker` | Versión del servidor Docker o `unavailable` si no responde |
| `Kubernetes` | Disponibilidad del cliente `kubectl` |
| `Platform` | Versión del estado instalado o indicación de no instalado |
| `Cluster` | Nombre e identificador del cluster (`[modo]`) |
| `Nodes` | Cantidad de nodos registrados en el estado |
| `Resources` | Total de recursos gestionados por la plataforma |

## Casos de Uso

### Verificar el entorno antes de un deploy
```bash
synctl status
synctl apply -f resources.yaml
```

### Health check en scripts de CI/CD
```bash
synctl status
if [ $? -ne 0 ]; then
  echo "Plataforma no disponible"
  exit 1
fi
```

### Diagnóstico rápido post-install
```bash
synctl install
synctl status
```

## Consideraciones

- **No modifica estado**: comando de solo lectura, sin efectos secundarios.
- **Sin privilegios**: no requiere sudo.
- **Runtimes**: la verificación de Docker y Kubernetes es independiente del estado almacenado; refleja disponibilidad real en el momento de la ejecución.
- **Complemento a `inspect`**: `status` ofrece una vista de salud rápida, mientras que `inspect` muestra el detalle del estado almacenado.
