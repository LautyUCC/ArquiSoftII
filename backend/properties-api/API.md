# Properties API - Documentación de Endpoints

API REST para gestión de propiedades inmobiliarias. Esta documentación describe todos los endpoints disponibles, sus parámetros, respuestas y posibles errores.

**Base URL:** `http://localhost:8081`

---

## 1. Crear Propiedad

Crea una nueva propiedad inmobiliaria con validación de usuario y cálculo automático de precio.

### Endpoint

```
POST /properties
```

### Descripción

Crea una nueva propiedad validando que el usuario propietario existe en users-api. El precio final se calcula automáticamente usando concurrencia (precio base + impuestos 21% + $50 por amenidad + $30 por persona de capacidad).

### Headers

```
Content-Type: application/json
```

### Request Body

```json
{
  "title": "Hermoso apartamento en el centro",
  "description": "Apartamento moderno con vista al mar, completamente amueblado. Ideal para vacaciones o estancias largas. Ubicado en el corazón de la ciudad, cerca de restaurantes y transporte público.",
  "price": 150000.00,
  "location": "Bogotá, Colombia",
  "ownerId": "user123",
  "amenities": ["wifi", "pool", "parking", "kitchen", "air-conditioning"],
  "capacity": 4,
  "available": true
}
```

### Response Success (201 Created)

```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "title": "Hermoso apartamento en el centro",
    "description": "Apartamento moderno con vista al mar, completamente amueblado. Ideal para vacaciones o estancias largas. Ubicado en el corazón de la ciudad, cerca de restaurantes y transporte público.",
    "price": 181750.00,
    "location": "Bogotá, Colombia",
    "ownerId": "user123",
    "amenities": ["wifi", "pool", "parking", "kitchen", "air-conditioning"],
    "capacity": 4,
    "available": true,
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  },
  "message": "Property created successfully"
}
```

### Posibles Errores

| Código | Descripción | Ejemplo |
|--------|-------------|---------|
| **400 Bad Request** | Request body inválido o campos faltantes | `{"error": "Invalid request body", "message": "Key: 'PropertyCreateDTO.Price' Error:Field validation for 'Price' failed on the 'required' tag"}` |
| **404 Not Found** | Usuario propietario no existe | `{"error": "User not found", "message": "usuario owner con ID 'user123' no existe"}` |
| **500 Internal Server Error** | Error interno del servidor | `{"error": "Internal server error", "message": "error creando propiedad en repositorio: connection timeout"}` |

---

## 2. Obtener Propiedad por ID

Obtiene los detalles de una propiedad específica por su ID.

### Endpoint

```
GET /properties/:id
```

### Descripción

Obtiene la información completa de una propiedad utilizando su ID único de MongoDB.

### Headers

```
Content-Type: application/json
```

### Path Parameters

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `id` | string | ID único de la propiedad (ObjectID de MongoDB) |

### Response Success (200 OK)

```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "title": "Hermoso apartamento en el centro",
    "description": "Apartamento moderno con vista al mar, completamente amueblado.",
    "price": 181750.00,
    "location": "Bogotá, Colombia",
    "ownerId": "user123",
    "amenities": ["wifi", "pool", "parking", "kitchen", "air-conditioning"],
    "capacity": 4,
    "available": true,
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  }
}
```

### Posibles Errores

| Código | Descripción | Ejemplo |
|--------|-------------|---------|
| **400 Bad Request** | ID no proporcionado | `{"error": "Invalid request", "message": "Property ID is required"}` |
| **404 Not Found** | Propiedad no encontrada | `{"error": "Property not found", "message": "propiedad con ID '507f1f77bcf86cd799439011' no encontrada"}` |
| **500 Internal Server Error** | Error interno del servidor | `{"error": "Internal server error", "message": "error obteniendo propiedad: database connection failed"}` |

---

## 3. Actualizar Propiedad

Actualiza una propiedad existente. Solo el propietario puede actualizar su propiedad.

### Endpoint

```
PUT /properties/:id
```

### Descripción

Actualiza los campos de una propiedad existente. Solo se actualizan los campos proporcionados (actualización parcial). El usuario debe ser el propietario de la propiedad para poder actualizarla.

### Headers

```
Content-Type: application/json
Authorization: Bearer <token>
```

