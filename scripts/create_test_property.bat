@echo off
REM Script para crear una propiedad de prueba
REM Requiere que el usuario admin exista y que los servicios estén corriendo

setlocal enabledelayedexpansion

echo ========================================
echo Creando Propiedad de Prueba
echo ========================================
echo.

REM Esperar un momento para asegurar que la API esté lista
timeout /t 2 /nobreak >nul

REM Verificar que los servicios estén corriendo
echo [1/4] Verificando que la API esté disponible...
curl -s -f http://localhost/api/users/login >nul 2>&1
if errorlevel 1 (
    echo ERROR: La API no está disponible en http://localhost/api
    echo Asegúrate de que los servicios estén ejecutándose con docker-compose up
    pause
    exit /b 1
)
echo OK: API disponible
echo.

REM Hacer login como admin
echo [2/4] Haciendo login como admin...
set LOGIN_RESPONSE=login_response.json
curl -s -X POST http://localhost/api/users/login ^
    -H "Content-Type: application/json" ^
    -d "{\"usernameOrEmail\":\"admin\",\"password\":\"admin123\"}" ^
    -o %LOGIN_RESPONSE%

REM Extraer el token del JSON usando PowerShell (más robusto)
for /f "delims=" %%i in ('powershell -Command "$json = Get-Content '%LOGIN_RESPONSE%' | ConvertFrom-Json; if ($json.token) { Write-Output $json.token } else { Write-Output '' }"') do set TOKEN=%%i

if "!TOKEN!"=="" (
    echo ERROR: No se pudo obtener el token. Verificando respuesta:
    type %LOGIN_RESPONSE%
    echo.
    echo Asegúrate de que el usuario admin exista y tenga las credenciales correctas.
    echo Ejecuta scripts\create_test_users.bat si aún no has creado los usuarios.
    del %LOGIN_RESPONSE% 2>nul
    pause
    exit /b 1
)

echo OK: Token obtenido
echo.

REM Crear propiedad de prueba
echo [3/4] Creando propiedad de prueba...
set PROPERTY_RESPONSE=property_response.json
curl -s -X POST http://localhost/api/properties ^
    -H "Content-Type: application/json" ^
    -H "Authorization: Bearer !TOKEN!" ^
    -d "{\"title\":\"Apartamento Moderno en el Centro\",\"description\":\"Hermoso apartamento completamente amueblado en el corazón de la ciudad. Cerca de restaurantes, transporte público y centros comerciales. Ideal para turistas o profesionales.\",\"price\":120000.00,\"location\":\"Bogotá, Colombia\",\"ownerId\":\"1\",\"amenities\":[\"wifi\",\"pool\",\"parking\",\"kitchen\",\"air-conditioning\"],\"capacity\":4,\"available\":true,\"images\":[]}" ^
    -o %PROPERTY_RESPONSE%

REM Verificar si se creó correctamente
findstr /C:"id" %PROPERTY_RESPONSE% >nul 2>&1
if errorlevel 1 (
    echo ERROR: No se pudo crear la propiedad. Respuesta:
    type %PROPERTY_RESPONSE%
    echo.
    echo Posibles causas:
    echo - El usuario no tiene rol ADMIN
    echo - El ownerId no es válido
    echo - Error de validación en los datos
    del %LOGIN_RESPONSE% %PROPERTY_RESPONSE% 2>nul
    pause
    exit /b 1
)

echo OK: Propiedad creada exitosamente
echo.

REM Mostrar respuesta
echo [4/4] Respuesta del servidor:
type %PROPERTY_RESPONSE%
echo.

REM Limpiar archivos temporales
del %LOGIN_RESPONSE% %PROPERTY_RESPONSE% 2>nul

echo ========================================
echo Propiedad de prueba creada exitosamente
echo ========================================
echo.
echo Puedes verificar la propiedad en el frontend:
echo http://localhost
echo.

pause
endlocal

