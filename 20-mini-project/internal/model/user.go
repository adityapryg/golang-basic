package model

import (
	"time"

	"gorm.io/gorm"
)

// User merepresentasikan entity pengguna di database
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null;size:50;index"`
	Email     string `gorm:"unique;not null;size:100;index"`
	Password  string `gorm:"not null"` // Hash password, jangan pernah expose ke luar
	FullName  string `gorm:"size:100"`
	Todos     []Todo `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName override nama tabel
func (User) TableName() string {
	return "users"
}
