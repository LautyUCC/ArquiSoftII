package config

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMQConn *amqp.Connection
var RabbitMQChannel *amqp.Channel

// ConnectRabbitMQ establece la conexión con RabbitMQ
func ConnectRabbitMQ() error {
	conn, err := amqp.Dial(AppConfig.RabbitMQ.URI)
	if err != nil {
		return fmt.Errorf("error conectando a RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("error abriendo canal de RabbitMQ: %w", err)
	}

	// Declarar el exchange
	err = ch.ExchangeDeclare(
		AppConfig.RabbitMQ.Exchange, // nombre
		"topic",                      // tipo
		true,                         // durable
		false,                        // auto-deleted
		false,                        // internal
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("error declarando exchange: %w", err)
	}

	RabbitMQConn = conn
	RabbitMQChannel = ch

	fmt.Println("✅ Conectado a RabbitMQ exitosamente")
	return nil
}

// DisconnectRabbitMQ cierra la conexión con RabbitMQ
func DisconnectRabbitMQ() error {
	if RabbitMQChannel != nil {
		RabbitMQChannel.Close()
	}
	if RabbitMQConn != nil {
		return RabbitMQConn.Close()
	}
	return nil
}

// PublishEvent publica un evento en RabbitMQ
func PublishEvent(routingKey string, body []byte) error {
	if RabbitMQChannel == nil {
		return fmt.Errorf("canal de RabbitMQ no está disponible")
	}

	err := RabbitMQChannel.Publish(
		AppConfig.RabbitMQ.Exchange, // exchange
		routingKey,                   // routing key
		false,                        // mandatory
		false,                        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("error publicando evento: %w", err)
	}

	return nil
}

