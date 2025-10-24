# Spotly Microservices - Estructura del Proyecto

## 📁 Estructura General

```
spotly-microservices/
│
├── docker-compose.yml              # Orquestación de todos los servicios
├── .env                            # Variables de entorno
│
├── users-api/                      # 🔐 Microservicio de Usuarios
│   ├── main.go                     # Punto de entrada
│   ├── go.mod                      # Dependencias
│   ├── Dockerfile                  # Imagen Docker
│   │
│   ├── controllers/                # Capa de controladores (HTTP handlers)
│   │   └── user_controller.go
│   │
│   ├── services/                   # Capa de lógica de negocio
│   │   └── user_service.go
│   │
│   ├── domain/                     # Modelos de dominio
│   │   └── user.go
│   │
│   ├── repositories/               # Capa de acceso a datos
│   │   └── user_repository.go
│   │
│   ├── dto/                        # Data Transfer Objects
│   │   └── user_dto.go
│   │
│   ├── utils/                      # Utilidades (hashing, etc)
│   │   ├── crypto.go
│   │   └── jwt.go
│   │
│   └── middleware/                 # Middlewares (auth, etc)
│       └── auth_middleware.go
│
├── properties-api/                 # 🏠 Microservicio de Propiedades
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   │
│   ├── controllers/                # Controladores HTTP
│   │   ├── property_controller.go
│   │   └── booking_controller.go
│   │
│   ├── services/                   # Lógica de negocio
│   │   ├── property_service.go
│   │   └── booking_service.go
│   │
│   ├── domain/                     # Modelos
│   │   ├── property.go
│   │   └── booking.go
│   │
│   ├── repositories/               # Acceso a MongoDB
│   │   ├── property_repository.go
│   │   └── booking_repository.go
│   │
│   ├── dto/                        # DTOs
│   │   ├── property_dto.go
│   │   └── booking_dto.go
│   │
│   ├── clients/                    # Clientes HTTP (para llamar a users-api)
│   │   └── user_client.go
│   │
│   └── queue/                      # Productor de RabbitMQ
│       └── rabbitmq_producer.go
│
├── search-api/                     # 🔍 Microservicio de Búsqueda
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile
│   │
│   ├── controllers/                # Controladores HTTP
│   │   └── search_controller.go
│   │
│   ├── services/                   # Lógica de búsqueda
│   │   └── search_service.go
│   │
│   ├── domain/                     # Modelos
│   │   └── property_search.go
│   │
│   ├── repositories/               # Acceso a Solr
│   │   └── solr_repository.go
│   │
│   ├── clients/                    # Cliente para properties-api
│   │   └── property_client.go
│   │
│   ├── queue/                      # Consumidor de RabbitMQ
│   │   └── rabbitmq_consumer.go
│   │
│   └── cache/                      # Capas de caché
│       ├── local_cache.go          # CCache
│       └── distributed_cache.go    # Memcached
│
└── frontend/                       # ⚛️ Aplicación React
    ├── package.json
    ├── Dockerfile
    │
    ├── public/                     # Archivos estáticos
    │   └── index.html
    │
    └── src/                        # Código fuente
        ├── App.js                  # Componente principal
        ├── index.js                # Punto de entrada
        │
        ├── pages/                  # Páginas/Vistas
        │   ├── Login.js
        │   ├── Register.js
        │   ├── Home.js
        │   ├── PropertyDetails.js
        │   ├── Congrats.js
        │   ├── MyBookings.js
        │   └── Admin.js
        │
        ├── components/             # Componentes reutilizables
        │   ├── Navbar.js
        │   ├── PropertyCard.js
        │   └── SearchBar.js
        │
        ├── services/               # Llamadas a API
        │   ├── authService.js
        │   ├── propertyService.js
        │   └── searchService.js
        │
        └── utils/                  # Utilidades
            └── auth.js
```

## 🔧 Tecnologías por Servicio

### users-api
- **Lenguaje:** Go
- **Base de datos:** MySQL
- **ORM:** GORM
- **Autenticación:** JWT
- **Hashing:** bcrypt

### properties-api
- **Lenguaje:** Go
- **Base de datos:** MongoDB
- **Driver:** mongo-go-driver
- **Mensajería:** RabbitMQ (producer)
- **Concurrencia:** Goroutines + Channels + WaitGroups

### search-api
- **Lenguaje:** Go
- **Motor de búsqueda:** Apache Solr
- **Mensajería:** RabbitMQ (consumer)
- **Caché local:** CCache
- **Caché distribuida:** Memcached

### frontend
- **Framework:** React
- **Comunicación:** HTTP/REST (fetch/axios)
- **Formato:** JSON

## 🐳 Servicios Docker

El `docker-compose.yml` orquestará:
1. **users-api** (puerto 8080)
2. **properties-api** (puerto 8081)
3. **search-api** (puerto 8082)
4. **frontend** (puerto 3000)
5. **MySQL** (puerto 3306)
6. **MongoDB** (puerto 27017)
7. **RabbitMQ** (puerto 5672, management: 15672)
8. **Solr** (puerto 8983)
9. **Memcached** (puerto 11211)

## 📝 Próximos Pasos

1. ✅ Estructura creada
2. ⏭️ Configurar docker-compose.yml
3. ⏭️ Implementar users-api
4. ⏭️ Implementar properties-api
5. ⏭️ Implementar search-api
6. ⏭️ Implementar frontend
7. ⏭️ Integración completa
