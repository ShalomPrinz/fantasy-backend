package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTeams(c *gin.Context) {
	teams := lib.GetAll("teams", entities.TeamAttributes[:])
	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

func GetTeam(c *gin.Context) {
	team := lib.GetSingle("teams", c.Param("id"), entities.TeamAttributes[:])
	c.JSON(http.StatusOK, gin.H{"team": team})
}

func NewTeam(c *gin.Context) {
	var input entities.AddTeam
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For now Team only has ID. Later I will replace this function call
	lib.InsertItemCustomID("teams", input.ID, map[string]interface{}{})

	c.JSON(http.StatusOK, gin.H{"addedTeam": true})
}
