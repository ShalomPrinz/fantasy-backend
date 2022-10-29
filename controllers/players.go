package controllers

import (
	"fantasy/database/db"
	"fantasy/database/entities"
	"fantasy/database/utils"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func GetPlayerData(doc *firestore.DocumentSnapshot) interface{} {
	return entities.Player{
		ID:   doc.Ref.ID,
		Name: utils.GetDocString(doc, "Name"),
		Role: utils.GetDocString(doc, "Role"),
	}
}

func GetPlayers(c *gin.Context) {
	players := db.GetAll("players", GetPlayerData)
	c.JSON(http.StatusOK, gin.H{"data": players})
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

	c.JSON(http.StatusOK, gin.H{"data": true})
}
