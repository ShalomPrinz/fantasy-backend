package tests

import (
	"fantasy/database/entities"
	testUtils "fantasy/database/tests/utils"
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
