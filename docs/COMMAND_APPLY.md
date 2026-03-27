# Comando: apply

## Descripción

El comando `apply` permite crear o actualizar recursos en la plataforma Syncloud basándose en archivos YAML declarativos. Este comando implementa el patrón de recursos declarativos donde el usuario define el estado deseado y el sistema se encarga de alcanzarlo.

## Uso

```bash
synctl apply -f <ruta_archivo_yaml>
```

### Parámetros

| Parámetro | Alias | Descripción |
|-----------|-------|-------------|
| `--file` | `-f` | Ruta al archivo YAML con definiciones de recursos (requerido) |

## Flujo de Ejecución

```
┌─────────────────────────────────────────────────────────────────────┐
│                     synctl apply -f <file>                         │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  1. Validación del Archivo                                          │
│     - Verificar que el archivo existe                               │
│     - Verificar que el archivo es legible                           │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  2. Inicialización de Componentes                                  │
│     - Logger: Inicializar sistema de logs                          │
│     - Repository: Cargar estado actual del sistema                 │
│     - Parser: Configurar parser de recursos YAML                  │
│     - Validator: Configurar validador de recursos                  │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  3. Parseo del Archivo YAML                                         │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  YAML Input                                              │   │
│     │  ─────────────────────────────────────────────────────   │   │
│     │  kind: Deployment                                        │   │
│     │  metadata:                                               │   │
│     │    name: my-app                                          │   │
│     │  spec:                                                   │   │
│     │    replicas: 3                                           │   │
│     └──────────────────────────────────────────────────────────┘   │
│                           │                                         │
│                           ▼                                         │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  ResourceParser                                          │   │
│     │  - Convertir YAML a objetos del dominio                  │   │
│     │  - Identificar tipo de recurso (K8s/Docker)              │   │
│     │  - Extraer metadatos y especificaciones                 │   │
│     └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  4. Validación de Recursos                                          │
│     - Validator: Validar estructura del recurso                    │
│     - Verificar campos requeridos                                  │
│     - Validar referencias entre recursos                           │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  5. Reconciliación (ApplyService)                                   │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  StateBuilder: Comparar estado desired vs actual         │   │
│     │  ┌────────────────┐      ┌────────────────┐            │   │
│     │  │  Estado Actual  │ ───► │ Estado Deseado │            │   │
│     │  └────────────────┘      └────────────────┘            │   │
│     │          │                        │                     │   │
│     │          └────────┬────────────────┘                     │   │
│     │                   ▼                                      │   │
│     │          ┌────────────────┐                             │   │
│     │          │     Diff/Plan   │                             │   │
│     │          └────────────────┘                             │   │
│     │                   │                                      │   │
│     │          ┌────────┴────────┐                             │   │
│     │          ▼                 ▼                             │   │
│     │  ┌────────────┐   ┌────────────┐                         │   │
│     │  │   Create   │   │   Update   │                         │   │
│     │  └────────────┘   └────────────┘                         │   │
│     └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  6. Ejecución en Runtime Destino                                    │
│     ┌──────────────────┐        ┌──────────────────┐              │
│     │  Docker Adapter  │        │  Kubernetes      │              │
│     │                  │        │  Adapter         │              │
│     │  - Containers   │        │  - Deployments   │              │
│     │  - Networks     │        │  - Services      │              │
│     │  - Images       │        │  - Ingress       │              │
│     └──────────────────┘        └──────────────────┘              │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  7. Actualización de Estado                                         │
│     - Persistir nuevo estado en syncloud-state.json                │
│     - Registrar recursos creados/actualizados                       │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Final: Recursos aplicados correctamente                            │
│  - Estado actualizado                                              │
│  - Recursos operativos                                              │
└─────────────────────────────────────────────────────────────────────┘
```

## Tipos de Recursos Soportados

### Kubernetes
- Deployment
- Service
- Ingress
- Namespace
- Secret
- ConfigMap
- ServiceAccount
- Role / RoleBinding
- ClusterRole / ClusterRoleBinding
- Pod
- Endpoints

### Docker
- Container
- Network
- Image

## Características

### Idempotencia
El comando es idempotente: ejecutarlo múltiples veces con el mismo archivo producirá el mismo resultado final. El sistema detectará recursos existentes y los actualizará según sea necesario.

### Reconciliación
El sistema compara el estado deseado (definido en el YAML) con el estado actual (recursos existentes) y determina las acciones necesarias:
- **Create**: Crear recursos que no existen
- **Update**: Actualizar recursos que existen pero tienen cambios
- **Delete**: (No realizado por apply, requiere comando específico)

### Validación
Antes de aplicar recursos, el sistema valida:
- Estructura del YAML
- Campos requeridos
- Referencias a otros recursos
- Compatibilidad con el runtime destino
