package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"search-api/dto"
	"search-api/services"
	"search-api/utils"
)

// SearchController maneja las peticiones HTTP relacionadas con búsqueda
type SearchController struct {
	service services.SearchService
}

// NewSearchController crea una nueva instancia del controlador de búsqueda
func NewSearchController(service services.SearchService) *SearchController {
	return &SearchController{
		service: service,
	}
}

// Search maneja GET /search
// Parsea query parameters, valida, llama al servicio y retorna resultados
func (c *SearchController) Search(w http.ResponseWriter, r *http.Request) {
	// Solo permitir método GET
	if r.Method != http.MethodGet {
		utils.SendErrorSimple(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parsear query parameters a SearchRequest
	request, err := parseSearchRequest(r)
	if err != nil {
		utils.LogError("GET /search", err, r)
		utils.SendErrorSimple(w, http.StatusBadRequest, fmt.Sprintf("Error parseando parámetros: %v", err))
		return
	}

	// Validar parámetros
	if err := validateSearchRequest(request); err != nil {
		utils.LogError("GET /search", err, r)
		utils.SendErrorSimple(w, http.StatusBadRequest, err.Error())
		return
	}

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Llamar al servicio
	response, err := c.service.Search(ctx, *request)
	if err != nil {
		utils.LogError("GET /search", err, r)
		utils.SendErrorSimple(w, http.StatusInternalServerError, fmt.Sprintf("Error en búsqueda: %v", err))
		return
	}

	// Escribir respuesta exitosa
	utils.SendSuccess(w, http.StatusOK, response)
	utils.LogInfo("GET /search", fmt.Sprintf("Búsqueda completada exitosamente: %d resultados", response.TotalResults), r)
}

// parseSearchRequest parsea los query parameters a SearchRequest
func parseSearchRequest(r *http.Request) (*dto.SearchRequest, error) {
	request := &dto.SearchRequest{}

	// Obtener query parameters
	query := r.URL.Query()

	// Query (término de búsqueda)
	request.Query = query.Get("query")

	// City
	request.City = query.Get("city")

	// Country
	request.Country = query.Get("country")

	// MinPrice
	if minPriceStr := query.Get("minPrice"); minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			return nil, fmt.Errorf("minPrice debe ser un número válido: %w", err)
		}
		request.MinPrice = minPrice
	}

	// MaxPrice
	if maxPriceStr := query.Get("maxPrice"); maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			return nil, fmt.Errorf("maxPrice debe ser un número válido: %w", err)
		}
		request.MaxPrice = maxPrice
	}

	// Bedrooms
	if bedroomsStr := query.Get("bedrooms"); bedroomsStr != "" {
		bedrooms, err := strconv.Atoi(bedroomsStr)
		if err != nil {
			return nil, fmt.Errorf("bedrooms debe ser un número entero válido: %w", err)
		}
		request.Bedrooms = bedrooms
	}

	// Bathrooms
	if bathroomsStr := query.Get("bathrooms"); bathroomsStr != "" {
		bathrooms, err := strconv.Atoi(bathroomsStr)
		if err != nil {
			return nil, fmt.Errorf("bathrooms debe ser un número entero válido: %w", err)
		}
		request.Bathrooms = bathrooms
	}

	// MinGuests
	if minGuestsStr := query.Get("minGuests"); minGuestsStr != "" {
		minGuests, err := strconv.Atoi(minGuestsStr)
		if err != nil {
			return nil, fmt.Errorf("minGuests debe ser un número entero válido: %w", err)
		}
		request.MinGuests = minGuests
	}

	// Page
	if pageStr := query.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, fmt.Errorf("page debe ser un número entero válido: %w", err)
		}
		request.Page = page
	} else {
		request.Page = 1 // Default
	}

	// PageSize
	if pageSizeStr := query.Get("pageSize"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			return nil, fmt.Errorf("pageSize debe ser un número entero válido: %w", err)
		}
		request.PageSize = pageSize
	} else {
		request.PageSize = 10 // Default
	}

	// SortBy (opcional - sin default)
	request.SortBy = query.Get("sortBy")
	// Si no se especifica, se deja vacío (sort opcional)

	// SortOrder (opcional - solo se usa si SortBy está especificado)
	request.SortOrder = strings.ToLower(query.Get("sortOrder"))
	if request.SortOrder == "" {
		request.SortOrder = "asc" // Default solo si SortBy está especificado
	}

	return request, nil
}

// validateSearchRequest valida los parámetros de búsqueda
func validateSearchRequest(request *dto.SearchRequest) error {
	// Validar Page
	if request.Page < 1 {
		return fmt.Errorf("page debe ser mayor o igual a 1")
	}

	// Validar PageSize
	if request.PageSize <= 0 {
		return fmt.Errorf("pageSize debe ser mayor a 0")
	}
	if request.PageSize > 100 {
		return fmt.Errorf("pageSize no puede ser mayor a 100")
	}

	// Validar rango de precio
	if request.MinPrice < 0 {
		return fmt.Errorf("minPrice no puede ser negativo")
	}
	if request.MaxPrice < 0 {
		return fmt.Errorf("maxPrice no puede ser negativo")
	}
	if request.MinPrice > 0 && request.MaxPrice > 0 && request.MinPrice > request.MaxPrice {
		return fmt.Errorf("minPrice no puede ser mayor que maxPrice")
	}

	// Validar SortOrder
	if request.SortOrder != "" && request.SortOrder != "asc" && request.SortOrder != "desc" {
		return fmt.Errorf("sortOrder debe ser 'asc' o 'desc'")
	}

	return nil
}

// Las funciones writeJSONResponse y writeErrorResponse han sido reemplazadas
// por utils.SendSuccess y utils.SendError

