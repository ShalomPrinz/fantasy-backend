package tests

import (
	testUtils "fantasy/database/tests/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testTokenUrl = "test-token"
)

func TestVerifyToken(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	var response any
	testUtils.GetWithToken(testTokenUrl, loginDetails, &response)

	assert.Equal(t,
		nil,
		response,
		"Should not return neither response nor error for verified token")
}

func TestVerifyToken_NoToken(t *testing.T) {
	beforeEach(t.FailNow)

	var response any
	testUtils.Get(testTokenUrl, &response)

	assert.Equal(t,
		map[string]any{"error": "id-token-missing"},
		response,
		"Should return missing token error if the request is for a token-protected route")
}

func TestVerifyToken_WrongToken(t *testing.T) {
	beforeEach(t.FailNow)

	var response any
	testUtils.GetWithCustomToken(testTokenUrl, "some_wrong_token", &response)

	assert.Equal(t,
		map[string]any{"error": "id-token-invalid"},
		response,
		"Should return invalid token error for malformed token")
}

func TestVerifyToken_UserNotExists(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)
	token, err := testUtils.GenerateIdToken(loginDetails)
	if err != nil {
		t.FailNow()
	}
	// clear all data, including user
	beforeEach(t.FailNow)

	var response any
	testUtils.GetWithCustomToken(testTokenUrl, token, &response)

	assert.Equal(t,
		map[string]any{"error": "user-not-exists"},
		response,
		"Should return user not exists error for nonexistent user")
}
