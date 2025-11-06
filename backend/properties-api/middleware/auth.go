package middleware

import (
	"net/http"
	"properties-api/clients"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware valida el token de autenticación
func AuthMiddleware(userClient *clients.UserClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token de autorización requerido"})
			c.Abort()
			return
		}

		// Extraer el token (formato: "Bearer <token>")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validar el token con users-api
		user, err := userClient.ValidateUserWithToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			c.Abort()
			return
		}

		// Guardar información del usuario en el contexto
		c.Set("userID", user.ID)
		c.Set("userEmail", user.Email)
		c.Set("username", user.Username)

		c.Next()
	}
}

// OptionalAuthMiddleware valida el token si está presente, pero no falla si no está
func OptionalAuthMiddleware(userClient *clients.UserClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		user, err := userClient.ValidateUserWithToken(token)
		if err == nil {
			c.Set("userID", user.ID)
			c.Set("userEmail", user.Email)
			c.Set("username", user.Username)
		}

		c.Next()
	}
}

