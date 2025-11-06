package clients

import (
	"fmt"
	"io"
	"net/http"
)

// UsersClient define la interfaz para la comunicación HTTP con users-api
// Implementa el patrón de cliente para abstraer la lógica de comunicación HTTP
type UsersClient interface {
	// ValidateUser valida si un usuario existe en users-api
	// Hace una petición GET a {baseURL}/users/{userID}
	// Retorna true si el usuario existe (status 200), false si no existe (status 404)
	// Retorna error en otros casos (errores de red, status codes inesperados, etc.)
	ValidateUser(userID string) (bool, error)
}

// usersClient es la implementación concreta de UsersClient
// Usa net/http estándar de Go para realizar peticiones HTTP
type usersClient struct {
	baseURL string
}

// NewUsersClient crea una nueva instancia del cliente de usuarios
// Recibe la URL base del servicio users-api como parámetro
// Retorna la interfaz UsersClient para permitir intercambiabilidad y testabilidad
func NewUsersClient(baseURL string) UsersClient {
	return &usersClient{
		baseURL: baseURL,
	}
}

// ValidateUser valida si un usuario existe en users-api
// Realiza una petición GET a {baseURL}/users/{userID}
func (c *usersClient) ValidateUser(userID string) (bool, error) {
	// Construir la URL completa para la petición
	url := fmt.Sprintf("%s/users/%s", c.baseURL, userID)

	// Crear la petición HTTP GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("error creando request HTTP: %w", err)
	}

	// Establecer headers apropiados
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Realizar la petición HTTP usando el cliente por defecto
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error haciendo petición HTTP a users-api: %w", err)
	}

	// Cerrar el body de la respuesta al finalizar (importante para liberar recursos)
	defer resp.Body.Close()

	// Manejar diferentes códigos de estado HTTP
	switch resp.StatusCode {
	case http.StatusOK:
		// Usuario existe y la petición fue exitosa
		return true, nil

	case http.StatusNotFound:
		// Usuario no encontrado (404)
		// Esto no es un error, simplemente retornamos false
		return false, nil

	default:
		// Cualquier otro código de estado es un error
		// Intentar leer el body del error para dar más contexto
		body, readErr := io.ReadAll(resp.Body)
		var errorMsg string
		if readErr != nil {
			errorMsg = fmt.Sprintf("status code %d", resp.StatusCode)
		} else {
			errorMsg = fmt.Sprintf("status code %d: %s", resp.StatusCode, string(body))
		}

		return false, fmt.Errorf("error validando usuario en users-api: %s", errorMsg)
	}
}

