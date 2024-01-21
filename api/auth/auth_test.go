package auth

import (
	"testing"
	"time"

	"github.com/jackmcguire1/UserService/dom/user"
	"github.com/stretchr/testify/assert"
)

var testToken = []byte("1234")

func TestSignClaims(t *testing.T) {

	h := &Handler{JWTSecret: testToken, Expiry: time.Minute}
	token, err := h.SignClaims(&user.User{
		ID:          "1234",
		FirstName:   "example",
		LastName:    "test",
		Email:       "test@example.com",
		NickName:    "test-user",
		CountryCode: "GB",
		Saved:       time.Now().Format(time.RFC3339),
		IsAdmin:     true,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestSignClaimsWithVerification(t *testing.T) {

	h := &Handler{JWTSecret: testToken, Expiry: time.Hour}
	token, err := h.SignClaims(&user.User{
		ID:          "1234",
		FirstName:   "example",
		LastName:    "test",
		Email:       "test@example.com",
		NickName:    "test-user",
		CountryCode: "GB",
		Saved:       time.Now().Format(time.RFC3339),
		IsAdmin:     true,
	})
	assert.NoError(t, err)

	claims, err := h.ValidateJWT(token)
	assert.NoError(t, err)
	assert.EqualValues(t, "1234", claims.Subject)
}
