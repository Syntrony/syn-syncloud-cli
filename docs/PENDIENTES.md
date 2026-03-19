# Roadmap de Desarrollo - Syncloud CLI v1.0

## Visión: Primera Versión Liberable

El objetivo es una **v1.0 completa** que permita gestionar todo el ciclo de vida de Syncloud. Esta versión debe ser lo suficientemente robusta para uso en producción, con todas las operaciones CRUD sobre los objetos principales.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                     syncloud-state.json                                 │
│                     (FUENTE DE VERDAD)                                 │
│  ┌─────────┐  ┌──────────┐  ┌──────────┐  ┌────────────────────┐  │
│  │ Cluster │  │  Nodes    │  │Resources │  │    DNS Records     │  │
│  └─────────┘  └──────────┘  └──────────┘  └────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Objetos del Sistema

### 1. Cluster
Configuración global del cluster Syncloud.

```json
{
  "id": "uuid",
  "name": "syncloud",
  "mode": "container|vps",
  "dns": "syncloud.local"
}
```

### 2. Nodes
Nodos físicos o virtuales del cluster.

```json
{
  "id": "node-1",
  "hostname": "server-01",
  "role": "control-plane|worker",
  "ip": "192.168.1.100"
}
```

### 3. Resources
Recursos abstractos que se reconcilian a Docker/Kubernetes.

```json
{
  "id": "uuid",
  "name": "postgres-db",
  "kind": "Container|Deployment|Service|...",
  "runtime": "docker|kubernetes",
  "nodeId": "node-1",
  "spec": {...},
  "status": {...},
  "ownership": {...}
}
```

### 4. DNS Records
Registros DNS para acceso a recursos.

```json
{
  "records": [
    {"name": "app", "server": "192.168.1.100"},
    {"name": "api", "server": "192.168.1.101"}
  ]
}
```

---

## Comandos MVP v1.0 - COMPLETOS

### Cluster

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl install` | Instala y configura el cluster | ✅ |
| `synctl inspect` | Muestra estado del cluster | ✅ |
| `synctl cluster status` | Estado detallado del cluster | 🔴 Pendiente |
| `synctl cluster config` | Ver/actualizar configuración | 🔴 Pendiente |

### Nodes

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl get nodes` | Lista nodos | ✅ |
| `synctl node describe <name>` | Detalle de nodo | 🔴 Pendiente |
| `synctl node add -f <file>` | Agregar nodo | 🔴 Pendiente |
| `synctl node remove <name>` | Remover nodo | 🔴 Pendiente |
| `synctl node cordon <name>` | Marcar nodo como no scheduling | 🔴 Pendiente |
| `synctl node uncordon <name>` | Reactivar nodo | 🔴 Pendiente |

### Resources

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl deploy` | Despliega plataforma completa | ✅ |
| `synctl apply -f <file>` | Aplica recursos desde YAML | ✅ |
| `synctl get resources` | Lista recursos | ✅ |
| `synctl resource describe <name>` | Detalle de recurso | 🔴 Pendiente |
| `synctl resource delete <name>` | Elimina recurso | 🔴 Pendiente |
| `synctl resource logs <name>` | Ver logs | 🔴 Pendiente |
| `synctl resource exec <name> -- <cmd>` | Ejecutar comando | 🔴 Pendiente |
| `synctl resource scale <name> --replicas=N` | Escalar recurso | 🔴 Pendiente |
| `synctl resource restart <name>` | Reiniciar recurso | 🔴 Pendiente |

### DNS Records

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl get dns` | Lista registros DNS | 🔴 Pendiente |
| `synctl dns add <name> --server <ip>` | Agregar registro | 🔴 Pendiente |
| `synctl dns delete <name>` | Eliminar registro | 🔴 Pendiente |
| `synctl dns update <name> --server <ip>` | Actualizar registro | 🔴 Pendiente |
| `synctl dns resolve <name>` | Resolver DNS | 🔴 Pendiente |

### Validate & Diff

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl validate -f <file>` | Validar recursos | 🔴 Pendiente |
| `synctl diff -f <file>` | Comparar con estado actual | 🔴 Pendiente |
| `synctl apply --dry-run -f <file>` | Validar sin aplicar | 🔴 Pendiente |
| `synctl apply --preview -f <file>` | Preview de cambios | 🔴 Pendiente |

### Output

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl get resources --output json` | Formato JSON | 🔴 Pendiente |
| `synctl get resources --output yaml` | Formato YAML | 🔴 Pendiente |
| `synctl get resources --output table` | Formato tabla (default) | ✅ |

