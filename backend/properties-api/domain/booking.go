package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookingStatus representa los estados posibles de una reserva
const (
	BookingStatusConfirmed = "CONFIRMED"
	BookingStatusCancelled = "CANCELLED"
)

// Booking representa una reserva de una propiedad
// Persistida en MongoDB en la misma base de datos que las propiedades
type Booking struct {
	// ID es el identificador único de MongoDB
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// PropertyID es el identificador de la propiedad reservada
	PropertyID string `bson:"propertyId" json:"propertyId"`

	// UserID es el identificador del usuario que realizó la reserva
	UserID string `bson:"userId" json:"userId"`

	// StartDate es la fecha de inicio de la reserva (sin hora o con hora 00:00)
	// Formato: YYYY-MM-DD con hora 00:00:00
	StartDate time.Time `bson:"startDate" json:"startDate"`

	// EndDate es la fecha de fin de la reserva
	// NOTA: La fecha de fin es EXCLUSIVA, es decir, el huésped debe salir antes de esta fecha
	// Ejemplo: Si StartDate = 2025-01-01 y EndDate = 2025-01-05,
	// la reserva es para las noches del 1, 2, 3 y 4 de enero (4 noches)
	// El huésped debe salir antes del 5 de enero (EndDate)
	// Formato: YYYY-MM-DD con hora 00:00:00
	EndDate time.Time `bson:"endDate" json:"endDate"`

	// TotalPrice es el precio total de la reserva al momento de crearla
	// Se calcula como: precio por noche * número de noches
	// Se guarda para mantener un registro histórico del precio (el precio de la propiedad puede cambiar)
	TotalPrice float64 `bson:"totalPrice" json:"totalPrice"`

	// Status es el estado de la reserva
	// Valores posibles: "CONFIRMED" o "CANCELLED"
	Status string `bson:"status" json:"status"`

	// CreatedAt es la fecha y hora de creación del registro
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`

	// UpdatedAt es la fecha y hora de última actualización
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// IsValidStatus verifica si el estado de la reserva es válido
func (b *Booking) IsValidStatus() bool {
	return b.Status == BookingStatusConfirmed || b.Status == BookingStatusCancelled
}

// CalculateNights calcula el número de noches de la reserva
// Considera que EndDate es exclusiva
func (b *Booking) CalculateNights() int {
	if b.EndDate.Before(b.StartDate) || b.EndDate.Equal(b.StartDate) {
		return 0
	}
	duration := b.EndDate.Sub(b.StartDate)
	return int(duration.Hours() / 24)
}

