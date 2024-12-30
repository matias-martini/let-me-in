package models

import (
	"time"
)

// User represents a user in the system.
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time
}

