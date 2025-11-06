package services

import (
	"encoding/json"
	"fmt"
	"properties-api/config"
	"time"
)

// EventService maneja la publicación de eventos en RabbitMQ
type EventService struct{}

// NewEventService crea una nueva instancia del servicio de eventos
func NewEventService() *EventService {
	return &EventService{}
}

// EventType representa el tipo de evento
type EventType string

const (
	EventPropertyCreated EventType = "property.created"
	EventPropertyUpdated EventType = "property.updated"
	EventPropertyDeleted EventType = "property.deleted"
)

// PropertyEvent representa un evento relacionado con propiedades
type PropertyEvent struct {
	Type      EventType   `json:"type"`
	Property  interface{} `json:"property"`
	Timestamp time.Time   `json:"timestamp"`
	Source    string      `json:"source"`
}

// PublishPropertyEvent publica un evento relacionado con propiedades
func (s *EventService) PublishPropertyEvent(eventType EventType, property interface{}) error {
	event := PropertyEvent{
		Type:      eventType,
		Property:  property,
		Timestamp: time.Now(),
		Source:    "properties-api",
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error serializando evento: %w", err)
	}

	routingKey := fmt.Sprintf("property.%s", string(eventType))
	if err := config.PublishEvent(routingKey, eventJSON); err != nil {
		return fmt.Errorf("error publicando evento: %w", err)
	}

	return nil
}

