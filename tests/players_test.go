package tests

import (
	"fantasy/database/entities"
	"fantasy/database/tests/utils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	player = entities.AddPlayer{
		Name: "Lionel Messi",
		Role: "ATT",
		Team: "Barcelona",
	}

	storedPlayer = map[string]any{
		"firstName": "Lionel",
		"id":        "EH7h6PVXQpA7WC1gA2jc",
		"lastName":  "Messi",
		"role":      "ATT",
		"team":      "Barcelona",
	}
)

func TestGetPlayer(t *testing.T) {
	var actualPlayer any
	err := utils.Get("players/EH7h6PVXQpA7WC1gA2jc", &actualPlayer)
	if err != nil {
		fmt.Println("Request failed for single player. Check player id")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"player": storedPlayer},
		actualPlayer,
		"Should get player by id.")
}

func TestPostPlayer(t *testing.T) {
	var response any
	err := utils.Post("players", player, &response)
	if err != nil {
		fmt.Println("Request failed for post player")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"addedPlayer": true},
		response,
		"Should respond with addedPlayer: true.")
}
