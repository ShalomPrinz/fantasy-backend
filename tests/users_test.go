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
	getUserInfoUrl = testUtils.Url{
		Path: "user",
	}
	userAddPlayerUrl = testUtils.Url{
		Path: "user/addplayer",
	}
	userRemovePlayerUrl = testUtils.Url{
		Path: "user/removeplayer",
	}

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
	postThisUser(failTest, user, &response)
}

func postThisUser(failTest func(), user entities.AddUser, response any) {
	if err := testUtils.Post("register", user, &response); err != nil {
		fmt.Println("Request failed for register user")
		failTest()
	}
}

func getUser(failTest func()) entities.DetailedAccount {
	var response map[string]any
	if err := testUtils.GetWithToken(getUserInfoUrl, loginDetails, &response); err != nil {
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

	type TestCase struct {
		user          entities.AddUser
		expected      map[string]any
		description   string
		postUserTwice bool
	}

	tests := []TestCase{
		{
			user:        user,
			expected:    map[string]any{"status": "success"},
			description: "Success: Should register new user",
		},
		{
			user:        entities.AddUser{},
			expected:    map[string]any{"error": "missing-request-data"},
			description: "NoData: Should return error for user without data",
		},
		{
			user: entities.AddUser{
				FullName: "Test User",
				Nickname: "Testy",
				Email:    "bad_email",
				Password: "testtest",
			},
			expected:    map[string]any{"error": "invalid-email"},
			description: "InvalidEmail: Should return error for user with invalid email",
		},
		{
			user: entities.AddUser{
				FullName: "Test User",
				Nickname: "Testy",
				Email:    "test@test.test",
				Password: "short",
			},
			expected:    map[string]any{"error": "password-too-short"},
			description: "ShortPassword: Should return error for user with a password shorter than 6 letters",
		},
		{
			user:          user,
			expected:      map[string]any{"error": "email-already-exists"},
			description:   "EmailExists: Should return error for user with an email already exists",
			postUserTwice: true,
		},
	}

	for _, test := range tests {
		beforeEach(t.FailNow)

		var response any
		postThisUser(t.FailNow, test.user, &response)
		if test.postUserTwice {
			postThisUser(t.FailNow, test.user, &response)
		}

		assert.Equal(t,
			test.expected,
			response,
			test.description)
	}
}

func TestGetUserInfo(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	var response any
	if err := testUtils.GetWithToken(getUserInfoUrl, loginDetails, &response); err != nil {
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
	if err := testUtils.GetWithToken(getUserInfoUrl, fakeUserDetails, &response); err != nil {
		fmt.Println("Request failed for get user info")
		t.FailNow()
	}

	assert.Contains(t,
		response,
		"error",
		"Should return error for nonexistent user email")
}

func TestGetUserInfo_WrongPassword(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	fakeUserDetails := testUtils.LoginUser{
		Email:    user.Email,
		Password: "not_the_password",
	}

	var response any
	if err := testUtils.GetWithToken(getUserInfoUrl, fakeUserDetails, &response); err != nil {
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
	playerId := postPlayer(t.FailNow).PlayerId
	data := entities.Entity{
		ID: playerId,
	}

	playersNumBefore := len(getUser(t.FailNow).Team)
	testUtils.PostWithToken(userAddPlayerUrl, data, loginDetails, nil)
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
	testUtils.PostWithToken(userAddPlayerUrl, data, loginDetails, &response)

	assert.Contains(t,
		response,
		"error",
		"Should return error if the given id is not in the players collection")
}

func TestAddTeamPlayer_NoData(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	var response any
	testUtils.PostWithToken(userAddPlayerUrl, entities.Entity{}, loginDetails, &response)

	assert.Contains(t,
		response,
		"error",
		"Should return error if no player id supplied")
}

func TestRemoveTeamPlayer(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)
	playerId := postPlayer(t.FailNow).PlayerId
	data := entities.Entity{
		ID: playerId,
	}
	testUtils.PostWithToken(userAddPlayerUrl, data, loginDetails, nil)

	playersNumBefore := len(getUser(t.FailNow).Team)
	testUtils.PostWithToken(userRemovePlayerUrl, data, loginDetails, nil)
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
	testUtils.PostWithToken(userRemovePlayerUrl, data, loginDetails, &response)

	assert.Contains(t,
		response,
		"error",
		"Should return error if the given id is not in the players collection")
}

func TestRemoveTeamPlayer_NoData(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	var response any
	testUtils.PostWithToken(userRemovePlayerUrl, entities.Entity{}, loginDetails, &response)

	assert.Contains(t,
		response,
		"error",
		"Should return error if no player id supplied")
}
