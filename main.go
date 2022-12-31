package main

import (
	"fantasy/database/controllers"
	"fantasy/database/lib"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env file. %v", err)
	}

	// Firestore online database
	// lib.InitClient()

	// Emulator
	lib.InitTestClient()

	defer lib.Client.Close()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Content-Type", os.Getenv("AUTHHEADER")},
		AllowCredentials: true,
	}))

	router.GET("/players/:id", controllers.GetPlayer)
	router.GET("/players/query", controllers.QueryPlayersByName)
	router.POST("/players", controllers.NewPlayer)

	router.GET("/league", controllers.VerifyIdToken, controllers.GetLeagueInfo)
	router.POST("/newleague", controllers.VerifyIdToken, controllers.NewLeague)

	router.POST("/register", controllers.NewUser)
	router.POST("/user/addplayer", controllers.VerifyIdToken, controllers.AddTeamPlayer)
	router.POST("/user/removeplayer", controllers.VerifyIdToken, controllers.RemoveTeamPlayer)
	router.GET("/user", controllers.VerifyIdToken, controllers.GetUserInfo)

	router.Run(":8080")
}
