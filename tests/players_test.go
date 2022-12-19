package tests

import (
	"fantasy/database/entities"
	testUtils "fantasy/database/tests/utils"
	"fantasy/database/utils"
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
)

func postPlayer(failTest func()) (map[string]any, string) {
	type postPlayerRes struct {
		PlayerId string
	}
	var response map[string]any
	err := testUtils.Post("players", player, &response)
	if err != nil {
		fmt.Println("Request failed for post player")
		failTest()
	}

	playerId := utils.MapToStruct[postPlayerRes](response).PlayerId
	return map[string]any{
		"firstName": "Lionel",
		"id":        playerId,
		"lastName":  "Messi",
		"role":      "ATT",
		"team":      "Barcelona",
	}, playerId
}

func resetData(failTest func()) {
	deleteDataUrl := "http://localhost:8090/emulator/v1/projects/demo-test-fantasy/databases/(default)/documents"
	err := testUtils.Delete(deleteDataUrl)
	if err != nil {
		fmt.Println("Reset data: Failed deleting data")
		failTest()
	}
}

func TestGetPlayer(t *testing.T) {
	playerSent, playerId := postPlayer(t.FailNow)

	var playerReceived any
	err := testUtils.Get("players/"+playerId, &playerReceived)
	if err != nil {
		fmt.Println("Request failed for single player")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"player": playerSent},
		playerReceived,
		"Should get player by id.")
}

func TestGetPlayer_NotExists(t *testing.T) {
	var actualPlayer any
	err := testUtils.Get("players/fake_id", &actualPlayer)
	if err != nil {
		fmt.Println("Request failed for fake id player.")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"error": "not-found"},
		actualPlayer,
		"Should get not found error.")
}

func TestNewPlayer(t *testing.T) {
	var response any
	err := testUtils.Post("players", player, &response)
	if err != nil {
		fmt.Println("Request failed for post player")
		t.FailNow()
	}

	assert.Contains(t,
		response,
		"playerId",
		"Should return the player ID in database")
}

func TestNewPlayer_NoData(t *testing.T) {
	var response any
	err := testUtils.Post("players", entities.AddPlayer{}, &response)
	if err != nil {
		fmt.Println("Request failed for post player")
		t.FailNow()
	}

	assert.Contains(t,
		response,
		"error",
		"Should return error when posting player without data")
}

func TestQueryPlayers_NoTerm(t *testing.T) {
	resetData(t.FailNow)

	var queryResult any
	err := testUtils.Get("players/query", &queryResult)
	if err != nil {
		fmt.Println("Empty players query failed")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"players": nil},
		queryResult,
		"Should get no players")
}

func TestQueryPlayers_TermExists(t *testing.T) {
	resetData(t.FailNow)
	storedPlayer, _ := postPlayer(t.FailNow)

	var queryResult any
	err := testUtils.Get("players/query?term=l", &queryResult)
	if err != nil {
		fmt.Println("Empty players query failed")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"players": []any{storedPlayer}},
		queryResult,
		"Should get all players matching the term")
}

func TestQueryPlayers_TermNotExists(t *testing.T) {
	resetData(t.FailNow)
	postPlayer(t.FailNow)

	var queryResult any
	err := testUtils.Get("players/query?term=P", &queryResult)
	if err != nil {
		fmt.Println("Empty players query failed")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"players": nil},
		queryResult,
		"Should get no players when no one matches the term")
}
