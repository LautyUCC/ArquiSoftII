package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"search-api/services"

	"github.com/streadway/amqp"
)

// PropertyMessage representa un mensaje sobre una propiedad
type PropertyMessage struct {
	// Operation indica la operación a realizar: "create", "update", "delete"
	Operation string `json:"operation"`

	// PropertyID es el identificador único de la propiedad
	PropertyID string `json:"propertyId"`
}

// RabbitMQConsumer maneja el consumo de mensajes de RabbitMQ
type RabbitMQConsumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queueName  string
	service    services.SearchService
}

// NewRabbitMQConsumer crea una nueva instancia del consumidor de RabbitMQ
// Conecta con RabbitMQ, crea un channel y declara la queue "property_events"
func NewRabbitMQConsumer(rabbitURL, queueName string, service services.SearchService) (*RabbitMQConsumer, error) {
	log.Printf("🔌 Conectando a RabbitMQ en: %s", rabbitURL)

	// Conectar con RabbitMQ
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("error conectando a RabbitMQ: %w", err)
	}

	log.Println("✅ Conectado a RabbitMQ exitosamente")

	// Crear un channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error creando channel de RabbitMQ: %w", err)
	}

	log.Println("✅ Channel de RabbitMQ creado exitosamente")

	// Declarar la queue "property_events"
	// durable=true significa que la queue sobrevive a reinicios del servidor RabbitMQ
	_, err = channel.QueueDeclare(
		queueName, // nombre de la queue
		true,      // durable - la queue sobrevive a reinicios del servidor
		false,     // delete when unused - no se elimina cuando no hay consumidores
		false,     // exclusive - solo accesible por la conexión que la crea
		false,     // no-wait - no espera confirmación del servidor
		nil,       // arguments - argumentos adicionales
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("error declarando queue '%s' en RabbitMQ: %w", queueName, err)
	}

	log.Printf("✅ Queue '%s' declarada exitosamente", queueName)

	return &RabbitMQConsumer{
		connection: conn,
		channel:    channel,
		queueName:  queueName,
		service:    service,
	}, nil
}

// Start inicia el consumo de mensajes de la queue
// Procesa cada mensaje según su Operation y hace ACK
func (c *RabbitMQConsumer) Start() error {
	log.Printf("🚀 Iniciando consumo de mensajes de la queue: %s", c.queueName)

	// Configurar QoS para procesar un mensaje a la vez
	err := c.channel.Qos(
		1,     // prefetch count - número de mensajes sin ACK que puede tener el consumidor
		0,     // prefetch size - tamaño en bytes (0 = ilimitado)
		false, // global - aplicar a todos los consumidores de esta conexión
	)
	if err != nil {
		return fmt.Errorf("error configurando QoS: %w", err)
	}

	// Consumir mensajes de la queue
	msgs, err := c.channel.Consume(
		c.queueName, // queue
		"",          // consumer tag - identificador único del consumidor (vacío = auto-generado)
		false,       // auto-ack - no hacer ACK automático (queremos hacerlo manualmente)
		false,       // exclusive - solo este consumidor puede acceder a la queue
		false,       // no-local - no rechazar mensajes publicados en la misma conexión
		false,       // no-wait - no esperar confirmación del servidor
		nil,         // arguments - argumentos adicionales
	)
	if err != nil {
		return fmt.Errorf("error registrando consumidor: %w", err)
	}

	log.Printf("✅ Consumidor registrado exitosamente en queue: %s", c.queueName)

	// Procesar mensajes en un loop infinito
	go func() {
		for msg := range msgs {
			c.processMessage(msg)
		}
	}()

	log.Println("✅ Consumidor de RabbitMQ iniciado y escuchando mensajes")
	return nil
}

