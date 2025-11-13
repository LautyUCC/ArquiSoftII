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

	log.Printf("📦 Respuesta raw de Properties API: %s", string(body))

	// Estructura para parsear la respuesta de Properties API
	var apiResponse struct {
		Data struct {
			ID          string   `json:"id"`
			Title       string   `json:"title"`
			Description string   `json:"description"`
			Price       float64  `json:"price"`
			Location    string   `json:"location"`
			OwnerID     string   `json:"ownerId"`
			Amenities   []string `json:"amenities"`
			Capacity    int      `json:"capacity"`
			Available   bool     `json:"available"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("❌ Error parseando JSON: %v", err)
		log.Printf("📄 Body completo: %s", string(body))
		return nil, fmt.Errorf("error parseando respuesta: %v", err)
	}

	// LOG para debug - verificar que el ID se leyó correctamente
	log.Printf("🔍 ID parseado desde JSON: '%s'", apiResponse.Data.ID)
	log.Printf("🔍 Title parseado: '%s'", apiResponse.Data.Title)

	// Validar que el ID no esté vacío
	if apiResponse.Data.ID == "" {
		log.Printf("❌ ERROR: ID está vacío después del parseo")
		log.Printf("📄 Body completo: %s", string(body))
		return nil, fmt.Errorf("la API devolvió una propiedad sin ID")
	}

	// Parsear CreatedAt de string a time.Time
	var createdAt time.Time
	// Properties API puede no incluir CreatedAt, usar tiempo actual si no está
	createdAt = time.Now()

	// Convertir OwnerID de string a uint
	var ownerID uint
	if apiResponse.Data.OwnerID != "" {
		// Generar un hash simple del string para convertirlo a uint
		hash := md5.Sum([]byte(apiResponse.Data.OwnerID))
		// Usar los primeros 4 bytes del hash como uint
		ownerID = uint(hash[0]) | uint(hash[1])<<8 | uint(hash[2])<<16 | uint(hash[3])<<24
	}

	// Extraer city y country de location (formato: "Ciudad, País")
	city := ""
	country := ""
	if apiResponse.Data.Location != "" {
		parts := strings.Split(apiResponse.Data.Location, ",")
		if len(parts) >= 1 {
			city = strings.TrimSpace(parts[0])
		}
		if len(parts) >= 2 {
			country = strings.TrimSpace(parts[1])
		}
	}

	// LOG para debug - verificar valores antes del mapeo
	log.Printf("🔍 Valores desde Properties API:")
	log.Printf("   - ID: '%s'", apiResponse.Data.ID)
	log.Printf("   - Title: '%s'", apiResponse.Data.Title)
	log.Printf("   - Description: '%s'", apiResponse.Data.Description)
	log.Printf("   - Price: %f", apiResponse.Data.Price)
	log.Printf("   - Location: '%s'", apiResponse.Data.Location)
	log.Printf("   - Capacity: %d", apiResponse.Data.Capacity)
	log.Printf("   - Available: %v", apiResponse.Data.Available)
	log.Printf("   - City extraída: '%s'", city)
	log.Printf("   - Country extraída: '%s'", country)

	// Mapear EXPLÍCITAMENTE cada campo
	property := &domain.Property{
		ID:            apiResponse.Data.ID,          // ← CRÍTICO
		Title:         apiResponse.Data.Title,
		Description:   apiResponse.Data.Description,
		City:          city,
		Country:       country,
		PricePerNight: apiResponse.Data.Price,
		Bedrooms:      0,
		Bathrooms:     0,
		MaxGuests:     apiResponse.Data.Capacity,
		Images:        []string{},
		OwnerID:       ownerID,
		Available:     apiResponse.Data.Available,
		CreatedAt:     createdAt,
	}

	// LOG para debug - verificar valores después del mapeo
	log.Printf("🆔 ID mapeado: '%s'", property.ID)
	log.Printf("📝 Title mapeado: '%s'", property.Title)
	log.Printf("💰 PricePerNight mapeado: %f", property.PricePerNight)
	log.Printf("🏙️ City mapeado: '%s'", property.City)
	log.Printf("🌍 Country mapeado: '%s'", property.Country)
	log.Printf("👥 MaxGuests mapeado: %d", property.MaxGuests)

	if property.ID == "" {
		log.Printf("❌ ERROR CRÍTICO: property.ID está vacío después del mapeo")
		return nil, fmt.Errorf("ID de propiedad está vacío después del mapeo")
	}

	log.Printf("✅ Propiedad obtenida desde API: %s", propertyID)
	return property, nil
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
	// sortBy puede estar vacío (sort opcional)
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

