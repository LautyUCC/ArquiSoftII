package services

import (
	"errors"
	"properties-api/dto"
	"properties-api/domain"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ============================================
// MOCKS
// ============================================

// mockRepository es un mock de PropertyRepository
// Permite controlar el comportamiento de las operaciones de repositorio en los tests
type mockRepository struct {
	CreateFunc        func(property domain.Property) (domain.Property, error)
	GetByIDFunc       func(id string) (domain.Property, error)
	UpdateFunc        func(id string, property domain.Property) error
	DeleteFunc        func(id string) error
	GetByOwnerIDFunc  func(ownerID string) ([]domain.Property, error)
}

// Create implementa PropertyRepository.Create
func (m *mockRepository) Create(property domain.Property) (domain.Property, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(property)
	}
	return domain.Property{}, errors.New("CreateFunc not set")
}

// GetByID implementa PropertyRepository.GetByID
func (m *mockRepository) GetByID(id string) (domain.Property, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return domain.Property{}, errors.New("GetByIDFunc not set")
}

// Update implementa PropertyRepository.Update
func (m *mockRepository) Update(id string, property domain.Property) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(id, property)
	}
	return errors.New("UpdateFunc not set")
}

// Delete implementa PropertyRepository.Delete
func (m *mockRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return errors.New("DeleteFunc not set")
}

// GetByOwnerID implementa PropertyRepository.GetByOwnerID
func (m *mockRepository) GetByOwnerID(ownerID string) ([]domain.Property, error) {
	if m.GetByOwnerIDFunc != nil {
		return m.GetByOwnerIDFunc(ownerID)
	}
	return nil, errors.New("GetByOwnerIDFunc not set")
}

// mockUsersClient es un mock de UsersClient
// Permite controlar el comportamiento de la validación de usuarios en los tests
type mockUsersClient struct {
	ValidateUserFunc func(userID string) (bool, error)
}

// ValidateUser implementa UsersClient.ValidateUser
func (m *mockUsersClient) ValidateUser(userID string) (bool, error) {
	if m.ValidateUserFunc != nil {
		return m.ValidateUserFunc(userID)
	}
	return false, errors.New("ValidateUserFunc not set")
}

// mockRabbitClient es un mock de RabbitMQClient
// Permite controlar el comportamiento de la publicación de eventos en los tests
type mockRabbitClient struct {
	PublishPropertyEventFunc func(operation string, propertyID string) error
}

// PublishPropertyEvent implementa RabbitMQClient.PublishPropertyEvent
func (m *mockRabbitClient) PublishPropertyEvent(operation string, propertyID string) error {
	if m.PublishPropertyEventFunc != nil {
		return m.PublishPropertyEventFunc(operation, propertyID)
	}
	return nil // Por defecto no retorna error para no bloquear tests
}

// ============================================
// HELPERS
// ============================================

