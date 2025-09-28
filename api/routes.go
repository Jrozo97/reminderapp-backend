package api

import (
	"github.com/Jrozo97/reminderapp-backend/internal/handler"
	"github.com/Jrozo97/reminderapp-backend/internal/repository"
	"github.com/Jrozo97/reminderapp-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	//Ro

	// users
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
}
