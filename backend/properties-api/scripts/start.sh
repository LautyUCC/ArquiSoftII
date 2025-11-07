#!/bin/bash

# Script para iniciar todos los servicios de la aplicación
# Inicia los contenedores con docker-compose, espera a que estén listos y muestra logs

set -e

echo "🚀 Iniciando servicios de Properties API..."
echo ""

# Navegar al directorio raíz del proyecto (donde está docker-compose.yml)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$PROJECT_ROOT"

# Construir e iniciar contenedores en modo detached
echo "📦 Construyendo e iniciando contenedores..."
docker-compose up --build -d

echo ""
echo "⏳ Esperando a que los servicios estén listos..."

# Esperar a que MongoDB esté listo
echo "   - Esperando MongoDB..."
until docker-compose exec -T mongodb mongosh --eval "db.adminCommand('ping')" > /dev/null 2>&1; do
  sleep 2
done
echo "   ✅ MongoDB está listo"

# Esperar a que RabbitMQ esté listo
echo "   - Esperando RabbitMQ..."
until docker-compose exec -T rabbitmq rabbitmqctl status > /dev/null 2>&1; do
  sleep 2
done
echo "   ✅ RabbitMQ está listo"

# Esperar a que properties-api esté listo (verificar health check)
echo "   - Esperando Properties API..."
max_attempts=30
attempt=0
while [ $attempt -lt $max_attempts ]; do
  if curl -s http://localhost:8081/health > /dev/null 2>&1; then
    break
  fi
  attempt=$((attempt + 1))
  sleep 2
done

if [ $attempt -eq $max_attempts ]; then
  echo "   ⚠️  Properties API no responde después de 60 segundos"
else
  echo "   ✅ Properties API está listo"
fi

echo ""
echo "✅ Todos los servicios están iniciados"
echo ""
echo "📋 Mostrando logs de los servicios..."
echo "   (Presiona Ctrl+C para salir de los logs)"
echo ""

# Mostrar logs de todos los servicios
docker-compose logs -f

