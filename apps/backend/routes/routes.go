package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/controllers"
	middlewares "github.com/zeusnotfound04/Tranza/middleware"
	"github.com/zeusnotfound04/Tranza/repositories"
	"github.com/zeusnotfound04/Tranza/services"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api/v1")


	api.POST("/register", controllers.SignupHandler)
	api.POST("/login", controllers.LoginHandler)
	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})


	cardRepo := repositories.NewCardRepository(db)
	txnRepo := repositories.NewTransactionRepository(db)
	cardService := services.NewCardService(cardRepo)
	paymentService := services.NewPaymentService(cardRepo, txnRepo)

	
	cardController := controllers.NewCardController(cardService)
	transactionController := controllers.NewTransactionController(paymentService)


	auth := api.Group("")
	auth.Use(middlewares.JWTAuthMiddleware())

	auth.POST("/cards", cardController.LinkCard)
	auth.GET("/cards", cardController.GetCards)
	auth.DELETE("/cards/:id", cardController.DeleteCard)
	auth.PUT("/cards/:id/limit", cardController.UpdateLimit)

	
	auth.POST("/transactions", transactionController.MakePayment)
	auth.GET("/transactions", transactionController.GetUserTransactions)
}
