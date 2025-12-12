# ✅ Endpoints de Bookings - Implementación Completa

## 📋 Resumen

Se ha implementado un sistema completo de endpoints HTTP para bookings con todas las reglas de negocio y autorización requeridas.

## 🔐 Reglas de Autorización

### Roles y Permisos

- **Cualquier usuario autenticado (USER o ADMIN):**
  - ✅ Crear reservas (`POST /bookings`)
  - ✅ Ver sus propias reservas (`GET /bookings/my`)
  - ✅ Editar sus propias reservas (`PATCH /bookings/:id`)
  - ✅ Eliminar sus propias reservas (`DELETE /bookings/:id`)
  - ✅ Ver reservas CONFIRMED de una propiedad (`GET /bookings/property/:id`)

- **Rol ADMIN:**
  - ✅ Puede ver y gestionar **todas** las reservas (no solo las propias)
  - ✅ Puede editar cualquier reserva (`PATCH /bookings/:id`)
  - ✅ Puede eliminar cualquier reserva (`DELETE /bookings/:id`)

## 📝 Reglas de Negocio

### Validaciones de Fechas

1. **startDate < endDate**: La fecha de inicio debe ser anterior a la fecha de fin
2. **startDate no puede ser en el pasado**: No se pueden crear reservas con fechas pasadas
3. **Normalización de fechas**: Todas las fechas se normalizan a hora 00:00:00

### Validación de Solapamiento

- **No se puede crear una reserva** si ya existe una reserva CONFIRMED para la misma propiedad con fechas que se solapen
- **No se puede cambiar las fechas** de una reserva a un rango que se solape con otra reserva CONFIRMED de esa propiedad
- **Al actualizar**: Se excluye la reserva actual del chequeo de solapamiento

### Soft Delete

- **OPCIÓN ELEGIDA**: Marca la reserva como `CANCELLED` (soft delete) en lugar de borrarla físicamente
- **Razón**: Mantener un registro histórico de reservas canceladas para auditoría y análisis
- **Implementación**: Se actualiza el campo `status` a `CANCELLED` y `updatedAt` al momento actual

## 🌐 Endpoints Implementados

### 1. POST /api/bookings

**Descripción:** Crea una nueva reserva

**Autorización:** Cualquier usuario autenticado (USER o ADMIN)

**Body:**
```json
{
  "propertyId": "string (requerido)",
  "startDate": "2025-01-15T00:00:00Z (requerido)",
  "endDate": "2025-01-20T00:00:00Z (requerido)"
}
```

**Compatibilidad Frontend:**
También acepta `checkIn`/`checkOut` que se mapean automáticamente a `startDate`/`endDate`:
```json
{
  "propertyId": "string",
  "checkIn": "2025-01-15T00:00:00Z",
  "checkOut": "2025-01-20T00:00:00Z"
}
```

**Validaciones:**
- ✅ Propiedad existe
- ✅ Propiedad está disponible
- ✅ startDate < endDate
- ✅ startDate no es en el pasado
- ✅ No hay solapamiento con reservas CONFIRMED

**Códigos HTTP:**
- `201 Created`: Reserva creada exitosamente
- `400 Bad Request`: Error de validación
- `409 Conflict`: Hay solapamiento con otra reserva CONFIRMED
- `500 Internal Server Error`: Error del servidor

**Evento RabbitMQ:**
Publica evento `booking.created` con el ID de la reserva

**Ejemplo de respuesta:**
```json
{
  "id": "693b59b7d790b976bd70fcf8",
  "propertyId": "693b56f8d790b976bd70fcf6",
  "userId": "1",
  "checkIn": "2025-01-15T00:00:00Z",
  "checkOut": "2025-01-20T00:00:00Z",
  "totalPrice": 12400.0,
  "status": "CONFIRMED",
  "createdAt": "2025-12-11T23:42:48Z",
  "updatedAt": "2025-12-11T23:42:48Z"
}
```

---

### 2. GET /api/bookings/my

**Descripción:** Obtiene todas las reservas del usuario autenticado

**Autorización:** Cualquier usuario autenticado (USER o ADMIN)

**Query Parameters:** Ninguno

**Ordenamiento:** Por fecha inicio descendente (más recientes primero)

**Códigos HTTP:**
- `200 OK`: Lista de reservas obtenida exitosamente
- `500 Internal Server Error`: Error del servidor

**Ejemplo de respuesta:**
```json
[
  {
    "id": "693b59b7d790b976bd70fcf8",
    "propertyId": "693b56f8d790b976bd70fcf6",
    "userId": "1",
    "checkIn": "2025-01-15T00:00:00Z",
    "checkOut": "2025-01-20T00:00:00Z",
    "totalPrice": 12400.0,
    "status": "CONFIRMED",
    "createdAt": "2025-12-11T23:42:48Z",
    "updatedAt": "2025-12-11T23:42:48Z"
  }
]
```

---

### 3. PATCH /api/bookings/:id

