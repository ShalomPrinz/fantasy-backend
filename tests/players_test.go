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

func TestGetPlayer_NotExists(t *testing.T) {
	var actualPlayer any
	err := utils.Get("players/fake_id", &actualPlayer)
	if err != nil {
		fmt.Println("Request failed for fake id player.")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"error": "not-found"},
		actualPlayer,
		"Should get not found error.")
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

func TestPostPlayer_NoData(t *testing.T) {
	var fakePlayer any
	var response any
	err := utils.Post("players", fakePlayer, &response)
	if err != nil {
		fmt.Println("Request failed for post player")
		t.FailNow()
	}

	assert.Contains(t,
		response,
		"error",
		"Should return error when posting player without data")
}
