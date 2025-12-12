package domain

import "time"

// Role define los roles de usuario disponibles en el sistema
type Role string

const (
	// RoleUser es el rol por defecto para usuarios normales
	// Permisos: leer y reservar propiedades
	RoleUser Role = "USER"

	// RoleAdmin es el rol para administradores
	// Permisos: todas las acciones de USER + crear/editar/eliminar propiedades y cambiar roles
	RoleAdmin Role = "ADMIN"

	// RoleOwner es el rol para propietarios de propiedades
	// Permisos: todas las acciones de USER + gestionar sus propias propiedades
	RoleOwner Role = "OWNER"
)

// IsValid valida que un rol sea válido
func (r Role) IsValid() bool {
	return r == RoleUser || r == RoleAdmin || r == RoleOwner
}

// User representa un usuario en el sistema
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // El "-" oculta el password en JSON
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	UserType  string    `gorm:"column:user_type" json:"-"` // Campo legacy, ignorado en JSON
	Role      string    `gorm:"column:role;default:'USER';type:varchar(20)" json:"role"` // USER, ADMIN, OWNER - mapea explícitamente a columna 'role'
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName especifica el nombre de la tabla en MySQL
func (User) TableName() string {
	return "users"
}
