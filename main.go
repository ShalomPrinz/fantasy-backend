package main

import (
	"fantasy/database/controllers"
	"fantasy/database/lib"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env file. %v", err)
	}

	lib.InitClient()
	defer lib.Client.Close()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
	}))

	router.GET("/players", controllers.GetPlayers)
	router.GET("/players/:id", controllers.GetPlayer)
	router.POST("/players", controllers.NewPlayer)

	router.GET("/teams", controllers.GetTeams)
	router.GET("/teams/:id", controllers.GetTeam)
	router.POST("/teams", controllers.NewTeam)

	router.POST("/users", controllers.NewUser)

	router.Run(":8080")
}
