# synctl - Syncloud CLI

`synctl` es una herramienta CLI para gestionar **Syncloud Resources** de forma declarativa.

## Filosofía

Syncloud **NO administra Docker o Kubernetes directamente**. Administra recursos abstractos que luego se reconcilian hacia el runtime apropiado:

```
synctl apply → Syncloud Resource → Runtime Adapter → Docker/Kubernetes
```

## Comandos

```bash
# Instalación
synctl install              # Instala Syncloud y dependencias

# Gestión de recursos
synctl apply -f <file>      # Aplica recursos desde YAML
synctl inspect              # Inspecciona estado actual
synctl get resources         # Lista recursos
synctl get nodes            # Lista nodos

# Utilidades
synctl version              # Muestra versión
```

### apply

Aplica recursos Syncloud desde archivos YAML:

```bash
synctl apply -f resources.yaml
synctl apply -f examples/resources.yaml
```

**Flag:**
- `-f, --file` - Ruta al archivo YAML (requerido)

### get resources

Lista recursos con filtros opcionales:

```bash
synctl get resources                    # Todos los recursos
synctl get resources --name myapp       # Por nombre
synctl get resources --kind Deployment # Por tipo
synctl get resources --runtime docker  # Por runtime
synctl get resources --namespace default # Por namespace
```

**Flags:**
- `-n, --name` - Filtrar por nombre
- `-k, --kind` - Filtrar por tipo de recurso
- `-r, --runtime` - Filtrar por runtime (docker, kubernetes)
- `--namespace` - Filtrar por namespace

### get nodes

Lista nodos del sistema:

```bash
synctl get nodes
```

### inspect

Inspecciona el estado actual de Syncloud:

```bash
synctl inspect
```

Muestra: versión, modo (Container/VPS), nodos y conteo de recursos.

### install

Instala Syncloud y sus dependencias:

```bash
synctl install
```

Instala automáticamente:
- Runtime detectado (K3s/K3d o Docker)
- DNS (dnsmasq)
- Configuración de red

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

## Modos de Ejecución

| Modo | Descripción | Runtime |
|------|-------------|---------|
| **Container** | Ejecutando en contenedor (k3d) | K3d |
| **VPS** | Ejecutando en servidor | K3s |

El modo se detecta automáticamente. Si es necesario sudo, se solicita la contraseña.

## Estado

El estado se guarda en: `helpers/syncloud-state.json`

## Ejemplos

Ver `examples/resources.yaml` para ejemplos completos de recursos.

```bash
# Aplicar ejemplos
synctl apply -f examples/resources.yaml

# Ver recursos aplicados
synctl get resources

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
