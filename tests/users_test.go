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
	user = entities.AddUser{
		FullName: "Test User",
		Nickname: "Testy",
		Email:    "test@test.test",
		Password: "testtest",
	}

	loginDetails = testUtils.LoginUser{
		Email:    user.Email,
		Password: user.Password,
	}
)

func postUser(failTest func(), response any) {
	err := testUtils.Post("register", user, &response)
	if err != nil {
		fmt.Println("Request failed for register user")
		failTest()
	}
}

func getUser(failTest func()) entities.DetailedAccount {
	var response map[string]any
	err := testUtils.GetWithToken("user", loginDetails, &response)
	if err != nil {
		fmt.Println("Request failed for get user info")
		failTest()
	}

	type GetUserResponse struct {
		User entities.DetailedAccount `json:"user"`
	}

	return utils.MapToStruct[GetUserResponse](response).User
}

func TestRegisterUser(t *testing.T) {
	beforeEach(t.FailNow)

	var response any
	postUser(t.FailNow, &response)

	assert.Equal(t,
		map[string]any{"status": "success"},
		response,
		"Should register new user.")
}

func TestRegisterUser_NoData(t *testing.T) {
	beforeEach(t.FailNow)

	var response any
	err := testUtils.Post("register", entities.AddUser{}, &response)
	if err != nil {
		fmt.Println("Request failed for register user")
		t.FailNow()
	}

	assert.Contains(t,
		response,
		"error",
		"Should return error when posting user without data")
}

func TestGetUserInfo(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	var response any
	err := testUtils.GetWithToken("user", loginDetails, &response)
	if err != nil {
		fmt.Println("Request failed for get user info")
		t.FailNow()
	}

	assert.Contains(t,
		response,
		"user",
		"Should return user object with all the relevant data")
}

func TestGetUserInfo_NotExists(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	fakeUserDetails := testUtils.LoginUser{
		Email:    "wrong_email",
		Password: "some_pass",
	}

	var response any
	err := testUtils.GetWithToken("user", fakeUserDetails, &response)
	if err != nil {
		fmt.Println("Request failed for get user info")
		t.FailNow()
	}

	assert.Contains(t,
		response,
		"error",
		"Should return error for non-exist user")
}

func TestGetUserInfo_WrongPassword(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	fakeUserDetails := testUtils.LoginUser{
		Email:    user.Email,
		Password: "not_the_password",
	}

	var response any
	err := testUtils.GetWithToken("user", fakeUserDetails, &response)
	if err != nil {
		fmt.Println("Request failed for get user info")
		t.FailNow()
	}

	assert.Contains(t,
		response,
		"error",
		"Should return error if password isn't correct")
}

func TestAddTeamPlayer(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)
	_, playerId := postPlayer(t.FailNow)
	data := entities.Entity{
		ID: playerId,
	}

	playersNumBefore := len(getUser(t.FailNow).Team)
	testUtils.PostWithToken("user/addplayer", data, loginDetails, nil)
	playersNumAfter := len(getUser(t.FailNow).Team)

	assert.Equal(t,
		playersNumBefore+1,
		playersNumAfter,
		"Should add one player to user team")
}

func TestAddTeamPlayer_NotExists(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)
	data := entities.Entity{
		ID: "fake_player_id",
	}

	var response any
	testUtils.PostWithToken("user/addplayer", data, loginDetails, &response)

	assert.Contains(t,
		response,
		"error",
		"Should return error if the given id is not in the players collection")
}

func TestAddTeamPlayer_NoData(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	var response any
	testUtils.PostWithToken("user/addplayer", entities.Entity{}, loginDetails, &response)

	assert.Contains(t,
		response,
		"error",
		"Should return error if no player id supplied")
}

func TestRemoveTeamPlayer(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)
	_, playerId := postPlayer(t.FailNow)
	data := entities.Entity{
		ID: playerId,
	}
	testUtils.PostWithToken("user/addplayer", data, loginDetails, nil)

	playersNumBefore := len(getUser(t.FailNow).Team)
	testUtils.PostWithToken("user/removeplayer", data, loginDetails, nil)
	playersNumAfter := len(getUser(t.FailNow).Team)

	assert.Equal(t,
		playersNumBefore-1,
		playersNumAfter,
		"Should remove one player from user team")
}

func TestRemoveTeamPlayer_NotExists(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)
	data := entities.Entity{
		ID: "fake_player_id",
	}

	var response any
	testUtils.PostWithToken("user/removeplayer", data, loginDetails, &response)

	assert.Contains(t,
		response,
		"error",
		"Should return error if the given id is not in the players collection")
}

func TestRemoveTeamPlayer_NoData(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	var response any
	testUtils.PostWithToken("user/removeplayer", entities.Entity{}, loginDetails, &response)

	assert.Contains(t,
		response,
		"error",
		"Should return error if no player id supplied")
}
