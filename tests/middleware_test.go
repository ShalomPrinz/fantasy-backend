package tests

import (
	testUtils "fantasy/database/tests/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyToken(t *testing.T) {
	beforeEach(t.FailNow)
	postUser(t.FailNow, nil)

	url := testUtils.Url{Path: "test-token"}
	var response any
	testUtils.GetWithToken(url, loginDetails, &response)

	assert.Equal(t,
		nil,
		response,
		"Should not return neither response nor error for verified token")
}

func TestVerifyToken_NoToken(t *testing.T) {
	beforeEach(t.FailNow)

	var response any
	testUtils.Get("test-token", &response)

	assert.Equal(t,
		map[string]any{"error": "id-token-missing"},
		response,
		"Should return missing token error if the request is for a token-protected route")
}

func TestVerifyToken_WrongToken(t *testing.T) {
	beforeEach(t.FailNow)

	url := testUtils.Url{Path: "test-token"}
	var response any
	testUtils.GetWithCustomToken(url, "some_wrong_token", &response)

	assert.Equal(t,
		map[string]any{"error": "id-token-invalid"},
		response,
		"Should return invalid token error for malformed token")
}
