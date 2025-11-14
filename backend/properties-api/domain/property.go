package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Property representa una propiedad disponible para reserva (tipo Airbnb)
type Property struct {
	// ID es el identificador único de MongoDB
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// Title es el título o nombre de la propiedad
	Title string `bson:"title" json:"title"`
	// Description contiene la descripción detallada de la propiedad
	Description string `bson:"description" json:"description"`
	// Location es la ubicación completa de la propiedad
	Location string `bson:"location" json:"location"`
	// Price es el precio por noche de la propiedad
	Price float64 `bson:"price" json:"price"`
	// Capacity es la cantidad máxima de huéspedes que puede alojar la propiedad
	Capacity int `bson:"capacity" json:"capacity"`
	// Amenities son las comodidades de la propiedad
	Amenities []string `bson:"amenities" json:"amenities"`
	// Images es una lista de URLs de imágenes de la propiedad
	Images []string `bson:"images" json:"images"`
	// OwnerID es el identificador del usuario propietario de la propiedad
	OwnerID string `bson:"ownerId" json:"ownerId"`
	// Available indica si la propiedad está disponible para reserva
	Available bool `bson:"available" json:"available"`
	// CreatedAt es la fecha y hora de creación del registro
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	// UpdatedAt es la fecha y hora de última actualización
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PropertyID string             `bson:"propertyId" json:"propertyId"`
	UserID     string             `bson:"userId" json:"userId"`
	CheckIn    time.Time          `bson:"checkIn" json:"checkIn"`
	CheckOut   time.Time          `bson:"checkOut" json:"checkOut"`
	TotalPrice float64            `bson:"totalPrice" json:"totalPrice"`
	Status     string             `bson:"status" json:"status"` // "pending", "confirmed", "cancelled"
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
}

// PropertyUpdate representa los campos actualizables de una propiedad
type PropertyUpdate struct {
	Title       *string   `json:"title,omitempty" bson:"title,omitempty"`
	Description *string   `json:"description,omitempty" bson:"description,omitempty"`
	Price       *float64  `json:"price,omitempty" bson:"price,omitempty"`
	Location    *string   `json:"location,omitempty" bson:"location,omitempty"`
	Amenities   *[]string `json:"amenities,omitempty" bson:"amenities,omitempty"`
	Capacity    *int      `json:"capacity,omitempty" bson:"capacity,omitempty"`
	Available   *bool     `json:"available,omitempty" bson:"available,omitempty"`
	UpdatedAt   time.Time `bson:"updatedAt" json:"updatedAt"`
}
