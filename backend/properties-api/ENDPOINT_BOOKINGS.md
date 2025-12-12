# 📅 Endpoint POST /bookings - Documentación

Este documento describe la implementación del endpoint POST /bookings en properties-api.

---

## 🔐 Regla de Autorización

**POST /bookings solo requiere autenticación (cualquier rol: USER, ADMIN, OWNER)**

- ✅ **Requiere:** Token JWT válido
- ❌ **NO requiere:** Rol específico (cualquier usuario autenticado puede crear reservas)
- 🔑 **userId:** Se obtiene automáticamente del JWT, NO del body del request

Esta regla asegura que:
1. Un usuario solo pueda crear reservas para sí mismo
2. Cualquier usuario autenticado (USER, ADMIN, OWNER) pueda reservar propiedades
3. No se pueda crear reservas para otros usuarios manipulando el body del request

---

## 📋 Endpoint

```
POST /api/bookings
```

### Headers Requeridos

```
Authorization: Bearer <token_jwt>
Content-Type: application/json
```

### Request Body

```json
{
  "propertyId": "507f1f77bcf86cd799439011",
  "checkIn": "2024-02-15T14:00:00Z",
  "checkOut": "2024-02-20T14:00:00Z"
}
```

**Nota:** El campo `userId` NO debe enviarse en el body (se ignora si se envía). El userId se obtiene automáticamente del JWT.

### Response Success (201 Created)

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

### Posibles Errores

| Código | Descripción |
|--------|-------------|
| 401 Unauthorized | No hay token o el token es inválido |
| 400 Bad Request | Propiedad no existe, no está disponible, fechas inválidas, etc. |

---

## 🔧 Implementación

### 1. Ruta (`main.go`)

```go
// Rutas protegidas (requieren autenticación - cualquier rol puede acceder)
protected := router.Group("/api")
protected.Use(middleware.RequireAuth(jwtSecret)) // requireAuth valida token válido (cualquier rol)
{
    // POST /bookings: Cualquier usuario autenticado puede crear reservas
    // El userId se obtiene del JWT, no del body del request
    // Esto asegura que un usuario solo pueda crear reservas para sí mismo
    protected.POST("/bookings", bookingController.CreateBooking)
}
```

### 2. Controlador (`controllers/booking_controller.go`)

```go
// CreateBooking maneja la creación de una nueva reserva
// REGLA DE AUTORIZACIÓN: Solo requiere autenticación (cualquier rol: USER, ADMIN, OWNER)
// El userId se obtiene del JWT, no del body del request
// Esto asegura que un usuario solo pueda crear reservas para sí mismo
func (c *BookingController) CreateBooking(ctx *gin.Context) {
    // Extraer userID del contexto (agregado por middleware RequireAuth)
    userIDValue, exists := ctx.Get("userID")
    // ... validación y conversión ...
    
    // Ignorar el userId del body si viene (por seguridad, siempre usamos el del JWT)
    createDTO.UserID = "" // Limpiar el userId del body
    
    // Llamar al servicio con el userId del JWT
    responseDTO, err := c.service.CreateBooking(requestCtx, createDTO, userID)
    // ...
}
```

### 3. Servicio (`services/booking_service.go`)

```go
// CreateBooking crea una nueva reserva
// REGLA DE NEGOCIO: Solo requiere autenticación (cualquier rol: USER, ADMIN, OWNER)
// El userId se obtiene del JWT, no del body del request
// Esto asegura que un usuario solo pueda crear reservas para sí mismo
func (s *bookingService) CreateBooking(ctx context.Context, createDTO dto.BookingCreateDTO, userID string) (dto.BookingDTO, error) {
    // Validaciones:
    // 1. Propiedad existe
    // 2. Propiedad está disponible
    // 3. Fechas válidas
    // 4. Calcular precio total
    // 5. Crear booking con userId del JWT
    // ...
}
```

### 4. DTO (`dto/booking_dto.go`)

