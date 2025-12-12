package utils

import (
	"fmt"
	"log"
	"net/http"
)

const (
	ServiceName = "search-api"
)

// LogError registra un error con contexto (servicio, endpoint)
func LogError(endpoint string, err error, r *http.Request) {
	logMsg := fmt.Sprintf("[%s] ERROR [%s]", ServiceName, endpoint)
	if r != nil {
		logMsg += fmt.Sprintf(" method=%s path=%s", r.Method, r.URL.Path)
	}
	logMsg += fmt.Sprintf(" error=%v", err)

	log.Println(logMsg)
}

// LogInfo registra información con contexto (servicio, endpoint)
func LogInfo(endpoint string, message string, r *http.Request) {
	logMsg := fmt.Sprintf("[%s] INFO [%s]", ServiceName, endpoint)
	if r != nil {
		logMsg += fmt.Sprintf(" method=%s path=%s", r.Method, r.URL.Path)
	}
	logMsg += fmt.Sprintf(" %s", message)

	log.Println(logMsg)
}

// LogInfoWithoutContext registra información sin contexto de request (útil para logs de inicialización)
func LogInfoWithoutContext(message string) {
	log.Printf("[%s] INFO %s", ServiceName, message)
}

// LogStartup registra mensajes de inicio del servicio
func LogStartup(message string) {
	log.Printf("[%s] STARTUP %s", ServiceName, message)
}

// LogConfig registra configuración (sin valores sensibles)
func LogConfig(key string, value string) {
	// Ocultar valores sensibles
	if key == "JWT_SECRET" || key == "DB_PASSWORD" || key == "SOLR_URL" {
		value = "***"
	}
	log.Printf("[%s] CONFIG %s=%s", ServiceName, key, value)
}

