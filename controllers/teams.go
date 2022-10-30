package controllers

import (
	"fantasy/database/db"
	"fantasy/database/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTeams(c *gin.Context) {
	teams := db.GetAll("teams", entities.TeamAttributes[:])
	c.JSON(http.StatusOK, gin.H{"data": teams})
}

func NewTeam(c *gin.Context) {
	var input entities.AddTeam
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.InsertItem("teams", entities.AddTeam{
		Name: input.Name,
	})

	c.JSON(http.StatusOK, gin.H{"data": true})
}
