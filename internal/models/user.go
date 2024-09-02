package models

import (
	"time"
)

type User struct {
	ID              uint       `gorm:"primaryKey;autoIncrement"` // ID del usuario, clave primaria auto-incremental
	Name            string     `gorm:"size:255"`                 // Nombre del usuario
	Email           string     `gorm:"size:255;unique"`          // Email del usuario, debe ser único
	EmailVerifiedAt *time.Time `gorm:"type:timestamp"`           // Marca de tiempo de verificación de correo
	Password        string     `gorm:"size:255"`                 // Contraseña del usuario
	RememberToken   string     `gorm:"size:100"`                 // Token de recordatorio para sesiones persistentes
	CreatedAt       time.Time  `gorm:"autoCreateTime"`           // Marca de tiempo de creación
	UpdatedAt       time.Time  `gorm:"autoUpdateTime"`           // Marca de tiempo de última actualización

}

func (User) TableName() string {
	return "users"
}
