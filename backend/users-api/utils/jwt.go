package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5" // ✅ CORRECTO - sin versión
	"os"
	"time"
)

// Esta es la "llave secreta" para firmar los tokens
// En producción debe estar en variables de entorno
var jwtSecret = []byte(getJWTSecret())

// Claims es la estructura de los datos que guardamos EN el token
// Cuando el usuario hace login, le damos un token con esta info
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"` // USER, ADMIN, OWNER - incluido como claim en JWT
	jwt.RegisteredClaims
}

// getJWTSecret obtiene el secret desde variables de entorno
// REQUIERE que JWT_SECRET esté configurado, no hay valor por defecto
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET no está configurado. Debe establecerse como variable de entorno.")
	}
	return secret
}

// GenerateToken genera un nuevo JWT token para un usuario
// Se llama después del login exitoso
// Incluye el role como claim en el JWT
// El token expira en 1 hora
func GenerateToken(userID uint, username, role string) (string, error) {
	// El token expira en 1 hora
	expirationTime := time.Now().Add(1 * time.Hour)

	// Creamos los "claims" (datos que va a tener el token)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role, // USER, ADMIN, OWNER - incluido como claim en JWT
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Creamos el token y lo firmamos con nuestro secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken valida un JWT token y retorna los claims
// Se usa en el middleware para verificar que el usuario esté autenticado
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// Parseamos el token y verificamos la firma
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// IsAdmin es una función helper que verifica si un role es ADMIN
func IsAdmin(role string) bool {
	return role == "ADMIN"
}

// IsOwner es una función helper que verifica si un role es OWNER
func IsOwner(role string) bool {
	return role == "OWNER"
}

// IsUser es una función helper que verifica si un role es USER
func IsUser(role string) bool {
	return role == "USER"
}
