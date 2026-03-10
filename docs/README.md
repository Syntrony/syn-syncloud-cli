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

## Devcontainer

Para ejecutar en un entorno de desarrollo con devcontainer:

### Levantar dnsmasq en contenedor

```bash
dnsmasq --no-daemon --log-queries --conf-dir=/etc/dnsmasq.d,*.conf
```

### Validar configuracion con dig

```bash
# Resolver un dominio configurado
dig syncloud.local @127.0.0.1

# Ver respuesta completa con consultas DNS
dig +trace syncloud.local @127.0.0.1
```

### Configuracion de red

Asegurate de que el archivo `/etc/dnsmasq.d/syncloud.conf` contenga la configuracion correcta:

```conf
address=/syncloud.local/192.168.1.100
server=8.8.8.8
server=8.8.4.4
listen-address=127.0.0.1,192.168.1.100
bind-interfaces
```