---

## Arquitectura - Requerimientos

### 🔴 CRÍTICO - v1.0

| # | Componente | Descripción |
|---|------------|-------------|
| 1 | CRUD Resources | Create, Read, Update, Delete sobre recursos |
| 2 | CRUD DNS | Create, Read, Update, Delete sobre DNS records |
| 3 | CRUD Nodes | Create, Read, Update, Delete sobre nodos |
| 4 | Observe Cycle | Leer estado actual del runtime |
| 5 | Diff Engine | Comparar desired vs actual state |
| 6 | Reconcile Engine | Aplicar cambios de forma idempotente |
| 7 | Status Updates | Actualizar Resource.status desde runtime |
| 8 | Error Handling | Manejo robusto de errores |
| 9 | Validation | Validación de recursos antes de aplicar |
| 10 | Rollback | Reversión ante fallos |

### 🟡 IMPORTANTE - v1.0

| # | Componente | Descripción |
|---|------------|-------------|
| 1 | Logging | Logging estructurado |
| 2 | Progress indicators | Feedback visual durante operaciones |
| 3 | Confirmation prompts | Confirmación para operaciones destructivas |
| 4 | Help & docs inline | Ayuda contextual en comandos |

---

## Roadmap de Implementación

### Fase 1: Fundamentos (Semana 1-2)

```bash
# Completar CRUD de recursos
synctl resource delete <name>          # Eliminar recurso
synctl resource describe <name>         # Describir recurso
synctl apply --dry-run -f <file>      # Dry run
```

### Fase 2: DNS Management (Semana 2-3)

```bash
# Gestionar DNS records
synctl get dns                         # Listar
synctl dns add <name> --server <ip>   # Agregar
synctl dns delete <name>              # Eliminar
synctl dns update <name> --server <ip> # Actualizar
```

### Fase 3: Node Management (Semana 3-4)

```bash
# Gestionar nodos
synctl node describe <name>           # Describir
synctl node add -f <file>            # Agregar
synctl node remove <name>             # Remover
synctl node cordon/uncordon <name>    # Scheduling
```

### Fase 4: Observability (Semana 4-5)

```bash
# Debug y diagnóstico
synctl resource logs <name>           # Logs
synctl resource exec <name> -- <cmd>  # Exec
synctl diff -f <file>                # Diff
```

### Fase 5: Polish (Semana 5-6)

```bash
# Output formatting
synctl get resources --output json    # JSON
synctl get resources --output yaml   # YAML

# Validation
synctl validate -f <file>            # Validar

# Scale & Restart
synctl resource scale <name> --replicas=3
synctl resource restart <name>
```

---

## Comandos por Categoría

### Currently Implemented (v0.0.x)
```
synctl install
synctl deploy
synctl apply -f
synctl inspect
synctl get resources
synctl get nodes
synctl version
```

### Phase 1 Target (v0.1.0)
```
synctl resource delete <name>
synctl resource describe <name>
synctl apply --dry-run -f
synctl apply --preview -f
```

### Phase 2 Target (v0.2.0)
```
synctl get dns
synctl dns add <name> --server <ip>
synctl dns delete <name>
synctl dns update <name> --server <ip>
```

### Phase 3 Target (v0.3.0)
```
synctl node describe <name>
synctl node add -f <file>
synctl node remove <name>
synctl node cordon <name>
synctl node uncordon <name>
```

### Phase 4 Target (v0.4.0)
```
synctl resource logs <name>
synctl resource exec <name> -- <cmd>
synctl diff -f <file>
```

### Phase 5 Target (v0.5.0)
```
synctl get resources --output json|yaml
synctl validate -f <file>
synctl resource scale <name> --replicas=N
synctl resource restart <name>
```

### v1.0.0 - RELEASE
```
Todos los comandos anteriores implementados
Documentación completa
Tests coverage > 80%
```

---

## Notas

- ✅ = Implementado
- 🔴 = Crítico para v1.0
- 🟡 = Importante para v1.0
- 🟢 = Nice to have post-v1.0
