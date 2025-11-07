package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"search-api/domain"
	"search-api/dto"
)

// SolrRepository define la interfaz para las operaciones de repositorio de Solr
type SolrRepository interface {
	// Search realiza una búsqueda de propiedades con filtros y paginación
	Search(ctx context.Context, request dto.SearchRequest) ([]domain.Property, int, error)

	// IndexProperty indexa una nueva propiedad en Solr
	IndexProperty(ctx context.Context, property domain.Property) error

	// UpdateProperty actualiza una propiedad existente en Solr
	UpdateProperty(ctx context.Context, property domain.Property) error

	// DeleteProperty elimina una propiedad de Solr por su ID
	DeleteProperty(ctx context.Context, propertyID string) error
}

// solrRepository es la implementación concreta de SolrRepository
type solrRepository struct {
	solrURL    string
	httpClient *http.Client
}

// NewSolrRepository crea una nueva instancia del repositorio de Solr
func NewSolrRepository(solrURL string) SolrRepository {
	return &solrRepository{
		solrURL: solrURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SolrResponse representa la estructura de respuesta de Solr
type SolrResponse struct {
	Response struct {
		NumFound int                      `json:"numFound"`
		Start    int                      `json:"start"`
		Docs     []map[string]interface{} `json:"docs"`
	} `json:"response"`
}

// SolrProperty representa una propiedad en formato Solr
type SolrProperty struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	City          string    `json:"city"`
	Country       string    `json:"country"`
	PricePerNight float64   `json:"price_per_night"`
	Bedrooms      int       `json:"bedrooms"`
	Bathrooms     int       `json:"bathrooms"`
	MaxGuests     int       `json:"max_guests"`
	Images        []string  `json:"images"`
	OwnerID       uint      `json:"owner_id"`
	Available     bool      `json:"available"`
	CreatedAt     time.Time `json:"created_at"`
}

// Search realiza una búsqueda de propiedades con filtros y paginación
func (r *solrRepository) Search(ctx context.Context, request dto.SearchRequest) ([]domain.Property, int, error) {
	// Construir la URL base de búsqueda
	baseURL := strings.TrimSuffix(r.solrURL, "/") + "/select"

	// Construir parámetros de la query
	params := url.Values{}
	params.Set("wt", "json") // Formato de respuesta JSON

	// Construir query de búsqueda por texto
	if request.Query != "" {
		// Búsqueda en title, city, country
		query := fmt.Sprintf("(title:*%s* OR city:*%s* OR country:*%s*)", 
			escapeSolrQuery(request.Query), 
			escapeSolrQuery(request.Query), 
			escapeSolrQuery(request.Query))
		params.Set("q", query)
	} else {
		params.Set("q", "*:*") // Buscar todo si no hay query
	}

	// Construir filtros (fq parameters)
	var filters []string

	// Filtro por ciudad
	if request.City != "" {
		filters = append(filters, fmt.Sprintf("city:\"%s\"", escapeSolrQuery(request.City)))
	}

	// Filtro por país
	if request.Country != "" {
		filters = append(filters, fmt.Sprintf("country:\"%s\"", escapeSolrQuery(request.Country)))
	}

	// Filtro por rango de precio
	if request.MinPrice > 0 || request.MaxPrice > 0 {
		minPrice := request.MinPrice
		maxPrice := request.MaxPrice
		if maxPrice == 0 {
			maxPrice = 999999 // Valor alto si no se especifica máximo
		}
		filters = append(filters, fmt.Sprintf("price_per_night:[%f TO %f]", minPrice, maxPrice))
	}

	// Filtro por número de habitaciones
	if request.Bedrooms > 0 {
		filters = append(filters, fmt.Sprintf("bedrooms:%d", request.Bedrooms))
	}

	// Filtro por número de baños
	if request.Bathrooms > 0 {
		filters = append(filters, fmt.Sprintf("bathrooms:%d", request.Bathrooms))
	}

	// Filtro por capacidad mínima de huéspedes
	if request.MinGuests > 0 {
		filters = append(filters, fmt.Sprintf("max_guests:[%d TO *]", request.MinGuests))
	}

	// Agregar filtros a los parámetros
	for _, filter := range filters {
		params.Add("fq", filter)
	}

	// Paginación
	page := request.Page
	if page < 1 {
		page = 1
	}
	pageSize := request.PageSize
	if pageSize < 1 {
		pageSize = 10
	}
	start := (page - 1) * pageSize
	params.Set("start", strconv.Itoa(start))
	params.Set("rows", strconv.Itoa(pageSize))

	// Ordenamiento
	sortBy := request.SortBy
	if sortBy == "" {
		sortBy = "price_per_night"
	}
	sortOrder := request.SortOrder
	if sortOrder == "" {
		sortOrder = "asc"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	params.Set("sort", fmt.Sprintf("%s %s", sortBy, sortOrder))

	// Construir URL completa
	fullURL := baseURL + "?" + params.Encode()

	// Crear request HTTP
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("error creando request HTTP: %w", err)
	}

	// Realizar petición
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("error realizando petición a Solr: %w", err)
	}
	defer resp.Body.Close()

	// Verificar código de estado
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("error en respuesta de Solr (status %d): %s", resp.StatusCode, string(body))
	}

	// Leer y parsear respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("error leyendo respuesta de Solr: %w", err)
	}

	var solrResp SolrResponse
	if err := json.Unmarshal(body, &solrResp); err != nil {
		return nil, 0, fmt.Errorf("error parseando respuesta JSON de Solr: %w", err)
	}

	// Convertir documentos de Solr a domain.Property
	properties := make([]domain.Property, 0, len(solrResp.Response.Docs))
	for _, doc := range solrResp.Response.Docs {
		property, err := r.solrDocToProperty(doc)
		if err != nil {
			// Log error pero continuar con otros documentos
			fmt.Printf("error convirtiendo documento de Solr: %v\n", err)
			continue
		}
		properties = append(properties, property)
	}

	return properties, solrResp.Response.NumFound, nil
}

