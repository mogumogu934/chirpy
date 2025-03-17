package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTCreationAndValidationSuccess(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	extractedID, err := ValidateJWT(token, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, extractedID)
}

func TestJWTValidationFailureWithExpiredJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := -time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	assert.NoError(t, err)

	_, err = ValidateJWT(token, secret)
	assert.Error(t, err)
}

func TestJWTValidationFailureWithWrongSecret(t *testing.T) {
	userID := uuid.New()
	correctSecret := "test-correct-secret"
	wrongSecret := "test-wrong-secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, correctSecret, expiresIn)
	assert.NoError(t, err)

	_, err = ValidateJWT(token, wrongSecret)
	assert.Error(t, err)
}

func TestGetBearerTokenSuccess(t *testing.T) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer test_token")
	token, err := GetBearerToken(headers)
	assert.Equal(t, "test_token", token)
	assert.NoError(t, err)
}

func TestGetBearerTokenFailureWithMissingHeader(t *testing.T) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	_, err := GetBearerToken(headers)
	assert.Error(t, err)
}

func TestGetBearerTokenFailureWithInvalidHeader(t *testing.T) {
	headers := http.Header{}
	headers.Add("Authorization", "This is my test_token")
	_, err := GetBearerToken(headers)
	assert.Error(t, err)
}

func TestGetBearerTokenSuccessWithExtraWhitespace(t *testing.T) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer         test_token")
	token, err := GetBearerToken(headers)
	assert.Equal(t, "test_token", token)
	assert.NoError(t, err)
}
