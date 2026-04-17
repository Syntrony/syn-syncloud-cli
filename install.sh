#!/bin/bash
set -e

# Configuración del binario
BINARY_NAME="synctl"
# Estas son constantes de distribución, no secretos.
GITHUB_OWNER="Syntrony" 
GITHUB_REPO="syn-syncloud-cli"

# --- Lógica de detección ---
OS=$(uname -s)
ARCH=$(uname -m)

# Normalizar arquitecturas (solo para el formato del archivo)
case $ARCH in
    x86_64) ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Arquitectura no soportada: $ARCH"; exit 1 ;;
esac

# Obtener la URL de la última versión usando la API de GitHub
# Esto evita que tengas que actualizar el script cada vez que saques una versión
LATEST_RELEASE_URL="https://api.github.com/repos/${GITHUB_OWNER}/${GITHUB_REPO}/releases/latest"
DOWNLOAD_URL=$(curl -s $LATEST_RELEASE_URL | grep "browser_download_url" | grep "${OS}_${ARCH}\.tar\.gz" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "No se encontró un binario para ${OS} y ${ARCH}."
    exit 1
fi

# --- Instalación ---
echo "Instalando ${BINARY_NAME} desde ${DOWNLOAD_URL}..."
curl -L "$DOWNLOAD_URL" -o "${BINARY_NAME}.tar.gz"
tar -xzf "${BINARY_NAME}.tar.gz"
chmod +x "${BINARY_NAME}"

# Intentar mover a bin, requiere sudo
sudo mv "${BINARY_NAME}" /usr/local/bin/

echo "✅ ${BINARY_NAME} se ha instalado en /usr/local/bin/${BINARY_NAME}"