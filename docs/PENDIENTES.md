# Roadmap de Desarrollo

## Contexto: El centro es syncloud-state.json

Todo gira alrededor de `syncloud-state.json`:
- **Estado deseado**: Recursos declarados por el usuario
- **Estado actual**: Recursos observados en los runtimes
- **Diff**: Diferencia entre desired y actual
- **Reconciler**: Aplica cambios para cerrar el gap

```
syncloud-state.json ──► Reconciler ──► Docker/K8s
         ▲                                    │
         │                                    ▼
         └──────────── Observe ──────────── Runtime
```

---

## MVP (0.1.0) - Flujo Completo

### ✅ Implementado

| Item | Comando/Feature | Descripción |
|------|-----------------|-------------|
| ✅ | `synctl install` | Instala Docker/K3s + dnsmasq |
| ✅ | `synctl apply -f` | Aplica recursos desde YAML |
| ✅ | `synctl deploy` | Despliega stack completo |
| ✅ | `synctl inspect` | Muestra estado actual |
| ✅ | `synctl get resources` | Lista recursos |
| ✅ | `synctl get nodes` | Lista nodos |
| ✅ | Runtime Reconciler | Adaptador Docker/K8s |

### 🔄 En Progreso

| # | Item | Descripción |
|---|------|-------------|
| 1 | Deploy resources | Recursos hardcodeados necesitan externalizarse |

### ❌ Pendientes para MVP Completo

#### Comandos

| # | Item | Prioridad | Descripción |
|---|------|-----------|-------------|
| 1 | `synctl delete <resource>` | 🔴 Alta | Eliminar recursos del state y runtime |
| 2 | `synctl deploy --dry-run` | 🔴 Alta | Validar sin aplicar |
| 3 | Output formatters | 🔴 Alta | `get resources --output json/yaml` |

#### Arquitectura

| # | Item | Prioridad | Descripción |
|---|------|-----------|-------------|
| 1 | Diff separated | 🔴 Alta | Diff → Plan → Execute (no directo) |
| 2 | Observe cycle | 🔴 Alta | Leer estado actual de runtime |
| 3 | State comparison | 🔴 Alta | desired vs actual en reconciler |

---

## Post-MVP 1 (0.2.0) - Developer Experience

### Comandos

| # | Item | Prioridad | Descripción |
|---|------|-----------|-------------|
| 1 | `synctl logs <resource>` | 🟡 Media | Ver logs de contenedor/pod |
| 2 | `synctl exec <resource> -- <cmd>` | 🟡 Media | Ejecutar en recurso |
| 3 | `synctl describe <resource>` | 🟡 Media | Detalle completo |
| 4 | `synctl port-forward <resource>` | 🟡 Media | Forward de puertos |

### Arquitectura

| # | Item | Prioridad | Descripción |
|---|------|-----------|-------------|
| 1 | Events/Audit | 🟡 Media | Historial de cambios |
| 2 | Hooks pre/post | 🟡 Media | Scripts antes/después de apply |

---

## Post-MVP 2 (0.3.0) - Orquestación

### Comandos

| # | Item | Prioridad | Descripción |
|---|------|-----------|-------------|
| 1 | `synctl scale <resource> --replicas=N` | 🟢 Baja | Escalar deployments |
| 2 | `synctl rollout undo <resource>` | 🟢 Baja | Rollback |
| 3 | `synctl diff -f <file>` | 🟢 Baja | Comparar con estado |

### Arquitectura

| # | Item | Prioridad | Descripción |
|---|------|-----------|-------------|
| 1 | Resource dependencies | 🟢 Baja | Orden de aplicación |
| 2 | Parallel apply | 🟢 Baja | Recursos independientes en paralelo |

---

## Post-MVP 3 (1.0.0) - Producción

| # | Item | Prioridad | Descripción |
|---|------|-----------|-------------|
| 1 | Plugin system | 🟢 Baja | Runtimes adicionales |
| 2 | MCP Protocol | 🟢 Baja | Integración IA |
| 3 | API Server mode | 🟢 Baja | Daemon con REST API |

---

## Roadmap de Comandos

```
v0.1.0 (MVP)                    v0.2.0                  v0.3.0
───────────────────            ──────────────           ───────────
✅ install                      ✅ logs                  ✅ scale
✅ deploy                       ✅ exec                  ✅ rollout
✅ apply                        ✅ describe              ✅ diff
✅ inspect                      ✅ port-forward         ✅ validate
✅ get resources                ✅ events               
✅ get nodes                    ✅ hooks                
🔴 delete (pendiente)           
🔴 deploy --dry-run (pendiente) 
🔴 output formats (pendiente)   
```

---

## Siguientes Pasos Inmediatos

### 1. Completar flujo apply completo
```bash
# Implementar
synctl apply --dry-run     # Validación sin ejecución
synctl delete <resource>   # Eliminación
```

### 2. Externalizar deploy resources
```bash
# De:
internal/application/services/apply/deploy/resources.yaml

# A:
~/.syncloud/deploy.yaml  # O configurable
```

### 3. Implementar ciclo observe
```go
# En reconciler
for _, resource := range state.Resources {
    actual := runtime.Observe(resource)
    diff := compare(resource.Spec, actual)
    if diff.HasChanges() {
        runtime.Reconcile(diff)
    }
}
```

---

## Notas

- 🔴 = Crítico para MVP
- 🟡 = Importante post-MVP
- 🟢 = Nice to have
- ✅ = Completado
- 🔄 = En progreso
