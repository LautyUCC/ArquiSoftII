package repositories

import (
	"context"
	"fmt"
	"properties-api/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BookingRepository define la interfaz para las operaciones de repositorio de reservas
// Implementa el patrón de repositorio para abstraer la lógica de acceso a datos
type BookingRepository interface {
	// Create crea una nueva reserva en la base de datos
	Create(ctx context.Context, booking *domain.Booking) error

	// GetByID obtiene una reserva por su ID
	GetByID(ctx context.Context, id string) (*domain.Booking, error)

	// GetByUser obtiene todas las reservas de un usuario específico
	GetByUser(ctx context.Context, userID string) ([]domain.Booking, error)

	// GetByProperty obtiene todas las reservas de una propiedad específica
	GetByProperty(ctx context.Context, propertyID string) ([]domain.Booking, error)

	// Update actualiza una reserva existente
	Update(ctx context.Context, booking *domain.Booking) error

	// Delete elimina una reserva por su ID
	Delete(ctx context.Context, id string) error

	// FindConfirmedByPropertyID obtiene solo las reservas confirmadas de una propiedad
	// Útil para verificar disponibilidad y solapamientos
	FindConfirmedByPropertyID(ctx context.Context, propertyID string) ([]domain.Booking, error)

	// GetAll obtiene todas las reservas del sistema (solo para admin)
	GetAll(ctx context.Context) ([]domain.Booking, error)
}

// bookingRepository es la implementación concreta de BookingRepository
// Usa MongoDB como almacenamiento en la misma base de datos que las propiedades
type bookingRepository struct {
	collection *mongo.Collection
}

// NewBookingRepository crea una nueva instancia del repositorio de reservas
// Recibe la base de datos de MongoDB como parámetro para inyección de dependencias
// Retorna la interfaz BookingRepository para permitir intercambiabilidad
func NewBookingRepository(db *mongo.Database) BookingRepository {
	return &bookingRepository{
		collection: db.Collection("bookings"),
	}
}

// Create crea una nueva reserva en la base de datos
// Genera automáticamente un nuevo ObjectID y establece las fechas de creación y actualización
func (r *bookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	// Generar nuevo ObjectID automáticamente si no existe
	if booking.ID.IsZero() {
		booking.ID = primitive.NewObjectID()
	}

	// Establecer timestamps
	now := time.Now()
	booking.CreatedAt = now
	booking.UpdatedAt = now

	// Establecer status por defecto si no está definido
	if booking.Status == "" {
		booking.Status = domain.BookingStatusConfirmed
	}

	// Insertar en la base de datos
	_, err := r.collection.InsertOne(ctx, booking)
	if err != nil {
		return fmt.Errorf("error creando reserva: %w", err)
	}

	return nil
}

// GetByID obtiene una reserva por su ID
func (r *bookingRepository) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	// Convertir string ID a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID de reserva inválido: %w", err)
	}

	var booking domain.Booking
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("reserva no encontrada con ID: %s", id)
		}
		return nil, fmt.Errorf("error obteniendo reserva: %w", err)
	}

	return &booking, nil
}

// GetByUser obtiene todas las reservas de un usuario específico
func (r *bookingRepository) GetByUser(ctx context.Context, userID string) ([]domain.Booking, error) {
	var bookings []domain.Booking

	cursor, err := r.collection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, fmt.Errorf("error buscando reservas del usuario: %w", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, fmt.Errorf("error decodificando reservas: %w", err)
	}

	return bookings, nil
}

// GetByProperty obtiene todas las reservas de una propiedad específica
func (r *bookingRepository) GetByProperty(ctx context.Context, propertyID string) ([]domain.Booking, error) {
	var bookings []domain.Booking

	cursor, err := r.collection.Find(ctx, bson.M{"propertyId": propertyID})
	if err != nil {
		return nil, fmt.Errorf("error buscando reservas de la propiedad: %w", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, fmt.Errorf("error decodificando reservas: %w", err)
	}

	return bookings, nil
}

// Update actualiza una reserva existente
func (r *bookingRepository) Update(ctx context.Context, booking *domain.Booking) error {
	// Validar que el booking tenga ID
	if booking.ID.IsZero() {
		return fmt.Errorf("no se puede actualizar una reserva sin ID")
	}

	// Actualizar timestamp de actualización
	booking.UpdatedAt = time.Now()

	// Actualizar en la base de datos
	filter := bson.M{"_id": booking.ID}
	update := bson.M{
		"$set": bson.M{
			"propertyId": booking.PropertyID,
			"userId":     booking.UserID,
			"startDate":  booking.StartDate,
			"endDate":    booking.EndDate,
			"totalPrice": booking.TotalPrice,
			"status":     booking.Status,
			"updatedAt":  booking.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error actualizando reserva: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("reserva no encontrada con ID: %s", booking.ID.Hex())
	}

	return nil
}

// Delete elimina una reserva por su ID
func (r *bookingRepository) Delete(ctx context.Context, id string) error {
	// Convertir string ID a ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID de reserva inválido: %w", err)
	}

	// Eliminar de la base de datos
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("error eliminando reserva: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("reserva no encontrada con ID: %s", id)
	}

	return nil
}

// FindConfirmedByPropertyID obtiene solo las reservas confirmadas de una propiedad
// Útil para verificar disponibilidad y solapamientos de fechas
func (r *bookingRepository) FindConfirmedByPropertyID(ctx context.Context, propertyID string) ([]domain.Booking, error) {
	var bookings []domain.Booking

	// Buscar solo reservas confirmadas
	filter := bson.M{
		"propertyId": propertyID,
		"status":     domain.BookingStatusConfirmed,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error buscando reservas confirmadas: %w", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, fmt.Errorf("error decodificando reservas: %w", err)
	}

	return bookings, nil
}

// GetAll obtiene todas las reservas del sistema (solo para admin)
func (r *bookingRepository) GetAll(ctx context.Context) ([]domain.Booking, error) {
	var bookings []domain.Booking

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error buscando todas las reservas: %w", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, fmt.Errorf("error decodificando reservas: %w", err)
	}

	return bookings, nil
}
