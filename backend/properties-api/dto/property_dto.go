package dto

// Este archivo contiene DTOs antiguos que no coinciden con el dominio actual
// Se mantiene por compatibilidad pero no se usa en el servicio actual
// Los DTOs nuevos están en property_create_dto.go, property_update_dto.go y response_dto.go

// CreatePropertyRequest representa el DTO para crear una propiedad (VERSIÓN ANTIGUA - NO USAR)
// type CreatePropertyRequest struct {
// 	Title       string   `json:"title" binding:"required"`
// 	Description string   `json:"description" binding:"required"`
// 	Address     string   `json:"address" binding:"required"`
// 	City        string   `json:"city" binding:"required"`
// 	State       string   `json:"state" binding:"required"`
// 	Country     string   `json:"country" binding:"required"`
// 	ZipCode     string   `json:"zipCode" binding:"required"`
// 	Price       float64  `json:"price" binding:"required,gt=0"`
// 	Bedrooms    int      `json:"bedrooms" binding:"required,gte=0"`
// 	Bathrooms   int      `json:"bathrooms" binding:"required,gte=0"`
// 	Area        float64  `json:"area" binding:"required,gt=0"`
// 	Type        string   `json:"type" binding:"required"`
// 	Status      string   `json:"status" binding:"required"`
// 	OwnerID     string   `json:"ownerId" binding:"required"`
// 	Images      []string `json:"images"`
// }

// UpdatePropertyRequest representa el DTO para actualizar una propiedad (VERSIÓN ANTIGUA - NO USAR)
// type UpdatePropertyRequest struct {
// 	Title       *string   `json:"title,omitempty"`
// 	Description *string   `json:"description,omitempty"`
// 	Address     *string   `json:"address,omitempty"`
// 	City        *string   `json:"city,omitempty"`
// 	State       *string   `json:"state,omitempty"`
// 	Country     *string   `json:"country,omitempty"`
// 	ZipCode     *string   `json:"zipCode,omitempty"`
// 	Price       *float64  `json:"price,omitempty"`
// 	Bedrooms    *int      `json:"bedrooms,omitempty"`
// 	Bathrooms   *int      `json:"bathrooms,omitempty"`
// 	Area        *float64  `json:"area,omitempty"`
// 	Type        *string   `json:"type,omitempty"`
// 	Status      *string   `json:"status,omitempty"`
// 	Images      *[]string `json:"images,omitempty"`
// }

// PropertyResponse representa el DTO de respuesta de una propiedad (VERSIÓN ANTIGUA - NO USAR)
// type PropertyResponse struct {
// 	ID          string   `json:"id"`
// 	Title       string   `json:"title"`
// 	Description string   `json:"description"`
// 	Address     string   `json:"address"`
// 	City        string   `json:"city"`
// 	State       string   `json:"state"`
// 	Country     string   `json:"country"`
// 	ZipCode     string   `json:"zipCode"`
// 	Price       float64  `json:"price"`
// 	Bedrooms    int      `json:"bedrooms"`
// 	Bathrooms   int      `json:"bathrooms"`
// 	Area        float64  `json:"area"`
// 	Type        string   `json:"type"`
// 	Status      string   `json:"status"`
// 	OwnerID     string   `json:"ownerId"`
// 	Images      []string `json:"images"`
// 	CreatedAt   string   `json:"createdAt"`
// 	UpdatedAt   string   `json:"updatedAt"`
// }
