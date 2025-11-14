package dto

import "time"

type BookingCreateDTO struct {
	PropertyID string    `json:"propertyId" binding:"required"`
	UserID     string    `json:"userId" binding:"required"`
	CheckIn    time.Time `json:"checkIn" binding:"required"`
	CheckOut   time.Time `json:"checkOut" binding:"required"`
}

type BookingDTO struct {
	ID         string    `json:"id"`
	PropertyID string    `json:"propertyId"`
	UserID     string    `json:"userId"`
	CheckIn    time.Time `json:"checkIn"`
	CheckOut   time.Time `json:"checkOut"`
	TotalPrice float64   `json:"totalPrice"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
}
