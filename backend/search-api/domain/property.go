package domain

import "time"

// Property representa una propiedad de alquiler tipo Airbnb
type Property struct {
	// ID es el identificador único de la propiedad
	ID string `json:"id"`

	// Title es el título o nombre de la propiedad
	Title string `json:"title"`

	// Description contiene la descripción detallada de la propiedad
	Description string `json:"description"`

	// City es la ciudad donde se encuentra la propiedad
	City string `json:"city"`

	// Country es el país donde se encuentra la propiedad
	Country string `json:"country"`

	// PricePerNight es el precio por noche de la propiedad
	PricePerNight float64 `json:"pricePerNight"`

	// Bedrooms es el número de habitaciones de la propiedad
	Bedrooms int `json:"bedrooms"`

	// Bathrooms es el número de baños de la propiedad
	Bathrooms int `json:"bathrooms"`

	// MaxGuests es la cantidad máxima de huéspedes que puede alojar la propiedad
	MaxGuests int `json:"maxGuests"`

	// Images es una lista de URLs de imágenes de la propiedad
	Images []string `json:"images"`

	// OwnerID es el identificador del usuario propietario de la propiedad
	OwnerID uint `json:"ownerID"`

	// Available indica si la propiedad está disponible para reserva
	Available bool `json:"available"`

	// CreatedAt es la fecha y hora de creación del registro
	CreatedAt time.Time `json:"createdAt"`
}
