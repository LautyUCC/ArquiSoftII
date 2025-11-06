# Properties API

API de microservicio para gestión de propiedades inmobiliarias desarrollada en Go utilizando el framework Gin.

## Estructura del Proyecto

```
properties-api/
├── main.go                    # Punto de entrada de la aplicación
├── go.mod                     # Definición del módulo y dependencias
├── controllers/               # Controladores HTTP (capa de presentación)
├── services/                  # Lógica de negocio
├── repositories/              # Acceso a datos (capa de persistencia)
├── domain/                    # Modelos de dominio (entidades del negocio)
├── dto/                       # Data Transfer Objects (DTOs para request/response)
├── clients/                    # Clientes HTTP para comunicación con otros servicios
├── config/                    # Configuración de la aplicación
└── utils/                     # Utilidades y helpers
```

## Descripción de Directorios

### `controllers/`
Contiene los controladores HTTP que manejan las peticiones entrantes y las respuestas. Los controladores se encargan de:
- Validar los datos de entrada
- Llamar a los servicios correspondientes
- Formatear las respuestas
- Manejar errores HTTP

### `services/`
Contiene la lógica de negocio de la aplicación. Los servicios:
- Implementan las reglas de negocio
- Coordinan las operaciones entre repositorios y clientes externos
- Manejan la publicación de eventos
- Validan datos de negocio

### `repositories/`
Contiene la capa de acceso a datos. Los repositorios:
- Interactúan directamente con la base de datos (MongoDB)
- Implementan operaciones CRUD
- Manejan consultas y filtros
- Abstraen la lógica de persistencia

### `domain/`
Contiene los modelos de dominio que representan las entidades del negocio. Estos modelos:
- Definen la estructura de datos principal
- Contienen las validaciones a nivel de dominio
- Son independientes de la capa de persistencia

### `dto/`
Contiene los Data Transfer Objects (DTOs) que se usan para:
- Serialización/deserialización de datos en HTTP
- Validación de entrada
- Transformación entre capas (domain ↔ HTTP)
- Separar la estructura interna de la API externa

### `clients/`
Contiene los clientes HTTP para comunicación con otros microservicios:
- `user_client.go`: Cliente para comunicarse con users-api
- Maneja la comunicación HTTP entre servicios
- Implementa retry logic y manejo de errores

### `config/`
Contiene la configuración de la aplicación:
- Variables de entorno
- Configuración de conexiones (MongoDB, RabbitMQ)
- Inicialización de servicios externos

### `utils/`
Contiene funciones de utilidad reutilizables:
- Helpers para respuestas HTTP
- Validaciones comunes
- Utilidades de paginación
- Funciones auxiliares

## Arquitectura

El proyecto sigue una arquitectura en capas (Layered Architecture) con separación clara de responsabilidades:

```
HTTP Request
    ↓
Controllers (validación HTTP, formateo)
    ↓
Services (lógica de negocio)
    ↓
Repositories (persistencia) / Clients (servicios externos)
    ↓
Database / External APIs
```

## Flujo de Datos

1. **Request HTTP** → Llega al `Controller`
2. **Controller** → Valida y convierte DTO → llama a `Service`
3. **Service** → Ejecuta lógica de negocio → usa `Repository` o `Client`
4. **Repository/Client** → Interactúa con DB/API externa → retorna datos
5. **Service** → Transforma datos → retorna al `Controller`
6. **Controller** → Convierte a DTO → retorna respuesta HTTP

## Tecnologías

- **Go 1.21+**: Lenguaje de programación
- **Gin**: Framework web
- **MongoDB**: Base de datos NoSQL
- **RabbitMQ**: Sistema de mensajería
- **Docker**: Contenedorización

## Variables de Entorno

```env
SERVER_PORT=8082
ENVIRONMENT=development
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=properties_db
RABBITMQ_URI=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=properties_exchange
USERS_API_URL=http://users-api:8081
```

## Principios de Diseño

- **Separation of Concerns**: Cada capa tiene una responsabilidad única
- **Dependency Injection**: Las dependencias se inyectan en lugar de crearse internamente
- **Interface-based Design**: Uso de interfaces para desacoplar componentes
- **Error Handling**: Manejo consistente de errores en todas las capas
- **Clean Code**: Código legible, mantenible y siguiendo convenciones de Go

## Próximos Pasos

1. Implementar la lógica en `main.go`
2. Configurar las rutas y middleware
3. Implementar los controladores
4. Desarrollar los servicios de negocio
5. Configurar los repositorios
6. Integrar con servicios externos

