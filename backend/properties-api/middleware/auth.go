package middleware

import (
	"strings"

	"properties-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims define la estructura de los claims del JWT
// Actualizado para usar "role" en lugar de "user_type"
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"` // USER, ADMIN, OWNER
	jwt.RegisteredClaims
}

// requireAuth valida el token JWT
// Si el token es válido, permite continuar y guarda la info del usuario en el contexto
// Si no hay token o es inválido, devuelve error 401 (Unauthorized)
func requireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.LogError("requireAuth", nil, c)
			utils.Unauthorized(c, "authorization header required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.LogError("requireAuth", nil, c)
			utils.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			utils.LogError("requireAuth", err, c)
			utils.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			utils.LogError("requireAuth", nil, c)
			utils.Unauthorized(c, "invalid claims")
			c.Abort()
			return
		}

		// Guardar la info del usuario en el contexto
		c.Set("userID", claims.UserID)
		c.Set("user_id", claims.UserID) // Compatibilidad con diferentes nombres
		c.Set("username", claims.Username)
		c.Set("role", claims.Role) // Role: USER, ADMIN, OWNER

		c.Next()
	}
}

// requireAdmin valida que el usuario tenga token válido Y rol ADMIN
// Si no hay token válido, responde 401 (Unauthorized)
// Si hay token pero no es ADMIN, responde 403 (Forbidden)
func requireAdmin(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Primero verificar que hay un token válido
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.LogError("requireAdmin", nil, c)
			utils.Unauthorized(c, "authorization header required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.LogError("requireAdmin", nil, c)
			utils.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			utils.LogError("requireAdmin", err, c)
			utils.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			utils.LogError("requireAdmin", nil, c)
			utils.Unauthorized(c, "invalid claims")
			c.Abort()
			return
		}

		// Guardar la info del usuario en el contexto
		c.Set("userID", claims.UserID)
		c.Set("user_id", claims.UserID) // Compatibilidad
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// Verificar que el role sea ADMIN
		if claims.Role != "ADMIN" {
			utils.LogError("requireAdmin", nil, c)
			utils.Forbidden(c, "admin privileges required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAuth es el export público de requireAuth
func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return requireAuth(jwtSecret)
}

// RequireAdmin es el export público de requireAdmin
func RequireAdmin(jwtSecret string) gin.HandlerFunc {
	return requireAdmin(jwtSecret)
}

// AuthMiddleware mantiene compatibilidad con código existente
// Deprecated: Usar RequireAuth() en su lugar
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return requireAuth(jwtSecret)
}

// AdminRequired mantiene compatibilidad con código existente
// Deprecated: Usar RequireAdmin() en su lugar
func AdminRequired() gin.HandlerFunc {
	// Esta función no puede funcionar sin jwtSecret, por lo que se mantiene
	// pero se recomienda usar RequireAdmin(jwtSecret) en su lugar
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			// Si no existe role, intentar obtener de userType (backward compatibility)
			userType, exists := c.Get("userType")
			if !exists {
				utils.LogError("AdminRequired", nil, c)
				utils.Unauthorized(c, "role not found in token")
				c.Abort()
				return
			}
			// Convertir userType antiguo a role
			if userType == "admin" {
				role = "ADMIN"
			} else {
				role = "USER"
			}
		}

		if role != "ADMIN" {
			utils.LogError("AdminRequired", nil, c)
			utils.Forbidden(c, "admin privileges required")
			c.Abort()
			return
		}
		c.Next()
	}
}
