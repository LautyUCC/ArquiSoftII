package main

import (
	"fmt"
	"os"
	"users-api/controllers"
	"users-api/domain"
	"users-api/middleware"
	"users-api/repositories"
	"users-api/services"
	"users-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// ============================================
	// 1. CONFIGURACIÓN - Leer variables de entorno
	// ============================================
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "spotly_user")
	dbPassword := getEnv("DB_PASSWORD", "spotly_password")
	dbName := getEnv("DB_NAME", "users_db")

	utils.LogStartup("Configuración cargada")
	utils.LogConfig("DB_HOST", dbHost)
	utils.LogConfig("DB_PORT", dbPort)
	utils.LogConfig("DB_NAME", dbName)

	// ============================================
	// 2. CONECTAR A MYSQL
	// ============================================
	// DSN = Data Source Name (string de conexión)
	// Formato: usuario:password@tcp(host:puerto)/base_de_datos?opciones
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	utils.LogStartup("Conectando a MySQL...")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.LogInfoWithoutContext(fmt.Sprintf("ERROR: Failed to connect to database: %v", err))
		panic(err)
	}
	utils.LogStartup("Conexión a MySQL exitosa")

	// ============================================
	// 3. AUTO-MIGRAR LAS TABLAS
	// ============================================
	// GORM crea automáticamente la tabla "users" si no existe
	utils.LogStartup("Ejecutando migraciones...")
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		utils.LogInfoWithoutContext(fmt.Sprintf("ERROR: Failed to migrate database: %v", err))
		panic(err)
	}
	utils.LogStartup("Tablas creadas/actualizadas")

	// ============================================
	// 4. INICIALIZAR CAPAS (Patrón MVC)
	// ============================================
	utils.LogStartup("Inicializando capas...")

	// Repository: acceso a datos
	userRepo := repositories.NewUserRepository(db)

	// Service: lógica de negocio
	userService := services.NewUserService(userRepo)

	// Controller: maneja HTTP
	userController := controllers.NewUserController(userService)

	utils.LogStartup("Capas inicializadas")

	// ============================================
	// 5. CONFIGURAR GIN (Framework web)
	// ============================================
	// Gin es como Express en Node.js
	router := gin.Default()

	// CORS - Permitir requests desde el frontend
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

	// ============================================
	// 6. DEFINIR RUTAS (Endpoints)
	// ============================================
	utils.LogStartup("Configurando rutas...")

	// Rutas PÚBLICAS (sin autenticación)
	router.GET("/health", userController.HealthCheck)
	router.POST("/users", userController.CreateUser)     // Registro
	router.POST("/users/login", userController.Login)    // Login
	router.GET("/users/:id", userController.GetUserByID) // Obtener usuario

	// Rutas PROTEGIDAS (requieren JWT - solo admin)
	admin := router.Group("/admin")
	admin.Use(middleware.RequireAdmin()) // requireAdmin valida token Y rol ADMIN
	{
		admin.GET("/users", userController.GetAllUsers)       // Listar todos
		admin.PUT("/users/:id", userController.UpdateUser)    // Actualizar
		admin.DELETE("/users/:id", userController.DeleteUser) // Eliminar
	}

	utils.LogStartup("Rutas configuradas")

	// ============================================
	// 7. ARRANCAR EL SERVIDOR
	// ============================================
	port := getEnv("SERVER_PORT", "8080")

	utils.LogStartup(fmt.Sprintf("Users API corriendo en puerto %s", port))

	if err := router.Run(":" + port); err != nil {
		utils.LogInfoWithoutContext(fmt.Sprintf("ERROR: Failed to start server: %v", err))
		panic(err)
	}
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
