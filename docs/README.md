# synctl - Syncloud CLI

`synctl` es una herramienta CLI en Go para gestionar recursos de Syncloud utilizando Cobra y Clean Architecture.

## Uso

```bash
# Instalar Syncloud en el cluster
./synctl install

# Inspeccionar recursos existentes
./synctl inspect

# Listar nodos
./synctl get nodes

# Listar recursos (con filtros)
./synctl get resources
./synctl get resources --name myapp
./synctl get resources --runtime k3s
./synctl get resources --kind deployment

# Version
./synctl version
```

## Comandos

| Comando | Descripcion |
|---------|-------------|
| `install` | Instala Syncloud en el cluster |
| `inspect` | Inspecciona recursos existentes |
| `get nodes` | Lista nodos del sistema |
| `get resources` | Lista recursos con filtros |
| `version` | Muestra version |

## Filtros para recursos

| Flag | Descripcion |
|------|-------------|
| `-n, --name` | Filtrar por nombre de recurso |
| `-k, --kind` | Filtrar por tipo de recurso |
| `-r, --runtime` | Filtrar por runtime (k3s, docker) |
| `-i, --id` | Filtrar por ID de recurso |

## Modos de ejecucion

- **Container (k3d)**: Detecta automaticamente cuando se ejecuta dentro de un contenedor
- **VPS (k3s)**: Ejecucion estandar en servidor
- **Sudo**: Solicita contrasena sudo si no es root

## Estado

El estado de Syncloud se guarda en: `helpers/syncloud-state.json`
