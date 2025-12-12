# 📅 Módulo de Bookings - Documentación Completa

Este documento describe la implementación completa del módulo de bookings en properties-api.

---

## 📋 Endpoints Implementados

### 1. **POST /api/bookings** - Crear Reserva

**Autorización:** Requiere autenticación (cualquier rol: USER, ADMIN, OWNER)

**Descripción:** Crea una nueva reserva para una propiedad. El userId se obtiene del JWT automáticamente.

**Request Body:**
```json
{
  "propertyId": "507f1f77bcf86cd799439011",
  "checkIn": "2024-02-15T14:00:00Z",
  "checkOut": "2024-02-20T14:00:00Z"
}
```

**Validaciones:**
- ✅ Propiedad existe
- ✅ Propiedad está disponible
- ✅ Fechas válidas (check-in < check-out, check-in no en el pasado)
- ✅ **Verifica solapamiento de reservas confirmadas** (409 si hay conflicto)

**Códigos HTTP:**
- `201 Created` - Reserva creada exitosamente
- `400 Bad Request` - Error de validación (propiedad no existe, fechas inválidas, etc.)
- `409 Conflict` - **Conflicto de reserva** (ya existe una reserva confirmada en esas fechas)
- `500 Internal Server Error` - Error interno del servidor

**Response (201):**
```json
{
  "id": "507f1f77bcf86cd799439012",
  "propertyId": "507f1f77bcf86cd799439011",
  "userId": "123",
  "checkIn": "2024-02-15T14:00:00Z",
  "checkOut": "2024-02-20T14:00:00Z",
  "totalPrice": 500.00,
  "status": "confirmed",
  "createdAt": "2024-01-15T10:30:00Z"
}
```

**Evento RabbitMQ:** Publica `booking.created` después de crear la reserva exitosamente.

---

### 2. **GET /api/bookings/my** - Mis Reservas

**Autorización:** Requiere autenticación (cualquier rol)

**Descripción:** Obtiene todas las reservas del usuario autenticado (userId del JWT).

**Códigos HTTP:**
- `200 OK` - Lista de reservas
- `500 Internal Server Error` - Error interno del servidor

**Response (200):**
```json
[
  {
    "id": "507f1f77bcf86cd799439012",
    "propertyId": "507f1f77bcf86cd799439011",
    "userId": "123",
    "checkIn": "2024-02-15T14:00:00Z",
    "checkOut": "2024-02-20T14:00:00Z",
    "totalPrice": 500.00,
    "status": "confirmed",
    "createdAt": "2024-01-15T10:30:00Z"
  }
]
```

---

### 3. **GET /api/bookings/property/:id** - Reservas de una Propiedad (Solo ADMIN)

**Autorización:** Requiere rol ADMIN (validado por middleware RequireAdmin)

**Descripción:** Obtiene todas las reservas de una propiedad específica. Solo ADMIN puede acceder.

**Códigos HTTP:**
- `200 OK` - Lista de reservas
- `400 Bad Request` - Propiedad no existe
- `401 Unauthorized` - No hay token o token inválido
- `403 Forbidden` - Token válido pero no es ADMIN
- `500 Internal Server Error` - Error interno del servidor

**Response (200):**
```json
[
  {
    "id": "507f1f77bcf86cd799439012",
    "propertyId": "507f1f77bcf86cd799439011",
    "userId": "123",
    "checkIn": "2024-02-15T14:00:00Z",
    "checkOut": "2024-02-20T14:00:00Z",
    "totalPrice": 500.00,
    "status": "confirmed",
    "createdAt": "2024-01-15T10:30:00Z"
  }
]
```

---

## 🔧 Implementación Técnica

### Service: CreateBooking

**Función:** `CreateBooking(ctx context.Context, createDTO dto.BookingCreateDTO, userID string) (dto.BookingDTO, error)`

**Pasos:**
1. ✅ Verifica que la propiedad existe
2. ✅ Verifica que la propiedad está disponible
3. ✅ Valida fechas (check-in < check-out, check-in no en el pasado)
4. ✅ **Verifica solapamiento de reservas confirmadas** (409 si hay conflicto)
5. ✅ Calcula precio total (precio por noche × número de noches)
6. ✅ Crea la reserva con status "confirmed"
7. ✅ Guarda en repositorio
8. ✅ **Publica evento `booking.created` en RabbitMQ**

**Verificación de Solapamiento:**
```go
// Busca todas las reservas confirmadas para la propiedad
existingBookings, err := s.bookingRepo.FindConfirmedByPropertyID(ctx, propertyID)

// Verifica solapamiento con cada reserva existente
for _, existing := range existingBookings {
    if s.hasOverlap(checkIn, checkOut, existing.CheckIn, existing.CheckOut) {
        return error 409 Conflict
    }
}
```

**Función hasOverlap:**
```go
// Dos rangos se solapan si: start1 < end2 AND start2 < end1
func hasOverlap(start1, end1, start2, end2 time.Time) bool {
    return start1.Before(end2) && start2.Before(end1)
}
```

---

### Repository

**Métodos agregados:**
- `FindByPropertyID(ctx, propertyID)` - Obtiene todas las reservas de una propiedad
- `FindConfirmedByPropertyID(ctx, propertyID)` - Obtiene solo reservas confirmadas de una propiedad (para verificar solapamientos)

