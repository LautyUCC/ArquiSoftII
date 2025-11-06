package clients

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

// PropertyEvent representa un evento relacionado con propiedades
// Se serializa a JSON para ser publicado en RabbitMQ
type PropertyEvent struct {
	// Operation indica la operación realizada: "create", "update", "delete"
	Operation string `json:"operation"`

	// PropertyID es el identificador único de la propiedad afectada
	PropertyID string `json:"propertyId"`
}

// RabbitMQClient define la interfaz para publicar eventos en RabbitMQ
// Implementa el patrón de cliente para abstraer la lógica de mensajería
type RabbitMQClient interface {
	// PublishPropertyEvent publica un evento de propiedad en la cola de RabbitMQ
	// Serializa el evento a JSON y lo publica en la cola "property_events"
	// Retorna error si falla la serialización o la publicación
	PublishPropertyEvent(operation string, propertyID string) error
}

// rabbitMQClient es la implementación concreta de RabbitMQClient
// Usa github.com/streadway/amqp para la comunicación con RabbitMQ
type rabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQClient crea una nueva instancia del cliente de RabbitMQ
// Se conecta a RabbitMQ usando la URL proporcionada
// Declara la cola "property_events" como durable (sobrevive reinicios del servidor)
// Retorna error si falla la conexión o la declaración de la cola
func NewRabbitMQClient(url string) (RabbitMQClient, error) {
	// Conectar a RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error conectando a RabbitMQ en %s: %w", url, err)
	}

	// Abrir un canal de comunicación
	channel, err := conn.Channel()
	if err != nil {
		// Cerrar la conexión si falla la apertura del canal
		conn.Close()
		return nil, fmt.Errorf("error abriendo canal de RabbitMQ: %w", err)
	}

	// Declarar la cola "property_events" como durable
	// durable=true significa que la cola sobrevive a reinicios del servidor RabbitMQ
	queueName := "property_events"
	_, err = channel.QueueDeclare(
		queueName, // nombre de la cola
		true,      // durable - la cola sobrevive a reinicios del servidor
		false,     // delete when unused - no se elimina cuando no hay consumidores
		false,     // exclusive - solo accesible por la conexión que la crea
		false,     // no-wait - no espera confirmación del servidor
		nil,       // arguments - argumentos adicionales
	)
	if err != nil {
		// Cerrar canal y conexión si falla la declaración de la cola
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("error declarando cola '%s' en RabbitMQ: %w", queueName, err)
	}

	return &rabbitMQClient{
		conn:    conn,
		channel: channel,
	}, nil
}

// PublishPropertyEvent publica un evento de propiedad en la cola de RabbitMQ
// Serializa el evento a JSON y lo publica en la cola "property_events"
func (c *rabbitMQClient) PublishPropertyEvent(operation string, propertyID string) error {
	// Crear el struct del evento
	event := PropertyEvent{
		Operation:  operation,
		PropertyID: propertyID,
	}

	// Serializar el evento a JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error serializando evento a JSON: %w", err)
	}

	// Nombre de la cola donde se publicará el evento
	queueName := "property_events"

	// Publicar el mensaje en la cola
	err = c.channel.Publish(
		"",        // exchange - usar exchange por defecto (cadena vacía)
		queueName, // routing key - nombre de la cola
		false,     // mandatory - no es obligatorio que haya un consumidor
		false,     // immediate - no es inmediato
		amqp.Publishing{
			ContentType: "application/json",
			Body:        eventJSON,
		},
	)
	if err != nil {
		return fmt.Errorf("error publicando evento en la cola '%s': %w", queueName, err)
	}

	return nil
}

// Close cierra la conexión y el canal de RabbitMQ
// Útil para liberar recursos cuando ya no se necesita el cliente
func (c *rabbitMQClient) Close() error {
	var errs []error

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error cerrando canal: %w", err))
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error cerrando conexión: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errores al cerrar RabbitMQ client: %v", errs)
	}

	return nil
}

