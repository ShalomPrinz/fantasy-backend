package tests

import (
	"fantasy/database/controllers"
	"fantasy/database/lib"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func initTestRouter() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/players/:id", controllers.GetPlayer)
	router.POST("/players", controllers.NewPlayer)

	go router.Run(":8080")
	time.Sleep(50 * time.Millisecond)
}

func TestMain(m *testing.M) {
	lib.InitTestClient()
	defer lib.Client.Close()
	initTestRouter()

	code := m.Run()

	os.Exit(code)
}
