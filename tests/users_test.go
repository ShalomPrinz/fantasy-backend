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
	getUserInfoUrl      = "user"
	userAddPlayerUrl    = "user/addplayer"
	userRemovePlayerUrl = "user/removeplayer"
)

var (
	user = entities.AddUser{
		FullName: "Test User",
		Username: "Testy",
		Email:    "test@test.test",
		Password: "testtest",
	}

	loginDetails = testUtils.LoginUser{
		Email:    user.Email,
		Password: user.Password,
	}
)

func queryUsersUrl(term string) string {
	if term == "" {
		return "users/query"
	} else {
		return "users/query?term=" + term
	}
}

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

	var response any
	postUser(t.FailNow, &response)

	assert.Contains(t,
		response,
		"userId",
		"Should register new user and return his id")
}

func TestRegisterUserErrors(t *testing.T) {
	type TestCase struct {
		user          entities.AddUser
		expected      map[string]any
		description   string
		postUserTwice bool
	}

	tests := []TestCase{
		{
			user:        entities.AddUser{},
			expected:    map[string]any{"error": "missing-request-data"},
			description: "NoData: Should return error for user without data",
		},
		{
			user: entities.AddUser{
				FullName: "Test User",
				Username: "Testy",
				Email:    "bad_email",
				Password: "testtest",
			},
			expected:    map[string]any{"error": "invalid-email"},
			description: "InvalidEmail: Should return error for user with invalid email",
		},
		{
			user: entities.AddUser{
				FullName: "Test User",
				Username: "Testy",
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

func TestQueryUsers_NoTerm(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)
	otherUser := entities.AddUser{
		FullName: "Other User",
		Username: user.Username + "2",
		Email:    "other@user.test",
		Password: "usertest",
	}
	postThisUser(t.FailNow, otherUser, nil)

	var queryResult any
	if err := testUtils.GetWithToken(queryUsersUrl(""), loginDetails, &queryResult); err != nil {
		fmt.Println("Users query failed")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"users": []any{}},
		queryResult,
		"Should return empty users list")
}

func TestQueryUsers_NoSelf(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	var queryResult any
	if err := testUtils.GetWithToken(queryUsersUrl(""), loginDetails, &queryResult); err != nil {
		fmt.Println("Users query failed")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"users": []any{}},
		queryResult,
		"Should return empty users list")
}

func TestQueryUsers_TermExists(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)
	storedUser := getUser(t.FailNow)
	otherUser := entities.AddUser{
		FullName: "Other User",
		Username: "Someone",
		Email:    "other@user.test",
		Password: "usertest",
	}
	otherLoginDetails := testUtils.LoginUser{
		Email:    otherUser.Email,
		Password: otherUser.Password,
	}
	postThisUser(t.FailNow, otherUser, nil)

	var queryResult map[string]any
	if err := testUtils.GetWithToken(queryUsersUrl(user.Username), otherLoginDetails, &queryResult); err != nil {
		fmt.Println("Users query failed")
		t.FailNow()
	}

	assert.Equal(t,
		[]any{
			map[string]any{
				"id":       storedUser.ID,
				"username": storedUser.Username,
			},
		},
		queryResult["users"],
		"Should return users list which contains otherUser query details")
}

func TestQueryUsers_TermNotExists(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	var queryResult any
	if err := testUtils.GetWithToken(queryUsersUrl("NoSuchTerm"), loginDetails, &queryResult); err != nil {
		fmt.Println("Users query failed")
		t.FailNow()
	}

	assert.Equal(t,
		map[string]any{"users": []any{}},
		queryResult,
		"Should return empty users list")
}

func TestQueryUsers_QueryLimit(t *testing.T) {
	beforeEach(t.FailNow)
	queryLimit := 3
	usersNum := 5
	postUser(t.FailNow, nil)
	usernameToQuery := "user" + user.Username
	for i := 0; i < usersNum-1; i++ {
		indexString := fmt.Sprintf("%d", i)
		user := entities.AddUser{
			FullName: user.FullName,
			Username: usernameToQuery + indexString,
			Email:    "user" + indexString + user.Email,
			Password: user.Password,
		}
		postThisUser(t.FailNow, user, nil)
	}

	var queryResult map[string][]any
	if err := testUtils.GetWithToken(queryUsersUrl(usernameToQuery), loginDetails, &queryResult); err != nil {
		fmt.Println("Users query failed")
		t.FailNow()
	}

	assert.Equal(t,
		queryLimit,
		len(queryResult["users"]),
		"Should return users list limitted by queryLimit")
}
