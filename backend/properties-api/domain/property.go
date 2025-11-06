package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Property representa una propiedad disponible para reserva (tipo Airbnb)
// Contiene toda la información necesaria para listar y reservar propiedades
type Property struct {
	// ID es el identificador único de MongoDB para la propiedad
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Title es el título o nombre de la propiedad
	// Ejemplo: "Hermoso apartamento en el centro", "Casa con piscina"
	Title string `json:"title" bson:"title" binding:"required"`

	// Description contiene la descripción detallada de la propiedad
	// Incluye información sobre características, ambiente, cercanías, etc.
	Description string `json:"description" bson:"description" binding:"required"`

	// Price es el precio por noche de la propiedad
	// Valor en la moneda local (ej: pesos, dólares, euros)
	Price float64 `json:"price" bson:"price" binding:"required,gt=0"`

	// Location es la ubicación completa de la propiedad
	// Puede incluir dirección, ciudad, país o coordenadas
	// Ejemplo: "Bogotá, Colombia", "Calle 123, Medellín"
	Location string `json:"location" bson:"location" binding:"required"`

	// OwnerID es el identificador del usuario propietario de la propiedad
	// Referencia al usuario en el sistema (users-api)
	OwnerID string `json:"ownerId" bson:"ownerId" binding:"required"`

	// Amenities es una lista de comodidades y servicios disponibles en la propiedad
	// Ejemplos: ["wifi", "pool", "parking", "kitchen", "air-conditioning", "tv"]
	// Puede incluir servicios como wifi, piscina, estacionamiento, cocina, etc.
	Amenities []string `json:"amenities" bson:"amenities"`

	// Capacity es la cantidad máxima de personas que pueden hospedarse
	// Define cuántas personas pueden ocupar la propiedad simultáneamente
	Capacity int `json:"capacity" bson:"capacity" binding:"required,gt=0"`

	// Available indica si la propiedad está disponible para reserva
	// true = disponible, false = no disponible (ya reservada o deshabilitada)
	Available bool `json:"available" bson:"available"`

	// CreatedAt es la fecha y hora de creación del registro en formato string
	// Formato recomendado: ISO 8601 (ej: "2024-01-15T10:30:00Z")
	CreatedAt string `json:"createdAt" bson:"createdAt"`

	// UpdatedAt es la fecha y hora de la última actualización en formato string
	// Se actualiza cada vez que se modifica la propiedad
	// Formato recomendado: ISO 8601 (ej: "2024-01-15T10:30:00Z")
	UpdatedAt string `json:"updatedAt" bson:"updatedAt"`
}

// PropertyUpdate representa los campos actualizables de una propiedad
// Usa punteros para permitir actualizaciones parciales (solo campos no nil se actualizan)
type PropertyUpdate struct {
	Title     *string   `json:"title,omitempty" bson:"title,omitempty"`
	Description *string `json:"description,omitempty" bson:"description,omitempty"`
	Price     *float64  `json:"price,omitempty" bson:"price,omitempty"`
	Location  *string   `json:"location,omitempty" bson:"location,omitempty"`
	Amenities *[]string `json:"amenities,omitempty" bson:"amenities,omitempty"`
	Capacity  *int      `json:"capacity,omitempty" bson:"capacity,omitempty"`
	Available *bool     `json:"available,omitempty" bson:"available,omitempty"`
	UpdatedAt string    `json:"updatedAt" bson:"updatedAt"`
}
