package dto

// SearchRequest representa los parámetros de búsqueda y filtrado de propiedades
// Se usa para recibir query parameters desde las peticiones HTTP
type SearchRequest struct {
	// Query es el término de búsqueda general
	Query string `json:"query" form:"query"`

	// City es un filtro opcional por ciudad
	City string `json:"city" form:"city"`

	// Country es un filtro opcional por país
	Country string `json:"country" form:"country"`

	// MinPrice es el precio mínimo por noche
	MinPrice float64 `json:"minPrice" form:"minPrice"`

	// MaxPrice es el precio máximo por noche
	MaxPrice float64 `json:"maxPrice" form:"maxPrice"`

	// Bedrooms es el número de habitaciones requerido
	Bedrooms int `json:"bedrooms" form:"bedrooms"`

	// Bathrooms es el número de baños requerido
	Bathrooms int `json:"bathrooms" form:"bathrooms"`

	// MinGuests es la capacidad mínima de huéspedes
	MinGuests int `json:"minGuests" form:"minGuests"`

	// Page es el número de página para paginación (default: 1)
	Page int `json:"page" form:"page"`

	// PageSize es el tamaño de página para paginación (default: 10)
	PageSize int `json:"pageSize" form:"pageSize"`

	// SortBy es el campo para ordenar los resultados (default: "price_per_night")
	// Opciones comunes: "price_per_night", "created_at", "bedrooms", etc.
	SortBy string `json:"sortBy" form:"sortBy"`

	// SortOrder es el orden de clasificación: "asc" o "desc" (default: "asc")
	SortOrder string `json:"sortOrder" form:"sortOrder"`
}