---

### Eventos RabbitMQ

**Evento:** `booking.created`

**Publicación:** Después de crear la reserva exitosamente

**Cola:** `property_events` (misma cola que eventos de propiedades)

**Datos:** 
- Operation: `"booking.created"`
- PropertyID: ID de la reserva (bookingID)

**Nota:** Si falla la publicación del evento, se registra un log de advertencia pero no se falla la operación (la reserva ya fue creada).

---

## 📊 Códigos HTTP y Manejo de Errores

| Código | Descripción | Cuándo se usa |
|--------|-------------|---------------|
| 201 Created | Reserva creada exitosamente | POST /bookings exitoso |
| 200 OK | Operación exitosa | GET /bookings/my, GET /bookings/property/:id |
| 400 Bad Request | Error de validación | Propiedad no existe, fechas inválidas, etc. |
| 409 Conflict | **Conflicto de reserva** | Ya existe una reserva confirmada en esas fechas |
| 401 Unauthorized | No autenticado | No hay token o token inválido |
| 403 Forbidden | Sin permisos | Token válido pero no es ADMIN (GET /bookings/property/:id) |
| 500 Internal Server Error | Error interno | Error en base de datos, RabbitMQ, etc. |

---

## 📝 Logs

El servicio incluye logs claros en cada paso:

**Ejemplos de logs:**
```
📅 [CreateBooking] Iniciando creación de reserva para usuario 123, propiedad 507f1f77bcf86cd799439011
✅ [CreateBooking] Propiedad 507f1f77bcf86cd799439011 encontrada: Apartamento en el centro
✅ [CreateBooking] No hay conflictos de reserva para las fechas 2024-02-15 - 2024-02-20
💰 [CreateBooking] Precio total calculado: $500.00 (5 noches x $100.00)
✅ [CreateBooking] Reserva creada exitosamente con ID: 507f1f77bcf86cd799439012
✅ [CreateBooking] Evento 'booking.created' publicado en RabbitMQ para reserva 507f1f77bcf86cd799439012
```

**Logs de error:**
```
❌ [CreateBooking] Propiedad 507f1f77bcf86cd799439011 no existe: ...
❌ [CreateBooking] Conflicto de reserva detectado: reserva 507f1f77bcf86cd799439012 solapa con fechas 2024-02-15 - 2024-02-20
⚠️ [CreateBooking] Error publicando evento 'booking.created' en RabbitMQ para reserva 507f1f77bcf86cd799439012: ...
```

---

## 🧪 Testing

### Test 1: Crear Reserva Exitosamente

```bash
curl -X POST http://localhost:8082/api/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "propertyId": "507f1f77bcf86cd799439011",
    "checkIn": "2024-02-15T14:00:00Z",
    "checkOut": "2024-02-20T14:00:00Z"
  }'
# Response: 201 Created
```

### Test 2: Conflicto de Reserva (409)

```bash
# Primera reserva
curl -X POST http://localhost:8082/api/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "propertyId": "507f1f77bcf86cd799439011",
    "checkIn": "2024-02-15T14:00:00Z",
    "checkOut": "2024-02-20T14:00:00Z"
  }'
# Response: 201 Created

# Segunda reserva con fechas solapadas
curl -X POST http://localhost:8082/api/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "propertyId": "507f1f77bcf86cd799439011",
    "checkIn": "2024-02-18T14:00:00Z",
    "checkOut": "2024-02-25T14:00:00Z"
  }'
# Response: 409 Conflict
```

### Test 3: Obtener Mis Reservas

```bash
curl -X GET http://localhost:8082/api/bookings/my \
  -H "Authorization: Bearer <token>"
# Response: 200 OK con lista de reservas
```

### Test 4: Obtener Reservas de Propiedad (Solo ADMIN)

```bash
curl -X GET http://localhost:8082/api/bookings/property/507f1f77bcf86cd799439011 \
  -H "Authorization: Bearer <token_admin>"
# Response: 200 OK con lista de reservas
```

---

## ✅ Verificación

Para verificar que todo funciona correctamente:

1. **Compilar el servicio:**
   ```bash
   cd backend/properties-api && go build
   ```

2. **Probar creación de reserva:**
   - Crear reserva exitosamente → 201
   - Crear reserva con fechas solapadas → 409
   - Crear reserva con propiedad inexistente → 400

3. **Probar endpoints de lectura:**
   - GET /bookings/my con token válido → 200
   - GET /bookings/property/:id con token ADMIN → 200
   - GET /bookings/property/:id con token USER → 403

4. **Verificar eventos RabbitMQ:**
   - Verificar que se publica `booking.created` después de crear reserva

---

## 📊 Resumen de Archivos

1. `backend/properties-api/services/booking_service.go` - Servicio con CreateBooking, GetMyBookings, GetPropertyBookings
2. `backend/properties-api/controllers/booking_controller.go` - Controlador con los 3 endpoints
3. `backend/properties-api/repositories/booking_repository.go` - Repositorio con FindByPropertyID y FindConfirmedByPropertyID
4. `backend/properties-api/main.go` - Rutas actualizadas
5. `backend/properties-api/MODULO_BOOKINGS.md` - Este documento

