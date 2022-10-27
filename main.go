package main

import (
	"fantasy/database/controllers"
	"fantasy/database/db"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env file. %v", err)
	}

	engine := gin.Default()
	db.InitConnection()

	engine.GET("/players", controllers.GetPlayers)
	engine.POST("/players", controllers.NewPlayer)
	engine.Run(":8080")
}
