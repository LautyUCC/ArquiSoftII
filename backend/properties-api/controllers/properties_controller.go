package controllers

import (
	"net/http"
	"properties-api/dto"
	"properties-api/services"

	"github.com/gin-gonic/gin"
)

// PropertyController maneja las peticiones HTTP relacionadas con propiedades
// Usa Gin framework para manejar las rutas y respuestas
type PropertyController struct {
	service services.PropertyService
}

// NewPropertyController crea una nueva instancia del controlador de propiedades
// Recibe el servicio como parámetro para inyección de dependencias
func NewPropertyController(service services.PropertyService) *PropertyController {
	return &PropertyController{
		service: service,
	}
}

// CreateProperty maneja la creación de una nueva propiedad
// POST /properties
// Bind JSON a PropertyCreateDTO, llama al servicio y retorna 201 Created o error apropiado
func (c *PropertyController) CreateProperty(ctx *gin.Context) {
	// Bind JSON a PropertyCreateDTO
	var createDTO dto.PropertyCreateDTO
	if err := ctx.ShouldBindJSON(&createDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	// Llamar service.CreateProperty
	responseDTO, err := c.service.CreateProperty(createDTO)
	if err != nil {
		// Verificar tipo de error para retornar código de estado apropiado
		if contains(err.Error(), "no existe") || contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "User not found",
				"message": err.Error(),
			})
			return
		}

		// Error interno del servidor
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	// Retornar 201 Created con la propiedad creada
	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    responseDTO,
		"message": "Property created successfully",
	})
}

// GetProperty maneja la obtención de una propiedad por ID
// GET /properties/:id
// Obtiene ID de params, llama al servicio y retorna 200 OK o 404 Not Found
func (c *PropertyController) GetProperty(ctx *gin.Context) {
	// Obtener ID de params
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Property ID is required",
		})
		return
	}

	// Llamar service.GetPropertyByID
	responseDTO, err := c.service.GetPropertyByID(id)
	if err != nil {
		// Retornar 404 Not Found si la propiedad no existe
		if contains(err.Error(), "no encontrada") || contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Property not found",
				"message": err.Error(),
			})
			return
		}

		// Error interno del servidor
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	// Retornar 200 OK con la propiedad
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    responseDTO,
	})
}

// UpdateProperty maneja la actualización de una propiedad
// PUT /properties/:id
// Obtiene ID de params y userID del contexto, bind JSON, llama al servicio
// Retorna 200 OK, 403 Forbidden o error apropiado
func (c *PropertyController) UpdateProperty(ctx *gin.Context) {
	// Obtener ID de params
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Property ID is required",
		})
		return
	}

	// Obtener userID del contexto (ctx.GetString("user_id"))
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User ID not found in context",
		})
		return
	}

	// Bind JSON a PropertyUpdateDTO
	var updateDTO dto.PropertyUpdateDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	// Llamar service.UpdateProperty
	err := c.service.UpdateProperty(id, updateDTO, userID)
	if err != nil {
		// Verificar si es error de permisos (403 Forbidden)
		if contains(err.Error(), "no tiene permisos") || contains(err.Error(), "permission") || contains(err.Error(), "forbidden") {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}

		// Verificar si es error de no encontrado (404 Not Found)
		if contains(err.Error(), "no encontrada") || contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Property not found",
				"message": err.Error(),
			})
			return
		}

		// Error interno del servidor
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	// Retornar 200 OK
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Property updated successfully",
	})
}

// DeleteProperty maneja la eliminación de una propiedad
// DELETE /properties/:id
// Similar a Update pero sin body
// Retorna 200 OK, 403 Forbidden o error apropiado
func (c *PropertyController) DeleteProperty(ctx *gin.Context) {
	// Obtener ID de params
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Property ID is required",
		})
		return
	}

	// Obtener userID del contexto (ctx.GetString("user_id"))
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User ID not found in context",
		})
		return
	}

	// Llamar service.DeleteProperty (sin body, solo ID y userID)
	err := c.service.DeleteProperty(id, userID)
	if err != nil {
		// Verificar si es error de permisos (403 Forbidden)
		if contains(err.Error(), "no tiene permisos") || contains(err.Error(), "permission") || contains(err.Error(), "forbidden") {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": err.Error(),
			})
			return
		}

		// Verificar si es error de no encontrado (404 Not Found)
		if contains(err.Error(), "no encontrada") || contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "Property not found",
				"message": err.Error(),
			})
			return
		}

		// Error interno del servidor
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	// Retornar 200 OK
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Property deleted successfully",
	})
}

// GetUserProperties maneja la obtención de todas las propiedades de un usuario
// GET /properties/user
// Obtiene userID del contexto y retorna lista de propiedades del usuario
func (c *PropertyController) GetUserProperties(ctx *gin.Context) {
	// Obtener userID del contexto
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User ID not found in context",
		})
		return
	}

	// Llamar service.GetUserProperties
	properties, err := c.service.GetUserProperties(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"message": err.Error(),
		})
		return
	}

	// Retornar 200 OK con la lista de propiedades
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    properties,
		"count":   len(properties),
	})
}

// contains es una función auxiliar para verificar si un string contiene otro
// Usada para determinar el tipo de error y retornar el código de estado apropiado
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

