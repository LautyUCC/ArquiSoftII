package services

import (
	"fmt"
	"time"

	"properties-api/clients"
	"properties-api/domain"
	"properties-api/dto"
	"properties-api/repositories"
	"properties-api/utils"
)

// PropertyService define la interfaz para la lógica de negocio de propiedades
// Implementa las reglas de negocio y coordina las operaciones entre repositorios y clientes
type PropertyService interface {
	// CreateProperty crea una nueva propiedad con validación de usuario y cálculo de precio
	CreateProperty(createDTO dto.PropertyCreateDTO) (dto.PropertyResponseDTO, error)

	// GetPropertyByID obtiene una propiedad por su ID
	GetPropertyByID(id string) (dto.PropertyResponseDTO, error)

	// UpdateProperty actualiza una propiedad existente con validación de ownership y admin
	UpdateProperty(id string, updateDTO dto.PropertyUpdateDTO, userID string, isAdmin bool) error

	// DeleteProperty elimina una propiedad con validación de ownership y admin
	DeleteProperty(id string, userID string, isAdmin bool) error

	// GetUserProperties obtiene todas las propiedades de un usuario específico
	GetUserProperties(userID string) ([]dto.PropertyResponseDTO, error)

	// GetAllProperties obtiene todas las propiedades (solo admin)
	GetAllProperties() ([]dto.PropertyResponseDTO, error)
}

// propertyService es la implementación concreta de PropertyService
// Coordina las operaciones entre repositorio, cliente de usuarios y cliente de RabbitMQ
type propertyService struct {
	repo         repositories.PropertyRepository
	usersClient  clients.UsersClient
	rabbitClient clients.RabbitMQClient
}

// NewPropertyService crea una nueva instancia del servicio de propiedades
// Recibe las dependencias como parámetros para inyección de dependencias
// Esto permite testear el servicio fácilmente y cambiar implementaciones
func NewPropertyService(
	repo repositories.PropertyRepository,
	usersClient clients.UsersClient,
	rabbitClient clients.RabbitMQClient,
) PropertyService {
	return &propertyService{
		repo:         repo,
		usersClient:  usersClient,
		rabbitClient: rabbitClient,
	}
}

// CreateProperty crea una nueva propiedad con validación de usuario y cálculo de precio
// Implementa los siguientes pasos:
// 1. Validar que el owner existe llamando a usersClient.ValidateUser
// 2. Calcular precio final usando CalculatePriceWithConcurrency
// 3. Crear property con timestamps actuales
// 4. Guardar en repository
// 5. Publicar evento "create" en RabbitMQ
// 6. Retornar DTO de respuesta
func (s *propertyService) CreateProperty(createDTO dto.PropertyCreateDTO) (dto.PropertyResponseDTO, error) {
	// 1. Validar que el owner existe llamando a usersClient.ValidateUser
	ownerExists, err := s.usersClient.ValidateUser(createDTO.OwnerID)
	if err != nil {
		return dto.PropertyResponseDTO{}, fmt.Errorf("error validando usuario owner: %w", err)
	}
	if !ownerExists {
		return dto.PropertyResponseDTO{}, fmt.Errorf("usuario owner con ID '%s' no existe", createDTO.OwnerID)
	}

	// 2. Calcular precio final usando CalculatePriceWithConcurrency
	// El precio base del DTO se usa como base para el cálculo
	finalPrice := utils.CalculatePriceWithConcurrency(
		createDTO.Price,     // precio base
		createDTO.Amenities, // lista de amenidades
		createDTO.Capacity,  // capacidad
	)

	// 3. Crear property con timestamps actuales
	now := time.Now()
	property := domain.Property{
		Title:       createDTO.Title,
		Description: createDTO.Description,
		Price:       finalPrice, // Usar el precio calculado con concurrencia
		Location:    createDTO.Location,
		OwnerID:     createDTO.OwnerID,
		Amenities:   createDTO.Amenities,
		Capacity:    createDTO.Capacity,
		Available:   createDTO.Available,
		Images:      createDTO.Images,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 4. Guardar en repository
	createdProperty, err := s.repo.Create(property)
	if err != nil {
		return dto.PropertyResponseDTO{}, fmt.Errorf("error creando propiedad en repositorio: %w", err)
	}

	// 5. Publicar evento "create" en RabbitMQ
	// Convertir ObjectID a string para el evento
	propertyID := createdProperty.ID.Hex()
	if err := s.rabbitClient.PublishPropertyEvent("create", propertyID); err != nil {
		// Log del error pero no fallar la operación si el evento no se publica
		// La propiedad ya fue creada exitosamente
		fmt.Printf("⚠️ Error publicando evento 'create' en RabbitMQ para propiedad %s: %v\n", propertyID, err)
	}

	// 6. Retornar DTO de respuesta
	return s.toDTO(createdProperty), nil
}

// GetPropertyByID obtiene una propiedad por su ID
// Retorna el DTO de respuesta o error si no se encuentra
func (s *propertyService) GetPropertyByID(id string) (dto.PropertyResponseDTO, error) {
	property, err := s.repo.GetByID(id)
	if err != nil {
		return dto.PropertyResponseDTO{}, fmt.Errorf("error obteniendo propiedad: %w", err)
	}

	return s.toDTO(property), nil
}

// UpdateProperty actualiza una propiedad existente con validación de ownership y admin
// Implementa los siguientes pasos:
// 1. Obtener propiedad existente
// 2. Validar que userID == ownerID O que sea admin
// 3. Actualizar solo campos no vacíos
// 4. Actualizar timestamp
// 5. Publicar evento "update"
func (s *propertyService) UpdateProperty(id string, updateDTO dto.PropertyUpdateDTO, userID string, isAdmin bool) error {
	// 1. Obtener propiedad existente
	property, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("error obteniendo propiedad para actualizar: %w", err)
	}

	// 2. Validar permisos: solo owner o admin pueden actualizar
	if property.OwnerID != userID && !isAdmin {
		return fmt.Errorf("forbidden: usuario con ID '%s' no tiene permisos para actualizar propiedad '%s' (owner: '%s')", userID, id, property.OwnerID)
	}

	// 3. Actualizar solo campos no vacíos (no nil)
	// Crear un nuevo objeto Property con los valores actualizados
	updatedProperty := property

	if updateDTO.Title != nil {
		updatedProperty.Title = *updateDTO.Title
	}
	if updateDTO.Description != nil {
		updatedProperty.Description = *updateDTO.Description
	}
	if updateDTO.Price != nil {
		// Si se actualiza el precio, recalcular con concurrencia
		// Usar los valores actuales de amenities y capacity
		updatedProperty.Price = utils.CalculatePriceWithConcurrency(
			*updateDTO.Price,
			updatedProperty.Amenities,
			updatedProperty.Capacity,
		)
	}
	if updateDTO.Location != nil {
		updatedProperty.Location = *updateDTO.Location
	}
	if updateDTO.Amenities != nil {
		updatedProperty.Amenities = *updateDTO.Amenities
		// Si se actualizan las amenidades y hay precio, recalcular
		if updateDTO.Price != nil {
			updatedProperty.Price = utils.CalculatePriceWithConcurrency(
				*updateDTO.Price,
				updatedProperty.Amenities,
				updatedProperty.Capacity,
			)
		}
	}
	if updateDTO.Capacity != nil {
		updatedProperty.Capacity = *updateDTO.Capacity
		// Si se actualiza la capacidad y hay precio, recalcular
		if updateDTO.Price != nil {
			updatedProperty.Price = utils.CalculatePriceWithConcurrency(
				*updateDTO.Price,
				updatedProperty.Amenities,
				updatedProperty.Capacity,
			)
		}
	}
	if updateDTO.Available != nil {
		updatedProperty.Available = *updateDTO.Available
	}
	if updateDTO.Images != nil {
		updatedProperty.Images = *updateDTO.Images
	}

	// 4. Actualizar timestamp
	updatedProperty.UpdatedAt = time.Now()

	// Guardar la actualización en el repositorio
	err = s.repo.Update(id, updatedProperty)
	if err != nil {
		return fmt.Errorf("error actualizando propiedad en repositorio: %w", err)
	}

	// 5. Publicar evento "update"
	if err := s.rabbitClient.PublishPropertyEvent("update", id); err != nil {
		// Log del error pero no fallar la operación
		fmt.Printf("⚠️ Error publicando evento 'update' en RabbitMQ para propiedad %s: %v\n", id, err)
	}

	return nil
}

