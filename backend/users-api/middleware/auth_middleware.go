package middleware

import (
	"net/http"
	"strings"
	"users-api/utils"

	"github.com/gin-gonic/gin"
)

// requireAuth valida el JWT token en cada request
// Si el token es válido, permite continuar y guarda la info del usuario en el contexto
// Si no hay token o es inválido, devuelve error 401 (Unauthorized)
func requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el header "Authorization"
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			utils.LogError("requireAuth", nil, c)
			utils.SendErrorSimple(c, http.StatusUnauthorized, "authorization header required")
			c.Abort() // Detiene la ejecución
			return
		}

		// Formato esperado: "Bearer <token>"
		// Ejemplo: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.LogError("requireAuth", nil, c)
			utils.SendErrorSimple(c, http.StatusUnauthorized, "invalid authorization header format")
			c.Abort()
			return
		}

		// Extraer el token
		tokenString := parts[1]

		// Validar el token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.LogError("requireAuth", err, c)
			utils.SendErrorSimple(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		// Guardar la info del usuario en el contexto
		// Así los endpoints pueden saber quién hizo la request
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role) // Role incluido en el contexto desde el JWT

		c.Next() // Continúa con el endpoint
	}
}

// requireAdmin valida que el usuario tenga token válido Y rol ADMIN
// Si no hay token válido, responde 401 (Unauthorized)
// Si hay token pero no es ADMIN, responde 403 (Forbidden)
// Este middleware debe usarse DESPUÉS de requireAuth o incluir la validación del token
func requireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Primero verificar que hay un token válido (requireAuth)
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			utils.LogError("requireAdmin", nil, c)
			utils.SendErrorSimple(c, http.StatusUnauthorized, "authorization header required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.LogError("requireAdmin", nil, c)
			utils.SendErrorSimple(c, http.StatusUnauthorized, "invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.LogError("requireAdmin", err, c)
			utils.SendErrorSimple(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		// Guardar la info del usuario en el contexto
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// Verificar que el role sea ADMIN
		if claims.Role != "ADMIN" {
			utils.LogError("requireAdmin", nil, c)
			utils.SendErrorSimple(c, http.StatusForbidden, "admin privileges required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAuth es el export público de requireAuth
// Mantiene compatibilidad con código existente
func RequireAuth() gin.HandlerFunc {
	return requireAuth()
}

// RequireAdmin es el export público de requireAdmin
// Mantiene compatibilidad con código existente
func RequireAdmin() gin.HandlerFunc {
	return requireAdmin()
}

// AuthMiddleware mantiene compatibilidad con código existente
// Deprecated: Usar RequireAuth() en su lugar
func AuthMiddleware() gin.HandlerFunc {
	return requireAuth()
}

// AdminMiddleware mantiene compatibilidad con código existente
// Deprecated: Usar RequireAdmin() en su lugar
func AdminMiddleware() gin.HandlerFunc {
	return requireAdmin()
}
