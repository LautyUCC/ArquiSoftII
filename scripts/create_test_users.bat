@echo off
REM Script para crear usuarios de prueba usando la API (Windows)
REM Este script crea los usuarios y luego actualiza uno a admin usando SQL directo

echo 🔧 Creando usuarios de prueba...
echo.

REM Esperar un momento para asegurar que la API esté lista
timeout /t 2 /nobreak >nul

REM URL base de la API (a través de nginx)
set API_URL=http://localhost/api/users

REM 1. Crear usuario ADMINISTRADOR
echo.
echo 📝 Creando usuario administrador...
curl -X POST "%API_URL%" ^
  -H "Content-Type: application/json" ^
  -d "{\"username\": \"admin\", \"email\": \"admin@spotly.com\", \"password\": \"admin123\", \"firstName\": \"Admin\", \"lastName\": \"Usuario\"}"

echo.
echo.

REM Esperar un momento
timeout /t 1 /nobreak >nul

REM 2. Crear usuario NORMAL
echo 📝 Creando usuario normal...
curl -X POST "%API_URL%" ^
  -H "Content-Type: application/json" ^
  -d "{\"username\": \"usuario\", \"email\": \"usuario@spotly.com\", \"password\": \"usuario123\", \"firstName\": \"Usuario\", \"lastName\": \"Normal\"}"

echo.
echo.

REM Esperar un momento
timeout /t 1 /nobreak >nul

REM 3. Actualizar el usuario admin a tipo "admin" usando SQL directo
echo 🔧 Actualizando usuario admin a tipo 'admin'...
docker-compose exec -T mysql mysql -u spotly -pspotlypass spotly_users -e "UPDATE users SET user_type = 'admin' WHERE username = 'admin'; SELECT id, username, email, user_type FROM users WHERE username IN ('admin', 'usuario');"

echo.
echo ✅ Usuarios creados exitosamente!
echo.
echo 📋 Credenciales:
echo    ADMIN:
echo      Username: admin
echo      Password: admin123
echo      Email: admin@spotly.com
echo.
echo    USUARIO NORMAL:
echo      Username: usuario
echo      Password: usuario123
echo      Email: usuario@spotly.com
echo.

pause

