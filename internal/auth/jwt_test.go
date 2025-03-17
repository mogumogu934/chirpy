package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTCreationAndValidation(t *testing.T) {
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

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	expiresIn := -time.Hour

	token, err := MakeJWT(userID, secret, expiresIn)
	assert.NoError(t, err)

	_, err = ValidateJWT(token, secret)
	assert.Error(t, err)
}

func TestWrongSecretJWT(t *testing.T) {
	userID := uuid.New()
	correctSecret := "test-correct-secret"
	wrongSecret := "test-wrong-secret"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, correctSecret, expiresIn)
	assert.NoError(t, err)

	_, err = ValidateJWT(token, wrongSecret)
	assert.Error(t, err)
}
