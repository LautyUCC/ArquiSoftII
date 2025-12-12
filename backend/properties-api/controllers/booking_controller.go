package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"properties-api/dto"
	"properties-api/services"
	"properties-api/utils"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	service services.BookingService
}

func NewBookingController(service services.BookingService) *BookingController {
	return &BookingController{
		service: service,
	}
}

// CreateBooking maneja la creación de una nueva reserva
// REGLA DE AUTORIZACIÓN: Solo requiere autenticación (cualquier rol: USER, ADMIN, OWNER)
// El userId se obtiene del JWT, no del body del request
// Esto asegura que un usuario solo pueda crear reservas para sí mismo
// Códigos HTTP: 201 (creado), 400 (bad request), 409 (conflict), 500 (error interno)
func (c *BookingController) CreateBooking(ctx *gin.Context) {
	var createDTO dto.BookingCreateDTO
	if err := ctx.ShouldBindJSON(&createDTO); err != nil {
		utils.LogError("POST /bookings", err, ctx)
		utils.BadRequest(ctx, "Error al parsear el cuerpo de la solicitud", nil)
		return
	}

	// Extraer userID del contexto (agregado por middleware RequireAuth)
	// El middleware RequireAuth valida el token JWT y guarda el user_id en el contexto
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		userIDValue, exists = ctx.Get("user_id") // Intentar con otro nombre
		if !exists {
			utils.LogError("POST /bookings", nil, ctx)
			utils.Unauthorized(ctx, "Usuario no autenticado")
			return
		}
	}

	// Convertir userID de any a string
	var userID string
	switch v := userIDValue.(type) {
	case uint:
		userID = fmt.Sprintf("%d", v)
	case int:
		userID = fmt.Sprintf("%d", v)
	case string:
		userID = v
	default:
		utils.LogError("POST /bookings", nil, ctx)
		utils.InternalServerError(ctx, "Tipo de userID inválido", nil)
		return
	}

	// Ignorar el userId del body si viene (por seguridad, siempre usamos el del JWT)
	// Esto previene que un usuario intente crear reservas para otros usuarios
	createDTO.UserID = "" // Limpiar el userId del body

	// Mapear CheckIn/CheckOut a StartDate/EndDate si vienen (compatibilidad con frontend)
	checkInTime := createDTO.CheckIn.Time()
	checkOutTime := createDTO.CheckOut.Time()
	startDateTime := createDTO.StartDate.Time()
	endDateTime := createDTO.EndDate.Time()

	if !checkInTime.IsZero() && startDateTime.IsZero() {
		createDTO.StartDate = createDTO.CheckIn
		startDateTime = checkInTime
	}
	if !checkOutTime.IsZero() && endDateTime.IsZero() {
		createDTO.EndDate = createDTO.CheckOut
		endDateTime = checkOutTime
	}

	// Validar que al menos una de las fechas esté presente
	if startDateTime.IsZero() && endDateTime.IsZero() {
		utils.LogError("POST /bookings", nil, ctx)
		utils.BadRequest(ctx, "Se requiere startDate/endDate o checkIn/checkOut", nil)
		return
	}

	// Crear contexto con timeout
	requestCtx, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	// Llamar al servicio con el userId del JWT
	responseDTO, err := c.service.CreateBooking(requestCtx, createDTO, userID)
	if err != nil {
		// Verificar si es un error de conflicto (409)
		if err.Error() == "conflict: la propiedad ya tiene una reserva confirmada en esas fechas" {
			utils.LogError("POST /bookings", err, ctx)
			utils.Conflict(ctx, err.Error())
			return
		}
		// Otros errores son 400 (bad request)
		utils.LogError("POST /bookings", err, ctx)
		utils.BadRequest(ctx, err.Error(), nil)
		return
	}

	utils.LogInfo("POST /bookings", "Reserva creada exitosamente", ctx)
	utils.SendSuccess(ctx, http.StatusCreated, responseDTO)
}

