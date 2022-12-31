package tests

import (
	"fantasy/database/entities"
	testUtils "fantasy/database/tests/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareLeagueInvitation(failTest func()) (string, string) {
	leagueId := postLeague(failTest).LeagueId

	otherUser := entities.AddUser{
		FullName: "Other User",
		Username: "Other name",
		Email:    "other@other.test",
		Password: "otherpass",
	}

	var response map[string]string
	postThisUser(failTest, otherUser, &response)
	recipientId := response["userId"]

	return leagueId, recipientId
}

func TestNewLeagueInvitation(t *testing.T) {
	beforeEach(t.FailNow)
	leagueId, recipientId := prepareLeagueInvitation(t.FailNow)

	var response any
	url := testUtils.Url{Path: "leagueinvite"}
	message := entities.AddLeagueInvitation{
		To:       recipientId,
		LeagueId: leagueId,
	}
	testUtils.PostWithToken(url, message, loginDetails, &response)

	assert.Equal(t,
		map[string]any{"status": "success"},
		response,
		"Should send the message to the recipient")
}

func TestNewLeagueInvitation_NoData(t *testing.T) {
	beforeEach(t.FailNow)
	prepareLeagueInvitation(t.FailNow)

	var response any
	url := testUtils.Url{Path: "leagueinvite"}
	testUtils.PostWithToken(url, entities.AddLeagueInvitation{}, loginDetails, &response)

	assert.Equal(t,
		map[string]any{"error": "missing-request-data"},
		response,
		"Should return missing data error if no data is supplied")
}
