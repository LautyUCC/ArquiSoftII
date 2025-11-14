package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"properties-api/clients"
	"properties-api/controllers"
	"properties-api/middleware"
	"properties-api/repositories"
	"properties-api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Configuración de MongoDB
	mongoURI := "mongodb://mongodb:27017"
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Error conectando a MongoDB:", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Verificar conexión a MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error haciendo ping a MongoDB:", err)
	}
	fmt.Println("✅ Conectado a MongoDB")

	// Obtener colección de propiedades
	propertiesCollection := mongoClient.Database("spotly").Collection("properties")

	// Inicializar clientes
	usersClient := clients.NewUsersClient("http://users-api:8081")
	rabbitClient, err := clients.NewRabbitMQClient("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal("Error conectando a RabbitMQ:", err)
	}

	// Inicializar repositorios
	propertyRepo := repositories.NewPropertyRepository(propertiesCollection)

	// Inicializar servicios
	propertyService := services.NewPropertyService(propertyRepo, usersClient, rabbitClient)

	// Inicializar controladores
	propertyController := controllers.NewPropertyController(propertyService)

	// Configurar Gin
	router := gin.Default()

	// Middleware CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Rutas públicas
	public := router.Group("/api")
	{
		public.GET("/properties/:id", propertyController.GetPropertyByID)
		public.GET("/properties/user/:userId", propertyController.GetUserProperties)
	}

	// Rutas protegidas (requieren autenticación)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware("your-super-secret-jwt-key-change-this-in-production"))
	{
		protected.POST("/properties", propertyController.CreateProperty)
		protected.PUT("/properties/:id", propertyController.UpdateProperty)
		protected.DELETE("/properties/:id", propertyController.DeleteProperty)
	}

	// Rutas de administrador
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware("your-super-secret-jwt-key-change-this-in-production"))
	admin.Use(middleware.AdminRequired())
	{
		admin.GET("/properties", propertyController.GetAllProperties)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "properties-api",
		})
	})

	// Iniciar servidor
	fmt.Println("🚀 Properties API corriendo en puerto 8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
