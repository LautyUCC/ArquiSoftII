# 🎓 Guía Docker para Principiantes - Spotly

## ✅ Pre-requisitos

1. **Docker Desktop instalado** 
   - Windows: https://www.docker.com/products/docker-desktop
   - Verificar: `docker --version`

2. **Archivos del proyecto descargados y descomprimidos**

---

## 📝 Paso a Paso: Primera vez con Docker

### Paso 1: Abrir Docker Desktop

1. Abre **Docker Desktop** desde el menú de inicio
2. Espera a que inicie (verás la ballena en la barra de tareas)
3. La ballena debe estar estable (no intermitente)

### Paso 2: Abrir PowerShell en el proyecto

1. Abre el **Explorador de Archivos**
2. Ve a tu carpeta del proyecto:
   ```
   C:\Users\virrr\Desktop\Spotly\ProyectoArqui2---Sesin-Vallino-y-Rodriguez-
   ```
3. En la barra de direcciones, escribe `powershell` y presiona Enter
4. Se abrirá PowerShell en esa carpeta

### Paso 3: Verificar archivos

Verifica que tengas estos archivos:

```powershell
dir
```

Deberías ver:
- ✅ `docker-compose.yml`
- ✅ Carpetas: `users-api`, `properties-api`, `search-api`, `frontend`

### Paso 4: Levantar SOLO la infraestructura primero

Primero vamos a levantar solo las bases de datos y servicios, SIN los microservicios:

```powershell
docker-compose up -d mysql mongodb rabbitmq solr memcached
```

**¿Qué hace esto?**
- `-d` = modo "detached" (en segundo plano)
- Levanta solo: MySQL, MongoDB, RabbitMQ, Solr, Memcached

**Espera 1-2 minutos** para que inicien.

### Paso 5: Verificar que los servicios estén corriendo

```powershell
docker-compose ps
```

Deberías ver algo como:
```
NAME                    STATUS
spotly-mysql            Up (healthy)
spotly-mongodb          Up (healthy)
spotly-rabbitmq         Up (healthy)
spotly-solr             Up (healthy)
spotly-memcached        Up
```

✅ Si todos dicen "Up" o "Up (healthy)", perfecto!
❌ Si alguno dice "Exit", hay un error (me avisas cuál)

### Paso 6: Ver logs (opcional)

Para ver qué están haciendo:

```powershell
docker-compose logs -f mysql
```

Presiona `Ctrl+C` para salir.

### Paso 7: Levantar los microservicios

Ahora levantamos nuestras APIs:

```powershell
docker-compose up --build users-api properties-api search-api
```

**¿Qué hace esto?**
- `--build` = construye las imágenes desde los Dockerfile
- Levanta: users-api, properties-api, search-api

**Primera vez puede tardar 5-10 minutos** (descarga Go, compila, etc.)

Verás mucho texto. Al final deberías ver:
```
🚀 Users API starting on port 8080...
🏠 Properties API starting on port 8081...
🔍 Search API starting on port 8082...
```

### Paso 8: Probar que funcionan

**Abre un NUEVO PowerShell** (deja el anterior corriendo) y ejecuta:

```powershell
# Probar users-api
curl http://localhost:8080/health

# Probar properties-api
curl http://localhost:8081/health

# Probar search-api
curl http://localhost:8082/health
```

O abre en el navegador:
- http://localhost:8080/health
- http://localhost:8081/health
- http://localhost:8082/health

Deberías ver JSON con `"status": "healthy"`

### Paso 9: Levantar el frontend

En otro PowerShell:

```powershell
docker-compose up --build frontend
```

**Puede tardar 10-15 minutos la primera vez** (descarga Node, dependencias, compila React)

### Paso 10: Ver tu aplicación

Abre el navegador en:
**http://localhost:3000**

Deberías ver:
- 🏠 Spotly
- Estado de las 3 APIs (users, properties, search)
- Si están en verde = ¡Todo funciona! 🎉

---

## 🛑 Comandos Importantes

### Detener todo
```powershell
docker-compose down
```

### Ver qué está corriendo
```powershell
docker-compose ps
```

### Ver logs de todos los servicios
```powershell
docker-compose logs -f
```

### Ver logs de un servicio específico
```powershell
docker-compose logs -f users-api
```

### Reiniciar un servicio
```powershell
docker-compose restart users-api
```

### Reconstruir un servicio (después de cambiar código)
```powershell
docker-compose up --build users-api
```

### Limpiar todo (¡CUIDADO! Borra datos)
```powershell
docker-compose down -v
```

---

## 🐛 Problemas Comunes

### Error: "port is already allocated"

Significa que ese puerto ya está en uso.

**Solución:**
```powershell
# Ver qué usa el puerto
netstat -ano | findstr :8080

# Matar el proceso (reemplaza <PID>)
taskkill /PID <PID> /F
```

### Error: "Cannot connect to Docker daemon"

Docker Desktop no está corriendo.

**Solución:**
1. Abre Docker Desktop
2. Espera a que la ballena esté estable
3. Intenta de nuevo

### Error al compilar Go: "go.mod file not found"

Falta el archivo go.mod en algún microservicio.

**Solución:**
Verifica que cada carpeta tenga su `go.mod`:
```powershell
dir users-api\go.mod
dir properties-api\go.mod
dir search-api\go.mod
```

### Frontend no construye

Error con npm o Node.

**Solución:**
1. Verifica que `package.json` exista en `frontend/`
2. Reconstruye:
   ```powershell
   docker-compose build --no-cache frontend
   docker-compose up frontend
   ```

### Servicio dice "Exited (1)"

Hay un error en el código o configuración.

**Solución:**
```powershell
# Ver el error
docker-compose logs users-api
```

---

## 📊 URLs Útiles

| Servicio | URL | Usuario/Password |
|----------|-----|------------------|
| Frontend | http://localhost:3000 | - |
| Users API | http://localhost:8080 | - |
| Properties API | http://localhost:8081 | - |
| Search API | http://localhost:8082 | - |
| RabbitMQ UI | http://localhost:15672 | spotly / spotly_password |
| Solr Admin | http://localhost:8983 | - |

---

## 🎯 Flujo Normal de Trabajo

### Primera vez (hoy):
```powershell
# 1. Levantar infraestructura
docker-compose up -d mysql mongodb rabbitmq solr memcached

# 2. Esperar 1 minuto

# 3. Levantar microservicios
docker-compose up --build users-api properties-api search-api

# 4. En otra terminal, levantar frontend
docker-compose up --build frontend

# 5. Abrir http://localhost:3000
```

### Próximas veces (mañana):
```powershell
# Todo junto (ya están construidas las imágenes)
docker-compose up
```

### Después de cambiar código:
```powershell
# Reconstruir solo lo que cambiaste
docker-compose up --build users-api
```

### Al terminar el día:
```powershell
# Detener todo (los datos persisten)
docker-compose down
```

---

## ✅ Checklist Final

- [ ] Docker Desktop instalado y corriendo
- [ ] Proyecto descargado y descomprimido
- [ ] PowerShell abierto en la carpeta del proyecto
- [ ] `docker-compose.yml` presente
- [ ] Infraestructura levantada (mysql, mongodb, etc.)
- [ ] Microservicios levantados (users, properties, search)
- [ ] Frontend levantado
- [ ] http://localhost:3000 muestra la app
- [ ] APIs responden en /health

---

## 🆘 Si algo no funciona

1. Copia el error completo
2. Ejecuta: `docker-compose ps`
3. Ejecuta: `docker-compose logs <servicio-con-error>`
4. Mándame esa información

¡Estoy para ayudarte! 🚀
