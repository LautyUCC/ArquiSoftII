@echo off
REM Script para iniciar todos los servicios de Spotly en Windows
REM Uso: start.bat

echo 🚀 Iniciando Spotly Microservices...
echo.

REM Verificar que Docker esté corriendo
docker info >nul 2>&1
if errorlevel 1 (
    echo ❌ Error: Docker no está corriendo. Por favor inicia Docker Desktop.
    exit /b 1
)

echo ✅ Docker está corriendo
echo.

REM Verificar si existe archivo .env
if not exist .env (
    echo ⚠️  No se encontró archivo .env
    echo 📝 Creando .env con valores por defecto...
    echo.
    (
        echo # Spotly Microservices - Variables de Entorno
        echo # Valores por defecto para desarrollo local
        echo.
        echo # MySQL
        echo MYSQL_ROOT_PASSWORD=rootpassword
        echo MYSQL_DATABASE=spotly_users
        echo MYSQL_USER=spotly
        echo MYSQL_PASSWORD=spotlypass
        echo MYSQL_HOST=spotly-mysql
        echo MYSQL_PORT=3306
        echo.
        echo # MongoDB
        echo MONGO_URI=mongodb://mongodb:27017/spotly_properties
        echo MONGO_DATABASE=spotly_properties
        echo.
        echo # RabbitMQ
        echo RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
        echo RABBITMQ_USER=guest
        echo RABBITMQ_PASSWORD=guest
        echo.
        echo # Solr
        echo SOLR_URL=http://solr:8983/solr/properties
        echo.
        echo # Memcached
        echo MEMCACHED_ADDR=memcached:11211
        echo.
        echo # APIs
        echo USERS_API_URL=http://users-api:8081
        echo PROPERTIES_API_URL=http://spotly-properties-api:8081
        echo JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-min-32-chars
        echo.
        echo # Nginx
        echo NGINX_PORT=80
    ) > .env
    echo ✅ Archivo .env creado
    echo.
) else (
    echo ✅ Archivo .env encontrado
    echo.
)

REM Construir e iniciar servicios
echo 🔨 Construyendo e iniciando servicios...
echo.

docker-compose up --build -d

echo.
echo ⏳ Esperando a que los servicios estén listos...
timeout /t 10 /nobreak >nul

REM Verificar estado de los servicios
echo.
echo 📊 Estado de los servicios:
docker-compose ps

echo.
echo ✅ Servicios iniciados!
echo.
echo 🌐 URLs disponibles:
echo    - Frontend:        http://localhost
echo    - Users API:       http://localhost/api/users
echo    - Properties API:  http://localhost/api/properties
echo    - Search API:      http://localhost/api/search
echo.
echo 📝 Para ver los logs:
echo    docker-compose logs -f [servicio]
echo.
echo 🛑 Para detener los servicios:
echo    docker-compose down
echo.

pause

