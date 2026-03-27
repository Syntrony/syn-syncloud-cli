# Contexto de Proyecto: Syncloud Platform CLI (synctl)

**Propietario:** Syntrony Technologies Inc.
**Lenguaje:** Go (Golang)
**Objetivo del MVP:** Proveer comandos funcionales para la administración general de recursos en Docker y Kubernetes mediante un modelo unificado.

---

## Contexto General

**Syncloud Platform (synctl)** es una herramienta CLI que abstrae la funcionalidad del runtime (Docker, Kubernetes) proporcionando una interfaz unificada para gestionar recursos de ambos ecosistemas. 

La plataforma está diseñada para:
- Simplificar la administración de infraestructura containerizada
- Proporcionar un modelo declarativo de recursos (CRUD)
- Gestionar múltiples runtimes desde una sola herramienta
- Mantener un estado centralizado de la plataforma

### Recursos Soportados

**Docker:** Container, Network, Image

**Kubernetes:** Deployment, Service, Ingress, Namespace, Secret, ConfigMap, ServiceAccount, Role, RoleBinding, ClusterRole, Pod, Endpoints, K8sResource

### Estado de la Plataforma

El archivo `helpers/syncloud-state.json` es el centro del sistema:
- **Estado deseado**: Recursos declarados
- **Estado actual**: Recursos observados
- **Diff**: Diferencia para reconciliar
- **DNS Records**: Registros para acceso
- **Nodes**: Nodos del cluster

---

## Estructura del Código

```
syn-sycloud-cli/
├── main.go                  # Punto de entrada de la aplicación
├── cmd/                     # Comandos CLI (Cobra)
│   ├── root.go              # Comando raíz
│   ├── apply/               # Comando apply
│   ├── daemon/              # Comando daemon
│   ├── get/                 # Comando get (recursos, nodos)
│   ├── inspect/             # Comando inspect
│   ├── install/             # Comando install
│   └── version/             # Comando version
├── internal/                # Lógica de negocio
│   ├── app/                 # Componentes de la aplicación
│   ├── application/         # Capa de aplicación (servicios)
│   │   ├── interfaces/      # Contratos y abstracciones
│   │   ├── outputs/        # Formateador de salida
│   │   └── services/        # Servicios de negocio
│   ├── domain/              # Núcleo del negocio
│   │   ├── dto/             # Objetos de transferencia
│   │   ├── filters/         # Filtros de consulta
│   │   ├── parser/          # Analizadores de recursos
│   │   ├── persistence/     # Repositorios de estado
│   │   ├── Resource/       # Definiciones de recursos
│   │   ├── docker/          # Modelos Docker
│   │   ├── k8s/            # Modelos Kubernetes
│   │   └── dns/            # Modelos DNS
│   ├── infrastructure/     # Implementaciones concretas
│   │   ├── docker/         # Adaptador Docker
│   │   ├── k8s/            # Adaptador Kubernetes
│   │   ├── dns/            # Adaptador DNS
│   │   └── system/         # Adaptador sistema
│   ├── executor/           # Ejecución de comandos
│   └── logger/             # Sistema de logging
└── helpers/                 # Utilidades y archivos de estado
```

### Arquitectura

El proyecto sigue **Clean Architecture / Hexagonal Architecture**:

1. **cmd/**: Punto de entrada. Maneja los comandos de la CLI usando Cobra. No contiene lógica de negocio.

2. **internal/application/**: Capa de orquestación. Define servicios y casos de uso. Depende de interfaces del dominio.

3. **internal/domain/**: Núcleo del negocio. Contiene los modelos, contextos de reconciliación y definiciones de recursos.

4. **internal/infrastructure/**: Implementaciones concretas de interfaces (Docker, Kubernetes, DNS, sistema).

5. **internal/executor/**: Abstracciones para la ejecución de comandos del sistema.

### Patrón de Reconciliación

La lógica de reconciliación sigue un ciclo de:
1. **Detección (Diffing)**: Identificar el estado deseado vs. estado actual
2. **Planificación**: Generar un `ReconcileContext` con acciones (Create/Update/Delete)
3. **Ejecución**: Aplicar las acciones a través de los adaptadores de infraestructura

---

## Listado de Comandos

| Comando | Descripción | Estado |
|---------|-------------|--------|
| `synctl install` | Instalación de runtime y servicios necesarios (Docker, K3s, dnsmasq) | ✅ Implementado |
| `synctl deploy` | Despliegue inicial de la plataforma (pendiente) | ⏳ Pendiente |
| `synctl apply -f <file>` | Creación/Actualización de recursos basados en YAML | ✅ Implementado |
| `synctl get resources` | Listado de recursos del clúster | ✅ Implementado |
| `synctl get nodes` | Listado de nodos del sistema | ✅ Implementado |
| `synctl describe` | Detalle profundo de un recurso específico | ⏳ Pendiente |
| `synctl delete` | Eliminación por nombre o por archivo YAML | ⏳ Pendiente |
| `synctl status` | Estado de salud del clúster de la plataforma | ⏳ Pendiente |
| `synctl version` | Información de versión de synctl y plataforma | ✅ Implementado |
| `synctl inspect` | Diagnóstico profundo del clúster | ✅ Implementado |
| `synctl daemon` | Gestión del proceso en segundo plano de la plataforma | ✅ Implementado |

### Flujo General de Uso

```
┌─────────────────────────────────────────────────────────────────────┐
│  synctl install → Configurar runtimes (Docker/K3s + DNS)           │
│                         ↓                                           │
│  synctl deploy  → Desplegar DB + Backend + WebApp                  │
│                         ↓                                           │
│  synctl apply -f <file> → Crear/actualizar recursos                 │
│                         ↓                                           │
│  synctl get resources/nodes → Gestionar plataforma                 │
│                         ↓                                           │
│  synctl daemon → Iniciar control plane en segundo plano            │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Reglas de Arquitectura

1. **Idempotencia:** Los comandos `apply` y `deploy` deben ser seguros de ejecutar múltiples veces.

2. **Modelado Unificado:** Abstraer la complejidad de K8s/Docker bajo el modelo de Syncloud.

3. **Logs:** Formato limpio para terminal (utilizar los estándares de Syntrony definidos en el código).

4. **Principios SOLID:**
   - **SRP:** Cada capa tiene una única responsabilidad.
   - **DIP:** La capa de `application` depende de interfaces definidas en el `domain`.
