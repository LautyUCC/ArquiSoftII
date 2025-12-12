# 🏠 Spotly - Sistema de Reservas de Alojamientos

Sistema de reservas de alojamientos con arquitectura de microservicios.  
**Proyecto Final - Arquitectura de Software II - UCC**

---

## 📋 Tabla de Contenidos

- [Arquitectura del Sistema](#-arquitectura-del-sistema)
- [Stack Tecnológico](#-stack-tecnológico)
- [Instalación y Despliegue](#-instalación-y-despliegue)
- [Funcionalidades Implementadas](#-funcionalidades-implementadas)
- [Endpoints de la API](#-endpoints-de-la-api)
- [Estructura del Proyecto](#-estructura-del-proyecto)
- [Herramientas y Tecnologías por Tarea](#-herramientas-y-tecnologías-por-tarea)

---

## 🏗️ Arquitectura del Sistema

### Backend (Microservicios en Go)

#### 1. **users-api** (Puerto 8080)
- **Base de datos:** MySQL con GORM
- **Funcionalidad:** Gestión de usuarios, autenticación JWT, roles y permisos
- **Características:**
  - Registro y login de usuarios
  - Generación de tokens JWT
  - Validación de roles (USER, ADMIN)
  - Middleware de autenticación

#### 2. **properties-api** (Puerto 8081)
- **Base de datos:** MongoDB
- **Funcionalidad:** CRUD de propiedades y sistema completo de reservas
- **Características:**
  - Gestión de propiedades (crear, editar, eliminar, listar)
  - Sistema de reservas (bookings) con validación de solapamiento
  - Upload de imágenes locales
  - Publicación de eventos a RabbitMQ
  - Manejo de concurrencia con goroutines
  - Validación de permisos por rol

#### 3. **search-api** (Puerto 8082)
- **Motor de búsqueda:** Apache Solr
- **Funcionalidad:** Búsqueda avanzada y paginada de propiedades
- **Características:**
  - Búsqueda full-text con Solr
  - Paginación de resultados
  - Sistema de caché de dos niveles (CCache local + Memcached)
  - Consumer de RabbitMQ para sincronización automática
  - Filtros avanzados (precio, ciudad, país, habitaciones, baños, huéspedes)

### Frontend (React + Vite)

- **Framework:** React 19 con Vite 7
- **Routing:** React Router DOM
- **HTTP Client:** Axios
- **UI:** Tailwind CSS + Lucide React (iconos)
- **Funcionalidades:**
  - Login y registro de usuarios
  - Búsqueda de propiedades con paginación
  - Detalles de propiedad
  - Sistema de reservas (crear, editar, cancelar)
  - Panel de administración
  - Gestión de imágenes locales

### Infraestructura

- **Reverse Proxy:** Nginx
- **Message Broker:** RabbitMQ
- **Caché:** Memcached
- **Bases de datos:** MySQL, MongoDB
- **Motor de búsqueda:** Apache Solr
- **Orquestación:** Docker Compose

---

## 🛠️ Stack Tecnológico

### Backend
- **Lenguaje:** Go 1.21+
- **Frameworks:**
  - Gin Gonic (HTTP router)
  - GORM (ORM para MySQL)
  - MongoDB Go Driver
- **Autenticación:** JWT (JSON Web Tokens)
- **Message Queue:** RabbitMQ
- **Búsqueda:** Apache Solr
- **Caché:** CCache (local) + Memcached (distribuido)

### Frontend
- **Framework:** React 19
- **Build Tool:** Vite 7
- **Routing:** React Router DOM 7
- **HTTP Client:** Axios
- **Estilos:** Tailwind CSS
- **Iconos:** Lucide React

### DevOps
- **Contenedores:** Docker
- **Orquestación:** Docker Compose
- **Reverse Proxy:** Nginx
- **Volúmenes:** Docker Volumes para persistencia

---

## 🚀 Instalación y Despliegue

### Prerrequisitos

1. **Docker Desktop** instalado y ejecutándose
   - Windows: https://www.docker.com/products/docker-desktop
   - Verificar instalación: `docker --version`
   - Verificar que Docker Desktop esté corriendo

2. **Git** (opcional, para clonar el repositorio)

### Pasos para Desplegar en una Nueva PC

#### 1. Clonar o Descargar el Proyecto

```bash
# Si tienes Git instalado
git clone https://github.com/LautyUCC/ArquiSoftII.git
cd ArquiSoftII

# O descargar el ZIP y descomprimirlo
```

#### 2. Verificar Archivos Necesarios

Asegúrate de tener estos archivos en la raíz del proyecto:
- ✅ `docker-compose.yml`
- ✅ Carpetas: `backend/`, `frontend/`, `nginx/`

#### 3. Levantar los Servicios

**Opción A: Levantar todo de una vez (recomendado)**
```bash
docker-compose up --build
```

**Opción B: Levantar paso a paso (para debugging)**

Primero, levantar solo la infraestructura:
```bash
docker-compose up -d mysql mongodb rabbitmq solr memcached nginx
```

Esperar a que todos los servicios estén saludables (30-60 segundos), luego:
```bash
docker-compose up -d users-api properties-api search-api frontend
```

#### 4. Verificar que Todo Esté Funcionando

Abre tu navegador y visita:
- **Frontend:** http://localhost:3000
- **RabbitMQ Management:** http://localhost:15672 (guest/guest)
- **Solr Admin:** http://localhost:8983

#### 5. Crear Usuarios de Prueba

**Opción A: Usar scripts incluidos**

Windows:
```bash
scripts\create_test_users.bat
```

Linux/Mac:
```bash
chmod +x scripts/create_test_users.sh
./scripts/create_test_users.sh
```

**Opción B: Manualmente vía API**

```bash
# Crear usuario normal
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","password":"password123","name":"Usuario Test","role":"USER"}'

# Crear administrador
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"password123","name":"Admin Test","role":"ADMIN"}'
```

#### 6. Detener los Servicios

```bash
# Detener todos los servicios
docker-compose down

# Detener y eliminar volúmenes (CUIDADO: borra datos)
docker-compose down -v
```

---

## ✨ Funcionalidades Implementadas

### Funcionalidades Base (Implementadas Previamente)

1. **Sistema de Autenticación**
   - Registro de usuarios
   - Login con JWT
   - Roles y permisos (USER, ADMIN)
   - Middleware de autenticación y autorización

2. **Gestión de Propiedades**
   - CRUD completo de propiedades
   - Validación de permisos (solo ADMIN puede crear/editar/eliminar)
   - Búsqueda avanzada con filtros
   - Paginación de resultados

3. **Sistema de Búsqueda**
   - Búsqueda full-text con Solr
   - Filtros por precio, ciudad, país, habitaciones, baños, huéspedes
   - Caché de dos niveles para optimización
   - Sincronización automática con RabbitMQ

### Funcionalidades Implementadas Hoy

#### 1. **Sistema Completo de Reservas (Bookings)**

**Backend:**
- ✅ Creación de reservas con validación de solapamiento
- ✅ Edición de fechas de reservas
- ✅ Cancelación de reservas (soft delete)
- ✅ Visualización de reservas propias y de propiedades
- ✅ Cálculo automático de precio total
- ✅ Validación de disponibilidad
- ✅ Publicación de eventos a RabbitMQ al crear reservas

**Frontend:**
- ✅ Página "Mis Reservas" (`/my-bookings`)
- ✅ Creación de reservas desde detalles de propiedad
- ✅ Edición de fechas con modal
- ✅ Cancelación de reservas
- ✅ Visualización de reservas confirmadas y canceladas
- ✅ Bloqueo de fechas ya reservadas en el calendario

**Herramientas utilizadas:**
- Go (Gin, MongoDB Driver)
- React (useState, useEffect, React Router)
- Axios para peticiones HTTP
- JWT para autenticación
- RabbitMQ para eventos asíncronos

#### 2. **Sistema de Upload de Imágenes Locales**

**Backend:**
- ✅ Endpoint `POST /api/upload` para subir imágenes
- ✅ Validación de formato (JPG, PNG, GIF, WEBP)
- ✅ Validación de tamaño (máximo 5MB)
- ✅ Almacenamiento en volumen compartido de Docker
- ✅ Generación de nombres únicos para archivos

**Frontend:**
- ✅ Input de tipo file para seleccionar imágenes
- ✅ Preview de imágenes antes de crear propiedad
- ✅ Subida automática al seleccionar archivo
- ✅ Eliminación de imágenes del preview
- ✅ Soporte para múltiples imágenes por propiedad

**Infraestructura:**
- ✅ Volumen Docker compartido entre `properties-api` y `nginx`
- ✅ Configuración de Nginx para servir imágenes estáticas
- ✅ Cache de imágenes configurado

**Herramientas utilizadas:**
- Go (multipart/form-data handling)
- React (File API, FormData)
- Docker Volumes
- Nginx (servir archivos estáticos)

#### 3. **Paginación en Búsqueda de Propiedades**

**Backend:**
- ✅ Soporte de parámetros `page` y `pageSize` en search-api
- ✅ Cálculo de total de páginas
- ✅ Integración con Solr para paginación

**Frontend:**
- ✅ 6 propiedades por página
- ✅ Controles de navegación completos:
  - Botón "Primera página"
  - Botón "Página anterior"
  - Números de página con elipsis
  - Botón "Página siguiente"
  - Botón "Última página"
- ✅ Información de paginación (página X de Y, total de resultados)
- ✅ Scroll automático al cambiar de página

**Herramientas utilizadas:**
- Solr (paginación nativa)
- React (estado de paginación)
- Tailwind CSS (estilos de controles)

#### 4. **Panel de Administración Extendido**

**Funcionalidades:**
- ✅ Gestión de usuarios (ver, editar roles)
- ✅ Gestión de propiedades (CRUD completo)
- ✅ Gestión de reservas (ver todas, cancelar cualquier reserva)
- ✅ Upload de imágenes al editar propiedades
- ✅ Eliminación de propiedades

**Herramientas utilizadas:**
- React (componentes modales, tabs)
- Axios (peticiones HTTP)
- JWT (validación de rol ADMIN)

#### 5. **Mejoras en UX/UI**

- ✅ Bloqueo visual de fechas ya reservadas (rojo)
- ✅ Banner informativo con fechas no disponibles
- ✅ Validación de solapamiento en cliente antes de enviar
- ✅ Mensajes de error claros y específicos
- ✅ Estados de carga y feedback visual
- ✅ Confirmaciones antes de acciones destructivas

**Herramientas utilizadas:**
- React (conditional rendering, state management)
- Tailwind CSS (estilos condicionales)
- Lucide React (iconos)

---

## 📡 Endpoints de la API

### users-api (Puerto 8080)

```
POST   /users              # Crear usuario
GET    /users/:id          # Obtener usuario por ID
POST   /users/login        # Login (retorna JWT)
```

### properties-api (Puerto 8081)

**Públicos:**
```
GET    /api/properties/:id              # Obtener propiedad por ID
GET    /api/properties/user/:userId    # Obtener propiedades de un usuario
```

**Protegidos (requieren JWT):**
```
POST   /api/upload                      # Subir imagen
POST   /api/bookings                    # Crear reserva
GET    /api/bookings/my                 # Mis reservas
GET    /api/bookings/property/:id       # Reservas de una propiedad
PATCH  /api/bookings/:id                # Editar reserva
DELETE /api/bookings/:id                # Cancelar reserva
```

**Solo ADMIN:**
```
POST   /api/properties                  # Crear propiedad
PUT    /api/properties/:id             # Actualizar propiedad
DELETE /api/properties/:id             # Eliminar propiedad
GET    /api/admin/bookings              # Ver todas las reservas
```

### search-api (Puerto 8082)

```
GET    /search?query=...&page=1&pageSize=6&city=...&minPrice=...&maxPrice=...
       # Búsqueda paginada con filtros
```

---

## 📁 Estructura del Proyecto

```
ArquiSoftII/
├── backend/
│   ├── users-api/          # Microservicio de usuarios
│   │   ├── controllers/    # Controladores HTTP
│   │   ├── services/      # Lógica de negocio
│   │   ├── repositories/   # Acceso a datos
│   │   ├── domain/         # Modelos de dominio
│   │   ├── dto/            # Data Transfer Objects
│   │   ├── middleware/     # Middleware de autenticación
│   │   └── utils/          # Utilidades
│   │
│   ├── properties-api/     # Microservicio de propiedades y reservas
│   │   ├── controllers/    # Controladores HTTP
│   │   ├── services/       # Lógica de negocio
│   │   ├── repositories/   # Acceso a datos (MongoDB)
│   │   ├── domain/         # Modelos de dominio
│   │   ├── dto/            # Data Transfer Objects
│   │   ├── middleware/     # Middleware de autenticación/autorización
│   │   ├── clients/        # Clientes HTTP (users-api, RabbitMQ)
│   │   ├── config/         # Configuración
│   │   └── utils/          # Utilidades
│   │
│   └── search-api/         # Microservicio de búsqueda
│       ├── controllers/    # Controladores HTTP
│       ├── services/       # Lógica de búsqueda
│       ├── repositories/   # Acceso a Solr y caché
│       ├── consumers/      # Consumer de RabbitMQ
│       ├── domain/         # Modelos de dominio
│       ├── dto/            # Data Transfer Objects
│       └── utils/          # Utilidades
│
├── frontend/               # Aplicación React
│   ├── src/
│   │   ├── pages/          # Páginas principales
│   │   │   ├── Login.jsx
│   │   │   ├── Register.jsx
│   │   │   ├── Search.jsx
│   │   │   ├── PropertyDetail.jsx
│   │   │   ├── MyBookings.jsx
│   │   │   └── Admin.jsx
│   │   ├── components/     # Componentes reutilizables
│   │   ├── services/      # Cliente API (Axios)
│   │   └── App.jsx         # Componente principal
│   ├── public/            # Archivos estáticos
│   ├── Dockerfile
│   └── package.json
│
├── nginx/                  # Configuración de Nginx
│   └── nginx.conf
│
├── scripts/                # Scripts de utilidad
│   ├── create_test_users.sh
│   ├── create_test_users.bat
│   └── create_test_property.sh
│
├── docker-compose.yml      # Orquestación de servicios
└── README.md              # Este archivo
```

---

## 🔧 Herramientas y Tecnologías por Tarea

### Sistema de Reservas (Bookings)

**Backend:**
- **Go + Gin Gonic:** Framework HTTP para crear endpoints REST
- **MongoDB Go Driver:** Persistencia de reservas en MongoDB
- **JWT:** Autenticación y extracción de userId del token
- **RabbitMQ:** Publicación de eventos `booking.created` de forma asíncrona
- **Goroutines:** Manejo de concurrencia en operaciones de base de datos

**Frontend:**
- **React + React Router:** Navegación y rutas protegidas
- **Axios:** Cliente HTTP para peticiones a la API
- **localStorage:** Almacenamiento de JWT y datos de usuario
- **useState/useEffect:** Gestión de estado y ciclo de vida
- **Tailwind CSS:** Estilos para modales y formularios

### Upload de Imágenes

**Backend:**
- **Go multipart/form-data:** Manejo de archivos subidos
- **Docker Volumes:** Almacenamiento persistente compartido
- **Validación de archivos:** Verificación de tipo MIME y tamaño

**Frontend:**
- **File API:** Acceso a archivos seleccionados por el usuario
- **FormData:** Construcción de requests multipart
- **URL.createObjectURL:** Preview de imágenes antes de subir

**Infraestructura:**
- **Docker Volumes:** Volumen compartido `property_images`
- **Nginx:** Servir archivos estáticos con cache

### Paginación

**Backend:**
- **Solr:** Paginación nativa con parámetros `start` y `rows`
- **Cálculo de totalPages:** Matemática para determinar páginas totales

**Frontend:**
- **React State:** Estado de página actual, total de páginas
- **Algoritmo de paginación:** Cálculo de números de página visibles con elipsis
- **Scroll automático:** `window.scrollTo()` para UX mejorada

### Panel de Administración

**Frontend:**
- **React Tabs:** Navegación entre secciones (usuarios, propiedades, reservas)
- **Conditional Rendering:** Mostrar/ocultar según rol de usuario
- **Modales:** Formularios para crear/editar
- **Tablas:** Visualización de datos en formato tabla

### Búsqueda y Filtros

**Backend:**
- **Solr Query Syntax:** Construcción de queries complejas
- **Filtros (fq):** Filtrado por múltiples criterios simultáneos
- **Caché:** CCache (local) + Memcached (distribuido) para optimización

**Frontend:**
- **Formularios controlados:** React controlled inputs
- **Query parameters:** Construcción de URLs con filtros

### Autenticación y Autorización

**Backend:**
- **JWT (JSON Web Tokens):** Tokens firmados para autenticación
- **Middleware:** Validación de tokens y roles
- **GORM:** Consultas a MySQL para validar usuarios

**Frontend:**
- **Protected Routes:** Componente que valida JWT antes de renderizar
- **Axios Interceptors:** Agregar JWT automáticamente a requests
- **localStorage:** Persistencia de token entre sesiones

### Sincronización de Datos

**Backend:**
- **RabbitMQ:** Message broker para eventos asíncronos
- **Consumer Pattern:** search-api consume eventos de creación/actualización
- **HTTP Client:** search-api consulta properties-api para obtener datos completos

### Optimización y Performance

**Backend:**
- **Caché de dos niveles:** CCache (rápido, local) + Memcached (distribuido)
- **Goroutines:** Procesamiento concurrente de operaciones
- **Índices de MongoDB:** Optimización de consultas

**Frontend:**
- **Vite:** Build tool rápido con HMR (Hot Module Replacement)
- **Code Splitting:** Carga lazy de componentes
- **Optimización de imágenes:** Servidas desde Nginx con cache

---

## 🌐 URLs de Acceso

Una vez desplegado, puedes acceder a:

- **Frontend:** http://localhost:3000
- **users-api:** http://localhost:8080
- **properties-api:** http://localhost:8081
- **search-api:** http://localhost:8082
- **RabbitMQ Management:** http://localhost:15672 (usuario: `guest`, contraseña: `guest`)
- **Solr Admin:** http://localhost:8983

---

## 🧪 Testing

### Probar el Sistema

1. **Crear un usuario administrador:**
   ```bash
   curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@test.com","password":"password123","name":"Admin","role":"ADMIN"}'
   ```

2. **Hacer login:**
   ```bash
   curl -X POST http://localhost:8080/users/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@test.com","password":"password123"}'
   ```

3. **Crear una propiedad (usar el token JWT obtenido):**
   ```bash
   curl -X POST http://localhost:8081/api/properties \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <TU_TOKEN_JWT>" \
     -d '{"title":"Casa Test","description":"Descripción","pricePerNight":100,"city":"Buenos Aires","country":"Argentina","maxGuests":4}'
   ```

4. **Buscar propiedades:**
   ```bash
   curl "http://localhost:8082/search?query=casa&page=1&pageSize=6"
   ```

---

## 📚 Recursos y Documentación

- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [GORM](https://gorm.io/docs/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [RabbitMQ Tutorials](https://www.rabbitmq.com/tutorials)
- [Apache Solr Guide](https://solr.apache.org/guide/)
- [React Documentation](https://react.dev/)
- [Vite Documentation](https://vitejs.dev/)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [Docker Compose](https://docs.docker.com/compose/)

---

## 👥 Equipo

Proyecto desarrollado para Arquitectura de Software II - UCC

---

## 📝 Notas Finales

Este proyecto implementa una arquitectura de microservicios completa con:
- ✅ Separación de responsabilidades
- ✅ Comunicación asíncrona con RabbitMQ
- ✅ Búsqueda avanzada con Solr
- ✅ Caché distribuido
- ✅ Autenticación y autorización robusta
- ✅ Sistema completo de reservas
- ✅ Upload de archivos
- ✅ Paginación y filtros
- ✅ Panel de administración

**Estado del Proyecto:** ✅ Listo para presentación final
