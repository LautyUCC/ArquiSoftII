package utils

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse estructura unificada para respuestas de error
type ErrorResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// SendError envía una respuesta de error JSON unificada
// status: código HTTP (401, 403, 400, 404, 500, etc.)
// message: mensaje de error principal
// details: campos opcionales con información adicional del error
func SendError(c *gin.Context, status int, message string, details map[string]interface{}) {
	response := ErrorResponse{
		Status:  status,
		Message: message,
		Details: details,
	}
	c.JSON(status, response)
}

// SendErrorSimple envía una respuesta de error sin detalles adicionales
func SendErrorSimple(c *gin.Context, status int, message string) {
	SendError(c, status, message, nil)
}

// SendSuccess envía una respuesta exitosa
func SendSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