// GetMyBookings maneja la obtención de reservas del usuario autenticado
// REGLA DE AUTORIZACIÓN: Solo requiere autenticación (cualquier rol)
// Códigos HTTP: 200 (ok), 500 (error interno)
func (c *BookingController) GetMyBookings(ctx *gin.Context) {
	// Extraer userID del contexto (agregado por middleware RequireAuth)
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		userIDValue, exists = ctx.Get("user_id")
		if !exists {
			utils.LogError("GET /bookings/my", nil, ctx)
			utils.Unauthorized(ctx, "Usuario no autenticado")
			return
		}
	}

	// Convertir userID de any a string
	var userID string
	switch v := userIDValue.(type) {
	case uint:
		userID = fmt.Sprintf("%d", v)
	case int:
		userID = fmt.Sprintf("%d", v)
	case string:
		userID = v
	default:
		utils.LogError("GET /bookings/my", nil, ctx)
		utils.InternalServerError(ctx, "Tipo de userID inválido", nil)
		return
	}

	// Crear contexto con timeout
	requestCtx, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	// Llamar al servicio
	bookings, err := c.service.GetMyBookings(requestCtx, userID)
	if err != nil {
		utils.LogError("GET /bookings/my", err, ctx)
		utils.InternalServerError(ctx, err.Error(), nil)
		return
	}

	utils.LogInfo("GET /bookings/my", "Reservas obtenidas exitosamente", ctx)
	utils.SendSuccess(ctx, http.StatusOK, bookings)
}

// GetPropertyBookings maneja la obtención de todas las reservas de una propiedad
// REGLA DE AUTORIZACIÓN: Solo ADMIN puede acceder (validado por middleware RequireAdmin)
// Códigos HTTP: 200 (ok), 400 (bad request), 500 (error interno)
func (c *BookingController) GetPropertyBookings(ctx *gin.Context) {
	propertyID := ctx.Param("id")
	if propertyID == "" {
		utils.LogError("GET /bookings/property/:id", nil, ctx)
		utils.BadRequest(ctx, "property ID es requerido", nil)
		return
	}

	// Crear contexto con timeout
	requestCtx, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	// Llamar al servicio
	bookings, err := c.service.GetPropertyBookings(requestCtx, propertyID)
	if err != nil {
		// Verificar si es un error de propiedad no encontrada
		if err.Error() == fmt.Sprintf("propiedad con ID '%s' no existe", propertyID) {
			utils.LogError("GET /bookings/property/:id", err, ctx)
			utils.BadRequest(ctx, err.Error(), nil)
			return
		}
		utils.LogError("GET /bookings/property/:id", err, ctx)
		utils.InternalServerError(ctx, err.Error(), nil)
		return
	}

	utils.LogInfo("GET /bookings/property/:id", "Reservas obtenidas exitosamente", ctx)
	utils.SendSuccess(ctx, http.StatusOK, bookings)
}

// UpdateBooking maneja la actualización de fechas de una reserva (PATCH)
// REGLA DE AUTORIZACIÓN: Solo el dueño de la reserva o un ADMIN puede actualizar
// Códigos HTTP: 200 (ok), 400 (bad request), 403 (forbidden), 409 (conflict), 500 (error interno)
func (c *BookingController) UpdateBooking(ctx *gin.Context) {
	bookingID := ctx.Param("id")
	if bookingID == "" {
		utils.LogError("PATCH /bookings/:id", nil, ctx)
		utils.BadRequest(ctx, "booking ID es requerido", nil)
		return
	}

	var updateDTO dto.BookingUpdateDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		utils.LogError("PATCH /bookings/:id", err, ctx)
		utils.BadRequest(ctx, "Error al parsear el cuerpo de la solicitud", nil)
		return
	}

	// Extraer userID y role del contexto (agregado por middleware RequireAuth)
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		userIDValue, exists = ctx.Get("user_id")
		if !exists {
			utils.LogError("PATCH /bookings/:id", nil, ctx)
			utils.Unauthorized(ctx, "Usuario no autenticado")
			return
		}
	}

	// Convertir userID de any a string
	var userID string
	switch v := userIDValue.(type) {
	case uint:
		userID = fmt.Sprintf("%d", v)
	case int:
		userID = fmt.Sprintf("%d", v)
	case string:
		userID = v
	default:
		utils.LogError("PATCH /bookings/:id", nil, ctx)
		utils.InternalServerError(ctx, "Tipo de userID inválido", nil)
		return
	}

	// Obtener role del contexto
	roleValue, exists := ctx.Get("role")
	userRole := "USER"
	if exists {
		if role, ok := roleValue.(string); ok {
			userRole = role
		}
	}

	// Crear contexto con timeout
	requestCtx, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	// Llamar al servicio
	responseDTO, err := c.service.UpdateBooking(requestCtx, bookingID, updateDTO, userID, userRole)
	if err != nil {
		// Verificar si es un error de conflicto (409)
		if err.Error() == "conflict: la propiedad ya tiene una reserva confirmada en esas fechas" {
			utils.LogError("PATCH /bookings/:id", err, ctx)
			utils.Conflict(ctx, err.Error())
			return
		}
		// Verificar si es un error de permisos (403)
		if err.Error() == "no tienes permisos para actualizar esta reserva" {
			utils.LogError("PATCH /bookings/:id", err, ctx)
			utils.Forbidden(ctx, err.Error())
			return
		}
		// Otros errores son 400 (bad request)
		utils.LogError("PATCH /bookings/:id", err, ctx)
		utils.BadRequest(ctx, err.Error(), nil)
		return
	}

	utils.LogInfo("PATCH /bookings/:id", "Reserva actualizada exitosamente", ctx)
	utils.SendSuccess(ctx, http.StatusOK, responseDTO)
}

