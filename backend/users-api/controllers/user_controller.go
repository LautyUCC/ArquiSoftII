package controllers

import (
	"net/http"
	"strconv"

	"users-api/dto"
	"users-api/services"
	"users-api/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// CreateUser maneja la creación de un nuevo usuario
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("POST /users", err, c)
		utils.SendErrorSimple(c, http.StatusBadRequest, "Error al parsear el cuerpo de la solicitud")
		return
	}

	user, err := ctrl.service.CreateUser(req)
	if err != nil {
		utils.LogError("POST /users", err, c)
		utils.SendErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.LogInfo("POST /users", "Usuario creado exitosamente", c)
	utils.SendSuccess(c, http.StatusCreated, user)
}

// Login maneja el inicio de sesión
func (ctrl *UserController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("POST /users/login", err, c)
		utils.SendErrorSimple(c, http.StatusBadRequest, "Error al parsear el cuerpo de la solicitud")
		return
	}

	response, err := ctrl.service.Login(req)
	if err != nil {
		utils.LogError("POST /users/login", err, c)
		utils.SendErrorSimple(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.LogInfo("POST /users/login", "Login exitoso", c)
	utils.SendSuccess(c, http.StatusOK, response)
}

// GetUserByID obtiene un usuario por ID
func (ctrl *UserController) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.LogError("GET /users/:id", err, c)
		utils.SendErrorSimple(c, http.StatusBadRequest, "ID inválido")
		return
	}

	user, err := ctrl.service.GetUserByID(uint(id))
	if err != nil {
		utils.LogError("GET /users/:id", err, c)
		utils.SendErrorSimple(c, http.StatusNotFound, err.Error())
		return
	}

	utils.LogInfo("GET /users/:id", "Usuario obtenido exitosamente", c)
	utils.SendSuccess(c, http.StatusOK, user)
}

// UpdateUser actualiza un usuario
// Solo ADMIN puede cambiar roles de otros usuarios
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.LogError("PUT /admin/users/:id", err, c)
		utils.SendErrorSimple(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.LogError("PUT /admin/users/:id", err, c)
		utils.SendErrorSimple(c, http.StatusBadRequest, "Error al parsear el cuerpo de la solicitud")
		return
	}

	// Validar que solo ADMIN puede cambiar roles
	if req.Role != nil {
		role, exists := c.Get("role")
		if !exists || role != "ADMIN" {
			utils.LogError("PUT /admin/users/:id", err, c)
			utils.SendErrorSimple(c, http.StatusForbidden, "solo ADMIN puede cambiar roles")
			return
		}
	}

	err = ctrl.service.UpdateUser(uint(id), req)
	if err != nil {
		utils.LogError("PUT /admin/users/:id", err, c)
		utils.SendErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.LogInfo("PUT /admin/users/:id", "Usuario actualizado exitosamente", c)
	utils.SendSuccess(c, http.StatusOK, gin.H{"message": "Usuario actualizado exitosamente"})
}

// DeleteUser elimina un usuario
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.LogError("DELETE /admin/users/:id", err, c)
		utils.SendErrorSimple(c, http.StatusBadRequest, "ID inválido")
		return
	}

	err = ctrl.service.DeleteUser(uint(id))
	if err != nil {
		utils.LogError("DELETE /admin/users/:id", err, c)
		utils.SendErrorSimple(c, http.StatusNotFound, err.Error())
		return
	}

	utils.LogInfo("DELETE /admin/users/:id", "Usuario eliminado exitosamente", c)
	utils.SendSuccess(c, http.StatusOK, gin.H{"message": "Usuario eliminado exitosamente"})
}

// GetAllUsers obtiene todos los usuarios (solo admin)
func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	users, err := ctrl.service.GetAllUsers()
	if err != nil {
		utils.LogError("GET /admin/users", err, c)
		utils.SendErrorSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.LogInfo("GET /admin/users", "Usuarios obtenidos exitosamente", c)
	utils.SendSuccess(c, http.StatusOK, users)
}

func (ctrl *UserController) HealthCheck(c *gin.Context) {
	utils.SendSuccess(c, http.StatusOK, gin.H{
		"status":  "OK",
		"service": "users-api",
	})
}
