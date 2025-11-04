package controllers

import (
	"net/http"
	"strconv"
	"users-api/dto"
	"users-api/services"

	"github.com/gin-gonic/gin"
)

// UserController maneja los endpoints de usuarios
type UserController struct {
	service services.UserService
}

// NewUserController crea una nueva instancia del controlador
func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// CreateUser maneja POST /users
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	user, err := ctrl.service.CreateUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "create_user_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "User created successfully",
		Data:    user,
	})
}

// GetUserByID maneja GET /users/:id
func (ctrl *UserController) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid user ID",
		})
		return
	}

	user, err := ctrl.service.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "user_not_found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Login maneja POST /users/login
func (ctrl *UserController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	response, err := ctrl.service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "login_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheck maneja GET /health
func (ctrl *UserController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"service": "users-api",
	})
}