**Nota:** El `user_id` debe estar disponible en el contexto (seteado por middleware de autenticación).

### Path Parameters

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `id` | string | ID único de la propiedad a actualizar |

### Request Body

Todos los campos son opcionales. Solo los campos proporcionados serán actualizados.

```json
{
  "title": "Apartamento actualizado con vista al mar",
  "price": 180000.00,
  "available": false,
  "amenities": ["wifi", "pool", "parking", "kitchen", "air-conditioning", "tv"]
}
```

### Response Success (200 OK)

```json
{
  "success": true,
  "message": "Property updated successfully"
}
```

### Posibles Errores

| Código | Descripción | Ejemplo |
|--------|-------------|---------|
| **400 Bad Request** | Request body inválido o ID faltante | `{"error": "Invalid request body", "message": "invalid JSON format"}` |
| **401 Unauthorized** | User ID no encontrado en contexto | `{"error": "Unauthorized", "message": "User ID not found in context"}` |
| **403 Forbidden** | Usuario no tiene permisos para actualizar | `{"error": "Forbidden", "message": "usuario con ID 'user456' no tiene permisos para actualizar propiedad '507f1f77bcf86cd799439011' (owner: 'user123')"}` |
| **404 Not Found** | Propiedad no encontrada | `{"error": "Property not found", "message": "propiedad con ID '507f1f77bcf86cd799439011' no encontrada"}` |
| **500 Internal Server Error** | Error interno del servidor | `{"error": "Internal server error", "message": "error actualizando propiedad en repositorio: database error"}` |

---

## 4. Eliminar Propiedad

Elimina una propiedad existente. Solo el propietario puede eliminar su propiedad.

### Endpoint

```
DELETE /properties/:id
```

### Descripción

Elimina permanentemente una propiedad del sistema. Solo el propietario de la propiedad puede eliminarla. Se publica un evento de eliminación en RabbitMQ.

### Headers

```
Content-Type: application/json
Authorization: Bearer <token>
```

**Nota:** El `user_id` debe estar disponible en el contexto (seteado por middleware de autenticación).

### Path Parameters

| Parámetro | Tipo | Descripción |
|-----------|------|-------------|
| `id` | string | ID único de la propiedad a eliminar |

### Request Body

No requiere body.

### Response Success (200 OK)

```json
{
  "success": true,
  "message": "Property deleted successfully"
}
```

### Posibles Errores

| Código | Descripción | Ejemplo |
|--------|-------------|---------|
| **400 Bad Request** | ID no proporcionado | `{"error": "Invalid request", "message": "Property ID is required"}` |
| **401 Unauthorized** | User ID no encontrado en contexto | `{"error": "Unauthorized", "message": "User ID not found in context"}` |
| **403 Forbidden** | Usuario no tiene permisos para eliminar | `{"error": "Forbidden", "message": "usuario con ID 'user456' no tiene permisos para eliminar propiedad '507f1f77bcf86cd799439011' (owner: 'user123')"}` |
| **404 Not Found** | Propiedad no encontrada | `{"error": "Property not found", "message": "propiedad con ID '507f1f77bcf86cd799439011' no encontrada"}` |
| **500 Internal Server Error** | Error interno del servidor | `{"error": "Internal server error", "message": "error eliminando propiedad en repositorio: database error"}` |

---

## 5. Listar Propiedades del Usuario

Obtiene todas las propiedades de un usuario específico.

### Endpoint

```
GET /properties/user
```

### Descripción

Retorna una lista de todas las propiedades que pertenecen al usuario autenticado. El `user_id` se obtiene del contexto de autenticación.

### Headers

```
Content-Type: application/json
Authorization: Bearer <token>
```

**Nota:** El `user_id` debe estar disponible en el contexto (seteado por middleware de autenticación).

### Response Success (200 OK)

