package dto

// PropertyUpdateDTO representa el DTO para actualizar una propiedad
// Usa punteros para permitir actualizaciones parciales (solo campos no nil se actualizan)
type PropertyUpdateDTO struct {
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Price       *float64  `json:"price,omitempty"`
	Location    *string   `json:"location,omitempty"`
	Amenities   *[]string `json:"amenities,omitempty"`
	Capacity    *int      `json:"capacity,omitempty"`
	Available   *bool     `json:"available,omitempty"`
}

