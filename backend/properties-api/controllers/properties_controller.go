package controllers

import (
	"net/http"

	"properties-api/dto"
	"properties-api/services"

	"github.com/gin-gonic/gin"
)

type PropertyController struct {
	service services.PropertyService
}

func NewPropertyController(service services.PropertyService) *PropertyController {
	return &PropertyController{
		service: service,
	}
}

// CreateProperty maneja la creación de una nueva propiedad
func (c *PropertyController) CreateProperty(ctx *gin.Context) {
	var createDTO dto.PropertyCreateDTO
	if err := ctx.ShouldBindJSON(&createDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responseDTO, err := c.service.CreateProperty(createDTO)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, responseDTO)
}

// GetPropertyByID maneja la obtención de una propiedad por ID
func (c *PropertyController) GetPropertyByID(ctx *gin.Context) {
	id := ctx.Param("id")

	responseDTO, err := c.service.GetPropertyByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, responseDTO)
}

// UpdateProperty maneja la actualización de una propiedad
func (c *PropertyController) UpdateProperty(ctx *gin.Context) {
	id := ctx.Param("id")

	var updateDTO dto.PropertyUpdateDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extraer userID del contexto (agregado por middleware)
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Convertir userID de any a string
	var userID string
	switch v := userIDValue.(type) {
	case uint:
		userID = string(rune(v))
	case string:
		userID = v
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Tipo de userID inválido"})
		return
	}

	// Extraer isAdmin del contexto
	isAdminValue, _ := ctx.Get("isAdmin")
	isAdmin, _ := isAdminValue.(bool)

	err := c.service.UpdateProperty(id, updateDTO, userID, isAdmin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Propiedad actualizada exitosamente"})
}

// DeleteProperty maneja la eliminación de una propiedad
func (c *PropertyController) DeleteProperty(ctx *gin.Context) {
	id := ctx.Param("id")

	// Extraer userID del contexto
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Convertir userID de any a string
	var userID string
	switch v := userIDValue.(type) {
	case uint:
		userID = string(rune(v))
	case string:
		userID = v
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Tipo de userID inválido"})
		return
	}

	// Extraer isAdmin del contexto
	isAdminValue, _ := ctx.Get("isAdmin")
	isAdmin, _ := isAdminValue.(bool)

	err := c.service.DeleteProperty(id, userID, isAdmin)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Propiedad eliminada exitosamente"})
}

// GetUserProperties maneja la obtención de propiedades de un usuario
func (c *PropertyController) GetUserProperties(ctx *gin.Context) {
	userID := ctx.Param("userId")

	responseDTOs, err := c.service.GetUserProperties(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, responseDTOs)
}

// GetAllProperties maneja la obtención de todas las propiedades (solo admin)
func (c *PropertyController) GetAllProperties(ctx *gin.Context) {
	responseDTOs, err := c.service.GetAllProperties()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, responseDTOs)
}
