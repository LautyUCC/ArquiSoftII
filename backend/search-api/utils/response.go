package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
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
func SendError(w http.ResponseWriter, status int, message string, details map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := ErrorResponse{
		Status:  status,
		Message: message,
		Details: details,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Si falla la codificación, escribir un error simple
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":500,"message":"Error serializando respuesta de error"}`)
	}
}

// SendErrorSimple envía una respuesta de error sin detalles adicionales
func SendErrorSimple(w http.ResponseWriter, status int, message string) {
	SendError(w, status, message, nil)
}

// SendSuccess envía una respuesta exitosa
func SendSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Si falla la codificación, escribir un error simple
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":500,"message":"Error serializando respuesta"}`)
	}
}

