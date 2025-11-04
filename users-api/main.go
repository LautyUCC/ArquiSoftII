package main

import (
	"fmt"
	"log"
	"os"
	"users-api/controllers"
	"users-api/domain"
	"users-api/repositories"
	"users-api/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Configuración de la base de datos desde variables de entorno
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "spotly_user")
	dbPassword := getEnv("DB_PASSWORD", "spotly_password")
	dbName := getEnv("DB_NAME", "users_db")

	// Conectar a MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrar las tablas
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")

	// Inicializar capas
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Configurar Gin
	router := gin.Default()

	// Rutas públicas
	router.GET("/health", userController.HealthCheck)
	router.POST("/users", userController.CreateUser)
	router.POST("/users/login", userController.Login)
	router.GET("/users/:id", userController.GetUserByID)

	// Obtener puerto del servidor
	port := getEnv("SERVER_PORT", "8080")

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
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
