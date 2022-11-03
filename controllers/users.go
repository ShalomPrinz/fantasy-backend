package controllers

import (
	"fantasy/database/entities"
	"fantasy/database/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewUser(c *gin.Context) {
	var input entities.AddUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	UID, err := lib.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lib.InsertItemCustomID("accounts", UID, entities.AddAccount{
		Team: []entities.Player{},
	})

	c.JSON(http.StatusOK, gin.H{"addedAccount": true})
}
