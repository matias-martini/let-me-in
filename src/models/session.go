package models

import (
	"time"
)

// Session represents a terminal session.
type Session struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null"`
	Status       string    `gorm:"default:active"` // active, inactive, terminated
	ContainerID  string    `gorm:"not null"`
	IPAddress    string    `gorm:"not null"`
	LastActivity time.Time `gorm:"autoUpdateTime"`
	CreatedAt    time.Time
}

