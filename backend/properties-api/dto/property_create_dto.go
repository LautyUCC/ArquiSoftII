package dto

// PropertyCreateDTO representa el DTO para crear una propiedad
// Usado en el servicio para recibir datos de creación
type PropertyCreateDTO struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	Location    string   `json:"location" binding:"required"`
	OwnerID     string   `json:"ownerId" binding:"required"`
	Amenities   []string `json:"amenities"`
	Capacity    int      `json:"capacity" binding:"required,gt=0"`
	Available   bool     `json:"available"`
}