// DeleteBooking maneja la eliminación de una reserva (DELETE)
// REGLA DE AUTORIZACIÓN: Solo el dueño de la reserva o un ADMIN puede eliminar
// OPCIÓN ELEGIDA: Marca la reserva como CANCELLED (soft delete) en lugar de borrarla físicamente
// Códigos HTTP: 200 (ok), 403 (forbidden), 404 (not found), 500 (error interno)
func (c *BookingController) DeleteBooking(ctx *gin.Context) {
	bookingID := ctx.Param("id")
	if bookingID == "" {
		utils.LogError("DELETE /bookings/:id", nil, ctx)
		utils.BadRequest(ctx, "booking ID es requerido", nil)
		return
	}

	// Extraer userID y role del contexto (agregado por middleware RequireAuth)
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		userIDValue, exists = ctx.Get("user_id")
		if !exists {
			utils.LogError("DELETE /bookings/:id", nil, ctx)
			utils.Unauthorized(ctx, "Usuario no autenticado")
			return
		}
	}

	// Convertir userID de any a string
	var userID string
	switch v := userIDValue.(type) {
	case uint:
		userID = fmt.Sprintf("%d", v)
	case int:
		userID = fmt.Sprintf("%d", v)
	case string:
		userID = v
	default:
		utils.LogError("DELETE /bookings/:id", nil, ctx)
		utils.InternalServerError(ctx, "Tipo de userID inválido", nil)
		return
	}

	// Obtener role del contexto
	roleValue, exists := ctx.Get("role")
	userRole := "USER"
	if exists {
		if role, ok := roleValue.(string); ok {
			userRole = role
		}
	}

	// Crear contexto con timeout
	requestCtx, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	// Llamar al servicio
	err := c.service.DeleteBooking(requestCtx, bookingID, userID, userRole)
	if err != nil {
		// Verificar si es un error de permisos (403)
		if err.Error() == "no tienes permisos para eliminar esta reserva" {
			utils.LogError("DELETE /bookings/:id", err, ctx)
			utils.Forbidden(ctx, err.Error())
			return
		}
		// Verificar si es un error de no encontrado (404)
		if err.Error() == fmt.Sprintf("reserva con ID '%s' no existe", bookingID) {
			utils.LogError("DELETE /bookings/:id", err, ctx)
			utils.NotFound(ctx, err.Error())
			return
		}
		// Otros errores son 500 (error interno)
		utils.LogError("DELETE /bookings/:id", err, ctx)
		utils.InternalServerError(ctx, err.Error(), nil)
		return
	}

	utils.LogInfo("DELETE /bookings/:id", "Reserva eliminada exitosamente", ctx)
	utils.SendSuccess(ctx, http.StatusOK, gin.H{"message": "Reserva cancelada exitosamente"})
}

// GetAllBookings obtiene todas las reservas del sistema (solo para admin)
// Requiere rol ADMIN (validado por middleware RequireAdmin)
func (c *BookingController) GetAllBookings(ctx *gin.Context) {
	// Crear contexto con timeout
	requestCtx, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	// Llamar al servicio
	bookings, err := c.service.GetAllBookings(requestCtx)
	if err != nil {
		utils.LogError("GET /admin/bookings", err, ctx)
		utils.InternalServerError(ctx, err.Error(), nil)
		return
	}

	utils.LogInfo("GET /admin/bookings", fmt.Sprintf("Obtenidas %d reservas", len(bookings)), ctx)
	utils.SendSuccess(ctx, http.StatusOK, bookings)
}

