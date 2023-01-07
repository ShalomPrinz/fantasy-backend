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

type postPlayerResponse struct {
	PlayerId string `json:"playerId"`
}

func postPlayer(failTest func()) postPlayerResponse {
	return postThisPlayer(failTest, player)
}

func postThisPlayer(failTest func(), player entities.AddPlayer) postPlayerResponse {
	var response map[string]any
	if err := testUtils.Post("players", player, &response); err != nil {
		fmt.Println("Request failed for post player")
		failTest()
	}

	return utils.MapToStruct[postPlayerResponse](response)
}

func getPlayer(failTest func(), playerId string, response any) {
	if err := testUtils.Get("players/"+playerId, &response); err != nil {
		fmt.Println("Request failed for single player")
		failTest()
	}
}

func getPlayerById(failTest func(), playerId string) any {
	var getPlayerResponse map[string]any
	getPlayer(failTest, playerId, &getPlayerResponse)
	return getPlayerResponse["player"]
}

func TestGetPlayer(t *testing.T) {
	beforeEach(t.FailNow)
	playerId := postPlayer(t.FailNow).PlayerId

	var response any
	getPlayer(t.FailNow, playerId, &response)

	assert.Contains(t,
		response,
		"player",
		"Should return player object with all the relevant data")
}

func TestGetPlayer_NotExists(t *testing.T) {
	beforeEach(t.FailNow)

	var response any
	getPlayer(t.FailNow, "fake_id", &response)

	assert.Contains(t,
		response,
		"error",
		"Should return error for nonexistent player id")
}

func TestNewPlayer(t *testing.T) {
	beforeEach(t.FailNow)

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
	beforeEach(t.FailNow)

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
	beforeEach(t.FailNow)

	var queryResult any
	if err := testUtils.Get("players/query", &queryResult); err != nil {
		fmt.Println("Players query failed")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"players": []any{}},
		queryResult,
		"Should return empty players list")
}

func TestQueryPlayers_TermLetterCase(t *testing.T) {
	beforeEach(t.FailNow)
	postPlayer(t.FailNow)

	var lowerCaseTerm any
	if err := testUtils.Get("players/query?term=l", &lowerCaseTerm); err != nil {
		fmt.Println("Players query failed")
		t.FailNow()
	}
	var upperCaseTerm any
	if err := testUtils.Get("players/query?term=L", &upperCaseTerm); err != nil {
		fmt.Println("Players query failed")
		t.FailNow()
	}

	assert.Equal(t,
		lowerCaseTerm,
		upperCaseTerm,
		"Should get same search result for lower case term and upper case term")
}

func TestQueryPlayers_TermExists(t *testing.T) {
	beforeEach(t.FailNow)
	playerId := postPlayer(t.FailNow).PlayerId
	storedPlayer := getPlayerById(t.FailNow, playerId)

	var queryResult any
	if err := testUtils.Get("players/query?term=l", &queryResult); err != nil {
		fmt.Println("Players query failed")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"players": []any{storedPlayer}},
		queryResult,
		"Should get all players matching the term")
}

func TestQueryPlayers_TermNotExists(t *testing.T) {
	beforeEach(t.FailNow)
	postPlayer(t.FailNow)

	var queryResult any
	if err := testUtils.Get("players/query?term=P", &queryResult); err != nil {
		fmt.Println("Players query failed")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"players": []any{}},
		queryResult,
		"Should return empty players list if no player name matches the term")
}

func TestQueryPlayers_QueryLimit(t *testing.T) {
	beforeEach(t.FailNow)
	queryLimit := 10
	playersNum := 12
	for i := 0; i < playersNum; i++ {
		postPlayer(t.FailNow)
	}

	var queryResult map[string][]any
	if err := testUtils.Get("players/query?term=l", &queryResult); err != nil {
		fmt.Println("Players query failed")
		t.FailNow()
	}

	assert.Equal(t,
		queryLimit,
		len(queryResult["players"]),
		"Should return players list, limitted by the query limit")
}

func TestQueryPlayers_DuplicatedName(t *testing.T) {
	beforeEach(t.FailNow)
	duplicatedNamePlayer := entities.AddPlayer{
		Name: "Shalom Shalom",
		Role: "MID",
		Team: "Ajax",
	}
	postThisPlayer(t.FailNow, duplicatedNamePlayer)

	var queryResult map[string][]any
	if err := testUtils.Get("players/query?term=s", &queryResult); err != nil {
		fmt.Println("Players query failed")
		t.FailNow()
	}

	assert.Equal(t,
		1,
		len(queryResult),
		"Should return players list with no duplicates, even if first and last name matches the term")
}

func TestQueryPlayers_FirstNameLastName(t *testing.T) {
	beforeEach(t.FailNow)

	firstPlayerId := postThisPlayer(t.FailNow, entities.AddPlayer{
		Name: "Kyle Walker",
		Role: "DEF",
		Team: "Manchester City",
	}).PlayerId
	firstPlayer := getPlayerById(t.FailNow, firstPlayerId)

	secondPlayerId := postThisPlayer(t.FailNow, entities.AddPlayer{
		Name: "Harry Kane",
		Role: "ATT",
		Team: "Tottenham",
	}).PlayerId
	secondPlayer := getPlayerById(t.FailNow, secondPlayerId)

	var queryResult map[string][]any
	if err := testUtils.Get("players/query?term=K", &queryResult); err != nil {
		fmt.Println("Players query failed")
		t.FailNow()
	}

	assert.Equal(t,
		2,
		len(queryResult["players"]),
		"Should return players list with all players who their first or last name starts with the term")

	assert.Equal(t,
		firstPlayer,
		queryResult["players"][0],
		"First player should be the player who his first name start with the term")

	assert.Equal(t,
		secondPlayer,
		queryResult["players"][1],
		"Second player should be the player who his last name start with the term")
}

func TestQueryPlayers_FirstNameLastNameLimit(t *testing.T) {
	beforeEach(t.FailNow)
	queryLimit := 10
	playersNum := 12

	for i := 0; i < playersNum/2; i++ {
		postThisPlayer(t.FailNow, entities.AddPlayer{
			Name: "Kyle Walker",
			Role: "DEF",
			Team: "Manchester City",
		})
		postThisPlayer(t.FailNow, entities.AddPlayer{
			Name: "Harry Kane",
			Role: "ATT",
			Team: "Tottenham",
		})
	}

	var queryResult map[string][]any
	if err := testUtils.Get("players/query?term=K", &queryResult); err != nil {
		fmt.Println("Players query failed")
		t.FailNow()
	}

	assert.Equal(t,
		queryLimit,
		len(queryResult["players"]),
		"Should return players list limitted by queryLimit number of elements, and not all the players matching the term")
}
