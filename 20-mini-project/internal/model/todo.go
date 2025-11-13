package model

import (
	"time"

	"gorm.io/gorm"
)

// Todo merepresentasikan entity todo di database
type Todo struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null;size:200"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(20);default:'pending';index"`
	Priority    string `gorm:"type:varchar(10);default:'medium'"`
	DueDate     *time.Time
	UserID      uint `gorm:"not null;index"`
	User        User `gorm:"foreignKey:UserID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// TableName override nama tabel
func (Todo) TableName() string {
	return "todos"
}
