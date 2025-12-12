package utils

import (
	"net/http"

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

// Funciones de conveniencia para códigos HTTP comunes

// BadRequest envía una respuesta de error 400
func BadRequest(c *gin.Context, message string, details map[string]interface{}) {
	SendError(c, http.StatusBadRequest, message, details)
}

// Unauthorized envía una respuesta de error 401
func Unauthorized(c *gin.Context, message string) {
	SendErrorSimple(c, http.StatusUnauthorized, message)
}

// Forbidden envía una respuesta de error 403
func Forbidden(c *gin.Context, message string) {
	SendErrorSimple(c, http.StatusForbidden, message)
}

// NotFound envía una respuesta de error 404
func NotFound(c *gin.Context, message string) {
	SendErrorSimple(c, http.StatusNotFound, message)
}

// Conflict envía una respuesta de error 409
func Conflict(c *gin.Context, message string) {
	SendErrorSimple(c, http.StatusConflict, message)
}

// InternalServerError envía una respuesta de error 500
func InternalServerError(c *gin.Context, message string, details map[string]interface{}) {
	SendError(c, http.StatusInternalServerError, message, details)
}