// IndexProperty indexa una nueva propiedad en Solr
func (r *solrRepository) IndexProperty(ctx context.Context, property domain.Property) error {
	// Convertir domain.Property a SolrProperty
	solrProp := r.propertyToSolr(property)

	// Serializar a JSON
	jsonData, err := json.Marshal(solrProp)
	if err != nil {
		return fmt.Errorf("error serializando propiedad a JSON: %w", err)
	}

	// Construir URL de actualización
	updateURL := strings.TrimSuffix(r.solrURL, "/") + "/update/json/docs"

	// Crear request HTTP POST
	req, err := http.NewRequestWithContext(ctx, "POST", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando request HTTP: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Realizar petición
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error realizando petición a Solr: %w", err)
	}
	defer resp.Body.Close()

	// Verificar código de estado
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error indexando propiedad en Solr (status %d): %s", resp.StatusCode, string(body))
	}

	// Hacer commit
	return r.commit(ctx)
}

// UpdateProperty actualiza una propiedad existente en Solr
func (r *solrRepository) UpdateProperty(ctx context.Context, property domain.Property) error {
	// En Solr, actualizar es igual que indexar (reemplaza el documento)
	return r.IndexProperty(ctx, property)
}

// DeleteProperty elimina una propiedad de Solr por su ID
func (r *solrRepository) DeleteProperty(ctx context.Context, propertyID string) error {
	// Construir comando de eliminación
	deleteCmd := map[string]interface{}{
		"delete": map[string]string{
			"id": propertyID,
		},
	}

	// Serializar a JSON
	jsonData, err := json.Marshal(deleteCmd)
	if err != nil {
		return fmt.Errorf("error serializando comando de eliminación: %w", err)
	}

	// Construir URL de actualización
	updateURL := strings.TrimSuffix(r.solrURL, "/") + "/update"

	// Crear request HTTP POST
	req, err := http.NewRequestWithContext(ctx, "POST", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando request HTTP: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Realizar petición
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error realizando petición a Solr: %w", err)
	}
	defer resp.Body.Close()

	// Verificar código de estado
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error eliminando propiedad en Solr (status %d): %s", resp.StatusCode, string(body))
	}

	// Hacer commit
	return r.commit(ctx)
}