**Descripción:** Actualiza las fechas de una reserva

**Autorización:** Solo el dueño de la reserva o un ADMIN

**Parámetros de URL:**
- `id`: ID de la reserva a actualizar

**Body:**
```json
{
  "startDate": "2025-01-16T00:00:00Z (opcional)",
  "endDate": "2025-01-21T00:00:00Z (opcional)"
}
```

**Compatibilidad Frontend:**
También acepta `checkIn`/`checkOut`:
```json
{
  "checkIn": "2025-01-16T00:00:00Z",
  "checkOut": "2025-01-21T00:00:00Z"
}
```

**Validaciones:**
- ✅ Reserva existe
- ✅ Usuario es dueño o ADMIN
- ✅ startDate < endDate
- ✅ No hay solapamiento con otras reservas CONFIRMED (excluyendo la reserva actual)

**Códigos HTTP:**
- `200 OK`: Reserva actualizada exitosamente
- `400 Bad Request`: Error de validación
- `403 Forbidden`: Usuario no tiene permisos
- `409 Conflict`: Hay solapamiento con otra reserva CONFIRMED
- `500 Internal Server Error`: Error del servidor

**Nota:** Si las fechas cambian, se recalcula el precio total automáticamente

---

### 4. DELETE /api/bookings/:id

**Descripción:** Elimina una reserva (marca como CANCELLED)

**Autorización:** Solo el dueño de la reserva o un ADMIN

**Parámetros de URL:**
- `id`: ID de la reserva a eliminar

**Body:** Ninguno

**Códigos HTTP:**
- `200 OK`: Reserva cancelada exitosamente
- `403 Forbidden`: Usuario no tiene permisos
- `404 Not Found`: Reserva no existe
- `500 Internal Server Error`: Error del servidor

**Implementación:**
- Marca la reserva como `CANCELLED` (soft delete)
- Actualiza `updatedAt` al momento actual
- **NO** borra físicamente la reserva (mantiene registro histórico)

**Ejemplo de respuesta:**
```json
{
  "message": "Reserva cancelada exitosamente"
}
```

---

### 5. GET /api/bookings/property/:id

**Descripción:** Obtiene todas las reservas CONFIRMED de una propiedad

**Autorización:** Cualquier usuario autenticado (USER o ADMIN)

**Parámetros de URL:**
- `id`: ID de la propiedad

**Query Parameters:** Ninguno

**Uso:** Útil para bloquear el calendario en el frontend

**Códigos HTTP:**
- `200 OK`: Lista de reservas obtenida exitosamente
- `400 Bad Request`: Propiedad no existe
- `500 Internal Server Error`: Error del servidor

**Nota:** Solo devuelve reservas con `status = "CONFIRMED"` (no incluye CANCELLED)

**Ejemplo de respuesta:**
```json
[
  {
    "id": "693b59b7d790b976bd70fcf8",
    "propertyId": "693b56f8d790b976bd70fcf6",
    "userId": "1",
    "checkIn": "2025-01-15T00:00:00Z",
    "checkOut": "2025-01-20T00:00:00Z",
    "totalPrice": 12400.0,
    "status": "CONFIRMED",
    "createdAt": "2025-12-11T23:42:48Z",
    "updatedAt": "2025-12-11T23:42:48Z"
  }
]
```

---

## 🔄 Mapeo DTO ↔ Dominio

### Al crear/actualizar (Frontend → Backend):
- Frontend envía: `checkIn`, `checkOut` (o `startDate`, `endDate`)
- Se normaliza a: `StartDate` (00:00:00), `EndDate` (00:00:00)
- Se guarda en MongoDB: `startDate`, `endDate`

### Al leer (Backend → Frontend):
- MongoDB tiene: `startDate`, `endDate`
- Se mapea a DTO: `checkIn`, `checkOut`
- Frontend recibe: `checkIn`, `checkOut`

## 📊 Estados de Reserva

- **CONFIRMED**: Reserva confirmada y activa
- **CANCELLED**: Reserva cancelada (soft delete)

## 🐰 Eventos RabbitMQ

### booking.created

Se publica cuando se crea una nueva reserva exitosamente.

**Cola:** `property_events`

**Contenido:**
```json
{
  "operation": "booking.created",
  "propertyId": "693b56f8d790b976bd70fcf6"
}
```

**Nota:** Por ahora se publica solo el ID de la reserva. En el futuro se podría extender para incluir más información (userId, startDate, endDate, etc.).

## ✅ Verificación

1. **Compilación:** ✅ Sin errores
2. **Endpoints:** ✅ Todos implementados
3. **Autorización:** ✅ Reglas de rol implementadas
4. **Validaciones:** ✅ Todas las reglas de negocio implementadas
5. **Soft Delete:** ✅ Implementado (marca como CANCELLED)
6. **Ordenamiento:** ✅ GetMyBookings ordenado por fecha descendente
7. **Eventos RabbitMQ:** ✅ booking.created publicado

