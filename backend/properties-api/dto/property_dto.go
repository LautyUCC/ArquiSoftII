package dto

// PropertyCreateDTO representa el DTO para crear una propiedad
type PropertyCreateDTO struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	Location    string   `json:"location" binding:"required"`
	OwnerID     string   `json:"ownerId" binding:"required"`
	Amenities   []string `json:"amenities"`
	Capacity    int      `json:"capacity" binding:"required,gte=1"`
	Available   bool     `json:"available"`
	Images      []string `json:"images"`
}

// PropertyUpdateDTO representa el DTO para actualizar una propiedad
// Todos los campos son opcionales (punteros)
type PropertyUpdateDTO struct {
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Price       *float64  `json:"price,omitempty"`
	Location    *string   `json:"location,omitempty"`
	Amenities   *[]string `json:"amenities,omitempty"`
	Capacity    *int      `json:"capacity,omitempty"`
	Available   *bool     `json:"available,omitempty"`
	Images      *[]string `json:"images,omitempty"`
}

// PropertyResponseDTO representa el DTO de respuesta de una propiedad
type PropertyResponseDTO struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Location    string   `json:"location"`
	OwnerID     string   `json:"ownerId"`
	Amenities   []string `json:"amenities"`
	Capacity    int      `json:"capacity"`
	Available   bool     `json:"available"`
	Images      []string `json:"images"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}
