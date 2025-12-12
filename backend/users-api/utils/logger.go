package utils

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

const (
	ServiceName = "users-api"
)

// LogError registra un error con contexto (servicio, endpoint, userId si está disponible)
func LogError(endpoint string, err error, ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	role := getRoleFromContext(ctx)

	logMsg := fmt.Sprintf("[%s] ERROR [%s]", ServiceName, endpoint)
	if userID != "" {
		logMsg += fmt.Sprintf(" userId=%s", userID)
	}
	if role != "" {
		logMsg += fmt.Sprintf(" role=%s", role)
	}
	logMsg += fmt.Sprintf(" error=%v", err)

	log.Println(logMsg)
}

// LogInfo registra información con contexto (servicio, endpoint, userId si está disponible)
func LogInfo(endpoint string, message string, ctx *gin.Context) {
	userID := getUserIDFromContext(ctx)
	role := getRoleFromContext(ctx)

	logMsg := fmt.Sprintf("[%s] INFO [%s]", ServiceName, endpoint)
	if userID != "" {
		logMsg += fmt.Sprintf(" userId=%s", userID)
	}
	if role != "" {
		logMsg += fmt.Sprintf(" role=%s", role)
	}
	logMsg += fmt.Sprintf(" %s", message)

	log.Println(logMsg)
}

// LogInfoWithoutContext registra información sin contexto de request (útil para logs de inicialización)
func LogInfoWithoutContext(message string) {
	log.Printf("[%s] INFO %s", ServiceName, message)
}

// getUserIDFromContext extrae el userID del contexto de Gin si está disponible
func getUserIDFromContext(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}

	// Intentar obtener user_id
	if userID, exists := ctx.Get("user_id"); exists {
		return fmt.Sprintf("%v", userID)
	}

	// Intentar obtener userID
	if userID, exists := ctx.Get("userID"); exists {
		return fmt.Sprintf("%v", userID)
	}

	return ""
}

// getRoleFromContext extrae el role del contexto de Gin si está disponible
func getRoleFromContext(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}

	if role, exists := ctx.Get("role"); exists {
		return fmt.Sprintf("%v", role)
	}

	return ""
}

// LogStartup registra mensajes de inicio del servicio
func LogStartup(message string) {
	log.Printf("[%s] STARTUP %s", ServiceName, message)
}

// LogConfig registra configuración (sin valores sensibles)
func LogConfig(key string, value string) {
	// Ocultar valores sensibles
	if key == "JWT_SECRET" || key == "DB_PASSWORD" {
		value = "***"
	}
	log.Printf("[%s] CONFIG %s=%s", ServiceName, key, value)
}

