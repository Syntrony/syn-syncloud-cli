# Comando: install

## Descripción

El comando `install` es responsable de instalar y configurar los runtime y servicios necesarios para que Syncloud Platform funcione correctamente. Este comando detecta el entorno de ejecución, instala las dependencias requeridas y configura el sistema de DNS.

## Servicios Instalados

### 1. Runtime Kubernetes (K3s)
- Instalación de K3s como cluster Kubernetes ligero
- Configuración de red del cluster
- Habilitación de componentes necesarios

### 2. Runtime Docker
- Verificación de instalación existente de Docker
- Configuración del daemon de Docker
- Preparación para gestión de contenedores

### 3. DNS (dnsmasq)
- Configuración de servidor DNS local
- Resolución de nombres internos (syncloud.local)
- Configuración de registros DNS para servicios

## Flujo de Ejecución

```
┌─────────────────────────────────────────────────────────────────────┐
│                         synctl install                              │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  1. Inicialización de Componentes                                  │
│     - Logger: Inicializar sistema de logs                          │
│     - Runner: Configurar ejecutor de comandos                      │
│     - Detector: Detectar entorno (Container vs VPS)              │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  2. Detección de Privilegios                                       │
│     - Verificar si se ejecuta como root                            │
│     - Configurar ejecutor con privilegios (sudo) si es necesario  │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  3. Detección de Runtimes Existentes                               │
│     ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│     │   Docker     │  │  Kubernetes  │  │     DNS     │          │
│     │   Detector   │  │   Detector   │  │   Detector  │          │
│     └──────────────┘  └──────────────┘  └──────────────┘          │
│           │                 │                 │                    │
│           ▼                 ▼                 ▼                    │
│     ┌──────────────────────────────────────────────┐              │
│     │  Verificar estado actual de cada servicio    │              │
│     └──────────────────────────────────────────────┘              │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  4. Instalación de Runtimes                                        │
│     - Kubernetes: Instalar K3s si no existe                       │
│     - Docker: Instalar/configurar Docker si no existe              │
│     - DNS: Instalar y configurar dnsmasq                          │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  5. Construcción del Estado Inicial                                │
│     - StateBuilder: Crear estado inicial                          │
│     - Persistir en syncloud-state.json                             │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
┌─────────────────────────────────────────────────────────────────────┐
│  Final: Sistema instalado y configurado                            │
│  - Kubernetes operativo                                            │
│  - Docker operativo                                                │
│  - DNS configurado                                                  │
└─────────────────────────────────────────────────────────────────────┘
```

## Dependencias entre Servicios

```
install (orquestador principal)
    │
    ├── k8sRuntime (K3s)
    │     ├── Detector: ¿K3s instalado?
    │     ├── Inspector: ¿K3s operativo?
    │     └── Installer: Instalar K3s si es necesario
    │
    ├── dockerRuntime (Docker)
    │     ├── Detector: ¿Docker instalado?
    │     ├── Inspector: ¿Docker operativo?
    │     └── Installer: Instalar Docker si es necesario
    │
    └── dnsRuntime (dnsmasq)
          ├── Detector: ¿dnsmasq instalado?
          ├── Inspector: ¿dnsmasq operativo?
          └── Installer: Instalar dnsmasq si es necesario
```

## Modos de Entorno

### Modo Contenedor (DEV)
- Detecta entorno container (Docker-in-Docker o K3d)
- Utiliza K3d para desarrollo
- Configuración de red especial para contenedores

### Modo VPS (PROD)
- Instala K3s directamente en el host
- Configuración de red estándar
- Requiere privilegios de root

## Consideraciones

- **Idempotencia**: El comando verifica el estado existente antes de instalar. Si un servicio ya está instalado, lo跳过.
- **Privilegios**: Dependiendo del entorno, puede requerir ejecución con sudo.
- **Estado**: El resultado de la instalación se almacena en el archivo de estado para referencia futura.
