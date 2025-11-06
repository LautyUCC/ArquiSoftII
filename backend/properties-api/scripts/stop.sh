#!/bin/bash

# Script para detener todos los servicios de la aplicación
# Detiene y elimina los contenedores de docker-compose

set -e

echo "🛑 Deteniendo servicios de Properties API..."
echo ""

# Navegar al directorio raíz del proyecto (donde está docker-compose.yml)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$PROJECT_ROOT"

# Detener y eliminar contenedores
echo "📦 Deteniendo contenedores..."
docker-compose down

echo ""
echo "✅ Servicios detenidos correctamente"

