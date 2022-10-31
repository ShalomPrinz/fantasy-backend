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

	db.InitClient()
	defer db.Client.Close()

	router := gin.Default()

	router.GET("/players", controllers.GetPlayers)
	router.GET("/players/:id", controllers.GetPlayer)
	router.POST("/players", controllers.NewPlayer)

	router.GET("/teams", controllers.GetTeams)
	router.GET("/teams/:id", controllers.GetTeam)
	router.POST("/teams", controllers.NewTeam)

	router.Run(":8080")
}