```json
{
  "success": true,
  "data": [
    {
      "id": "507f1f77bcf86cd799439011",
      "title": "Hermoso apartamento en el centro",
      "description": "Apartamento moderno con vista al mar.",
      "price": 181750.00,
      "location": "Bogotá, Colombia",
      "ownerId": "user123",
      "amenities": ["wifi", "pool", "parking"],
      "capacity": 4,
      "available": true,
      "createdAt": "2024-01-15T10:30:00Z",
      "updatedAt": "2024-01-15T10:30:00Z"
    },
    {
      "id": "507f1f77bcf86cd799439012",
      "title": "Casa con jardín en las afueras",
      "description": "Casa espaciosa con jardín privado, ideal para familias.",
      "price": 250000.00,
      "location": "Medellín, Colombia",
      "ownerId": "user123",
      "amenities": ["wifi", "kitchen", "parking", "garden"],
      "capacity": 6,
      "available": true,
      "createdAt": "2024-01-20T14:15:00Z",
      "updatedAt": "2024-01-20T14:15:00Z"
    }
  ],
  "count": 2
}
```

### Posibles Errores

| Código | Descripción | Ejemplo |
|--------|-------------|---------|
| **401 Unauthorized** | User ID no encontrado en contexto | `{"error": "Unauthorized", "message": "User ID not found in context"}` |
| **500 Internal Server Error** | Error interno del servidor | `{"error": "Internal server error", "message": "error obteniendo propiedades del usuario: database connection failed"}` |

---

## Códigos de Estado HTTP

| Código | Descripción | Uso |
|--------|-------------|-----|
| **200 OK** | Operación exitosa | GET, PUT, DELETE exitosos |
| **201 Created** | Recurso creado exitosamente | POST exitoso |
| **400 Bad Request** | Solicitud inválida | Request body o parámetros incorrectos |
| **401 Unauthorized** | No autenticado | Falta token o user ID en contexto |
| **403 Forbidden** | Sin permisos | Usuario no es propietario de la propiedad |
| **404 Not Found** | Recurso no encontrado | Propiedad o usuario no existe |
| **500 Internal Server Error** | Error del servidor | Errores internos de base de datos o servicios |

---

## Formato de Respuestas

### Respuesta Exitosa

Todas las respuestas exitosas incluyen un campo `success: true`:

```json
{
  "success": true,
  "data": { ... },
  "message": "Operación exitosa"
}
```

### Respuesta de Error

Todas las respuestas de error incluyen campos `error` y `message`:

```json
{
  "error": "Error type",
  "message": "Descripción detallada del error"
}
```

---

## Notas Adicionales

### Cálculo de Precio

El precio final de una propiedad se calcula automáticamente usando goroutines:

- **Precio base con impuestos:** `precio base × 1.21` (21% de impuestos)
- **Costo por amenidad:** `$50 × cantidad de amenidades`
- **Costo por capacidad:** `$30 × capacidad de personas`

**Ejemplo:**
- Precio base: $100,000
- Amenidades: 3 (wifi, pool, parking)
- Capacidad: 4 personas

**Cálculo:**
- Precio con impuestos: $100,000 × 1.21 = $121,000
- Costo amenidades: 3 × $50 = $150
- Costo capacidad: 4 × $30 = $120
- **Precio final: $121,270**

### Eventos en RabbitMQ

Todas las operaciones de creación, actualización y eliminación publican eventos en RabbitMQ:

- **Cola:** `property_events`
- **Eventos:** `create`, `update`, `delete`
- **Formato:** JSON con `operation` y `propertyId`

### Validación de Usuarios

El servicio valida que los usuarios existan en `users-api` antes de crear propiedades. La comunicación se realiza mediante HTTP GET a `{users-api-url}/users/{userID}`.

---

## Ejemplos de Uso

### Crear una Propiedad

```bash
curl -X POST http://localhost:8081/properties \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Apartamento moderno",
    "description": "Apartamento con todas las comodidades",
    "price": 120000.00,
    "location": "Bogotá, Colombia",
    "ownerId": "user123",
    "amenities": ["wifi", "pool"],
    "capacity": 3,
    "available": true
  }'
```

### Obtener una Propiedad

```bash
curl -X GET http://localhost:8081/properties/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json"
```

### Actualizar una Propiedad

```bash
curl -X PUT http://localhost:8081/properties/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "title": "Título actualizado",
    "available": false
  }'
```

### Eliminar una Propiedad

```bash
curl -X DELETE http://localhost:8081/properties/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>"
```

### Listar Propiedades del Usuario

```bash
curl -X GET http://localhost:8081/properties/user \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>"
```

