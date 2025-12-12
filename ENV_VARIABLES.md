# Variables de Entorno - Spotly Microservices

Este documento describe todas las variables de entorno necesarias para el proyecto.

## 📋 Crear archivo .env

Crea un archivo `.env` en la raíz del proyecto con las siguientes variables:

```bash
# ============================================
# MySQL Configuration
# ============================================
MYSQL_ROOT_PASSWORD=rootpassword
MYSQL_DATABASE=spotly_users
MYSQL_USER=spotly
MYSQL_PASSWORD=spotlypass
MYSQL_HOST=spotly-mysql
MYSQL_PORT=3306

# ============================================
# MongoDB Configuration
# ============================================
MONGO_HOST=mongodb
MONGO_PORT=27017
MONGO_DATABASE=spotly_properties
MONGO_URI=mongodb://mongodb:27017/spotly_properties

# ============================================
# RabbitMQ Configuration
# ============================================
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
RABBITMQ_QUEUE=property_events

# ============================================
# Solr Configuration
# ============================================
SOLR_HOST=solr
SOLR_PORT=8983
SOLR_CORE=properties
SOLR_URL=http://solr:8983/solr/properties

# ============================================
# Memcached Configuration
# ============================================
MEMCACHED_HOST=memcached
MEMCACHED_PORT=11211
MEMCACHED_ADDR=memcached:11211

# ============================================
# Users API Configuration
# ============================================
USERS_API_HOST=users-api
USERS_API_PORT=8081
USERS_API_URL=http://users-api:8081
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-min-32-chars

# ============================================
# Properties API Configuration
# ============================================
PROPERTIES_API_HOST=spotly-properties-api
PROPERTIES_API_PORT=8081
PROPERTIES_API_EXTERNAL_PORT=8082
PROPERTIES_API_URL=http://spotly-properties-api:8081

# ============================================
# Search API Configuration
# ============================================
SEARCH_API_HOST=spotly-search-api
SEARCH_API_PORT=8083
SEARCH_API_URL=http://spotly-search-api:8083

# ============================================
# Nginx Configuration
# ============================================
NGINX_PORT=80
```

## 🔐 Valores Recomendados para Entorno Local

Los valores mostrados arriba son razonables para desarrollo local. 

### ⚠️ Importante para Producción

1. **JWT_SECRET**: Debe ser una cadena aleatoria de al menos 32 caracteres. Genera uno con:
   ```bash
   openssl rand -base64 32
   ```

2. **MYSQL_ROOT_PASSWORD**: Usa una contraseña fuerte en producción

3. **MYSQL_PASSWORD**: Usa una contraseña fuerte diferente a la root

4. **RABBITMQ_USER y RABBITMQ_PASSWORD**: Cambia de `guest/guest` en producción

## 📝 Uso

1. Copia el contenido de este archivo a un nuevo archivo `.env` en la raíz del proyecto
2. Ajusta los valores según tu entorno
3. El `docker-compose.yml` leerá automáticamente estas variables

## 🔄 Variables con Valores por Defecto

Todas las variables tienen valores por defecto en `docker-compose.yml` usando la sintaxis `${VARIABLE:-default}`, por lo que el proyecto funcionará incluso sin el archivo `.env`, pero se recomienda crearlo para mayor control.

