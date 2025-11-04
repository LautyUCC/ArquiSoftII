package services

import (
	"errors"
	"strings"
	"users-api/domain"
	"users-api/dto"
	"users-api/repositories"
	"users-api/utils"
)

// UserService define la interfaz del servicio
type UserService interface {
	CreateUser(req dto.CreateUserRequest) (*domain.User, error)
	GetUserByID(id uint) (*domain.User, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
}

// userService implementa UserService
type userService struct {
	repo repositories.UserRepository
}

// NewUserService crea una nueva instancia del servicio
func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

// CreateUser crea un nuevo usuario
func (s *userService) CreateUser(req dto.CreateUserRequest) (*domain.User, error) {
	// Verificar si el username ya existe
	existingUser, _ := s.repo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Verificar si el email ya existe
	existingUser, _ = s.repo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hashear la contraseña
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("error hashing password")
	}

	// Crear el usuario
	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UserType:  domain.UserTypeNormal, // Por defecto es normal
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID obtiene un usuario por su ID
func (s *userService) GetUserByID(id uint) (*domain.User, error) {
	return s.repo.GetByID(id)
}

// Login autentica un usuario y genera un token JWT
func (s *userService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	var user *domain.User
	var err error

	// Intentar login por username o email
	if strings.Contains(req.UsernameOrEmail, "@") {
		user, err = s.repo.GetByEmail(req.UsernameOrEmail)
	} else {
		user, err = s.repo.GetByUsername(req.UsernameOrEmail)
	}

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verificar contraseña
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generar token JWT
	token, err := utils.GenerateToken(user.ID, user.Username, string(user.UserType))
	if err != nil {
		return nil, errors.New("error generating token")
	}

	return &dto.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}
