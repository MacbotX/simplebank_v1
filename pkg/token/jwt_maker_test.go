package token

import (
	"testing"
	"time"

	"github.com/MacbotX/simplebank_v1/util"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err :=  NewJWTMaker(util.RandomString(36))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredJWTToken(t *testing.T)  {
	maker, err :=  NewJWTMaker(util.RandomString(36))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)

}

func TestInvalideJWTTokenALgNone(t *testing.T)  {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err :=  NewJWTMaker(util.RandomString(36))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalideToken.Error())
	require.Nil(t, payload)

}

func TestTamperedJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(36))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	// Tamper token by changing a character
	tamperedToken := token[:len(token)-1] + "X"

	payload, err := maker.VerifyToken(tamperedToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalideToken.Error())
	require.Nil(t, payload)
}

