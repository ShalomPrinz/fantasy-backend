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
	newLeagueUrl = testUtils.Url{
		Path: "newleague",
	}

	league = entities.AddLeague{
		Name: "Test League",
	}
)

func getLeagueInfoUrl(leagueId string) testUtils.Url {
	return testUtils.Url{
		Path: "league",
		Params: []testUtils.UrlParameter{
			{
				Key:   "id",
				Value: leagueId,
			},
		},
	}
}

type postLeagueResponse struct {
	LeagueId string `json:"leagueId"`
}

func postLeague(failTest func()) postLeagueResponse {
	postUser(failTest, nil)

	var response map[string]any
	err := testUtils.PostWithToken(newLeagueUrl, league, loginDetails, &response)
	if err != nil {
		fmt.Println("Request failed for new league")
		failTest()
	}

	return utils.MapToStruct[postLeagueResponse](response)
}

func getLeague(failTest func(), leagueId string, response any) {
	actualLeagueId := postLeague(failTest).LeagueId
	var url testUtils.Url
	if leagueId == "use-actual" {
		url = getLeagueInfoUrl(actualLeagueId)
	} else {
		url = getLeagueInfoUrl(leagueId)
	}

	err := testUtils.GetWithToken(url, loginDetails, &response)
	if err != nil {
		fmt.Println("Request failed for get league info")
		failTest()
	}
}

func TestNewLeague(t *testing.T) {
	beforeEach(t.FailNow)

	response := postLeague(t.FailNow)

	assert.NotEmpty(t,
		response.LeagueId,
		"Should add new league and return its id")
}

func TestNewLeague_NoData(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	var response any
	testUtils.PostWithToken(newLeagueUrl, entities.AddLeague{}, loginDetails, &response)

	assert.Contains(t,
		response,
		"error",
		"Should return error when posting league without data")
}

func TestNewLeague_AddMember(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	userLeaguesNumBefore := len(getUser(t.FailNow).Leagues)
	var response any
	testUtils.PostWithToken(newLeagueUrl, league, loginDetails, &response)
	userLeaguesNumAfter := len(getUser(t.FailNow).Leagues)

	assert.Equal(t,
		userLeaguesNumBefore+1,
		userLeaguesNumAfter,
		"Should add the logged user to Test League members")
}

func TestGetLeagueInfo(t *testing.T) {
	beforeEach(t.FailNow)

	var response map[string]any
	getLeague(t.FailNow, "use-actual", &response)

	assert.Contains(t,
		response,
		"league",
		"Should return league object with all the relevant data")
}

func TestGetLeagueInfo_NoData(t *testing.T) {
	beforeEach(t.FailNow)

	var response map[string]any
	getLeague(t.FailNow, "", &response)

	assert.Contains(t,
		response,
		"error",
		"Should return error if no league id supplied")
}

func TestGetLeagueInfo_NotExists(t *testing.T) {
	beforeEach(t.FailNow)

	var response map[string]any
	getLeague(t.FailNow, "fake_league_id", &response)

	assert.Contains(t,
		response,
		"error",
		"Should return error if the given id is not in the leagues collection")
}

func TestGetLeagueInfo_NotMember(t *testing.T) {
	beforeEach(t.FailNow)

	postOtherUser(t.FailNow, nil)
	leagueId := postLeague(t.FailNow).LeagueId
	url := getLeagueInfoUrl(leagueId)

	var response any
	err := testUtils.GetWithToken(url, otherLoginDetails, &response)
	if err != nil {
		fmt.Println("Request failed for get league info")
		t.FailNow()
	}

	assert.Contains(t,
		response,
		"error",
		"Should return error if the logged user is not a member of the league")
}
