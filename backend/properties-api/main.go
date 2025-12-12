package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"properties-api/clients"
	"properties-api/controllers"
	"properties-api/middleware"
	"properties-api/repositories"
	"properties-api/services"
	"properties-api/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// ============================================
	// CONFIGURACIÓN - Leer variables de entorno
	// ============================================
	mongoURI := getEnv("MONGO_URI", "mongodb://mongodb:27017/spotly_properties")
	mongoDatabase := getEnv("MONGO_DATABASE", "spotly_properties")
	usersAPIURL := getEnv("USERS_API_URL", "http://users-api:8081")
	rabbitMQURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	jwtSecret := getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production")
	serverPort := getEnv("SERVER_PORT", "8081")

	utils.LogStartup("Configuración cargada")
	utils.LogConfig("MONGO_URI", mongoURI)
	utils.LogConfig("MONGO_DATABASE", mongoDatabase)
	utils.LogConfig("USERS_API_URL", usersAPIURL)
	utils.LogConfig("RABBITMQ_URL", rabbitMQURL)
	utils.LogConfig("JWT_SECRET", jwtSecret)
	utils.LogConfig("SERVER_PORT", serverPort)

	// ============================================
	// CONECTAR A MONGODB
	// ============================================
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		utils.LogInfoWithoutContext(fmt.Sprintf("ERROR: Error conectando a MongoDB: %v", err))
		panic(err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Verificar conexión a MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		utils.LogInfoWithoutContext(fmt.Sprintf("ERROR: Error haciendo ping a MongoDB: %v", err))
		panic(err)
	}
	utils.LogStartup("Conectado a MongoDB")

	// Obtener colección de propiedades
	propertiesCollection := mongoClient.Database(mongoDatabase).Collection("properties")

	// ============================================
	// INICIALIZAR CLIENTES
	// ============================================
	usersClient := clients.NewUsersClient(usersAPIURL)
	rabbitClient, err := clients.NewRabbitMQClient(rabbitMQURL)
	if err != nil {
		utils.LogInfoWithoutContext(fmt.Sprintf("ERROR: Error conectando a RabbitMQ: %v", err))
		panic(err)
	}

	// Inicializar repositorios
	propertyRepo := repositories.NewPropertyRepository(propertiesCollection)
	bookingRepo := repositories.NewBookingRepository(mongoClient.Database(mongoDatabase))

	// Inicializar servicios
	propertyService := services.NewPropertyService(propertyRepo, usersClient, rabbitClient)
	bookingService := services.NewBookingService(bookingRepo, propertyRepo, rabbitClient)

	// Inicializar controladores
	propertyController := controllers.NewPropertyController(propertyService)
	bookingController := controllers.NewBookingController(bookingService)
	uploadController := controllers.NewUploadController()

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

	// Rutas públicas (lectura - cualquier usuario puede ver)
	public := router.Group("/api")
	{
		public.GET("/properties/:id", propertyController.GetPropertyByID)
		public.GET("/properties/user/:userId", propertyController.GetUserProperties)
	}

	// Rutas protegidas (requieren autenticación - cualquier rol puede acceder)
	// REGLA: Solo requiere token JWT válido, sin restricción de rol
	// Esto permite que USER, ADMIN y OWNER puedan crear reservas y ver sus propias reservas
	protected := router.Group("/api")
	protected.Use(middleware.RequireAuth(jwtSecret)) // requireAuth valida token válido (cualquier rol)
	{
		// POST /upload: Cualquier usuario autenticado puede subir imágenes
		protected.POST("/upload", uploadController.UploadImage)

		// POST /bookings: Cualquier usuario autenticado puede crear reservas
		// El userId se obtiene del JWT, no del body del request
		// Esto asegura que un usuario solo pueda crear reservas para sí mismo
		// Verifica solapamiento de reservas confirmadas (409 si hay conflicto)
		protected.POST("/bookings", bookingController.CreateBooking)

		// GET /bookings/my: Cualquier usuario autenticado puede ver sus propias reservas
		// Ordenadas por fecha inicio descendente
		protected.GET("/bookings/my", bookingController.GetMyBookings)

		// PATCH /bookings/:id: Solo el dueño de la reserva o un ADMIN puede actualizar
		// Valida solapamiento con otras reservas CONFIRMED
		protected.PATCH("/bookings/:id", bookingController.UpdateBooking)

		// DELETE /bookings/:id: Solo el dueño de la reserva o un ADMIN puede eliminar
		// Marca la reserva como CANCELLED (soft delete)
		protected.DELETE("/bookings/:id", bookingController.DeleteBooking)

		// GET /bookings/property/:id: Cualquier usuario autenticado puede ver reservas CONFIRMED de una propiedad
		// Útil para bloquear el calendario en el frontend
		protected.GET("/bookings/property/:id", bookingController.GetPropertyBookings)
	}

	// Rutas de administrador (solo ADMIN puede crear/editar/eliminar propiedades)
	admin := router.Group("/api")
	admin.Use(middleware.RequireAdmin(jwtSecret)) // requireAdmin valida token Y rol ADMIN
	{
		admin.POST("/properties", propertyController.CreateProperty)
		admin.PUT("/properties/:id", propertyController.UpdateProperty)
		admin.DELETE("/properties/:id", propertyController.DeleteProperty)
	}

	// Rutas de administrador (endpoints adicionales)
	adminGroup := router.Group("/api/admin")
	adminGroup.Use(middleware.RequireAdmin(jwtSecret)) // requireAdmin valida token Y rol ADMIN
	{
		adminGroup.GET("/properties", propertyController.GetAllProperties)
		adminGroup.GET("/bookings", bookingController.GetAllBookings)
	}

	// NOTA: GET /bookings/property/:id ahora está en el grupo "protected" (cualquier usuario autenticado)
	// para permitir que el frontend bloquee el calendario. Si se necesita un endpoint ADMIN
	// que muestre TODAS las reservas (incluyendo CANCELLED), se puede agregar aquí.

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "properties-api",
		})
	})

	// Iniciar servidor
	utils.LogStartup(fmt.Sprintf("Properties API corriendo en puerto %s", serverPort))
	if err := router.Run(":" + serverPort); err != nil {
		utils.LogInfoWithoutContext(fmt.Sprintf("ERROR: Error iniciando servidor: %v", err))
		panic(err)
	}
}
