package config

import "os"

// Config contiene toda la configuración de la aplicación
type Config struct {
	// SolrURL es la URL del servidor Solr para búsquedas
	SolrURL string

	// MemcachedHost es la dirección del servidor Memcached para caché
	MemcachedHost string

	// RabbitMQURL es la URL de conexión a RabbitMQ
	RabbitMQURL string

	// PropertiesAPIURL es la URL base de la API de propiedades
	PropertiesAPIURL string

	// Port es el puerto en el que escuchará el servidor
	Port string
}

// LoadConfig carga la configuración desde variables de entorno
// Si una variable no está definida, usa los valores por defecto
func LoadConfig() *Config {
	return &Config{
		SolrURL:         getEnv("SOLR_URL", "http://localhost:8983/solr/properties"),
		MemcachedHost:   getEnv("MEMCACHED_HOST", "localhost:11211"),
		RabbitMQURL:     getEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
		PropertiesAPIURL: getEnv("PROPERTIES_API_URL", "http://localhost:8081"),
		Port:             getEnv("SERVER_PORT", "8083"),
	}
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

