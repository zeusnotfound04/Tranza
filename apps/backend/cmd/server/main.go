package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zeusnotfound04/Tranza/config"
	"github.com/zeusnotfound04/Tranza/routes"
)

func main() {
	config.LoadEnv()
	db := config.ConnectDB()
	defer config.CloseDB(db)

	router := gin.Default()
	routes.SetupRoutes(router , db)
	router.Run(":8080")
}