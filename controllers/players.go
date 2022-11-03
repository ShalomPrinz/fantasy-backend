package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPlayers(c *gin.Context) {
	players := lib.GetAll("players", entities.PlayerAttributes[:])
	c.JSON(http.StatusOK, gin.H{"players": players})
}

func GetPlayer(c *gin.Context) {
	player := lib.GetSingle("players", c.Param("id"), entities.PlayerAttributes[:])
	c.JSON(http.StatusOK, gin.H{"player": player})
}

func NewPlayer(c *gin.Context) {
	var input entities.AddPlayer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lib.InsertItem("players", entities.AddPlayer{
		Name: input.Name,
		Role: input.Role,
	})

	c.JSON(http.StatusOK, gin.H{"addedPlayer": true})
}