```go
// BookingCreateDTO DTO para crear una nueva reserva
// NOTA: userId NO debe venir en el body, se obtiene del JWT automáticamente
// Esto asegura que un usuario solo pueda crear reservas para sí mismo
type BookingCreateDTO struct {
    PropertyID string    `json:"propertyId" binding:"required"`
    UserID     string    `json:"userId,omitempty"` // OPCIONAL: se ignora, se usa el userId del JWT
    CheckIn    time.Time `json:"checkIn" binding:"required"`
    CheckOut   time.Time `json:"checkOut" binding:"required"`
}
```

---

## 🧪 Testing

### Test 1: Crear Reserva sin Token (debe retornar 401)

```bash
curl -X POST http://localhost:8082/api/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "propertyId": "507f1f77bcf86cd799439011",
    "checkIn": "2024-02-15T14:00:00Z",
    "checkOut": "2024-02-20T14:00:00Z"
  }'
# Response: {"error": "authorization header required"}
```

### Test 2: Crear Reserva con Token de USER (debe funcionar)

```bash
curl -X POST http://localhost:8082/api/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_user>" \
  -d '{
    "propertyId": "507f1f77bcf86cd799439011",
    "checkIn": "2024-02-15T14:00:00Z",
    "checkOut": "2024-02-20T14:00:00Z"
  }'
# Response: 201 Created con la reserva creada
# Nota: El userId en la respuesta será el del token, no el del body (si se envió)
```

### Test 3: Intentar Crear Reserva para Otro Usuario (debe ignorar userId del body)

```bash
curl -X POST http://localhost:8082/api/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_user_123>" \
  -d '{
    "propertyId": "507f1f77bcf86cd799439011",
    "userId": "999",  // Este userId se ignora
    "checkIn": "2024-02-15T14:00:00Z",
    "checkOut": "2024-02-20T14:00:00Z"
  }'
# Response: 201 Created
# El userId en la respuesta será "123" (del token), no "999" (del body)
```

### Test 4: Crear Reserva con Token de ADMIN (debe funcionar)

```bash
curl -X POST http://localhost:8082/api/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_admin>" \
  -d '{
    "propertyId": "507f1f77bcf86cd799439011",
    "checkIn": "2024-02-15T14:00:00Z",
    "checkOut": "2024-02-20T14:00:00Z"
  }'
# Response: 201 Created con la reserva creada
```

---

## ⚠️ Notas Importantes

1. **Seguridad:** El userId siempre se obtiene del JWT, nunca del body. Esto previene que un usuario cree reservas para otros usuarios.

2. **Cualquier Rol:** USER, ADMIN y OWNER pueden crear reservas. No hay restricción de rol.

3. **Validaciones:** El servicio valida:
   - Que la propiedad existe
   - Que la propiedad está disponible
   - Que las fechas son válidas (check-in < check-out, check-in no en el pasado)
   - Calcula el precio total automáticamente

4. **Middleware:** Se usa `RequireAuth()` que valida el token JWT pero NO restringe por rol.

---

## 📊 Resumen de Archivos

1. `backend/properties-api/services/booking_service.go` - **NUEVO** - Servicio de bookings
2. `backend/properties-api/controllers/booking_controller.go` - **NUEVO** - Controlador de bookings
3. `backend/properties-api/main.go` - Actualizado para incluir la ruta POST /bookings
4. `backend/properties-api/dto/booking_dto.go` - Actualizado para documentar que userId es opcional
5. `backend/properties-api/ENDPOINT_BOOKINGS.md` - Este documento

---

## ✅ Verificación

Para verificar que todo funciona correctamente:

1. **Compilar el servicio:**
   ```bash
   cd backend/properties-api && go build
   ```

2. **Probar con diferentes roles:**
   - Crear reserva con token de USER → debe funcionar
   - Crear reserva con token de ADMIN → debe funcionar
   - Crear reserva sin token → debe retornar 401

3. **Verificar seguridad:**
   - Intentar crear reserva con userId diferente en el body → debe usar el userId del token

