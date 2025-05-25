package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine , db *gorm.DB) {
	api := r.Group("/api")


	api.GET("/ping" , func(ctx *gin.Context) {
		ctx.JSON(200 , gin.H{"message" : "pong"})
	})
}