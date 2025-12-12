#!/bin/bash

# Script para iniciar todos los servicios de Spotly
# Uso: ./start.sh

set -e

echo "🚀 Iniciando Spotly Microservices..."
echo ""

# Verificar que Docker esté corriendo
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker no está corriendo. Por favor inicia Docker Desktop."
    exit 1
fi

echo "✅ Docker está corriendo"
echo ""

# Verificar si existe archivo .env
if [ ! -f .env ]; then
    echo "⚠️  No se encontró archivo .env"
    echo "📝 Creando .env con valores por defecto..."
    echo ""
    cat > .env << EOF
# Spotly Microservices - Variables de Entorno
# Valores por defecto para desarrollo local

# MySQL
MYSQL_ROOT_PASSWORD=rootpassword
MYSQL_DATABASE=spotly_users
MYSQL_USER=spotly
MYSQL_PASSWORD=spotlypass
MYSQL_HOST=spotly-mysql
MYSQL_PORT=3306

# MongoDB
MONGO_URI=mongodb://mongodb:27017/spotly_properties
MONGO_DATABASE=spotly_properties

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest

# Solr
SOLR_URL=http://solr:8983/solr/properties

# Memcached
MEMCACHED_ADDR=memcached:11211

# APIs
USERS_API_URL=http://users-api:8081
PROPERTIES_API_URL=http://spotly-properties-api:8081
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-min-32-chars

# Nginx
NGINX_PORT=80
EOF
    echo "✅ Archivo .env creado"
    echo ""
else
    echo "✅ Archivo .env encontrado"
    echo ""
fi

# Construir e iniciar servicios
echo "🔨 Construyendo e iniciando servicios..."
echo ""

docker-compose up --build -d

echo ""
echo "⏳ Esperando a que los servicios estén listos..."
sleep 10

# Verificar estado de los servicios
echo ""
echo "📊 Estado de los servicios:"
docker-compose ps

echo ""
echo "✅ Servicios iniciados!"
echo ""
echo "🌐 URLs disponibles:"
echo "   - Frontend:        http://localhost"
echo "   - Users API:       http://localhost/api/users"
echo "   - Properties API:  http://localhost/api/properties"
echo "   - Search API:      http://localhost/api/search"
echo ""
echo "📝 Para ver los logs:"
echo "   docker-compose logs -f [servicio]"
echo ""
echo "🛑 Para detener los servicios:"
echo "   docker-compose down"
echo ""

