package dto

// CreateUserRequest DTO para crear un nuevo usuario
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

// LoginRequest DTO para login
type LoginRequest struct {
	UsernameOrEmail string `json:"usernameOrEmail" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

// UpdateUserRequest DTO para actualizar usuario
// Solo ADMIN puede cambiar el role de otros usuarios
type UpdateUserRequest struct {
	Email     *string `json:"email"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Password  *string `json:"password"`
	Role      *string `json:"role"` // USER, ADMIN, OWNER - solo ADMIN puede cambiar
}

// UserResponse DTO de respuesta
type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"` // USER, ADMIN, OWNER
}

// LoginResponse DTO de respuesta del login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ErrorResponse DTO de respuesta de error
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse DTO de respuesta exitosa
type SuccessResponse struct {
	Message string `json:"message"`
}
