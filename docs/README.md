# synctl - Syncloud CLI

`synctl` es una herramienta CLI para gestionar **Syncloud Resources** de forma declarativa, orquestando contenedores en Docker y Kubernetes a través de un modelo de recursos universal.

## Flujo MVP v1.0

```bash
# 1. Instalar runtimes (Docker/K3s, dnsmasq)
synctl install

# 2. Desplegar plataforma completa
synctl deploy

# 3. Gestionar
synctl get resources          # Ver recursos
synctl get nodes              # Ver nodos
synctl get dns                # Ver DNS records

# 4. Acceder
# http://syncloud.local
```

## Arquitectura Conceptual

Syncloud **NO administra Docker o Kubernetes directamente**. Administra **Syncloud Resources** que luego se reconcilian hacia los runtimes apropiados.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         synctl CLI                                      │
│  ┌─────────┐     ┌─────────────────────┐     ┌─────────────────────┐  │
│  │ install │     │      deploy         │     │      apply          │  │
│  └────┬────┘     └──────────┬──────────┘     └──────────┬──────────┘  │
│  ┌────┴────┐     ┌──────────┴──────────┐     ┌──────────┴──────────┐ │
│  │  dns    │     │ resource describe   │     │ resource delete    │ │
│  │  node   │     │ resource logs      │     │ resource exec      │ │
│  └────┬────┘     └──────────┬──────────┘     └──────────┬──────────┘ │
└───────┼─────────────────────┼───────────────────────────┼─────────────┘
        │                     │                           │
        ▼                     ▼                           ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    syncloud-state.json                                 │
│                    (FUENTE DE VERDAD - MÉDULA DEL SISTEMA)            │
│  ┌──────────┐  ┌──────────┐  ┌───────────┐  ┌──────────────────────┐  │
│  │ Cluster  │  │  Nodes   │  │ Resources │  │    DNS Records       │  │
│  └──────────┘  └──────────┘  └───────────┘  └──────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
        │                     │                           │
        ▼                     ▼                           ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      Reconciler Layer                                   │
│  ┌─────────────────────┐         ┌─────────────────────┐             │
│  │  DockerReconciler   │         │   K8sReconciler     │             │
│  │  docker run / APIs  │         │  kubectl apply/APIs │             │
│  └─────────────────────┘         └─────────────────────┘             │
└─────────────────────────────────────────────────────────────────────────┘
```

## Comandos

### Instalación y Despliegue

```bash
synctl install              # Instala Syncloud y dependencias
synctl deploy              # Despliega la plataforma completa
synctl inspect              # Inspecciona estado actual
```

### Gestión de Recursos

```bash
synctl apply -f <file>                  # Aplicar recursos desde YAML
synctl get resources                     # Lista recursos
synctl get resources --name myapp       # Filtrar por nombre
synctl get resources --kind Deployment  # Filtrar por tipo
synctl get resources --runtime docker   # Filtrar por runtime
synctl resource describe <name>         # Ver detalle de recurso
synctl resource delete <name>           # Eliminar recurso
synctl resource logs <name>            # Ver logs
synctl resource exec <name> -- <cmd>   # Ejecutar comando
synctl resource scale <name> --replicas=3 # Escalar
synctl resource restart <name>          # Reiniciar
```

### Gestión de DNS

```bash
synctl get dns                           # Lista registros DNS
synctl dns add <name> --server <ip>     # Agregar registro
synctl dns delete <name>                # Eliminar registro
synctl dns update <name> --server <ip>  # Actualizar registro
```

### Gestión de Nodos

```bash
synctl get nodes                         # Lista nodos
synctl node describe <name>             # Ver detalle de nodo
synctl node add -f <file>              # Agregar nodo
synctl node remove <name>               # Remover nodo
synctl node cordon <name>               # Deshabilitar scheduling
synctl node uncordon <name>             # Habilitar scheduling
```

### Validación y Diagnóstico

```bash
synctl validate -f <file>               # Validar recursos
synctl diff -f <file>                   # Comparar con estado
synctl apply --dry-run -f <file>       # Validar sin aplicar
synctl apply --preview -f <file>       # Preview de cambios
```

### Output

```bash
synctl get resources --output json       # Formato JSON
synctl get resources --output yaml      # Formato YAML
synctl get resources --output table     # Formato tabla
```

### Utilidades

```bash
synctl version                           # Muestra versión
```

## Modelo de Recursos

Syncloud usa un modelo declarativo universal basado en `apiVersion: syncloud/v1`:

```yaml
apiVersion: syncloud/v1
kind: Container           # Kind abstracto
metadata:
  name: myapp
  namespace: production
spec:
  image: nginx:latest
  ports:
    - hostPort: 8080
      containerPort: 80
```

### Recursos Disponibles

**Docker:**
- `Container` - Contenedor
- `Network` - Red
- `Image` - Imagen

**Kubernetes:**
- `Deployment` - Deployment
- `Service` - Servicio
- `Ingress` - Ingress
- `Namespace` - Namespace
- `Secret` - Secret
- `ConfigMap` - ConfigMap
- `ServiceAccount` - ServiceAccount
- `Role` / `ClusterRole` - RBAC
- `Pod` - Pod individual
- `K8sResource` - Recurso raw (passthrough)

## syncloud-state.json

El estado se guarda en: `helpers/syncloud-state.json`

```json
{
  "version": "1.0.0",
  "cluster": {
    "id": "uuid",
    "name": "syncloud",
    "mode": "container|vps",
    "dns": "syncloud.local"
  },
  "nodes": [...],
  "resources": [...],
  "dns": {
    "records": [
      {"name": "app", "server": "192.168.1.100"}
    ]
  }
}
```

## Modos de Ejecución

| Modo | Descripción | Runtime |
|------|-------------|---------|
| **Container** | Ejecutando en contenedor (k3d) | K3d |
| **VPS** | Ejecutando en servidor | K3s |

El modo se detecta automáticamente. Si es necesario sudo, se solicita la contraseña.

## Ejemplos

Ver `examples/resources.yaml` para ejemplos completos de recursos.

```bash
# Aplicar ejemplos
synctl apply -f examples/resources.yaml

# Ver recursos aplicados
synctl get resources

# Ver DNS
synctl get dns

# Inspeccionar estado
synctl inspect
```

## Desarrollo

### Build

```bash
go build -o synctl .          # Compilar
go run .                      # Ejecutar sin compilar
```

### Linting

```bash
go fmt ./... && go vet ./... && go build -o /dev/null .
```

Ver `AGENTS.md` para guía completa de desarrollo.
