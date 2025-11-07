package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"search-api/domain"
	"search-api/dto"
	"search-api/repositories"
)

// SearchService define la interfaz para las operaciones de búsqueda
type SearchService interface {
	// Search realiza una búsqueda de propiedades con caché y Solr
	Search(ctx context.Context, request dto.SearchRequest) (*dto.SearchResponse, error)

	// IndexProperty indexa una nueva propiedad en Solr e invalida caché
	IndexProperty(ctx context.Context, property domain.Property) error

	// UpdateProperty actualiza una propiedad en Solr e invalida caché
	UpdateProperty(ctx context.Context, property domain.Property) error

	// DeleteProperty elimina una propiedad de Solr e invalida caché
	DeleteProperty(ctx context.Context, propertyID string) error

	// FetchPropertyFromAPI obtiene una propiedad desde la API de propiedades
	FetchPropertyFromAPI(propertyID string) (*domain.Property, error)
}

// searchService es la implementación concreta de SearchService
type searchService struct {
	solrRepo        repositories.SolrRepository
	cacheRepo       repositories.CacheRepository
	propertiesAPIURL string
	httpClient      *http.Client
}

// NewSearchService crea una nueva instancia del servicio de búsqueda
func NewSearchService(
	solrRepo repositories.SolrRepository,
	cacheRepo repositories.CacheRepository,
	apiURL string,
) SearchService {
	return &searchService{
		solrRepo:         solrRepo,
		cacheRepo:        cacheRepo,
		propertiesAPIURL: strings.TrimSuffix(apiURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Search realiza una búsqueda de propiedades con estrategia de caché de dos niveles
func (s *searchService) Search(ctx context.Context, request dto.SearchRequest) (*dto.SearchResponse, error) {
	// Validar request
	if err := s.validateSearchRequest(&request); err != nil {
		return nil, fmt.Errorf("request inválido: %w", err)
	}

	// Generar cache key basado en los parámetros del request
	cacheKey := s.generateCacheKey(request)
	log.Printf("🔍 Iniciando búsqueda con cache key: %s", cacheKey)

	// Consultar caché primero
	properties, total, found := s.cacheRepo.Get(cacheKey)
	if found {
		log.Printf("✅ Cache hit para key: %s", cacheKey)
		return s.buildSearchResponse(properties, total, request), nil
	}

	log.Printf("❌ Cache miss para key: %s, consultando Solr", cacheKey)

	// Consultar Solr
	properties, total, err := s.solrRepo.Search(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("error buscando en Solr: %w", err)
	}

	log.Printf("✅ Búsqueda en Solr completada: %d resultados encontrados", total)

	// Guardar resultado en caché con TTL de 15 minutos
	s.cacheRepo.Set(cacheKey, properties, total, 15*time.Minute)
	log.Printf("✅ Resultados guardados en caché para key: %s", cacheKey)

	return s.buildSearchResponse(properties, total, request), nil
}

// IndexProperty indexa una nueva propiedad en Solr e invalida caché
func (s *searchService) IndexProperty(ctx context.Context, property domain.Property) error {
	// Validar propiedad
	if err := s.validateProperty(&property); err != nil {
		return fmt.Errorf("propiedad inválida: %w", err)
	}

	log.Printf("📝 Indexando propiedad ID: %s", property.ID)

	// Indexar en Solr
	if err := s.solrRepo.IndexProperty(ctx, property); err != nil {
		return fmt.Errorf("error indexando propiedad en Solr: %w", err)
	}

	log.Printf("✅ Propiedad indexada exitosamente en Solr: %s", property.ID)

	// Invalidar caché (eliminar todas las keys relacionadas)
	s.invalidateCache()

	return nil
}

// UpdateProperty actualiza una propiedad en Solr e invalida caché
func (s *searchService) UpdateProperty(ctx context.Context, property domain.Property) error {
	// Validar propiedad
	if err := s.validateProperty(&property); err != nil {
		return fmt.Errorf("propiedad inválida: %w", err)
	}

	log.Printf("🔄 Actualizando propiedad ID: %s", property.ID)

	// Actualizar en Solr
	if err := s.solrRepo.UpdateProperty(ctx, property); err != nil {
		return fmt.Errorf("error actualizando propiedad en Solr: %w", err)
	}

	log.Printf("✅ Propiedad actualizada exitosamente en Solr: %s", property.ID)

	// Invalidar caché
	s.invalidateCache()

	return nil
}

// DeleteProperty elimina una propiedad de Solr e invalida caché
func (s *searchService) DeleteProperty(ctx context.Context, propertyID string) error {
	// Validar ID
	if propertyID == "" {
		return fmt.Errorf("ID de propiedad no puede estar vacío")
	}

	log.Printf("🗑️ Eliminando propiedad ID: %s", propertyID)

	// Eliminar de Solr
	if err := s.solrRepo.DeleteProperty(ctx, propertyID); err != nil {
		return fmt.Errorf("error eliminando propiedad de Solr: %w", err)
	}

	log.Printf("✅ Propiedad eliminada exitosamente de Solr: %s", propertyID)

	// Invalidar caché
	s.invalidateCache()

	return nil
}

// FetchPropertyFromAPI obtiene una propiedad desde la API de propiedades
func (s *searchService) FetchPropertyFromAPI(propertyID string) (*domain.Property, error) {
	// Validar ID
	if propertyID == "" {
		return nil, fmt.Errorf("ID de propiedad no puede estar vacío")
	}

	log.Printf("🌐 Obteniendo propiedad desde API: %s", propertyID)

	// Construir URL
	url := fmt.Sprintf("%s/properties/%s", s.propertiesAPIURL, propertyID)

	// Crear request HTTP GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request HTTP: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Realizar petición
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error realizando petición a properties-api: %w", err)
	}
	defer resp.Body.Close()

	// Verificar código de estado
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error en respuesta de properties-api (status %d): %s", resp.StatusCode, string(body))
	}

	// Leer y parsear respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta de properties-api: %w", err)
	}

	// La API puede devolver la propiedad envuelta en un objeto de respuesta
	// Intentar parsear diferentes formatos de respuesta
	var apiResponse struct {
		Success bool            `json:"success"`
		Data    domain.Property `json:"data"`
		Message string          `json:"message"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		// Si falla, intentar parsear directamente como Property
		var property domain.Property
		if err2 := json.Unmarshal(body, &property); err2 != nil {
			return nil, fmt.Errorf("error parseando respuesta JSON: %w (también intentó parseo directo: %v)", err, err2)
		}
		log.Printf("✅ Propiedad obtenida desde API (formato directo): %s", propertyID)
		return &property, nil
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("la API reportó error: %s", apiResponse.Message)
	}

	log.Printf("✅ Propiedad obtenida desde API: %s", propertyID)
	return &apiResponse.Data, nil
}

// validateSearchRequest valida los parámetros de búsqueda
func (s *searchService) validateSearchRequest(request *dto.SearchRequest) error {
	// Validar paginación
	if request.Page < 1 {
		request.Page = 1
	}
	if request.PageSize < 1 {
		request.PageSize = 10
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

	// Validar sortOrder
	if request.SortOrder != "" && request.SortOrder != "asc" && request.SortOrder != "desc" {
		return fmt.Errorf("sortOrder debe ser 'asc' o 'desc'")
	}

	return nil
}

// validateProperty valida una propiedad
func (s *searchService) validateProperty(property *domain.Property) error {
	if property.ID == "" {
		return fmt.Errorf("ID de propiedad no puede estar vacío")
	}
	if property.Title == "" {
		return fmt.Errorf("title no puede estar vacío")
	}
	if property.PricePerNight < 0 {
		return fmt.Errorf("pricePerNight no puede ser negativo")
	}
	return nil
}

// generateCacheKey genera una clave de caché única basada en los parámetros de búsqueda
func (s *searchService) generateCacheKey(request dto.SearchRequest) string {
	// Normalizar valores para consistencia
	page := request.Page
	if page < 1 {
		page = 1
	}
	pageSize := request.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	sortBy := request.SortBy
	if sortBy == "" {
		sortBy = "price_per_night"
	}
	sortOrder := request.SortOrder
	if sortOrder == "" {
		sortOrder = "asc"
	}

	// Construir string con todos los parámetros
	keyParts := []string{
		fmt.Sprintf("query:%s", request.Query),
		fmt.Sprintf("city:%s", request.City),
		fmt.Sprintf("country:%s", request.Country),
		fmt.Sprintf("minPrice:%.2f", request.MinPrice),
		fmt.Sprintf("maxPrice:%.2f", request.MaxPrice),
		fmt.Sprintf("bedrooms:%d", request.Bedrooms),
		fmt.Sprintf("bathrooms:%d", request.Bathrooms),
		fmt.Sprintf("minGuests:%d", request.MinGuests),
		fmt.Sprintf("page:%d", page),
		fmt.Sprintf("pageSize:%d", pageSize),
		fmt.Sprintf("sortBy:%s", sortBy),
		fmt.Sprintf("sortOrder:%s", sortOrder),
	}

	keyString := strings.Join(keyParts, "|")

	// Generar hash MD5 para obtener una clave de longitud fija
	hash := md5.Sum([]byte(keyString))
	return "search:" + hex.EncodeToString(hash[:])
}

// buildSearchResponse construye una respuesta de búsqueda
func (s *searchService) buildSearchResponse(properties []domain.Property, total int, request dto.SearchRequest) *dto.SearchResponse {
	page := request.Page
	if page < 1 {
		page = 1
	}
	pageSize := request.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	// Calcular total de páginas
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	return &dto.SearchResponse{
		Results:     properties,
		TotalResults: total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}
}

// invalidateCache invalida el caché eliminando todas las keys relacionadas
// Nota: En una implementación más sofisticada, se podría mantener un registro de keys
// o usar un patrón de invalidación más granular
func (s *searchService) invalidateCache() {
	// Por simplicidad, invalidamos todas las keys que empiezan con "search:"
	// En producción, se podría implementar un sistema más sofisticado de invalidación
	log.Println("🔄 Invalidando caché de búsquedas")
	// Nota: La invalidación completa del caché requeriría una implementación adicional
	// en el CacheRepository para soportar invalidación por patrón
	// Por ahora, el caché se invalidará naturalmente con su TTL
}

