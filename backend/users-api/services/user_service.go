package services

import (
	"errors"
	"users-api/domain"
	"users-api/dto"
	"users-api/repositories"
	"users-api/utils"
)

type UserService interface {
	CreateUser(userDTO dto.CreateUserRequest) (dto.UserResponse, error)
	Login(loginDTO dto.LoginRequest) (dto.LoginResponse, error)
	GetUserByID(id uint) (dto.UserResponse, error)
	UpdateUser(id uint, updateDTO dto.UpdateUserRequest) error
	DeleteUser(id uint) error
	GetAllUsers() ([]dto.UserResponse, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// CreateUser crea un nuevo usuario con contraseña hasheada
func (s *userService) CreateUser(userDTO dto.CreateUserRequest) (dto.UserResponse, error) {
	// Validar que el username no exista
	existingUser, _ := s.repo.GetByUsername(userDTO.Username)
	if existingUser != nil {
		return dto.UserResponse{}, errors.New("el username ya existe")
	}

	// Validar que el email no exista
	existingEmail, _ := s.repo.GetByEmail(userDTO.Email)
	if existingEmail != nil {
		return dto.UserResponse{}, errors.New("el email ya existe")
	}

	// Hashear la contraseña
	hashedPassword, err := utils.HashPassword(userDTO.Password)
	if err != nil {
		return dto.UserResponse{}, errors.New("error hasheando contraseña")
	}

	// Crear el usuario con user_type por defecto "normal"
	user := domain.User{
		Username:  userDTO.Username,
		Email:     userDTO.Email,
		Password:  hashedPassword,
		FirstName: userDTO.FirstName,
		LastName:  userDTO.LastName,
		UserType:  "normal", // Por defecto todos los usuarios son normales
	}

	// Guardar en la base de datos
	err = s.repo.Create(&user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// Retornar el DTO de respuesta (sin la contraseña)
	return s.toDTO(user), nil
}

// Login valida credenciales y genera token JWT
func (s *userService) Login(loginDTO dto.LoginRequest) (dto.LoginResponse, error) {
	// Buscar usuario por username o email
	var user *domain.User
	var err error

	// Intentar buscar por username primero
	user, err = s.repo.GetByUsername(loginDTO.UsernameOrEmail)
	if err != nil || user == nil {
		// Si no se encuentra, intentar por email
		user, err = s.repo.GetByEmail(loginDTO.UsernameOrEmail)
		if err != nil || user == nil {
			return dto.LoginResponse{}, errors.New("credenciales inválidas")
		}
	}

	// Verificar la contraseña
	if !utils.CheckPasswordHash(loginDTO.Password, user.Password) {
		return dto.LoginResponse{}, errors.New("credenciales inválidas")
	}

	// Generar token JWT
	token, err := utils.GenerateToken(user.ID, user.Username, user.UserType)
	if err != nil {
		return dto.LoginResponse{}, errors.New("error generando token")
	}

	// Retornar respuesta con token y datos del usuario
	return dto.LoginResponse{
		Token: token,
		User:  s.toDTO(*user),
	}, nil
}

// GetUserByID obtiene un usuario por su ID
func (s *userService) GetUserByID(id uint) (dto.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return dto.UserResponse{}, err
	}
	if user == nil {
		return dto.UserResponse{}, errors.New("usuario no encontrado")
	}

	return s.toDTO(*user), nil
}

// UpdateUser actualiza los datos de un usuario
func (s *userService) UpdateUser(id uint, updateDTO dto.UpdateUserRequest) error {
	// Obtener usuario existente
	user, err := s.repo.GetByID(id)
	if err != nil || user == nil {
		return errors.New("usuario no encontrado")
	}

	// Actualizar solo los campos que vienen en el DTO
	if updateDTO.Email != nil {
		// Verificar que el nuevo email no esté en uso por otro usuario
		existingEmail, _ := s.repo.GetByEmail(*updateDTO.Email)
		if existingEmail != nil && existingEmail.ID != id {
			return errors.New("el email ya está en uso")
		}
		user.Email = *updateDTO.Email
	}

	if updateDTO.FirstName != nil {
		user.FirstName = *updateDTO.FirstName
	}

	if updateDTO.LastName != nil {
		user.LastName = *updateDTO.LastName
	}

	if updateDTO.Password != nil {
		// Hashear la nueva contraseña
		hashedPassword, err := utils.HashPassword(*updateDTO.Password)
		if err != nil {
			return errors.New("error hasheando nueva contraseña")
		}
		user.Password = hashedPassword
	}

	// Guardar cambios
	return s.repo.Update(user)
}

// DeleteUser elimina un usuario por su ID
func (s *userService) DeleteUser(id uint) error {
	// Verificar que el usuario existe
	user, err := s.repo.GetByID(id)
	if err != nil || user == nil {
		return errors.New("usuario no encontrado")
	}

	return s.repo.Delete(id)
}

// GetAllUsers obtiene todos los usuarios (solo para admin)
func (s *userService) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Convertir cada usuario a DTO
	userDTOs := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userDTOs[i] = s.toDTO(user)
	}

	return userDTOs, nil
}

// toDTO convierte un domain.User a dto.UserResponse
func (s *userService) toDTO(user domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserType:  user.UserType,
	}
}
