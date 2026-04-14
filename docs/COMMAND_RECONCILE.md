# Comando: reconcile

## Descripción

El comando `reconcile` ejecuta una reconciliación puntual (one-shot) del estado deseado contra el estado actual de los runtimes (Docker y Kubernetes). A diferencia del `daemon`, que corre en loop continuo, `reconcile` aplica las diferencias detectadas una sola vez y termina.

## Uso

```bash
synctl reconcile
```

> No requiere parámetros. Utiliza el estado almacenado en `syncloud-state.json` como fuente de verdad.

## Flujo de Ejecución

```
┌─────────────────────────────────────────────────────────────────────┐
│                     synctl reconcile                                │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  1. Inicialización de Componentes                                   │
│     - Logger: Inicializar sistema de logs                           │
│     - Repository: Cargar estado actual del sistema                  │
│     - Runner / Sudo: Inicializar ejecutores de comandos             │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  2. Carga del Estado                                                │
│     - Leer syncloud-state.json                                      │
│     - Obtener lista de recursos deseados                            │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  3. Observación por Runtime                                         │
│     ┌──────────────────┐   ┌──────────────────┐                    │
│     │  DockerReconciler│   │  K8sReconciler   │                    │
│     │  - docker inspect│   │  - kubectl get   │                    │
│     └──────────────────┘   └──────────────────┘                    │
│     - Por cada recurso: comparar estado deseado vs actual           │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  4. Planificación de Acciones (diff)                                │
│     - ActionCreate  → recurso no existe en runtime                  │
│     - ActionUpdate  → recurso existe pero difiere                   │
│     - ActionDelete  → recurso fue removido del estado               │
│     - ActionNoop    → recurso en sync, sin cambios                  │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  5. Ejecución de Acciones                                           │
│     - Aplicar creates / updates / deletes en cada runtime           │
│     - Reconciliar registros DNS (DnsReconciler)                     │
└─────────────────────────────────────────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│  6. Resultado                                                       │
│     - Log de confirmación o error por recurso                       │
│     - Finaliza el proceso (exit 0 / exit 1)                         │
└─────────────────────────────────────────────────────────────────────┘
```

## Diferencia con `daemon`

| Aspecto           | `reconcile`              | `daemon`                        |
|-------------------|--------------------------|---------------------------------|
| Ejecución         | Una sola vez             | Loop continuo (cada 5s)         |
| Uso típico        | CI/CD, cron, manual      | Proceso de control plane        |
| Bloqueo           | No (termina solo)        | Sí (proceso en foreground)      |

## Runtimes Soportados

- **Docker**: containers, networks, images
- **Kubernetes**: deployments, services, namespaces, ingress, configmaps, secrets, etc.
- **DNS**: registros `dnsmasq` para acceso interno

## Casos de Uso

### Reconciliación manual tras un cambio
```bash
synctl apply -f resources.yaml
synctl reconcile
```

### Ejecutar desde un cron job
```bash
# /etc/cron.d/synctl-reconcile
*/5 * * * * root /usr/local/bin/synctl reconcile >> /var/log/synctl/reconcile.log 2>&1
```

### Verificar convergencia del cluster
```bash
synctl reconcile && echo "Cluster en sync"
```

## Consideraciones

- **Idempotente**: seguro de ejecutar múltiples veces; operaciones con `ActionNoop` no generan cambios.
- **Sudo**: requiere privilegios si el daemon de Docker exige permisos de root.
- **Estado previo**: si `syncloud-state.json` no existe, el comando termina con error indicando que la plataforma no está instalada.
