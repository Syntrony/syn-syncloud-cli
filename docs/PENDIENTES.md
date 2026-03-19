# Roadmap de Desarrollo - Syncloud CLI

## Contexto: syncloud-state.json como Médula

El archivo `syncloud-state.json` es el contrato central del sistema:
- **Estado deseado**: Recursos declarados por el usuario
- **Persistencia**: Sobrevive reinicios
- **Sincronización**: Reconciler compara desired vs actual
- **DNS**: Campo `cluster.dns` define la URL de acceso

```
syncloud-state.json ──► Reconciler ──► Docker/K8s
         ▲                                    │
         │                                    ▼
         └──────────── Observe ──────────── Runtime
```

---

## Objetivo MVP: Completar Flujo Despliegue

El MVP permite:
1. `synctl install` → Genera state con cluster.dns
2. `synctl deploy` → Lee state, reconcilia, muestra DNS
3. Usuario accede a `http://syncloud.local`

### ✅ Implementado

| Item | Descripción |
|------|-------------|
| `synctl install` | Instala Docker/K3s + dnsmasq, genera state |
| `synctl deploy` | Aplica recursos de plataforma, muestra DNS |
| `synctl apply -f` | Aplica recursos desde YAML |
| `synctl inspect` | Muestra estado actual |
| `synctl get resources` | Lista recursos |
| `synctl get nodes` | Lista nodos |
| Runtime Reconciler | Adaptador Docker/K8s |
| syncloud-state.json | Fuente de verdad central |

---

## Pendientes MVP

### 🔴 CRÍTICO - Requeridos para MVP funcional

#### Comandos

| # | Item | Descripción | Impacto |
|---|------|-------------|---------|
| 1 | `synctl delete <resource>` | Elimina recurso del state y runtime | Completar ciclo CRUD |
| 2 | `synctl deploy --dry-run` | Valida recursos sin aplicar | Evitar errores en deploy |
| 3 | Output formatters | `get resources --output json/yaml` | Depuración y scripting |

#### Arquitectura

| # | Item | Descripción | Impacto |
|---|------|-------------|---------|
| 1 | Ciclo Observe | Leer estado actual del runtime | Sincronizar desired vs actual |
| 2 | Diff separado | Diff → Plan → Execute | No ejecutar directo |
| 3 | Actualización de status | RuntimeResourceState → Resource.status | Reflejar estado real en state |

---

## Post-MVP 1 - Developer Experience

### 🟡 MEDIA - Mejoras de UX

| # | Comando | Descripción |
|---|---------|-------------|
| 1 | `synctl logs <resource>` | Ver logs de contenedor/pod |
| 2 | `synctl exec <resource> -- <cmd>` | Ejecutar comando en recurso |
| 3 | `synctl describe <resource>` | Descripción detallada |
| 4 | `synctl validate -f <file>` | Solo validación sin aplicar |

---

## Post-MVP 2 - Orquestación

### 🟢 BAJA - Nice to have

| # | Comando | Descripción |
|---|---------|-------------|
| 1 | `synctl scale <resource> --replicas=N` | Escalar deployments |
| 2 | `synctl rollout undo <resource>` | Rollback |
| 3 | `synctl diff -f <file>` | Comparar con estado actual |

---

## Roadmap Visual

```
MVP (v0.1.0)                    v0.2.0                  v0.3.0
───────────────────            ──────────────           ───────────
✅ install                      ✅ logs                  ✅ scale
✅ deploy                       ✅ exec                  ✅ rollout
✅ apply                        ✅ describe              ✅ diff
✅ inspect                      ✅ validate             
✅ get resources                🔴 delete (pendiente)   
✅ get nodes                    🔴 dry-run (pendiente)  
🔴 delete (pendiente)          🔴 output formats       
🔴 dry-run (pendiente)                                   
🔴 output formats                                          
🔴 observe cycle                                         
🔴 diff separated                                       
```

---

## Siguientes Pasos Inmediatos

### 1. Implementar delete
```bash
synctl delete postgres-db
# 1. Leer state.json
# 2. Remover recurso del array resources
# 3. Ejecutar docker rm / kubectl delete
# 4. Guardar state.json actualizado
```

### 2. Implementar deploy --dry-run
```bash
synctl deploy --dry-run
# 1. Parsear recursos
# 2. Validar
# 3. Mostrar diff sin aplicar
```

### 3. Implementar observe cycle
```go
// En reconciler
for _, resource := range state.Resources {
    observed := runtime.Observe(resource)
    resource.Status = observed
}
repo.Save(state)
```

### 4. Output formatters
```bash
synctl get resources --output json
synctl get resources --output yaml
```

---

## Notas

- 🔴 = Crítico para MVP
- 🟡 = Importante post-MVP  
- 🟢 = Nice to have
- ✅ = Completado
- 🔄 = En progreso
