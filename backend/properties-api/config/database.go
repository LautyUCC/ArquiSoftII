package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient es la variable global que almacena el cliente de MongoDB
var MongoClient *mongo.Client

// PropertiesCollection es la variable global que almacena la colección de propiedades
var PropertiesCollection *mongo.Collection

// ConnectMongoDB establece la conexión con MongoDB y configura la colección
// Se conecta a mongodb://mongodb:27017 con un timeout de 10 segundos
// Inicializa la colección "properties" en la base de datos "properties_db"
func ConnectMongoDB() error {
	// Crear contexto con timeout de 10 segundos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// URI de conexión a MongoDB
	mongoURI := "mongodb://mongodb:27017"

	// Configurar opciones del cliente
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Intentar conectar al servidor de MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Error conectando a MongoDB en %s: %v", mongoURI, err)
		fmt.Println(errorMsg)
		return fmt.Errorf("error conectando a MongoDB: %w", err)
	}

	// Verificar la conexión haciendo ping
	if err := client.Ping(ctx, nil); err != nil {
		errorMsg := fmt.Sprintf("❌ Error verificando conexión a MongoDB: %v", err)
		fmt.Println(errorMsg)
		// Cerrar la conexión si el ping falla
		if disconnectErr := client.Disconnect(ctx); disconnectErr != nil {
			fmt.Printf("⚠️ Error cerrando conexión después de ping fallido: %v\n", disconnectErr)
		}
		return fmt.Errorf("error haciendo ping a MongoDB: %w", err)
	}

	// Asignar el cliente a la variable global
	MongoClient = client

	// Obtener la base de datos "properties_db"
	database := client.Database("properties_db")

	// Inicializar la colección "properties" en la variable global
	PropertiesCollection = database.Collection("properties")

	// Mensaje de éxito
	fmt.Println("✅ Conectado a MongoDB exitosamente")
	fmt.Printf("   URI: %s\n", mongoURI)
	fmt.Printf("   Base de datos: properties_db\n")
	fmt.Printf("   Colección: properties\n")

	return nil
}

// DisconnectMongoDB cierra la conexión con MongoDB
// Maneja errores y usa un timeout de 10 segundos
func DisconnectMongoDB() error {
	if MongoClient == nil {
		fmt.Println("⚠️ MongoClient es nil, no hay conexión que cerrar")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := MongoClient.Disconnect(ctx); err != nil {
		fmt.Printf("❌ Error cerrando conexión a MongoDB: %v\n", err)
		return fmt.Errorf("error cerrando conexión a MongoDB: %w", err)
	}

	fmt.Println("✅ Conexión a MongoDB cerrada correctamente")
	return nil
}