// processMessage procesa un mensaje individual
func (c *RabbitMQConsumer) processMessage(msg amqp.Delivery) {
	log.Printf("📨 Mensaje recibido: %s", string(msg.Body))

	// Deserializar el JSON a PropertyMessage
	var propertyMsg PropertyMessage
	if err := json.Unmarshal(msg.Body, &propertyMsg); err != nil {
		log.Printf("❌ Error deserializando mensaje: %v. Body: %s", err, string(msg.Body))
		// Rechazar el mensaje y no reintentarlo
		msg.Nack(false, false)
		return
	}

	// Validar que el mensaje tenga Operation y PropertyID
	if propertyMsg.Operation == "" {
		log.Printf("❌ Mensaje inválido: Operation está vacío. Body: %s", string(msg.Body))
		msg.Nack(false, false)
		return
	}
	if propertyMsg.PropertyID == "" {
		log.Printf("❌ Mensaje inválido: PropertyID está vacío. Body: %s", string(msg.Body))
		msg.Nack(false, false)
		return
	}

	log.Printf("🔄 Procesando mensaje - Operation: %s, PropertyID: %s", propertyMsg.Operation, propertyMsg.PropertyID)

	// Crear contexto con timeout para las operaciones
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Procesar según el Operation
	var err error
	switch propertyMsg.Operation {
	case "create":
		err = c.handleCreate(ctx, propertyMsg.PropertyID)
	case "update":
		err = c.handleUpdate(ctx, propertyMsg.PropertyID)
	case "delete":
		err = c.handleDelete(ctx, propertyMsg.PropertyID)
	default:
		log.Printf("⚠️ Operation desconocido: %s. Ignorando mensaje.", propertyMsg.Operation)
		// ACK el mensaje aunque no sepamos qué hacer con él
		msg.Ack(false)
		return
	}

	// Si hay error, loguearlo pero hacer ACK del mensaje para no reintentarlo infinitamente
	// En producción, podrías querer implementar un sistema de reintentos o dead letter queue
	if err != nil {
		log.Printf("❌ Error procesando mensaje (Operation: %s, PropertyID: %s): %v", propertyMsg.Operation, propertyMsg.PropertyID, err)
		// Hacer ACK para no reintentar (o implementar lógica de reintentos)
		msg.Ack(false)
		return
	}

	// Hacer ACK del mensaje si todo salió bien
	msg.Ack(false)
	log.Printf("✅ Mensaje procesado exitosamente - Operation: %s, PropertyID: %s", propertyMsg.Operation, propertyMsg.PropertyID)
}

// handleCreate maneja la acción "create"
// Obtiene la propiedad desde la API y la indexa en Solr
func (c *RabbitMQConsumer) handleCreate(ctx context.Context, propertyID string) error {
	log.Printf("📝 Creando/Indexando propiedad: %s", propertyID)

	// Obtener propiedad desde la API
	property, err := c.service.FetchPropertyFromAPI(propertyID)
	if err != nil {
		return fmt.Errorf("error obteniendo propiedad desde API: %w", err)
	}

	// Indexar en Solr
	if err := c.service.IndexProperty(ctx, *property); err != nil {
		return fmt.Errorf("error indexando propiedad en Solr: %w", err)
	}

	log.Printf("✅ Propiedad indexada exitosamente: %s", propertyID)
	return nil
}

// handleUpdate maneja la acción "update"
// Obtiene la propiedad actualizada desde la API y la actualiza en Solr
func (c *RabbitMQConsumer) handleUpdate(ctx context.Context, propertyID string) error {
	log.Printf("🔄 Actualizando propiedad: %s", propertyID)

	// Obtener propiedad actualizada desde la API
	property, err := c.service.FetchPropertyFromAPI(propertyID)
	if err != nil {
		return fmt.Errorf("error obteniendo propiedad desde API: %w", err)
	}

	// Actualizar en Solr
	if err := c.service.UpdateProperty(ctx, *property); err != nil {
		return fmt.Errorf("error actualizando propiedad en Solr: %w", err)
	}

	log.Printf("✅ Propiedad actualizada exitosamente: %s", propertyID)
	return nil
}

// handleDelete maneja la acción "delete"
// Elimina la propiedad de Solr
func (c *RabbitMQConsumer) handleDelete(ctx context.Context, propertyID string) error {
	log.Printf("🗑️ Eliminando propiedad: %s", propertyID)

	// Eliminar de Solr
	if err := c.service.DeleteProperty(ctx, propertyID); err != nil {
		return fmt.Errorf("error eliminando propiedad de Solr: %w", err)
	}

	log.Printf("✅ Propiedad eliminada exitosamente: %s", propertyID)
	return nil
}

// Close cierra las conexiones de RabbitMQ
func (c *RabbitMQConsumer) Close() error {
	log.Println("🔌 Cerrando conexiones de RabbitMQ...")

	var errs []error

	// Cerrar channel
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error cerrando channel: %w", err))
		} else {
			log.Println("✅ Channel cerrado exitosamente")
		}
	}

	// Cerrar conexión
	if c.connection != nil {
		if err := c.connection.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error cerrando conexión: %w", err))
		} else {
			log.Println("✅ Conexión de RabbitMQ cerrada exitosamente")
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errores al cerrar conexiones: %v", errs)
	}

	log.Println("✅ Todas las conexiones de RabbitMQ cerradas exitosamente")
	return nil
}
