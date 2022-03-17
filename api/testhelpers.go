package api

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/require"
	"testing"
)

//  genJWTAuthToken generates jwtauth token from user id
func genJWTAuthToken(t *testing.T, userID uint64) string {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": userID})
	require.NoError(t, err)

	return tokenString
}
