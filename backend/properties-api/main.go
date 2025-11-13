package main

import (
	"fmt"
	"log"
	"net/http"

	"properties-api/clients"
	"properties-api/config"
	"properties-api/controllers"
	"properties-api/repositories"
	"properties-api/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// ============================================
	// SECCIÓN 0: CARGAR CONFIGURACIÓN
	// ============================================

	// Cargar configuración desde variables de entorno
	if err := config.Load(); err != nil {
		log.Fatalf("❌ Error cargando configuración: %v", err)
	}

	// ============================================
	// SECCIÓN 1: CONEXIÓN A BASE DE DATOS
	// ============================================

	// Conectar a MongoDB
	// Esto inicializa la conexión y configura la colección "properties"
	if err := config.ConnectMongoDB(); err != nil {
		log.Fatalf("❌ Error conectando a MongoDB: %v", err)
	}
	defer func() {
		if err := config.DisconnectMongoDB(); err != nil {
			log.Printf("⚠️ Error cerrando conexión a MongoDB: %v", err)
		}
	}()

	// ============================================
	// SECCIÓN 2: INICIALIZACIÓN DE CLIENTES
	// ============================================

	// Inicializar cliente de usuarios con URL del servicio users-api
	// Este cliente se comunica con el microservicio de usuarios para validar usuarios
	// Usar la URL desde configuración (variable de entorno USERS_API_URL)
	usersClientURL := config.AppConfig.UsersAPI.BaseURL
	if usersClientURL == "" {
		usersClientURL = "http://spotly-users-api:8081" // Valor por defecto si no está configurado
	}
	usersClient := clients.NewUsersClient(usersClientURL)

	// Inicializar cliente de RabbitMQ con URL del broker
	// Este cliente se usa para publicar eventos de propiedades
	rabbitMQURL := "amqp://guest:guest@rabbitmq:5672/"
	rabbitClient, err := clients.NewRabbitMQClient(rabbitMQURL)
	if err != nil {
		log.Fatalf("❌ Error creando cliente de RabbitMQ: %v", err)
	}
	defer func() {
		// Cerrar conexión de RabbitMQ al finalizar
		if closer, ok := rabbitClient.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				log.Printf("⚠️ Error cerrando conexión a RabbitMQ: %v", err)
			}
		}
	}()

	// ============================================
	// SECCIÓN 3: INICIALIZACIÓN DE REPOSITORIO
	// ============================================

	// Crear repositorio pasando la colección de MongoDB
	// El repositorio maneja todas las operaciones de base de datos
	propertyRepo := repositories.NewPropertyRepository(config.PropertiesCollection)

	// ============================================
	// SECCIÓN 4: INICIALIZACIÓN DE SERVICIO
	// ============================================

	// Crear servicio con inyección de dependencias
	// El servicio contiene la lógica de negocio y coordina repositorio y clientes
	propertyService := services.NewPropertyService(
		propertyRepo, // Repositorio de propiedades
		usersClient,  // Cliente de usuarios
		rabbitClient, // Cliente de RabbitMQ
	)

	// ============================================
	// SECCIÓN 5: INICIALIZACIÓN DE CONTROLADOR
	// ============================================

	// Crear controlador pasando el servicio
	// El controlador maneja las peticiones HTTP y las respuestas
	propertyController := controllers.NewPropertyController(propertyService)

	// ============================================
	// SECCIÓN 6: CONFIGURACIÓN DE ROUTER GIN
	// ============================================

	// Configurar router Gin
	// Gin es el framework web que maneja las rutas HTTP
	router := gin.Default()
        router.Use(func(c *gin.Context) {
            c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
            c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
            if c.Request.Method == "OPTIONS" {
                c.AbortWithStatus(204)
                return
            }
    
            c.Next()
         })

	// Middleware global para logging y recovery
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "properties-api",
		})
	})

	// ============================================
	// SECCIÓN 7: REGISTRO DE RUTAS
	// ============================================

	// Grupo de rutas para propiedades
	properties := router.Group("/properties")
	{
		// POST /properties → CreateProperty
		// Crea una nueva propiedad
		properties.POST("", propertyController.CreateProperty)

		// GET /properties/:id → GetProperty
		// Obtiene una propiedad por su ID
		properties.GET("/:id", propertyController.GetProperty)

		// PUT /properties/:id → UpdateProperty
		// Actualiza una propiedad existente
		// Requiere userID en el contexto (seteado por middleware de autenticación)
		properties.PUT("/:id", propertyController.UpdateProperty)

		// DELETE /properties/:id → DeleteProperty
		// Elimina una propiedad
		// Requiere userID en el contexto (seteado por middleware de autenticación)
		properties.DELETE("/:id", propertyController.DeleteProperty)

		// GET /properties/user → GetUserProperties
		// Obtiene todas las propiedades de un usuario
		// Requiere userID en el contexto (seteado por middleware de autenticación)
		properties.GET("/user", propertyController.GetUserProperties)
	}

	// ============================================
	// SECCIÓN 8: INICIO DEL SERVIDOR
	// ============================================

	// Obtener puerto desde configuración (variable de entorno SERVER_PORT)
	port := config.AppConfig.ServerPort
	if port == "" {
		port = "8082" // Valor por defecto si no está configurado
	}
	serverPort := ":" + port

	fmt.Println("🚀 Iniciando servidor Properties API...")
	fmt.Printf("   Puerto: %s\n", serverPort)
	fmt.Printf("   MongoDB: Conectado\n")
	fmt.Printf("   RabbitMQ: Conectado\n")
	fmt.Printf("   Users API: %s\n", usersClientURL)
	fmt.Println("✅ Servidor listo para recibir peticiones")

	// Iniciar servidor HTTP
	if err := router.Run(serverPort); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}
