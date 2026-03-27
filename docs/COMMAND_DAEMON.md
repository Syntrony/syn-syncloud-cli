# Comando: daemon

## Descripción

El comando `daemon` inicia el control plane de Syncloud en modo segundo plano. Este proceso es responsable de mantener el estado deseado de la plataforma reconciliándose continuamente con el estado actual de los recursos.

## Uso

```bash
synctl daemon
```

## Propósito

El daemon es el componente de reconciliación continua que:
- Monitorea el estado de los recursos
- Detecta drifts entre estado deseado y actual
- Aplica correcciones automáticas
- Mantiene la plataforma en el estado declarativo definido

## Flujo de Ejecución

```
┌─────────────────────────────────────────────────────────────────────┐
│                         synctl daemon                               │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  1. Inicialización de Contexto                                     │
│     - Crear contexto de Go para cancelación                         │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  2. Inicialización de Componentes                                  │
│     - Logger: Sistema de logs                                       │
│     - Repository: Acceso al estado                                  │
│     - Runner: Ejecutor de comandos                                 │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  3. Configuración de Reconcilers                                  │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  Runtime Reconcilers                                     │   │
│     │  ┌────────────────────┐    ┌────────────────────┐        │   │
│     │  │ DockerReconciler │    │  K8sReconciler     │        │   │
│     │  └────────────────────┘    └────────────────────┘        │   │
│     │         │                          │                     │   │
│     │         └──────────┬───────────────┘                     │   │
│     │                    ▼                                     │   │
│     │          ┌────────────────────┐                           │   │
│     │          │ ReconcileService  │                           │   │
│     │          └────────────────────┘                           │   │
│     └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  4. Inicialización del Controller                                  │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  Controller                                               │   │
│     │  - Frecuencia de reconcile: 5 segundos                   │   │
│     │  - Servicios: Repository + ReconcileService + Logger    │   │
│     └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  5. Inicio del Loop de Reconciliación                             │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │                    BUCLE INFINITO                         │   │
│     │  ┌─────────────────────────────────────────────────┐     │   │
│     │  │  1. Obtener estado deseado (Repository)        │     │   │
│     │  └─────────────────────────────────────────────────┘     │   │
│     │                         │                                 │   │
│     │                         ▼                                 │   │
│     │  ┌─────────────────────────────────────────────────┐     │   │
│     │  │  2. Obtener estado actual (Docker/K8s APIs)     │     │   │
│     │  └─────────────────────────────────────────────────┘     │   │
│     │                         │                                 │   │
│     │                         ▼                                 │   │
│     │  ┌─────────────────────────────────────────────────┐     │   │
│     │  │  3. Calcular diferencias (Diff)                 │     │   │
│     │  └─────────────────────────────────────────────────┘     │   │
│     │                         │                                 │   │
│     │                         ▼                                 │   │
│     │  ┌─────────────────────────────────────────────────┐     │   │
│     │  │  4. Generar plan de acciones                   │     │   │
│     │  │     - Create/Update/Delete recursos            │     │   │
│     │  └─────────────────────────────────────────────────┘     │   │
│     │                         │                                 │   │
│     │                         ▼                                 │   │
│     │  ┌─────────────────────────────────────────────────┐     │   │
│     │  │  5. Ejecutar acciones (Runtime Reconcilers)    │     │   │
│     │  └─────────────────────────────────────────────────┘     │   │
│     │                         │                                 │   │
│     │                         ▼                                 │   │
│     │  ┌─────────────────────────────────────────────────┐     │   │
│     │  │  6. Actualizar estado (Repository)              │     │   │
│     │  └─────────────────────────────────────────────────┘     │   │
│     │                         │                                 │   │
│     │                         ▼                                 │   │
│     │              ┌─────────────────────┐                     │   │
│     │              │ Esperar 5 segundos   │                     │   │
│     │              └─────────────────────┘                     │   │
│     │                         │                                 │   │
│     └─────────────────────────┼─────────────────────────────────┘   │
│                               │                                    │
│                   (Hasta cancelación)                              │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Final:Daemon detenido                                             │
│  - El proceso termina cuando se cancela el contexto               │
└─────────────────────────────────────────────────────────────────────┘
```

## Componentes del Reconciler

### Docker Reconciler
- Detecta contenedores, redes e imágenes
- Compara estado deseado vs actual
- Ejecuta acciones de creación, actualización y eliminación

### Kubernetes Reconciler
- Detecta recursos K8s (Deployments, Services, etc.)
- Compara estado deseado vs actual
- Ejecuta acciones mediante kubectl o API nativa

## Características

### Intervalo de Reconciliación
- **Frecuencia**: 5 segundos
- Configurable en el código

### Reconciliación Continua
El daemon nunca se detiene activamente; espera hasta que:
- Se reciba señal de cancelación
- Ocurra un error crítico

### Idempotencia
Cada ciclo de reconciliación es idempotente:
- Múltiples ejecuciones producen el mismo resultado
- No hay efectos secundarios no deseados

## Casos de Uso

### Iniciar el control plane
```bash
synctl daemon
```

### Detener el daemon
Presionar `Ctrl+C` para enviar señal de cancelación.

## Relación con apply

| Comando | Modo | Descripción |
|---------|------|-------------|
| `apply` | Único | Aplica recursos una vez |
| `daemon` | Continuo | Mantiene recursos en estado deseado |

El daemon complementa el comando `apply`:
1. `apply` define el estado deseado
2. `daemon` mantiene ese estado continuamente

## Consideraciones

- **Proceso de larga duración**: El daemon corre indefinidamente
- **Consumo de recursos**: Usa CPU y memoria continuamente
- **Logs extensos**: Genera logs de cada ciclo de reconciliación
- **Requiere Syncloud inicializado**: Necesita el archivo de estado
