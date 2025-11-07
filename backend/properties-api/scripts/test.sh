#!/bin/bash

# Script para ejecutar todos los tests unitarios de la aplicación
# Ejecuta tests con cobertura y modo verbose

set -e

echo "🧪 Ejecutando tests de Properties API..."
echo ""

# Navegar al directorio de la API
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$API_DIR"

# Ejecutar tests con cobertura y modo verbose
echo "📊 Ejecutando tests con cobertura..."
go test ./... -v -cover

echo ""
echo "✅ Tests completados"