// createTestProperty crea una propiedad de prueba para usar en los tests
func createTestProperty(id string, ownerID string) domain.Property {
	objectID := primitive.NewObjectID()
	if id != "" {
		var err error
		objectID, err = primitive.ObjectIDFromHex(id)
		if err != nil {
			objectID = primitive.NewObjectID()
		}
	}

	now := time.Now().Format(time.RFC3339)
	return domain.Property{
		ID:          objectID,
		Title:       "Test Property",
		Description: "Test Description",
		Price:       1000.0,
		Location:    "Test Location",
		OwnerID:     ownerID,
		Amenities:   []string{"wifi", "pool"},
		Capacity:    4,
		Available:  true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// createTestCreateDTO crea un DTO de creación de prueba
func createTestCreateDTO(ownerID string) dto.PropertyCreateDTO {
	return dto.PropertyCreateDTO{
		Title:       "Test Property",
		Description: "Test Description",
		Price:       1000.0,
		Location:    "Test Location",
		OwnerID:     ownerID,
		Amenities:   []string{"wifi", "pool"},
		Capacity:    4,
		Available:  true,
	}
}

// ============================================
// TESTS
// ============================================

// TestCreateProperty_Success testa el caso exitoso de creación de propiedad
func TestCreateProperty_Success(t *testing.T) {
	// Arrange
	ownerID := "user123"
	propertyID := primitive.NewObjectID().Hex()

	mockRepo := &mockRepository{
		CreateFunc: func(property domain.Property) (domain.Property, error) {
			// El servicio calcula el precio antes de crear, así que retornamos la propiedad con el precio calculado
			createdProperty := property
			createdProperty.ID, _ = primitive.ObjectIDFromHex(propertyID)
			return createdProperty, nil
		},
	}

	mockUsersClient := &mockUsersClient{
		ValidateUserFunc: func(userID string) (bool, error) {
			return true, nil // Usuario existe
		},
	}

	mockRabbitClient := &mockRabbitClient{
		PublishPropertyEventFunc: func(operation string, propertyID string) error {
			return nil
		},
	}

	service := NewPropertyService(mockRepo, mockUsersClient, mockRabbitClient)
	createDTO := createTestCreateDTO(ownerID)

	// Act
	result, err := service.CreateProperty(createDTO)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.ID == "" {
		t.Error("Expected property ID to be set")
	}

	if result.OwnerID != ownerID {
		t.Errorf("Expected OwnerID %s, got %s", ownerID, result.OwnerID)
	}

	// Verificar que el precio fue calculado con concurrencia (incluye impuestos y extras)
	// El cálculo incluye: precio base * 1.21 + ($50 * amenidades) + ($30 * capacidad)
	// Con precio base 1000, 2 amenidades, capacidad 4:
	// 1000 * 1.21 + (2 * 50) + (4 * 30) = 1210 + 100 + 120 = 1430
	expectedMinPrice := createDTO.Price * 1.21 // Al menos precio base con impuestos
	if result.Price < expectedMinPrice {
		t.Errorf("Expected price to be calculated with concurrency (at least %.2f), got %.2f", expectedMinPrice, result.Price)
	}
}

// TestCreateProperty_UserNotFound testa el caso cuando el owner no existe
func TestCreateProperty_UserNotFound(t *testing.T) {
	// Arrange
	ownerID := "nonexistent-user"
	mockRepo := &mockRepository{}

	mockUsersClient := &mockUsersClient{
		ValidateUserFunc: func(userID string) (bool, error) {
			return false, nil // Usuario no existe
		},
	}

	mockRabbitClient := &mockRabbitClient{}

	service := NewPropertyService(mockRepo, mockUsersClient, mockRabbitClient)
	createDTO := createTestCreateDTO(ownerID)

	// Act
	result, err := service.CreateProperty(createDTO)

	// Assert
	if err == nil {
		t.Fatal("Expected error when user does not exist")
	}

	if result.ID != "" {
		t.Error("Expected empty result when creation fails")
	}

	// Verificar que el error contiene información sobre el usuario no encontrado
	if !contains(err.Error(), "no existe") && !contains(err.Error(), "not found") {
		t.Errorf("Expected error message about user not found, got: %v", err)
	}
}

// TestGetPropertyByID_Success testa obtener una propiedad existente
func TestGetPropertyByID_Success(t *testing.T) {
	// Arrange
	propertyID := primitive.NewObjectID().Hex()
	ownerID := "user123"
	expectedProperty := createTestProperty(propertyID, ownerID)

	mockRepo := &mockRepository{
		GetByIDFunc: func(id string) (domain.Property, error) {
			if id == propertyID {
				return expectedProperty, nil
			}
			return domain.Property{}, errors.New("property not found")
		},
	}

	mockUsersClient := &mockUsersClient{}
	mockRabbitClient := &mockRabbitClient{}

	service := NewPropertyService(mockRepo, mockUsersClient, mockRabbitClient)

	// Act
	result, err := service.GetPropertyByID(propertyID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.ID != propertyID {
		t.Errorf("Expected ID %s, got %s", propertyID, result.ID)
	}

	if result.Title != expectedProperty.Title {
		t.Errorf("Expected Title %s, got %s", expectedProperty.Title, result.Title)
	}
}

// TestGetPropertyByID_NotFound testa obtener una propiedad inexistente
func TestGetPropertyByID_NotFound(t *testing.T) {
	// Arrange
	propertyID := "nonexistent-id"

	mockRepo := &mockRepository{
		GetByIDFunc: func(id string) (domain.Property, error) {
			return domain.Property{}, errors.New("property not found")
		},
	}

	mockUsersClient := &mockUsersClient{}
	mockRabbitClient := &mockRabbitClient{}

	service := NewPropertyService(mockRepo, mockUsersClient, mockRabbitClient)

	// Act
	result, err := service.GetPropertyByID(propertyID)

	// Assert
	if err == nil {
		t.Fatal("Expected error when property does not exist")
	}

	if result.ID != "" {
		t.Error("Expected empty result when property not found")
	}

	// Verificar que el error contiene información sobre la propiedad no encontrada
	if !contains(err.Error(), "no encontrada") && !contains(err.Error(), "not found") {
		t.Errorf("Expected error message about property not found, got: %v", err)
	}
}

// TestUpdateProperty_Unauthorized testa actualización sin permisos
func TestUpdateProperty_Unauthorized(t *testing.T) {
	// Arrange
	propertyID := primitive.NewObjectID().Hex()
	ownerID := "owner123"
	unauthorizedUserID := "user456"
	existingProperty := createTestProperty(propertyID, ownerID)

	mockRepo := &mockRepository{
		GetByIDFunc: func(id string) (domain.Property, error) {
			if id == propertyID {
				return existingProperty, nil
			}
			return domain.Property{}, errors.New("property not found")
		},
	}

	mockUsersClient := &mockUsersClient{}
	mockRabbitClient := &mockRabbitClient{}

	service := NewPropertyService(mockRepo, mockUsersClient, mockRabbitClient)

	updateDTO := dto.PropertyUpdateDTO{
		Title: stringPtr("Updated Title"),
	}

	// Act
	err := service.UpdateProperty(propertyID, updateDTO, unauthorizedUserID)

	// Assert
	if err == nil {
		t.Fatal("Expected error when user is not authorized")
	}

	// Verificar que el error contiene información sobre permisos
	if !contains(err.Error(), "no tiene permisos") && !contains(err.Error(), "permission") {
		t.Errorf("Expected error message about unauthorized access, got: %v", err)
	}
}

// TestDeleteProperty_Success testa eliminación exitosa
func TestDeleteProperty_Success(t *testing.T) {
	// Arrange
	propertyID := primitive.NewObjectID().Hex()
	ownerID := "owner123"
	existingProperty := createTestProperty(propertyID, ownerID)

	mockRepo := &mockRepository{
		GetByIDFunc: func(id string) (domain.Property, error) {
			if id == propertyID {
				return existingProperty, nil
			}
			return domain.Property{}, errors.New("property not found")
		},
		DeleteFunc: func(id string) error {
			if id == propertyID {
				return nil
			}
			return errors.New("property not found")
		},
	}

	mockUsersClient := &mockUsersClient{}
	mockRabbitClient := &mockRabbitClient{
		PublishPropertyEventFunc: func(operation string, propID string) error {
			if operation == "delete" && propID == propertyID {
				return nil
			}
			return errors.New("invalid event")
		},
	}

	service := NewPropertyService(mockRepo, mockUsersClient, mockRabbitClient)

	// Act
	err := service.DeleteProperty(propertyID, ownerID)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// TestDeleteProperty_Unauthorized testa eliminación sin permisos usando tabla de tests
func TestDeleteProperty_Unauthorized(t *testing.T) {
	// Tabla de tests para diferentes escenarios de no autorización
	tests := []struct {
		name           string
		propertyID     string
		ownerID        string
		requestingUser string
		expectedError  string
	}{
		{
			name:           "Different user tries to delete",
			propertyID:     "property123",
			ownerID:        "owner123",
			requestingUser: "user456",
			expectedError:  "no tiene permisos",
		},
		{
			name:           "Empty user ID",
			propertyID:     "property123",
			ownerID:        "owner123",
			requestingUser: "",
			expectedError:  "no tiene permisos",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			existingProperty := createTestProperty(tt.propertyID, tt.ownerID)

			mockRepo := &mockRepository{
				GetByIDFunc: func(id string) (domain.Property, error) {
					if id == tt.propertyID {
						return existingProperty, nil
					}
					return domain.Property{}, errors.New("property not found")
				},
			}

			mockUsersClient := &mockUsersClient{}
			mockRabbitClient := &mockRabbitClient{}

			service := NewPropertyService(mockRepo, mockUsersClient, mockRabbitClient)

			// Act
			err := service.DeleteProperty(tt.propertyID, tt.requestingUser)

			// Assert
			if err == nil {
				t.Fatal("Expected error when user is not authorized")
			}

			if !contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error to contain '%s', got: %v", tt.expectedError, err)
			}
		})
	}
}

// ============================================
// HELPER FUNCTIONS
// ============================================

// contains verifica si un string contiene otro
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// stringPtr crea un puntero a un string
func stringPtr(s string) *string {
	return &s
}

