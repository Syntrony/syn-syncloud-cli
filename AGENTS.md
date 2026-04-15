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

---

### Estado de la Plataforma

El archivo `helpers/syncloud-state.json` es el centro del sistema:
- **Estado deseado**: Recursos declarados
- **Estado actual**: Recursos observados
- **Diff**: Diferencia para reconciliar
- **DNS Records**: Registros para acceso
- **Nodes**: Nodos del cluster

---

### Patrón de Reconciliación

La lógica de reconciliación sigue un ciclo de:
1. **Detección (Diffing)**: Identificar el estado deseado vs. estado actual
2. **Planificación**: Generar un `ReconcileContext` con acciones (Create/Update/Delete)
3. **Ejecución**: Aplicar las acciones a través de los adaptadores de infraestructura

### Abstracción de Mutaciones

El sistema utiliza una abstracción de mutaciones para compartir código entre comandos:

- **MutationService**: Servicio genérico que maneja parser, validación y ejecución
- **ApplyMutation**: Implementación específica para aplicar recursos (Upsert)
- **DeleteMutation**: Implementación específica para eliminar recursos (Remove)

Esta abstracción permite que apply y delete compartan la misma lógica de:
- Parseo de archivos YAML
- Validación de recursos
- Construcción de estado
- Persistencia

Los componentes compartidos están en `internal/application/services/common/`:
- **parser/**: Parseo de recursos desde YAML
- **validator/**: Validación de estructura y reglas de recursos
- **diff/**: Comparación de estado deseado vs actual

---

## Listado de Comandos

| Comando | Descripción | Estado | |
|---------|-------------|--------|-|
| `synctl install` | Instalación de runtime y servicios necesarios (Docker, K3s, dnsmasq) | ✅ Implementado | [INSTALL_CMD](/docs/COMMAND_INSTALL.md)
| `synctl deploy` | Despliegue inicial de la plataforma (pendiente) | ⏳ Pendiente |
| `synctl apply -f <file>` | Creación/Actualización de recursos basados en YAML | ✅ Implementado | [APPLY_CMD](/docs/COMMAND_APPLY.md)
| `synctl get resources` | Listado de recursos del clúster | ✅ Implementado | [GET_CMD](/docs/COMMAND_GET.md)
| `synctl get nodes` | Listado de nodos del sistema | ✅ Implementado | [GET_CMD](/docs/COMMAND_GET.md)
| `synctl describe` | Detalle profundo de un recurso específico | ✅ Implementado | [COMMAND_DESCRIBE](/docs/COMMAND_DESCRIBE.md)
| `synctl reconcile` | Sincronización bidireccional estado ↔ runtime (importar recursos) | ⏳ Pendiente |
| `synctl logs` | Logs de algun recurso en especifico | ⏳ Pendiente |
| `synctl delete resource` | Eliminación por nombre, id o archivo YAML | ✅ Implementado | [DELETE_CMD](/docs/COMMAND_DELETE.md)
| `synctl status` | Estado de salud del clúster de la plataforma | ⏳ Pendiente |
| `synctl version` | Información de versión de synctl y plataforma | ✅ Implementado | [VERSION_CMD](/docs/COMMAND_VERSION.md)
| `synctl inspect` | Diagnóstico profundo del clúster | ✅ Implementado | [INSPECT_CMD](/docs/COMMAND_INSPECT.md)
| `synctl daemon` | Gestión del proceso en segundo plano de la plataforma | ✅ Implementado | [DAEMON_CMD](/docs/COMMAND_DAEMON.md)

## Reglas de Arquitectura

1. **Idempotencia:** Los comandos `apply` y `deploy` deben ser seguros de ejecutar múltiples veces.

2. **Modelado Unificado:** Abstraer la complejidad de K8s/Docker bajo el modelo de Syncloud.

3. **Logs:** Formato limpio para terminal (utilizar los estándares de Syntrony definidos en el código).

4. **Principios SOLID:**
   - **SRP:** Cada capa tiene una única responsabilidad.
   - **DIP:** La capa de `application` depende de interfaces definidas en el `domain`.