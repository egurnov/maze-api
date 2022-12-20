package jwtservice

import (
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTService_GenerateToken(t *testing.T) {
	jwtService := New([]byte("secret"), time.Minute)

	tokenString, err := jwtService.GenerateToken(42)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	claims, err := jwtService.ValidateToken(tokenString)
	require.NoError(t, err)
	assert.Equal(t, int64(42), claims.UserID)

	_, err = jwtService.ValidateToken(tokenString + "a")
	require.EqualError(t, err, "signature is invalid")
}

func TestJWTService_TokenExpiration(t *testing.T) {
	jwtService := New([]byte("secret"), time.Minute)

	jwt.TimeFunc = func() time.Time {
		return time.Now().Add(time.Hour)
	}
	defer func() {
		jwt.TimeFunc = time.Now
	}()

	tokenString, err := jwtService.GenerateToken(42)
	require.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	_, err = jwtService.ValidateToken(tokenString)
	require.Error(t, err)
	require.True(t, strings.HasPrefix(err.Error(), "token is expired"))
}
