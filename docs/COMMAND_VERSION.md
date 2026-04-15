# Comando: version

## Descripción

El comando `version` muestra información sobre la versión del CLI de Syncloud. Es un comando simple que permite identificar rápidamente qué versión del cliente está instalada.

## Uso

```bash
synctl version
```

## Flujo de Ejecución

```
┌─────────────────────────────────────────────────────────────────────┐
│                       synctl version                                │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  1. Inicialización del Logger                                      │
│     - Crear logger de consola para salida                          │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  2. Lectura de Versión                                             │
│     - La versión está hardcodeada en el código fuente             │
│     - Formato: semver (X.Y.Z)                                       │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  3. Presentación                                                   │
│     ┌──────────────────────────────────────────────────────────┐   │
│     │  Salida formateada:                                       │   │
│     │  "synctl version X.Y.Z"                                   │   │
│     └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

## Salida

```
synctl version 0.0.1
```

## Información de Versión

| Campo | Descripción |
|-------|-------------|
| **Versión actual** | 0.0.1 (versión inicial del MVP) |
| **Formato** | Semver (Major.Minor.Patch) |

## Consideraciones

- **Comando de solo lectura**: No modifica ningún estado
- **Sin dependencias**: No requiere que Syncloud esté instalado
- **Salida simple**: Muestra solo la versión del CLI
- **Versión hardcodeada**: La versión se define en tiempo de compilación
