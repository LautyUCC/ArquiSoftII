package domain

import "time"

// UserType define los tipos de usuario
type UserType string

const (
	UserTypeNormal UserType = "normal"
	UserTypeAdmin  UserType = "admin"
)

// User representa un usuario en el sistema
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"` // No se expone en JSON
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	UserType  UserType  `gorm:"type:varchar(20);default:'normal'" json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName especifica el nombre de la tabla
func (User) TableName() string {
	return "users"
}
