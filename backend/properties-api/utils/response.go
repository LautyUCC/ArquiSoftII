package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response estructura para respuestas JSON estandarizadas
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SuccessResponse envía una respuesta exitosa
func SuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ErrorResponse envía una respuesta de error
func ErrorResponse(c *gin.Context, statusCode int, err error, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   err.Error(),
		Message: message,
	})
}

// BadRequest envía una respuesta de error 400
func BadRequest(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusBadRequest, err, "Solicitud inválida")
}

// NotFound envía una respuesta de error 404
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error:   "Recurso no encontrado",
		Message: message,
	})
}

// InternalServerError envía una respuesta de error 500
func InternalServerError(c *gin.Context, err error) {
	ErrorResponse(c, http.StatusInternalServerError, err, "Error interno del servidor")
}

// Unauthorized envía una respuesta de error 401
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Error:   "No autorizado",
		Message: message,
	})
}

