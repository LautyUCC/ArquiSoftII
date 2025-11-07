package dto

// PaginatedResponse representa una respuesta paginada
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"totalPages"`
}

// PropertyStatsResponse representa el DTO de respuesta de estadísticas
type PropertyStatsResponse struct {
	PropertyID   string  `json:"propertyId"`
	PricePerM2   float64 `json:"pricePerM2"`
	TotalValue   float64 `json:"totalValue"`
	RoomRatio    float64 `json:"roomRatio"`
	CalculatedAt string  `json:"calculatedAt"`
}

// FromDomainStats convierte estadísticas del dominio a DTO de respuesta
// Nota: Esta función está comentada porque domain.PropertyStats no existe en el dominio actual
// Si se necesita estadísticas, se debe agregar al dominio primero
// func FromDomainStats(stats *domain.PropertyStats) *PropertyStatsResponse {
// 	return &PropertyStatsResponse{
// 		PropertyID:   stats.PropertyID,
// 		PricePerM2:   stats.PricePerM2,
// 		TotalValue:   stats.TotalValue,
// 		RoomRatio:    stats.RoomRatio,
// 		CalculatedAt: stats.CalculatedAt.Format("2006-01-02T15:04:05Z07:00"),
// 	}
// }

// PropertyResponseDTO representa el DTO de respuesta de una propiedad
// Usado para serializar la respuesta HTTP
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
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

// ErrorResponse representa un error en la respuesta
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

