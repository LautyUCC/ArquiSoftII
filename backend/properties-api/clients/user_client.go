package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"properties-api/config"
)

// UserClient maneja la comunicación con users-api
type UserClient struct {
	baseURL string
	client  *http.Client
}

// NewUserClient crea una nueva instancia del cliente de usuarios
func NewUserClient() *UserClient {
	return &UserClient{
		baseURL: config.AppConfig.UsersAPI.BaseURL,
		client:  &http.Client{},
	}
}

// UserResponse representa la respuesta de users-api
type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Active   bool   `json:"active"`
	Exists   bool   `json:"exists"`
}

// ValidateUser valida si un usuario existe y está activo en users-api
func (c *UserClient) ValidateUser(userID string) (bool, error) {
	url := fmt.Sprintf("%s/api/users/%s/validate", c.baseURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error haciendo request a users-api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("error validando usuario: %s", string(body))
	}

	var userResponse UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return false, fmt.Errorf("error decodificando respuesta: %w", err)
	}

	return userResponse.Exists && userResponse.Active, nil
}

// GetUser obtiene información de un usuario desde users-api
func (c *UserClient) GetUser(userID string) (*UserResponse, error) {
	url := fmt.Sprintf("%s/api/users/%s", c.baseURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error haciendo request a users-api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error obteniendo usuario: %s", string(body))
	}

	var userResponse UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %w", err)
	}

	return &userResponse, nil
}

// ValidateUserWithToken valida un usuario usando un token de autenticación
func (c *UserClient) ValidateUserWithToken(token string) (*UserResponse, error) {
	url := fmt.Sprintf("%s/api/users/me", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error haciendo request a users-api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("token inválido o expirado")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error validando usuario: %s", string(body))
	}

	var userResponse UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %w", err)
	}

	return &userResponse, nil
}

