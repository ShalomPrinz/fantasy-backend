package controllers

import (
	"fantasy/database/db"
	"fantasy/database/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPlayers(c *gin.Context) {
	var players []entities.Player
	db.DB.Find(&players)
	c.JSON(http.StatusOK, gin.H{"data": players})
}

func NewPlayer(c *gin.Context) {
	var input entities.AddPlayer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player := entities.Player{
		Name: input.Name,
		Role: input.Role,
		Team: input.Team,
	}

	db.DB.Create(&player)
	c.JSON(http.StatusOK, gin.H{"data": player})
}
