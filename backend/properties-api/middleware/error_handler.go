package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware maneja errores de forma centralizada
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Verificar si hay errores
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Determinar el código de estado HTTP apropiado
			statusCode := http.StatusInternalServerError
			if err.Type == gin.ErrorTypeBind {
				statusCode = http.StatusBadRequest
			}

			c.JSON(statusCode, gin.H{
				"error": err.Error(),
			})
		}
	}
}