// DeleteProperty elimina una propiedad con validación de ownership y admin
// Valida que el usuario tenga permisos (owner o admin) y publica evento "delete"
func (s *propertyService) DeleteProperty(id string, userID string, isAdmin bool) error {
	// Obtener propiedad existente para validar ownership
	property, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("error obteniendo propiedad para eliminar: %w", err)
	}

	// Validar permisos: solo owner o admin pueden eliminar
	if property.OwnerID != userID && !isAdmin {
		return fmt.Errorf("forbidden: usuario con ID '%s' no tiene permisos para eliminar propiedad '%s' (owner: '%s')", userID, id, property.OwnerID)
	}

	// Eliminar la propiedad
	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("error eliminando propiedad en repositorio: %w", err)
	}

	// Publicar evento "delete"
	if err := s.rabbitClient.PublishPropertyEvent("delete", id); err != nil {
		// Log del error pero no fallar la operación
		fmt.Printf("⚠️ Error publicando evento 'delete' en RabbitMQ para propiedad %s: %v\n", id, err)
	}

	return nil
}

// GetUserProperties obtiene todas las propiedades de un usuario específico
// Retorna un slice de DTOs de respuesta o error
func (s *propertyService) GetUserProperties(userID string) ([]dto.PropertyResponseDTO, error) {
	properties, err := s.repo.GetByOwnerID(userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo propiedades del usuario: %w", err)
	}

	// Convertir cada propiedad del dominio a DTO
	responseDTOs := make([]dto.PropertyResponseDTO, len(properties))
	for i, property := range properties {
		responseDTOs[i] = s.toDTO(property)
	}

	return responseDTOs, nil
}

// GetAllProperties obtiene todas las propiedades del sistema (solo para admin)
// Retorna un slice de DTOs de respuesta o error
func (s *propertyService) GetAllProperties() ([]dto.PropertyResponseDTO, error) {
	properties, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo todas las propiedades: %w", err)
	}

	// Convertir cada propiedad del dominio a DTO
	responseDTOs := make([]dto.PropertyResponseDTO, len(properties))
	for i, property := range properties {
		responseDTOs[i] = s.toDTO(property)
	}

	return responseDTOs, nil
}

// toDTO es una función privada que convierte un Property del dominio a PropertyResponseDTO
// Centraliza la lógica de conversión para evitar duplicación de código
func (s *propertyService) toDTO(property domain.Property) dto.PropertyResponseDTO {
	return dto.PropertyResponseDTO{
		ID:          property.ID.Hex(),
		Title:       property.Title,
		Description: property.Description,
		Price:       property.Price,
		Location:    property.Location,
		OwnerID:     property.OwnerID,
		Amenities:   property.Amenities,
		Capacity:    property.Capacity,
		Available:   property.Available,
		Images:      property.Images,
		CreatedAt:   property.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   property.UpdatedAt.Format(time.RFC3339),
	}
}
