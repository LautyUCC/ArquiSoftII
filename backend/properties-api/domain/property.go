package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

import (
	"time"
)

// Property representa una propiedad disponible para reserva (tipo Airbnb)
// Contiene toda la información necesaria para listar y reservar propiedades
type Property struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Location    string             `bson:"location" json:"location"`
	OwnerID     string             `bson:"ownerId" json:"ownerId"`
	Amenities   []string           `bson:"amenities" json:"amenities"`
	Capacity    int                `bson:"capacity" json:"capacity"`
	Available   bool               `bson:"available" json:"available"`
	Images      []string           `bson:"images" json:"images"` // ← AGREGAR ESTO
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
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
// Usa punteros para permitir actualizaciones parciales (solo campos no nil se actualizan)
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
