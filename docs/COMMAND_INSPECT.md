# Comando: inspect

## Descripción

El comando `inspect` proporciona una visión general del estado de la plataforma Syncloud. Es una herramienta de diagnóstico que permite verificar rápidamente si la plataforma está inicializada y operativa.

## Uso

```bash
synctl inspect
```

## Flujo de Ejecución

```
┌─────────────────────────────────────────────────────────────────────┐
│                        synctl inspect                               │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  1. Inicialización de Componentes                                  │
│     - Logger: Inicializar sistema de logs                          │
│     - Repository: Cargar estado actual del sistema                 │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  2. Ejecución del Servicio de Inspección                          │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  InspectService                                         │   │
│     │  - Consultar repositorio de estado                       │   │
│     │  - Verificar si Syncloud está inicializado               │   │
│     │  - Extraer información de versión y modo                 │   │
│     │  - Contar nodos y recursos registrados                  │   │
│     └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  3. Evaluación de Resultado                                        │
│     ┌──────────────────┐      ┌──────────────────┐                 │
│     │  Found = true   │      │  Found = false   │                 │
│     └──────────────────┘      └──────────────────┘                 │
│           │                           │                             │
│           ▼                           ▼                             │
│     ┌──────────────┐          ┌──────────────────┐                 │
│     │ Mostrar      │          │ Mostrar mensaje  │                 │
│     │ información  │          │ "Syncloud not    │                 │
│     │ del estado   │          │  initialized"    │                 │
│     └──────────────┘          └──────────────────┘                 │
└─────────────────────────────────────────────────────────────────────┘
```

## Información Reportada

Cuando Syncloud está inicializado, el comando muestra:

| Campo | Descripción |
|-------|-------------|
| **Syncloud State** | Estado de inicialización (FOUND/NOT FOUND) |
| **Version** | Versión de Syncloud Platform |
| **Mode** | Modo de operación (DEV/PROD) |
| **Nodes** | Cantidad de nodos registrados |
| **Resources** | Cantidad de recursos gestionados |

## Casos de Uso

### Verificar instalación
```bash
synctl inspect
```
Salida esperada si está instalado:
```
Syncloud State: FOUND
Version: 1.0.0
Mode: VPS
Nodes: 1
Resources: 5
```

### Verificar si Syncloud está inicializado
```bash
synctl inspect
```
Salida esperada si NO está instalado:
```
Syncloud not initialized
```

## Diferencia con otros comandos

| Comando | Propósito |
|---------|-----------|
| `inspect` | Resumen general del estado de la plataforma |
| `get nodes` | Lista detallada de nodos con sus propiedades |
| `get resources` | Lista detallada de recursos con filtros |
| `status` | Estado de salud de los servicios (pendiente) |

## Consideraciones

- **Solo lectura**: No modifica ningún estado del sistema
- **Lectura del estado**: Consulta el archivo de estado local
- **Diagnóstico rápido**: Diseñado para verificación rápida
