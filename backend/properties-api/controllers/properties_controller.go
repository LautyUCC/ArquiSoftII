package controllers

import (
	"fmt"
	"net/http"

	"properties-api/dto"
	"properties-api/services"
	"properties-api/utils"

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
// Solo ADMIN puede crear propiedades (validado por middleware RequireAdmin)
func (c *PropertyController) CreateProperty(ctx *gin.Context) {
	var createDTO dto.PropertyCreateDTO
	if err := ctx.ShouldBindJSON(&createDTO); err != nil {
		utils.LogError("POST /properties", err, ctx)
		utils.BadRequest(ctx, "Error al parsear el cuerpo de la solicitud", nil)
		return
	}

	responseDTO, err := c.service.CreateProperty(createDTO)
	if err != nil {
		utils.LogError("POST /properties", err, ctx)
		utils.BadRequest(ctx, err.Error(), nil)
		return
	}

	utils.LogInfo("POST /properties", "Propiedad creada exitosamente", ctx)
	utils.SendSuccess(ctx, http.StatusCreated, responseDTO)
}

// GetPropertyByID maneja la obtención de una propiedad por ID
func (c *PropertyController) GetPropertyByID(ctx *gin.Context) {
	id := ctx.Param("id")

	responseDTO, err := c.service.GetPropertyByID(id)
	if err != nil {
		utils.LogError("GET /properties/:id", err, ctx)
		utils.NotFound(ctx, err.Error())
		return
	}

	utils.LogInfo("GET /properties/:id", "Propiedad obtenida exitosamente", ctx)
	utils.SendSuccess(ctx, http.StatusOK, responseDTO)
}

// UpdateProperty maneja la actualización de una propiedad
// Solo ADMIN puede actualizar propiedades (validado por middleware RequireAdmin)
// La validación de owner se mantiene en el servicio para consistencia de datos,
// pero ADMIN puede editar cualquier propiedad sin importar el owner
func (c *PropertyController) UpdateProperty(ctx *gin.Context) {
	id := ctx.Param("id")

	var updateDTO dto.PropertyUpdateDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		utils.LogError("PUT /properties/:id", err, ctx)
		utils.BadRequest(ctx, "Error al parsear el cuerpo de la solicitud", nil)
		return
	}

	// Extraer userID del contexto (agregado por middleware RequireAdmin)
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		userIDValue, exists = ctx.Get("user_id") // Intentar con otro nombre
		if !exists {
			utils.LogError("PUT /properties/:id", nil, ctx)
			utils.Unauthorized(ctx, "Usuario no autenticado")
			return
		}
	}

	// Convertir userID de any a string
	var userID string
	switch v := userIDValue.(type) {
	case uint:
		userID = fmt.Sprintf("%d", v)
	case int:
		userID = fmt.Sprintf("%d", v)
	case string:
		userID = v
	default:
		utils.LogError("PUT /properties/:id", nil, ctx)
		utils.InternalServerError(ctx, "Tipo de userID inválido", nil)
		return
	}

	// Extraer role del contexto (siempre será ADMIN si pasó el middleware)
	roleValue, _ := ctx.Get("role")
	role, _ := roleValue.(string)
	isAdmin := role == "ADMIN"

	err := c.service.UpdateProperty(id, updateDTO, userID, isAdmin)
	if err != nil {
		utils.LogError("PUT /properties/:id", err, ctx)
		utils.BadRequest(ctx, err.Error(), nil)
		return
	}

	utils.LogInfo("PUT /properties/:id", "Propiedad actualizada exitosamente", ctx)
	utils.SendSuccess(ctx, http.StatusOK, gin.H{"message": "Propiedad actualizada exitosamente"})
}

// DeleteProperty maneja la eliminación de una propiedad
// Solo ADMIN puede eliminar propiedades (validado por middleware RequireAdmin)
// La validación de owner se mantiene en el servicio para consistencia de datos,
// pero ADMIN puede eliminar cualquier propiedad sin importar el owner
func (c *PropertyController) DeleteProperty(ctx *gin.Context) {
	id := ctx.Param("id")

	// Extraer userID del contexto (agregado por middleware RequireAdmin)
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		userIDValue, exists = ctx.Get("user_id") // Intentar con otro nombre
		if !exists {
			utils.LogError("DELETE /properties/:id", nil, ctx)
			utils.Unauthorized(ctx, "Usuario no autenticado")
			return
		}
	}

	// Convertir userID de any a string
	var userID string
	switch v := userIDValue.(type) {
	case uint:
		userID = fmt.Sprintf("%d", v)
	case int:
		userID = fmt.Sprintf("%d", v)
	case string:
		userID = v
	default:
		utils.LogError("DELETE /properties/:id", nil, ctx)
		utils.InternalServerError(ctx, "Tipo de userID inválido", nil)
		return
	}

	// Extraer role del contexto (siempre será ADMIN si pasó el middleware)
	roleValue, _ := ctx.Get("role")
	role, _ := roleValue.(string)
	isAdmin := role == "ADMIN"

	err := c.service.DeleteProperty(id, userID, isAdmin)
	if err != nil {
		utils.LogError("DELETE /properties/:id", err, ctx)
		utils.NotFound(ctx, err.Error())
		return
	}

	utils.LogInfo("DELETE /properties/:id", "Propiedad eliminada exitosamente", ctx)
	utils.SendSuccess(ctx, http.StatusOK, gin.H{"message": "Propiedad eliminada exitosamente"})
}

// GetUserProperties maneja la obtención de propiedades de un usuario
func (c *PropertyController) GetUserProperties(ctx *gin.Context) {
	userID := ctx.Param("userId")

	responseDTOs, err := c.service.GetUserProperties(userID)
	if err != nil {
		utils.LogError("GET /properties/user/:userId", err, ctx)
		utils.InternalServerError(ctx, err.Error(), nil)
		return
	}

	utils.LogInfo("GET /properties/user/:userId", "Propiedades obtenidas exitosamente", ctx)
	utils.SendSuccess(ctx, http.StatusOK, responseDTOs)
}

// GetAllProperties maneja la obtención de todas las propiedades (solo admin)
func (c *PropertyController) GetAllProperties(ctx *gin.Context) {
	responseDTOs, err := c.service.GetAllProperties()
	if err != nil {
		utils.LogError("GET /admin/properties", err, ctx)
		utils.InternalServerError(ctx, err.Error(), nil)
		return
	}

	utils.LogInfo("GET /admin/properties", "Propiedades obtenidas exitosamente", ctx)
	utils.SendSuccess(ctx, http.StatusOK, responseDTOs)
}
