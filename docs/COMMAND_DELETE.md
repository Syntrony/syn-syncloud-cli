# Comando: delete

## Descripción

El comando `delete` elimina recursos de la plataforma Syncloud. Soporta tres modos de selección: por archivo YAML declarativo, por nombre o por ID.

## Uso

```bash
synctl delete -f <ruta_archivo_yaml>
synctl delete -n <nombre_recurso>
synctl delete -I <id_recurso>
```

### Parámetros

| Parámetro | Alias | Descripción |
|-----------|-------|-------------|
| `--file`  | `-f`  | Ruta al archivo YAML con definiciones de recursos a eliminar |
| `--name`  | `-n`  | Nombre del recurso a eliminar |
| `--id`    | `-I`  | ID del recurso a eliminar |

Al menos uno de los tres flags es requerido.

## Flujo de Ejecución

```
┌─────────────────────────────────────────────────────────────────────┐
│          synctl delete -f <file> | -n <name> | -I <id>             │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  1. Validación de Flags                                             │
│     - Verificar que al menos un flag fue provisto                   │
│     - Si -f: verificar que el archivo existe y es legible           │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  2. Resolución de Recursos                                          │
│     ┌──────────────────────┐     ┌──────────────────────────────┐  │
│     │  Modo: -f (archivo)  │     │  Modo: -n/-I (filtro)        │  │
│     │  ResourceParser      │     │  GetResourceService          │  │
│     │  - Parsear YAML      │     │  - Buscar por nombre o ID    │  │
│     │  - Extraer recursos  │     │  - Cargar del estado actual  │  │
│     └──────────────────────┘     └──────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  3. MutationService → DeleteMutation                                │
│     - Validar recursos resueltos                                    │
│     - Ejecutar DeleteMutation (Remove del estado deseado)          │
│     - Persistir cambios en syncloud-state.json                     │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  4. Reconciliación en Runtime Destino                               │
│     ┌──────────────────┐        ┌──────────────────┐              │
│     │  Docker Adapter  │        │  Kubernetes       │              │
│     │  - Containers   │        │  Adapter          │              │
│     │  - Networks     │        │  - Deployments    │              │
│     │  - Images       │        │  - Services       │              │
│     └──────────────────┘        └──────────────────┘              │
│     Acción: ActionDelete por cada recurso del runtime              │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Final: Recursos eliminados del runtime y del estado               │
└─────────────────────────────────────────────────────────────────────┘
```

## Modos de Selección

### Por archivo YAML (`-f`)
Elimina todos los recursos definidos en el archivo. Usa el mismo formato YAML que `apply`.

```bash
synctl delete -f deployment.yaml
```

### Por nombre (`-n`)
Busca un recurso en el estado actual por su campo `metadata.name`.

```bash
synctl delete -n my-app
```

### Por ID (`-I`)
Busca un recurso en el estado actual por su ID interno.

```bash
synctl delete -I abc-123
```

## Notas

- Si no se encuentran recursos coincidentes, el comando termina sin error.
- La eliminación opera sobre Docker y Kubernetes de forma independiente según el runtime del recurso.
- El estado en `helpers/syncloud-state.json` se actualiza tras la eliminación.
