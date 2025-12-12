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
	solrRepo         repositories.SolrRepository
	cacheRepo        repositories.CacheRepository
	propertiesAPIURL string
	httpClient       *http.Client
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
	// PROPERTIES_API_URL puede incluir o no /api, así que lo manejamos
	baseURL := s.propertiesAPIURL
	if !strings.HasSuffix(baseURL, "/api") && !strings.HasSuffix(baseURL, "/api/") {
		baseURL = strings.TrimSuffix(baseURL, "/") + "/api"
	}
	url := fmt.Sprintf("%s/properties/%s", baseURL, propertyID)
	log.Printf("🌐 URL construida: %s", url)

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
	// Puede venir envuelta en un objeto con "data" o directamente
	var wrappedResponse struct {
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
			Images      []string `json:"images"`
		} `json:"data"`
		// También puede venir directamente sin wrapper
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
	}

	if err := json.Unmarshal(body, &wrappedResponse); err != nil {
		log.Printf("❌ Error parseando JSON: %v", err)
		log.Printf("📄 Body completo: %s", string(body))
		return nil, fmt.Errorf("error parseando respuesta: %v", err)
	}

	// Determinar si la respuesta viene envuelta en "data" o directamente
	var apiResponse struct {
		ID          string
		Title       string
		Description string
		Price       float64
		Location    string
		OwnerID     string
		Amenities   []string
		Capacity    int
		Available   bool
		Images      []string
	}

	if wrappedResponse.Data.ID != "" {
		// Respuesta envuelta en "data"
		apiResponse.ID = wrappedResponse.Data.ID
		apiResponse.Title = wrappedResponse.Data.Title
		apiResponse.Description = wrappedResponse.Data.Description
		apiResponse.Price = wrappedResponse.Data.Price
		apiResponse.Location = wrappedResponse.Data.Location
		apiResponse.OwnerID = wrappedResponse.Data.OwnerID
		apiResponse.Amenities = wrappedResponse.Data.Amenities
		apiResponse.Capacity = wrappedResponse.Data.Capacity
		apiResponse.Available = wrappedResponse.Data.Available
		apiResponse.Images = wrappedResponse.Data.Images
	} else {
		// Respuesta directa
		apiResponse.ID = wrappedResponse.ID
		apiResponse.Title = wrappedResponse.Title
		apiResponse.Description = wrappedResponse.Description
		apiResponse.Price = wrappedResponse.Price
		apiResponse.Location = wrappedResponse.Location
		apiResponse.OwnerID = wrappedResponse.OwnerID
		apiResponse.Amenities = wrappedResponse.Amenities
		apiResponse.Capacity = wrappedResponse.Capacity
		apiResponse.Available = wrappedResponse.Available
		apiResponse.Images = wrappedResponse.Images
	}

	// LOG para debug
	log.Printf("🔍 ID parseado desde JSON: '%s'", apiResponse.ID)
	log.Printf("🔍 Title parseado: '%s'", apiResponse.Title)
	log.Printf("🖼️ Images parseado: %v (len: %d)", apiResponse.Images, len(apiResponse.Images))

	// Validar que el ID no esté vacío
	if apiResponse.ID == "" {
		log.Printf("❌ ERROR: ID está vacío después del parseo")
		log.Printf("📄 Body completo: %s", string(body))
		return nil, fmt.Errorf("la API devolvió una propiedad sin ID")
	}

	// Parsear CreatedAt de string a time.Time
	var createdAt time.Time
	createdAt = time.Now()

	// Convertir OwnerID de string a uint
	var ownerID uint
	if apiResponse.OwnerID != "" {
		hash := md5.Sum([]byte(apiResponse.OwnerID))
		ownerID = uint(hash[0]) | uint(hash[1])<<8 | uint(hash[2])<<16 | uint(hash[3])<<24
	}

	// Extraer city y country de location
	city := ""
	country := ""
	if apiResponse.Location != "" {
		parts := strings.Split(apiResponse.Location, ",")
		if len(parts) >= 1 {
			city = strings.TrimSpace(parts[0])
		}
		if len(parts) >= 2 {
			country = strings.TrimSpace(parts[1])
		}
	}

	// Mapear cada campo
	property := &domain.Property{
		ID:            apiResponse.ID,
		Title:         apiResponse.Title,
		Description:   apiResponse.Description,
		City:          city,
		Country:       country,
		PricePerNight: apiResponse.Price,
		Bedrooms:      0,
		Bathrooms:     0,
		MaxGuests:     apiResponse.Capacity,
		Images:        apiResponse.Images, // ✅ CORRECTO
		OwnerID:       ownerID,
		Available:     apiResponse.Available,
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
		Results:      properties,
		TotalResults: total,
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
	}
}

// invalidateCache invalida el caché eliminando todas las keys relacionadas
// Limpia todo el caché local para forzar nuevas búsquedas en Solr
func (s *searchService) invalidateCache() {
	log.Println("🔄 Invalidando caché de búsquedas")
	s.cacheRepo.InvalidateAll()
	log.Println("✅ Caché invalidado exitosamente")
}
