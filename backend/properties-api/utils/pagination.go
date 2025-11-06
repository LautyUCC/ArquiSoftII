package utils

import (
	"math"
	"strconv"
)

// PaginationParams representa los parámetros de paginación
type PaginationParams struct {
	Page  int
	Limit int
}

// ParsePaginationParams parsea los parámetros de paginación desde query params
func ParsePaginationParams(pageStr, limitStr string, defaultLimit int) PaginationParams {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = defaultLimit
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// CalculateTotalPages calcula el número total de páginas
func CalculateTotalPages(total int64, limit int) int {
	if total == 0 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}

// CalculateSkip calcula el offset para la paginación
func CalculateSkip(page, limit int) int {
	return (page - 1) * limit
}

