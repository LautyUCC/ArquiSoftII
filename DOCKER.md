# 🐳 Guía Docker - Spotly Microservices

## 📋 Servicios Incluidos

El `docker-compose.yml` levanta **9 servicios**:

### Microservicios (Go)
1. **users-api** (8080) - Gestión de usuarios
2. **properties-api** (8081) - Propiedades y reservas
3. **search-api** (8082) - Motor de búsqueda

### Frontend
4. **frontend** (3000) - Aplicación React

### Infraestructura
5. **mysql** (3306) - Base de datos para usuarios
6. **mongodb** (27017) - Base de datos para propiedades
7. **rabbitmq** (5672, 15672) - Cola de mensajes
8. **solr** (8983) - Motor de búsqueda
9. **memcached** (11211) - Caché distribuida

---

## 🚀 Comandos Básicos

### Levantar todos los servicios
```bash
docker-compose up --build
```

### Levantar en segundo plano (detached)
```bash
docker-compose up -d --build
```

### Ver logs de todos los servicios
```bash
docker-compose logs -f
```

### Ver logs de un servicio específico
```bash
docker-compose logs -f users-api
docker-compose logs -f properties-api
docker-compose logs -f search-api
```

### Ver estado de los servicios
```bash
docker-compose ps
```

### Detener todos los servicios
```bash
docker-compose down
```

### Detener y eliminar volúmenes (borra los datos)
```bash
docker-compose down -v
```

### Reconstruir un servicio específico
```bash
docker-compose up --build users-api
```

---

## 🔍 Acceso a los Servicios

| Servicio | URL | Descripción |
|----------|-----|-------------|
| Frontend | http://localhost:3000 | Aplicación web |
| users-api | http://localhost:8080 | API REST usuarios |
| properties-api | http://localhost:8081 | API REST propiedades |
| search-api | http://localhost:8082 | API REST búsqueda |
| RabbitMQ UI | http://localhost:15672 | Panel admin (user: spotly, pass: spotly_password) |
| Solr Admin | http://localhost:8983 | Panel admin Solr |
| MySQL | localhost:3306 | Base de datos (user: spotly_user, pass: spotly_password) |
| MongoDB | localhost:27017 | Base de datos (user: admin, pass: adminpassword) |

---

## ⚙️ Configuración

### Variables de Entorno

1. Copia el archivo de ejemplo:
```bash
cp .env.example .env
```

2. Edita `.env` con tus configuraciones (opcional para desarrollo)

### Cambiar Puertos

Si algún puerto está ocupado, edita `docker-compose.yml`:

```yaml
services:
  users-api:
    ports:
      - "8080:8080"  # Cambiar el primer número: "PUERTO_HOST:PUERTO_CONTAINER"
```

---

## 🔧 Troubleshooting

### Puerto ya en uso
```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Mac/Linux
lsof -i :8080
kill -9 <PID>
```

### Reconstruir desde cero
```bash
# Detener y limpiar todo
docker-compose down -v
docker system prune -a

# Levantar de nuevo
docker-compose up --build
```

### Servicio no inicia
```bash
# Ver logs del servicio con error
docker-compose logs users-api

# Revisar health check
docker-compose ps
```

### Base de datos no conecta
```bash
# Verificar que MySQL/MongoDB estén healthy
docker-compose ps

# Restart del servicio
docker-compose restart mysql
docker-compose restart mongodb
```

---

## 📝 Orden de Inicio

Docker Compose inicia los servicios en este orden (gracias a `depends_on`):

1. **mysql**, **mongodb**, **rabbitmq**, **solr**, **memcached**
2. **users-api** (espera a mysql)
3. **properties-api** (espera a mongodb, rabbitmq, users-api)
4. **search-api** (espera a solr, rabbitmq, properties-api)
5. **frontend** (espera a las 3 APIs)

---

## 🧪 Testing

### Verificar que users-api funciona
```bash
curl http://localhost:8080/health
```

### Verificar que properties-api funciona
```bash
curl http://localhost:8081/health
```

### Verificar que search-api funciona
```bash
curl http://localhost:8082/health
```

### RabbitMQ - Ver colas
Ve a: http://localhost:15672
- User: `spotly`
- Pass: `spotly_password`

### Solr - Ver índice
Ve a: http://localhost:8983

---

## 💾 Datos Persistentes

Los siguientes datos persisten entre reinicios:
- MySQL: `mysql_data`
- MongoDB: `mongodb_data`
- RabbitMQ: `rabbitmq_data`
- Solr: `solr_data`

Para eliminar los datos:
```bash
docker-compose down -v
```

---

## 🐛 Debugging

### Entrar a un contenedor
```bash
# MySQL
docker exec -it spotly-mysql mysql -u spotly_user -pspotly_password users_db

# MongoDB
docker exec -it spotly-mongodb mongosh -u admin -p adminpassword

# RabbitMQ
docker exec -it spotly-rabbitmq rabbitmqctl list_queues

# Ver contenido de un contenedor
docker exec -it spotly-users-api sh
```

### Ver uso de recursos
```bash
docker stats
```

---

## 📦 Producción

**NO uses este docker-compose en producción directamente.**

Para producción:
- Cambia todas las contraseñas
- Usa secrets en lugar de variables de entorno
- Configura SSL/TLS
- Usa volúmenes externos
- Configura backups automáticos
- Añade rate limiting
- Configura monitoring (Prometheus, Grafana)
