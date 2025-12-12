@echo off
REM Script para reconstruir el frontend después de cambios
echo ========================================
echo Reconstruyendo Frontend
echo ========================================
echo.

echo [1/3] Deteniendo frontend...
docker-compose stop frontend
if errorlevel 1 (
    echo ERROR: No se pudo detener el frontend
    pause
    exit /b 1
)
echo OK: Frontend detenido
echo.

echo [2/3] Reconstruyendo frontend (esto puede tardar 1-2 minutos)...
docker-compose build --no-cache frontend
if errorlevel 1 (
    echo ERROR: No se pudo reconstruir el frontend
    pause
    exit /b 1
)
echo OK: Frontend reconstruido
echo.

echo [3/3] Iniciando frontend y nginx...
docker-compose up -d frontend nginx
if errorlevel 1 (
    echo ERROR: No se pudieron iniciar los servicios
    pause
    exit /b 1
)
echo OK: Servicios iniciados
echo.

echo ========================================
echo Frontend reconstruido exitosamente!
echo ========================================
echo.
echo IMPORTANTE: Limpia la caché del navegador:
echo   - Presiona Ctrl + Shift + R
echo   - O abre en modo incógnito
echo.
echo URL: http://localhost
echo.

pause

