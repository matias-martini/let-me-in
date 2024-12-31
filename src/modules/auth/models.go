package auth

import (
	"gorm.io/gorm"
)

type UserCredentials struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex"`
	Password string
	Salt     string
	UserID   uint
}

type User struct {
	gorm.Model
	DisplayName string
}

type RefreshToken struct {
	gorm.Model
	Token     string `gorm:"uniqueIndex"`
	UserID    uint
	User      User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE,foreignKey:UserID;"`
	ExpiresAt int64
	Active    bool
}
