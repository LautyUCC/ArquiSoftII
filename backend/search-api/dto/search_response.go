package dto

import "search-api/domain"

// SearchResponse representa la respuesta de una búsqueda de propiedades
// Incluye los resultados y la información de paginación
type SearchResponse struct {
	// Results es el array de propiedades encontradas
	Results []domain.Property `json:"results"`

	// TotalResults es el total de resultados que coinciden con la búsqueda
	TotalResults int `json:"totalResults"`

	// Page es la página actual de resultados
	Page int `json:"page"`

	// PageSize es el tamaño de página utilizado
	PageSize int `json:"pageSize"`

	// TotalPages es el total de páginas disponibles
	TotalPages int `json:"totalPages"`
}

// ErrorResponse representa una respuesta de error
// Se usa para devolver errores estructurados en las respuestas HTTP
type ErrorResponse struct {
	// Error es el mensaje de error descriptivo
	Error string `json:"error"`

	// Code es el código de error HTTP o código de error personalizado
	Code int `json:"code"`
}

