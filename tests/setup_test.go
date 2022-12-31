package tests

import (
	"fantasy/database/controllers"
	"fantasy/database/lib"
	"fantasy/database/tests/utils"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func resetAuthenticaion(failTest func()) {
	deleteDataUrl := "http://localhost:8110/emulator/v1/projects/demo-test-fantasy/accounts"
	err := utils.Delete(deleteDataUrl)
	if err != nil {
		fmt.Println("Reset data: Failed deleting authentication data")
		failTest()
	}
}

func resetFirestore(failTest func()) {
	deleteDataUrl := "http://localhost:8090/emulator/v1/projects/demo-test-fantasy/databases/(default)/documents"
	err := utils.Delete(deleteDataUrl)
	if err != nil {
		fmt.Println("Reset data: Failed deleting firestore data")
		failTest()
	}
}

func beforeEach(failTest func()) {
	resetFirestore(func() { panic("") })
	resetAuthenticaion(func() { panic("") })
}

func initTestRouter() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/players/:id", controllers.GetPlayer)
	router.GET("/players/query", controllers.QueryPlayersByName)
	router.POST("/players", controllers.NewPlayer)

	router.GET("/league", controllers.VerifyIdToken, controllers.GetLeagueInfo)
	router.POST("/newleague", controllers.VerifyIdToken, controllers.NewLeague)

	router.POST("/register", controllers.NewUser)
	router.POST("/user/addplayer", controllers.VerifyIdToken, controllers.AddTeamPlayer)
	router.POST("/user/removeplayer", controllers.VerifyIdToken, controllers.RemoveTeamPlayer)
	router.GET("/user", controllers.VerifyIdToken, controllers.GetUserInfo)

	router.POST("/leagueinvite", controllers.VerifyIdToken, controllers.NewLeagueInvitation)

	// Test Routes
	router.GET("/test-token", controllers.VerifyIdToken)

	go router.Run(":8080")
	time.Sleep(50 * time.Millisecond)
}

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	lib.InitTestClient()
	defer lib.Client.Close()
	initTestRouter()

	code := m.Run()

	os.Exit(code)
}
