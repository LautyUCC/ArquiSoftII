package dto

import (
	"strings"
	"time"
)

// DateOnly es un tipo personalizado para parsear fechas en formato YYYY-MM-DD
// Acepta tanto formato YYYY-MM-DD como formato ISO completo
type DateOnly time.Time

// UnmarshalJSON implementa la interfaz json.Unmarshaler para parsear fechas
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	// Remover comillas del JSON
	s := strings.Trim(string(b), `"`)

	// Intentar parsear formato YYYY-MM-DD primero
	if t, err := time.Parse("2006-01-02", s); err == nil {
		*d = DateOnly(t)
		return nil
	}

	// Si falla, intentar formato ISO completo (RFC3339)
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		*d = DateOnly(t)
		return nil
	}

	// Si falla, intentar formato ISO sin zona horaria
	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		*d = DateOnly(t)
		return nil
	}

	// Si todo falla, retornar error
	_, err := time.Parse("2006-01-02", s)
	return err // Retornar el error del primer intento
}

// Time convierte DateOnly a time.Time
func (d DateOnly) Time() time.Time {
	return time.Time(d)
}

// BookingCreateDTO DTO para crear una nueva reserva
// NOTA: userId NO debe venir en el body, se obtiene del JWT automáticamente
// Esto asegura que un usuario solo pueda crear reservas para sí mismo
// Usa startDate/endDate según requerimientos, pero mantiene CheckIn/CheckOut para compatibilidad con frontend
type BookingCreateDTO struct {
	PropertyID string   `json:"propertyId" binding:"required"` // ID de la propiedad a reservar
	UserID     string   `json:"userId,omitempty"`                // OPCIONAL: se ignora, se usa el userId del JWT
	StartDate  DateOnly `json:"startDate"`     // Fecha de inicio (00:00:00) - se valida manualmente
	EndDate    DateOnly `json:"endDate"`       // Fecha de fin EXCLUSIVA (00:00:00) - se valida manualmente
	// Campos legacy para compatibilidad con frontend (se mapean a StartDate/EndDate)
	CheckIn    DateOnly `json:"checkIn,omitempty"`      // Mapeado a StartDate
	CheckOut   DateOnly `json:"checkOut,omitempty"`      // Mapeado a EndDate
}

// BookingUpdateDTO DTO para actualizar fechas de una reserva (PATCH)
// Solo permite actualizar startDate y endDate
type BookingUpdateDTO struct {
	StartDate *DateOnly `json:"startDate,omitempty"` // Fecha de inicio (opcional)
	EndDate   *DateOnly `json:"endDate,omitempty"`   // Fecha de fin EXCLUSIVA (opcional)
	// Campos legacy para compatibilidad con frontend
	CheckIn  *DateOnly `json:"checkIn,omitempty"`  // Mapeado a StartDate
	CheckOut *DateOnly `json:"checkOut,omitempty"` // Mapeado a EndDate
}

// BookingDTO representa el DTO de respuesta de una reserva
// NOTA: Usa CheckIn/CheckOut para compatibilidad con el frontend
// Internamente el dominio usa StartDate/EndDate
type BookingDTO struct {
	ID         string    `json:"id"`
	PropertyID string    `json:"propertyId"`
	UserID     string    `json:"userId"`
	CheckIn    time.Time `json:"checkIn"`    // Mapeado desde StartDate del dominio
	CheckOut   time.Time `json:"checkOut"`   // Mapeado desde EndDate del dominio
	TotalPrice float64   `json:"totalPrice"` // Precio total de la reserva
	Status     string    `json:"status"`      // "CONFIRMED" o "CANCELLED"
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`  // Timestamp de última actualización
}