// commit realiza un commit en Solr para hacer persistentes los cambios
func (r *solrRepository) commit(ctx context.Context) error {
	commitCmd := map[string]interface{}{
		"commit": map[string]interface{}{},
	}

	jsonData, err := json.Marshal(commitCmd)
	if err != nil {
		return fmt.Errorf("error serializando comando de commit: %w", err)
	}

	updateURL := strings.TrimSuffix(r.solrURL, "/") + "/update"
	req, err := http.NewRequestWithContext(ctx, "POST", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creando request HTTP: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error realizando commit en Solr: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error en commit de Solr (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// propertyToSolr convierte una domain.Property a SolrProperty
func (r *solrRepository) propertyToSolr(property domain.Property) SolrProperty {
	return SolrProperty{
		ID:            property.ID,
		Title:         property.Title,
		Description:   property.Description,
		City:          property.City,
		Country:       property.Country,
		PricePerNight: property.PricePerNight,
		Bedrooms:      property.Bedrooms,
		Bathrooms:     property.Bathrooms,
		MaxGuests:     property.MaxGuests,
		Images:        property.Images,
		OwnerID:       property.OwnerID,
		Available:     property.Available,
		CreatedAt:     property.CreatedAt,
	}
}

// solrDocToProperty convierte un documento de Solr a domain.Property
func (r *solrRepository) solrDocToProperty(doc map[string]interface{}) (domain.Property, error) {
	property := domain.Property{}

	// Extraer y convertir campos
	if id, ok := doc["id"].(string); ok {
		property.ID = id
	}

	if title, ok := doc["title"].(string); ok {
		property.Title = title
	}

	if description, ok := doc["description"].(string); ok {
		property.Description = description
	}

	if city, ok := doc["city"].(string); ok {
		property.City = city
	}

	if country, ok := doc["country"].(string); ok {
		property.Country = country
	}

	if price, ok := doc["price_per_night"].(float64); ok {
		property.PricePerNight = price
	}

	if bedrooms, ok := doc["bedrooms"].(float64); ok {
		property.Bedrooms = int(bedrooms)
	}

	if bathrooms, ok := doc["bathrooms"].(float64); ok {
		property.Bathrooms = int(bathrooms)
	}

	if maxGuests, ok := doc["max_guests"].(float64); ok {
		property.MaxGuests = int(maxGuests)
	}

	if images, ok := doc["images"].([]interface{}); ok {
		property.Images = make([]string, 0, len(images))
		for _, img := range images {
			if imgStr, ok := img.(string); ok {
				property.Images = append(property.Images, imgStr)
			}
		}
	}

	if ownerID, ok := doc["owner_id"].(float64); ok {
		property.OwnerID = uint(ownerID)
	}

	if available, ok := doc["available"].(bool); ok {
		property.Available = available
	}

	if createdAt, ok := doc["created_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			property.CreatedAt = t
		}
	}

	return property, nil
}

// escapeSolrQuery escapa caracteres especiales en queries de Solr
func escapeSolrQuery(query string) string {
	// Escapar caracteres especiales de Solr
	specialChars := []string{"+", "-", "&", "|", "!", "(", ")", "{", "}", "[", "]", "^", "\"", "~", "*", "?", ":", "\\", "/"}
	escaped := query
	for _, char := range specialChars {
		escaped = strings.ReplaceAll(escaped, char, "\\"+char)
	}
	return escaped
}

