package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/controllers"
	"github.com/zeusnotfound04/Tranza/middleware" // This should map to the 'middlewares' package
	"github.com/zeusnotfound04/Tranza/repositories"
	"github.com/zeusnotfound04/Tranza/services"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")

	api.POST("/register", controllers.SignupHandler)
	api.POST("/login", controllers.LoginHandler)


	authGroup := api.Group("")
	authGroup.Use(middlewares.JWTAuthMiddleware())

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	cardRepo := repositories.NewCardRepository(db)
	cardService := services.NewCardService(cardRepo)
	cardController := controllers.NewCardController(cardService)

	authGroup.POST("/cards", cardController.LinkCard)
	authGroup.GET("/cards", cardController.GetCards)
	authGroup.DELETE("/cards/:id", cardController.DeleteCard)
	authGroup.PUT("/cards/:id/limit", cardController.UpdateLimit)
}
