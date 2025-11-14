package repositories

import (
	"context"
	"properties-api/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) error
	FindByUserID(ctx context.Context, userID string) ([]domain.Booking, error)
	FindByID(ctx context.Context, id string) (*domain.Booking, error)
}

type bookingRepository struct {
	collection *mongo.Collection
}

func NewBookingRepository(db *mongo.Database) BookingRepository {
	return &bookingRepository{
		collection: db.Collection("bookings"),
	}
}

func (r *bookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	booking.ID = primitive.NewObjectID()
	booking.CreatedAt = time.Now()
	booking.Status = "confirmed"

	_, err := r.collection.InsertOne(ctx, booking)
	return err
}

func (r *bookingRepository) FindByUserID(ctx context.Context, userID string) ([]domain.Booking, error) {
	var bookings []domain.Booking
	cursor, err := r.collection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) FindByID(ctx context.Context, id string) (*domain.Booking, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var booking domain.Booking
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}
