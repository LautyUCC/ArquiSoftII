package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config contiene toda la configuración de la aplicación
type Config struct {
	ServerPort   string
	MongoDB      MongoDBConfig
	RabbitMQ     RabbitMQConfig
	UsersAPI     UsersAPIConfig
	Environment  string
}

// MongoDBConfig contiene la configuración de MongoDB
type MongoDBConfig struct {
	URI      string
	Database string
}

// RabbitMQConfig contiene la configuración de RabbitMQ
type RabbitMQConfig struct {
	URI      string
	Exchange string
}

// UsersAPIConfig contiene la configuración para comunicarse con users-api
type UsersAPIConfig struct {
	BaseURL string
}

var AppConfig *Config

// Load carga la configuración desde variables de entorno
func Load() error {
	// Intentar cargar .env si existe (no es crítico si no existe)
	_ = godotenv.Load()

	AppConfig = &Config{
		ServerPort:  getEnv("SERVER_PORT", "8082"),
		Environment: getEnv("ENVIRONMENT", "development"),
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "properties_db"),
		},
		RabbitMQ: RabbitMQConfig{
			URI:      getEnv("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/"),
			Exchange: getEnv("RABBITMQ_EXCHANGE", "properties_exchange"),
		},
		UsersAPI: UsersAPIConfig{
			BaseURL: getEnv("USERS_API_URL", "http://users-api:8081"),
		},
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// Validate valida que la configuración sea correcta
func (c *Config) Validate() error {
	if c.MongoDB.URI == "" {
		return fmt.Errorf("MONGODB_URI no puede estar vacío")
	}
	if c.MongoDB.Database == "" {
		return fmt.Errorf("MONGODB_DATABASE no puede estar vacío")
	}
	if c.RabbitMQ.URI == "" {
		return fmt.Errorf("RABBITMQ_URI no puede estar vacío")
	}
	return nil
}

