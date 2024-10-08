package tests

import (
	"fantasy/database/entities"
	testUtils "fantasy/database/tests/utils"
	"fantasy/database/utils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	newLeagueUrl          = "newleague"
	acceptLeagueInviteUrl = "accept/leagueinvite"
	rejectLeagueInviteUrl = "reject/leagueinvite"
)

var (
	league = entities.AddLeague{
		Name: "Test League",
	}
)

func getLeagueUrl(leagueId string) string {
	return "league?id=" + leagueId
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
	if err := testUtils.GetWithToken(getLeagueUrl(leagueId), loginDetails, &response); err != nil {
		fmt.Println("Request failed for get league info")
		failTest()
	}
}

func TestNewLeague(t *testing.T) {
	beforeEach(t.FailNow)

	response := postLeague(t.FailNow)

	assert.NotEmpty(t,
		response.LeagueId,
		"Should return the league ID in database")
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

	leagueId := postLeague(t.FailNow).LeagueId
	var response map[string]any
	getLeague(t.FailNow, leagueId, &response)

	assert.Contains(t,
		response,
		"league",
		"Should return league object with all the relevant data")
}

func TestGetLeagueInfo_NoData(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	var response map[string]any
	getLeague(t.FailNow, "", &response)

	assert.Equal(t,
		map[string]any{"error": "missing-request-data"},
		response,
		"Should return missing data error if no league id supplied")
}

func TestGetLeagueInfo_NotExists(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	var response map[string]any
	getLeague(t.FailNow, "fake_league_id", &response)

	assert.Equal(t,
		map[string]any{"error": "not-found"},
		response,
		"Should return not found error if the given id is not in the leagues collection")
}

func TestGetLeagueInfo_NotMember(t *testing.T) {
	beforeEach(t.FailNow)

	notInLeagueUser := entities.AddUser{
		FullName: "Other User",
		Username: "Someone",
		Email:    "other@user.test",
		Password: "usertest",
	}
	loginDetails := testUtils.LoginUser{
		Email:    notInLeagueUser.Email,
		Password: notInLeagueUser.Password,
	}

	postThisUser(t.FailNow, notInLeagueUser, nil)
	leagueId := postLeague(t.FailNow).LeagueId

	var response any
	if err := testUtils.GetWithToken(getLeagueUrl(leagueId), loginDetails, &response); err != nil {
		fmt.Println("Request failed for get league info")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"error": "not-league-member"},
		response,
		"Should return error if the logged user is not a member of the league")
}

func TestAcceptLeagueInvitation(t *testing.T) {
	beforeEach(t.FailNow)
	var res map[string]string
	sendLeagueInvitation(t.FailNow, &res)
	messageId := res["messageId"]

	var response any
	addMemberData := entities.LeagueInviteResponse{MessageId: messageId}
	testUtils.PostWithToken(acceptLeagueInviteUrl, addMemberData, invitedUserDetails, &response)

	assert.Equal(t,
		map[string]any{"status": "success"},
		response,
		"Should add other user as a regular member to the league user created")
}

func TestAcceptLeagueInvitation_NoData(t *testing.T) {
	beforeEach(t.FailNow)
	sendLeagueInvitation(t.FailNow, nil)

	var response any
	inviteResponseData := entities.LeagueInviteResponse{MessageId: ""}
	testUtils.PostWithToken(acceptLeagueInviteUrl, inviteResponseData, invitedUserDetails, &response)

	assert.Equal(t,
		map[string]any{"error": "missing-request-data"},
		response,
		"Should return error if no data is sent")
}

func TestAcceptLeagueInvitation_MessageNotExists(t *testing.T) {
	beforeEach(t.FailNow)
	sendLeagueInvitation(t.FailNow, nil)

	var response any
	inviteResponseData := entities.LeagueInviteResponse{MessageId: "fake_message_id"}
	testUtils.PostWithToken(acceptLeagueInviteUrl, inviteResponseData, invitedUserDetails, &response)

	assert.Equal(t,
		map[string]any{"error": "no-such-message"},
		response,
		"Should return error if no data is sent")
}

func TestAcceptLeagueInvitation_AcceptTwice(t *testing.T) {
	beforeEach(t.FailNow)
	var res map[string]string
	sendLeagueInvitation(t.FailNow, &res)
	messageId := res["messageId"]

	var response any
	inviteResponseData := entities.LeagueInviteResponse{MessageId: messageId}
	testUtils.PostWithToken(acceptLeagueInviteUrl, inviteResponseData, invitedUserDetails, &response)
	testUtils.PostWithToken(acceptLeagueInviteUrl, inviteResponseData, invitedUserDetails, &response)

	assert.Equal(t,
		map[string]any{"error": "no-such-message"},
		response,
		"Should remove message if already accepted")
}

func TestRejectLeagueInvitation(t *testing.T) {
	beforeEach(t.FailNow)
	var res map[string]string
	sendLeagueInvitation(t.FailNow, &res)
	messageId := res["messageId"]

	var response any
	inviteResponseData := entities.LeagueInviteResponse{MessageId: messageId}
	testUtils.PostWithToken(rejectLeagueInviteUrl, inviteResponseData, invitedUserDetails, &response)

	assert.Equal(t,
		map[string]any{"status": "success"},
		response,
		"Should reject the league invitation by a given message id")
}

func TestRejectLeagueInvitation_NoData(t *testing.T) {
	beforeEach(t.FailNow)
	sendLeagueInvitation(t.FailNow, nil)

	var response any
	inviteResponseData := entities.LeagueInviteResponse{MessageId: ""}
	testUtils.PostWithToken(rejectLeagueInviteUrl, inviteResponseData, invitedUserDetails, &response)

	assert.Equal(t,
		map[string]any{"error": "missing-request-data"},
		response,
		"Should return error if no data is sent")
}

func TestRejectLeagueInvitation_MessageNotExists(t *testing.T) {
	beforeEach(t.FailNow)
	sendLeagueInvitation(t.FailNow, nil)

	var response any
	inviteResponseData := entities.LeagueInviteResponse{MessageId: "fake_message_id"}
	testUtils.PostWithToken(rejectLeagueInviteUrl, inviteResponseData, invitedUserDetails, &response)

	assert.Equal(t,
		map[string]any{"error": "no-such-message"},
		response,
		"Should return error if no data is sent")
}

func TestRejectLeagueInvitation_AcceptTwice(t *testing.T) {
	beforeEach(t.FailNow)
	var res map[string]string
	sendLeagueInvitation(t.FailNow, &res)
	messageId := res["messageId"]

	var response any
	inviteResponseData := entities.LeagueInviteResponse{MessageId: messageId}
	testUtils.PostWithToken(rejectLeagueInviteUrl, inviteResponseData, invitedUserDetails, &response)
	testUtils.PostWithToken(rejectLeagueInviteUrl, inviteResponseData, invitedUserDetails, &response)

	assert.Equal(t,
		map[string]any{"error": "no-such-message"},
		response,
		"Should remove message if already rejected")
}
