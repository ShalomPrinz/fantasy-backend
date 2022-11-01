package controllers

import (
	"fantasy/database/db"
	"fantasy/database/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPlayers(c *gin.Context) {
	players := db.GetAll("players", entities.PlayerAttributes[:])
	c.JSON(http.StatusOK, gin.H{"players": players})
}

func GetPlayer(c *gin.Context) {
	player := db.GetSingle("players", c.Param("id"), entities.PlayerAttributes[:])
	c.JSON(http.StatusOK, gin.H{"player": player})
}

func NewPlayer(c *gin.Context) {
	var input entities.AddPlayer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.InsertItem("players", entities.AddPlayer{
		Name: input.Name,
		Role: input.Role,
	})

	c.JSON(http.StatusOK, gin.H{"addedPlayer": true})
}
