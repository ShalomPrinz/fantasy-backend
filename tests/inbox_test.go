package tests

import (
	"fantasy/database/entities"
	testUtils "fantasy/database/tests/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	leagueInviteUrl = "leagueinvite"
)

var (
	invitedUser = entities.AddUser{
		FullName: "Other User",
		Username: "Other name",
		Email:    "other@other.test",
		Password: "otherpass",
	}

	invitedUserDetails = testUtils.LoginUser{
		Email:    "other@other.test",
		Password: "otherpass",
	}
)

func sendLeagueInvitation(failTest func(), response any) {
	leagueId, invitedUserId := prepareLeagueInvitation(failTest)

	message := entities.AddLeagueInvitation{
		To:       invitedUserId,
		LeagueId: leagueId,
	}
	testUtils.PostWithToken(leagueInviteUrl, message, loginDetails, &response)
}

func prepareLeagueInvitation(failTest func()) (string, string) {
	leagueId := postLeague(failTest).LeagueId

	var response map[string]string
	postThisUser(failTest, invitedUser, &response)
	invitedUserId := response["userId"]

	return leagueId, invitedUserId
}

func TestNewLeagueInvitation(t *testing.T) {
	beforeEach(t.FailNow)

	var response any
	sendLeagueInvitation(t.FailNow, &response)

	assert.Contains(t,
		response,
		"messageId",
		"Should send the message to the recipient and return the message id")
}

func TestNewLeagueInvitation_NoData(t *testing.T) {
	beforeEach(t.FailNow)
	prepareLeagueInvitation(t.FailNow)

	var response any
	testUtils.PostWithToken(leagueInviteUrl, entities.AddLeagueInvitation{}, loginDetails, &response)

	assert.Equal(t,
		map[string]any{"error": "missing-request-data"},
		response,
		"Should return missing data error if no data is supplied")
}
