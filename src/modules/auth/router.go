package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup) {
	router.POST("/register", RegisterUserHandler)
	router.POST("/login", LoginUserHandler)
	router.POST("/refresh", RefreshTokenHandler)
}
