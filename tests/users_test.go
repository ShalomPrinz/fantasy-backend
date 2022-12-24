package tests

import (
	"fantasy/database/entities"
	testUtils "fantasy/database/tests/utils"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	beforeEach(t.FailNow)

	var response any
	var registerInfo = entities.AddUser{
		FullName: "Test User",
		Nickname: "Testy",
		Email:    "test@test.test",
		Password: "password",
	}

	err := testUtils.Post("register", registerInfo, &response)
	if err != nil {
		fmt.Println("Request failed for register user")
		t.FailNow()
	}

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
