# Pendientes de Refinación

## Resumen

Este documento lista las tareas pendientes para alinear la implementación actual con la arquitectura de Syncloud Resources.

---

## Prioridad ALTA

### Comandos

| # | Item | Estado | Descripción |
|---|------|--------|-------------|
| 1 | `apply --dry-run` | ❌ Pendiente | Ejecutar parseo y validación sin aplicar |
| 2 | `apply --preview` | ❌ Pendiente | Mostrar diff antes de aplicar |
| 3 | `delete <resource>` | ❌ Pendiente | Eliminar recursos por nombre/kind |
| 4 | `get resources --output json/yaml` | ❌ Pendiente | Formato de salida configurable |
| 5 | `apply -f -` (stdin) | ❌ Pendiente | Leer desde stdin |

### Arquitectura

| # | Item | Estado | Descripción |
|---|------|--------|-------------|
| 1 | RuntimeReconciler observado | ⚠️ Parcial | El reconciler existe pero no hay ciclo de observación implementado |
| 2 | Plan de reconciliación separado | ❌ Pendiente | Diff → Plan → Execute (no ejecutar directo) |
| 3 | Estado deseado vs actual | ❌ Pendiente | Persistir y comparar estado antes de aplicar |

---

## Prioridad MEDIA

### Comandos

| # | Item | Estado | Descripción |
|---|------|--------|-------------|
| 1 | `logs <resource>` | ❌ Pendiente | Ver logs de contenedor/pod |
| 2 | `exec <resource> -- <cmd>` | ❌ Pendiente | Ejecutar comando en recurso |
| 3 | `port-forward <resource>` | ❌ Pendiente | Forward de puertos |
| 4 | `describe <resource>` | ❌ Pendiente | Descripción detallada |
| 5 | `validate -f <file>` | ❌ Pendiente | Solo validación sin aplicar |

### Arquitectura

| # | Item | Estado | Descripción |
|---|------|--------|-------------|
| 1 | Contenedor de políticas | ❌ Pendiente | Docker y K8s reconcilers en contenedores separados |
| 2 | Métricas/Observabilidad | ❌ Pendiente | Instrumentación para monitoreo |
| 3 | Eventos de recurso | ❌ Pendiente | Audit trail de cambios |
| 4 | Hooks pre/post apply | ❌ Pendiente | Extensibilidad para scripts |

---

## Prioridad BAJA

### Comandos

| # | Item | Estado | Descripción |
|---|------|--------|-------------|
| 1 | `scale <resource> --replicas=N` | ❌ Pendiente | Escalar deployments |
| 2 | `rollout status <resource>` | ❌ Pendiente | Ver estado de rollout |
| 3 | `rollout undo <resource>` | ❌ Pendiente | Revertir cambios |
| 4 | `top resources` | ❌ Pendiente | Uso de recursos (CPU/mem) |
| 5 | `diff -f <file>` | ❌ Pendiente | Diff entre YAML y estado actual |

### Arquitectura

| # | Item | Estado | Descripción |
|---|------|--------|-------------|
| 1 | Plugin system | ❌ Pendiente | Cargar runtimes adicionales |
| 2 | Resource dependencies | ❌ Pendiente | Orden de aplicación basado en dependencias |
| 3 | Parallel apply | ❌ Pendiente | Aplicar recursos independientes en paralelo |
| 4 | MCP Protocol | ❌ Pendiente | Integración con Model Context Protocol |
| 5 | API Server mode | ❌ Pendiente | Modo daemon con API REST/GRPC |

---

## Próximos Pasos Recomendados

### Fase 1: Completar ciclo de reconciliación

```bash
# Implementar
1. apply --dry-run    # Validación sin ejecución
2. apply --preview    # Mostrar diff
3. Separar Diff → Plan → Execute
```

### Fase 2: Completar operaciones CRUD

```bash
# Implementar
1. get resources --output json/yaml
2. delete <resource>
3. validate -f <file>
```

### Fase 3: Debug y diagnóstico

```bash
# Implementar
1. logs <resource>
2. describe <resource>
3. exec <resource> -- <cmd>
```

### Fase 4: Orquestación

```bash
# Implementar
1. Resource dependencies (orden)
2. Hooks pre/post apply
3. Parallel apply
```

---

## Notas

- ❌ = No implementado
- ⚠️ = Parcialmente implementado
- ✅ = Implementado
- 🔄 = En progreso
