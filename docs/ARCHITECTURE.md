# Arquitectura de Syncloud CLI

## Visión General
Syncloud CLI sigue una arquitectura basada en capas (Clean Architecture / Hexagonal Architecture) diseñada para separar las preocupaciones de infraestructura, lógica de dominio y comandos de usuario.

### Estructura de Capas
- **`cmd/`**: Punto de entrada. Maneja los comandos de la CLI (usando `cobra`). No contiene lógica de negocio.
- **`internal/application/`**: Capa de orquestación. Define servicios y casos de uso.
- **`internal/domain/`**: Núcleo del negocio. Contiene los modelos, contextos de reconciliación y definiciones de recursos.
- **`internal/infrastructure/`**: Implementaciones concretas de interfaces (Docker, Kubernetes, DNS, etc.).
- **`internal/executor/`**: Abstracciones para la ejecución de comandos de sistema.

## Patrón de Reconciliación
La lógica de reconciliación se basa en un ciclo de:
1. **Detección (Diffing):** Identificar el estado deseado vs. estado actual.
2. **Planificación:** Generar un `ReconcileContext` que contiene una lista de `Actions` (Create/Update/Delete).
3. **Ejecución:** Aplicar las acciones a través de los adaptadores de infraestructura (`internal/infrastructure`).

## Principios SOLID Aplicados
- **SRP:** Cada capa tiene una única responsabilidad.
- **DIP:** La capa de `application` depende de interfaces definidas en el `domain`, y la capa de `infrastructure` implementa esas interfaces.
