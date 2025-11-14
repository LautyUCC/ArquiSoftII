package repositories

import (
	"context"
	"fmt"
	"time"

	"properties-api/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PropertyRepository define la interfaz para las operaciones de repositorio de propiedades
// Implementa el patrón de repositorio para abstraer la lógica de acceso a datos
type PropertyRepository interface {
	Create(property domain.Property) (domain.Property, error)
	GetByID(id string) (domain.Property, error)
	GetByOwnerID(ownerID string) ([]domain.Property, error)
	Update(id string, property domain.Property) error
	Delete(id string) error
	GetAll() ([]domain.Property, error) // ← AGREGAR ESTA LÍNEA
}

// propertyRepository es la implementación concreta de PropertyRepository
// Usa MongoDB como almacenamiento
type propertyRepository struct {
	collection *mongo.Collection
}

// NewPropertyRepository crea una nueva instancia del repositorio de propiedades
// Recibe la colección de MongoDB como parámetro para inyección de dependencias
// Retorna la interfaz PropertyRepository para permitir intercambiabilidad
func NewPropertyRepository(collection *mongo.Collection) PropertyRepository {
	return &propertyRepository{
		collection: collection,
	}
}

// Create crea una nueva propiedad en la base de datos
// Genera automáticamente un nuevo ObjectID y establece las fechas de creación y actualización
func (r *propertyRepository) Create(property domain.Property) (domain.Property, error) {
	// Crear contexto con timeout para la operación
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generar nuevo ObjectID automáticamente si no existe
	if property.ID.IsZero() {
		property.ID = primitive.NewObjectID()
	}

	// Establecer fechas de creación y actualización en formato ISO 8601
	now := time.Now()
	property.CreatedAt = now
	property.UpdatedAt = now

	// Insertar el documento en la colección
	result, err := r.collection.InsertOne(ctx, property)
	if err != nil {
		return domain.Property{}, fmt.Errorf("error insertando propiedad en MongoDB: %w", err)
	}

	// Verificar que se insertó correctamente
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		property.ID = oid
	}

	return property, nil
}

// GetByID obtiene una propiedad por su ID (string)
// Convierte el string a ObjectID y realiza la búsqueda en MongoDB
func (r *propertyRepository) GetByID(id string) (domain.Property, error) {
	// Crear contexto con timeout para la operación
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convertir string a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Property{}, fmt.Errorf("ID inválido '%s': %w", id, err)
	}

	// Crear filtro BSON para buscar por _id
	filter := bson.M{"_id": objectID}

	// Buscar el documento
	var property domain.Property
	err = r.collection.FindOne(ctx, filter).Decode(&property)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Property{}, fmt.Errorf("propiedad con ID '%s' no encontrada", id)
		}
		return domain.Property{}, fmt.Errorf("error obteniendo propiedad de MongoDB: %w", err)
	}

	return property, nil
}

// Update actualiza una propiedad existente por su ID
// Convierte el string a ObjectID y actualiza todos los campos
func (r *propertyRepository) Update(id string, property domain.Property) error {
	// Crear contexto con timeout para la operación
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convertir string a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID inválido '%s': %w", id, err)
	}

	// Actualizar el campo UpdatedAt con la fecha actual en formato ISO 8601
	property.UpdatedAt = time.Now()

	// Asegurar que el ID de la propiedad coincida con el parámetro
	property.ID = objectID

	// Crear filtro BSON para buscar por _id
	filter := bson.M{"_id": objectID}

	// Crear documento de actualización usando $set para actualizar todos los campos
	update := bson.M{
		"$set": bson.M{
			"title":       property.Title,
			"description": property.Description,
			"price":       property.Price,
			"location":    property.Location,
			"ownerId":     property.OwnerID,
			"amenities":   property.Amenities,
			"capacity":    property.Capacity,
			"available":   property.Available,
			"updatedAt":   property.UpdatedAt,
		},
	}

	// Actualizar el documento
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error actualizando propiedad en MongoDB: %w", err)
	}

	// Verificar que se actualizó al menos un documento
	if result.MatchedCount == 0 {
		return fmt.Errorf("propiedad con ID '%s' no encontrada para actualizar", id)
	}

	return nil
}

// Delete elimina una propiedad por su ID
// Convierte el string a ObjectID y elimina el documento de MongoDB
func (r *propertyRepository) Delete(id string) error {
	// Crear contexto con timeout para la operación
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convertir string a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID inválido '%s': %w", id, err)
	}

	// Crear filtro BSON para buscar por _id
	filter := bson.M{"_id": objectID}

	// Eliminar el documento
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error eliminando propiedad de MongoDB: %w", err)
	}

	// Verificar que se eliminó al menos un documento
	if result.DeletedCount == 0 {
		return fmt.Errorf("propiedad con ID '%s' no encontrada para eliminar", id)
	}

	return nil
}

// GetByOwnerID obtiene todas las propiedades de un propietario específico
// Usa cursor para obtener múltiples resultados eficientemente
func (r *propertyRepository) GetByOwnerID(ownerID string) ([]domain.Property, error) {
	// Crear contexto con timeout para la operación
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Crear filtro BSON para buscar por ownerId
	filter := bson.M{"ownerId": ownerID}

	// Crear cursor para iterar sobre los resultados
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error buscando propiedades por ownerId '%s' en MongoDB: %w", ownerID, err)
	}

	// Cerrar el cursor al finalizar (importante para liberar recursos)
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			// Log del error pero no lo retornamos para no enmascarar errores anteriores
			fmt.Printf("⚠️ Error cerrando cursor: %v\n", err)
		}
	}()

	// Decodificar todos los resultados del cursor
	var properties []domain.Property
	if err = cursor.All(ctx, &properties); err != nil {
		return nil, fmt.Errorf("error decodificando propiedades del cursor: %w", err)
	}

	// Si no se encontraron propiedades, retornar slice vacío en lugar de error
	// Esto permite diferenciar entre "no hay propiedades" y "error en la consulta"
	if properties == nil {
		properties = []domain.Property{}
	}

	return properties, nil
}

// GetAll obtiene todas las propiedades del sistema (solo admin)
func (r *propertyRepository) GetAll() ([]domain.Property, error) {
	var properties []domain.Property
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error buscando todas las propiedades: %w", err)
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &properties); err != nil {
		return nil, fmt.Errorf("error decodificando propiedades: %w", err)
	}

	return properties, nil
}
